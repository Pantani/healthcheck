package alert

import (
	"fmt"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/Pantani/healthcheck/internal/config"
	"log/slog"
	"strings"
)

type PagerDuty struct {
	client           *pagerduty.Client
	service          string
	escalationPolicy string
}

func NewPagerDuty(cfg config.PagerDuty) (*PagerDuty, error) {
	if cfg.Key == "" {
		return nil, fmt.Errorf("PAGERDUTY_KEY is required")
	}
	if cfg.Service == "" {
		return nil, fmt.Errorf("PAGERDUTY_SERVICE is required")
	}
	if cfg.EscalationPolicy == "" {
		return nil, fmt.Errorf("PAGERDUTY_ESCALATION_POLICY is required")
	}
	return &PagerDuty{
		client:           pagerduty.NewClient(cfg.Key),
		service:          cfg.Service,
		escalationPolicy: cfg.EscalationPolicy,
	}, nil
}

func SendEvent(namespace, name, path string) error {
	pagerDuty, err := NewPagerDuty(config.Active.PagerDuty)
	if err != nil {
		return err
	}
	return pagerDuty.SendEvent(namespace, name, path)
}

func (p *PagerDuty) SendEvent(namespace, name, path string) error {
	incident, err := p.client.CreateIncident("Health Check Application", &pagerduty.CreateIncidentOptions{
		Urgency:     "high",
		Type:        "incident",
		IncidentKey: namespace + "_" + name,
		Title:       namespace + " - " + name,
		Service: &pagerduty.APIReference{
			ID:   p.service,
			Type: "service",
		},
		Body: &pagerduty.APIDetails{
			Type:    name,
			Details: getDescription(namespace, name, path),
		},
		EscalationPolicy: &pagerduty.APIReference{
			ID:   p.escalationPolicy,
			Type: "escalation_policy",
		},
	})
	if err != nil {
		return fmt.Errorf("create PagerDuty incident for %s.%s: %w", namespace, name, err)
	}
	slog.Info("PagerDuty incident created", "namespace", namespace, "name", name, "id", incident.ID, "number", incident.IncidentNumber, "url", incident.HTMLURL)
	return nil
}

func getDescription(namespace, name, path string) string {
	s := fmt.Sprintf("%s test %s %s", namespace, name, path)
	return strings.Join(strings.Fields(s), " ")
}
