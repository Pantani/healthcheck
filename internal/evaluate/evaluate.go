package evaluate

import (
	"github.com/antonmedv/expr"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
)

const (
	lastValueKey = "lastValue"
	newValueKey  = "newValue"
)

func Evaluate(exp string, lastValue, newValue interface{}) (bool, error) {
	var environment = map[string]interface{}{
		lastValueKey: lastValue,
		newValueKey:  newValue,
	}
	logParams := logger.Params{"exp": exp, "lastValue": lastValue, "newValue": newValue}
	program, err := expr.Compile(exp, expr.Env(environment))
	if err != nil {
		return false, errors.E(err, "cannot compile the expression", logParams)
	}

	output, err := expr.Run(program, environment)
	if err != nil {
		return false, errors.E(err, "cannot run the expression", logParams)
	}
	result, ok := output.(bool)
	if !ok {
		return false, errors.E(err, "evaluate result is not a boolean", logParams)
	}
	return result, nil
}
