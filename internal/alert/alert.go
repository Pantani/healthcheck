package alert

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/Pantani/healthcheck/internal/config"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"strings"
)

func SendEvent(namespace, name, path string) error {
	logParams := logger.Params{"namespace": namespace, "name": name, "path": path}
	client := pagerduty.NewClient(config.Configuration.PagerDuty.Key)
	incident, err := client.CreateIncident("Health Check Application", &pagerduty.CreateIncidentOptions{
		Urgency:     "high",
		Type:        "incident",
		IncidentKey: namespace + "_" + name,
		Title:       namespace + " - " + name,
		Service: &pagerduty.APIReference{
			ID:   config.Configuration.PagerDuty.Service,
			Type: "service",
		},
		Body: &pagerduty.APIDetails{
			Type:    name,
			Details: getDescription(namespace, name, path),
		},
		EscalationPolicy: &pagerduty.APIReference{
			ID:   config.Configuration.PagerDuty.Escalation_Policy,
			Type: "escalation_policy",
		},
	})
	if err != nil {
		return errors.E(err, "cannot create PagerDuty event", logParams)
	}
	logger.Info("PagerDuty incident created", logParams,
		logger.Params{"id": incident.ID, "number": incident.IncidentNumber, "url": incident.HTMLURL})
	return nil
}

func getDescription(namespace, name, path string) string {
	s := fmt.Sprintf("%s test %s %s", namespace, name, path)
	return strings.Join(strings.Fields(s), " ")
}
