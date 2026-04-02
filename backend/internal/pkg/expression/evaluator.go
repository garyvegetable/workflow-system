package expression

import (
	"github.com/Knetic/govaluate"
)

type Evaluator struct {
	expression *govaluate.EvaluableExpression
}

func NewEvaluator(expressionStr string) (*Evaluator, error) {
	expression, err := govaluate.NewEvaluableExpression(expressionStr)
	if err != nil {
		return nil, err
	}
	return &Evaluator{expression: expression}, nil
}

func (e *Evaluator) Evaluate(params map[string]interface{}) (bool, error) {
	result, err := e.expression.Evaluate(params)
	if err != nil {
		return false, err
	}
	if b, ok := result.(bool); ok {
		return b, nil
	}
	return false, nil
}
