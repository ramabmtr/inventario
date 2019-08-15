package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Environment struct {
	key      string
	value    string
	required bool
	oneOf    []string
}

func GetEnv(key string) *Environment {
	env := Environment{
		key:      key,
		required: false,
	}

	if value, ok := os.LookupEnv(key); ok {
		env.value = value
	}

	return &env
}

func (e *Environment) OneOf(values ...string) *Environment {
	e.oneOf = values
	return e
}

// if env var mark as required, will fail to parse if value is empty
func (e *Environment) Required() *Environment {
	e.required = true
	return e
}

func (e *Environment) SetDefault(defaultValue string) *Environment {
	if e.value != "" {
		return e
	}

	e.value = defaultValue
	return e
}

func (e *Environment) checkVal() {
	if e.required && e.value == "" {
		log.Fatal("\"", e.key, "\" marked as required but got an empty value")
	}
}

func (e *Environment) checkErr(err error) {
	if err != nil {
		log.Fatal("fail to parse \"", e.key, "\". Error message: ", err.Error())
	}
}

func (e *Environment) checkOneOf() {
	if len(e.oneOf) == 0 {
		return
	}

	found := false

	for _, v := range e.oneOf {
		if e.value == v {
			found = true
			break
		}
	}

	if !found {
		log.Fatal("\"", e.key, "\" must be one of ", e.oneOf)
	}
}

func (e *Environment) ToString() string {
	e.checkVal()
	e.checkOneOf()
	return e.value
}

func (e *Environment) ToInt() int {
	e.checkVal()
	e.checkOneOf()
	v, err := strconv.Atoi(e.value)
	e.checkErr(err)
	return v
}

func (e *Environment) ToFloat64() float64 {
	e.checkVal()
	e.checkOneOf()
	v, err := strconv.ParseFloat(e.value, 64)
	e.checkErr(err)
	return v
}

func (e *Environment) ToBool() bool {
	e.checkVal()
	e.checkOneOf()
	v, err := strconv.ParseBool(e.value)
	e.checkErr(err)
	return v
}

func (e *Environment) ToDuration() time.Duration {
	e.checkVal()
	e.checkOneOf()
	var v time.Duration
	var err error

	switch {
	case strings.ContainsAny(e.value, "nsuµmh"):
		v, err = time.ParseDuration(e.value)
	default:
		v, err = time.ParseDuration(e.value + "ms")
	}

	e.checkErr(err)

	return v
}
