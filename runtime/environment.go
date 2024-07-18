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
	env := &Environment{
		variables: make(map[string]RuntimeValue),
		parent:    p,
	}

	if p == nil {
		env.init()
	}

	return env
}

func (env *Environment) init() {
	env.Set("null", NewNumberValue(0))
	env.Set("true", NewNumberValue(1))
	env.Set("false", NewNumberValue(0))
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
