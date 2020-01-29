package collector

import (
	"fmt"
	"github.com/Pantani/healthcheck/internal/alert"
	"github.com/Pantani/healthcheck/internal/client"
	"github.com/Pantani/healthcheck/internal/config"
	"github.com/Pantani/healthcheck/internal/database"
	"github.com/Pantani/healthcheck/internal/evaluate"
	"github.com/Pantani/healthcheck/internal/fixtures"
	"github.com/robfig/cron"
	"github.com/tidwall/gjson"
	"github.com/trustwallet/blockatlas/pkg/logger"
)

func MetricsCollector() {
	db, err := database.Init(config.Configuration.Redis.Url)
	if err != nil {
		logger.Fatal(err)
	}
	c := cron.New()
	fxs, err := fixtures.GeFixtures()
	if err != nil {
		logger.Fatal(err)
	}
	for _, f := range fxs {
		cl := client.InitClient(f.Host)
		for _, t := range f.Tests {
			spec := fmt.Sprintf("@every %s", t.UpdateTime)
			err := AddFunc(c, f.Name, spec, t, cl, db)
			if err != nil {
				logger.Error(err)
			}
		}
	}
	c.Start()
	defer c.Stop()
	<-make(chan bool)
}

func AddFunc(c *cron.Cron, namespace, spec string, t fixtures.Test, cl client.Request, db *database.Database) error {
	err := c.AddFunc(spec, func() {
		collect(namespace, t, cl, db)
	})
	return err
}

func collect(namespace string, t fixtures.Test, c client.Request, db *database.Database) {
	result, err := c.Execute(t.Method, t.UrlPath, t.Body)
	if err != nil {
		return
	}

	value := gjson.Get(result, t.JsonPath)
	if !value.Exists() {
		logger.Error("gjson path doesn't exist", logger.Params{"namespace": namespace, "name": t.Name})
		return
	}
	logParams := logger.Params{"namespace": namespace, "name": t.Name, "value": value.Value()}

	var lastData interface{}
	err = db.GetData(namespace, t.Name, &lastData)
	if err != nil {
		lastData = 0
		logger.Error(err, logParams)
	}

	err = db.SaveData(namespace, t.Name, value.Value())
	if err != nil {
		logger.Error(err, "failed to save data", logParams)
		return
	}

	ev, err := evaluate.Evaluate(t.Expression, lastData, value.Value())
	if err != nil {
		logger.Error(err, logParams)
		return
	}
	if ev {
		return
	}

	err = alert.SendEvent(namespace, t.Name, t.UrlPath)
	if err != nil {
		logger.Error(err, logParams)
		return
	}
}
