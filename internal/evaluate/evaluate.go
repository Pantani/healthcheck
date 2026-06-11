package evaluate

import (
	"fmt"

	"github.com/expr-lang/expr"
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
	program, err := expr.Compile(exp, expr.Env(environment))
	if err != nil {
		return false, fmt.Errorf("compile expression %q: %w", exp, err)
	}

	output, err := expr.Run(program, environment)
	if err != nil {
		return false, fmt.Errorf("run expression %q: %w", exp, err)
	}
	result, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("expression %q returned %T, want bool", exp, output)
	}
	return result, nil
}
