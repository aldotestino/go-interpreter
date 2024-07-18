package runtime

import (
	"fmt"
	"go-interpreter/shared"
)

type Environment struct {
	variables map[string]RuntimeValue
	parent    *Environment
}

func NewEnvironment(p *Environment) *Environment {
	return &Environment{
		variables: make(map[string]RuntimeValue),
		parent:    p,
	}
}

func (env *Environment) Get(varName string) (RuntimeValue, error) {
	if v, found := env.variables[varName]; found {
		return v, nil
	} else if env.parent != nil {
		return env.parent.Get(varName)
	}

	return nil, shared.RuntimeError(fmt.Sprintf("'%s' is not defined", varName))
}

func (env *Environment) Set(varName string, value RuntimeValue) {
	env.variables[varName] = value
}

func (env *Environment) Unset(varName string) {
	delete(env.variables, varName)
}
