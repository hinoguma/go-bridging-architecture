package crosscutting

import (
	"errors"
	"fmt"
)

type deferFunc struct {
	name     string
	function func() error
}

type DeferFunctions struct {
	functions []deferFunc
}

func NewDeferFunctions() DeferFunctions {
	return DeferFunctions{
		functions: make([]deferFunc, 0),
	}
}

func (model *DeferFunctions) Append(name string, f func() error) *DeferFunctions {
	model.functions = append(model.functions, deferFunc{
		name:     name,
		function: f,
	})
	return model
}

func (model DeferFunctions) Execute() error {
	var err error
	for _, df := range model.functions {
		tmpErr := df.function()
		if tmpErr != nil {
			tmpErr = errors.Join(
				tmpErr, fmt.Errorf("defer function name:%s", df.name),
			)
			if err == nil {
				err = tmpErr
				continue
			}
			err = errors.Join(err, tmpErr)
			continue
		}
	}
	return err
}
