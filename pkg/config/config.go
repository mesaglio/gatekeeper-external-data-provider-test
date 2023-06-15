package config

import (
	"mesaglio/gatekeeper-external-data-provider-test/pkg/validators"
	"mesaglio/gatekeeper-external-data-provider-test/pkg/validators/naming"
)

type Config struct {
	Validators []validators.Validator
}

func Get() {
	config := Config{}
	nv := naming.NamingValidator{}
	config.Validators = append(config.Validators, nv)
}
