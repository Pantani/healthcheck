package collector

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Pantani/healthcheck/internal/alert"
	"github.com/Pantani/healthcheck/internal/client"
	"github.com/Pantani/healthcheck/internal/config"
	"github.com/Pantani/healthcheck/internal/database"
	"github.com/Pantani/healthcheck/internal/evaluate"
	"github.com/Pantani/healthcheck/internal/fixtures"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
)

type Requester interface {
	Execute(ctx context.Context, method string, path string, body interface{}) (string, error)
}

type Store interface {
	SaveData(ctx context.Context, entity, key string, value interface{}) error
	GetData(ctx context.Context, entity, key string, result interface{}) error
}

type AlertSender interface {
	SendEvent(namespace, name, path string) error
}

type AlertFunc func(namespace, name, path string) error

func (f AlertFunc) SendEvent(namespace, name, path string) error {
	return f(namespace, name, path)
}

type Runner struct {
	Store Store
	Alert AlertSender
}

func MetricsCollector(ctx context.Context, cfg config.Configuration) error {
	if err := cfg.ValidateRuntime(); err != nil {
		return err
	}

	db, err := database.Init(ctx, cfg.Redis.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	fxs, err := fixtures.GetFixtures(cfg.Fixtures.Path)
	if err != nil {
		return err
	}

	pagerDuty, err := alert.NewPagerDuty(cfg.PagerDuty)
	if err != nil {
		return err
	}

	c := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	runner := Runner{Store: db, Alert: pagerDuty}
	for _, f := range fxs {
		cl := client.InitClient(f.Host, cfg.HTTP.Timeout)
		for _, t := range f.Tests {
			spec := fmt.Sprintf("@every %s", t.UpdateTime)
			if err := AddFunc(ctx, c, runner, f.Namespace, spec, t, &cl); err != nil {
				return err
			}
		}
	}

	c.Start()
	slog.Info("healthcheck scheduler started", "fixtures", cfg.Fixtures.Path, "checks", fxs.CheckCount())
	<-ctx.Done()
	stopCtx := c.Stop()
	<-stopCtx.Done()
	slog.Info("healthcheck scheduler stopped")
	return nil
}

func AddFunc(ctx context.Context, c *cron.Cron, runner Runner, namespace, spec string, t fixtures.Test, cl Requester) error {
	_, err := c.AddFunc(spec, func() {
		runner.Collect(ctx, namespace, t, cl)
	})
	return err
}

func (r Runner) Collect(ctx context.Context, namespace string, t fixtures.Test, c Requester) {
	result, err := c.Execute(ctx, t.Method, t.URLPath, t.Body)
	if err != nil {
		r.sendAlert(namespace, t, err)
		return
	}

	value := gjson.Get(result, t.JSONPath)
	if !value.Exists() {
		slog.Error("json path does not exist", "namespace", namespace, "name", t.Name, "json_path", t.JSONPath)
		return
	}

	var lastData interface{}
	if err := r.Store.GetData(ctx, namespace, t.Name, &lastData); err != nil {
		lastData = 0
		slog.Info("previous value unavailable; using zero", "namespace", namespace, "name", t.Name, "error", err)
	}

	if err := r.Store.SaveData(ctx, namespace, t.Name, value.Value()); err != nil {
		slog.Error("failed to save data", "namespace", namespace, "name", t.Name, "error", err)
		return
	}

	passed, err := evaluate.Evaluate(t.Expression, lastData, value.Value())
	if err != nil {
		slog.Error("failed to evaluate expression", "namespace", namespace, "name", t.Name, "expression", t.Expression, "error", err)
		return
	}
	if passed {
		slog.Info("health check passed", "namespace", namespace, "name", t.Name, "value", value.Value())
		return
	}

	r.sendAlert(namespace, t, fmt.Errorf("expression evaluated to false"))
}

func (r Runner) sendAlert(namespace string, t fixtures.Test, cause error) {
	if err := r.Alert.SendEvent(namespace, t.Name, t.URLPath); err != nil {
		slog.Error("failed to create incident", "namespace", namespace, "name", t.Name, "cause", cause, "error", err)
		return
	}
	slog.Warn("health check failed; incident created", "namespace", namespace, "name", t.Name, "cause", cause)
}
