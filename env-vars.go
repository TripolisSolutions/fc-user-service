package main

import (
	"github.com/kelseyhightower/envconfig"
)

// envSettings struct
type envSettings struct {
	EventUsersExchange     string `envconfig:"event_users" default:"event.users"`
	EventUsersCreatedQueue string `envconfig:"event_users" default:"event.users.created"`
	EventUsersUpdatedQueue string `envconfig:"event_users" default:"event.users.updated"`
	EventUsersDeletedQueue string `envconfig:"event_users" default:"event.users.deleted"`
	Buffer                 int    `envconfig:"buffer" default:10`
}

// ProjectEnvSettings variable
var ProjectEnvSettings envSettings

// readEnvironmentVariables ...
func (settings *envSettings) readEnvironmentVariables() error {
	err := envconfig.Process("settings", settings)
	if err != nil {
		return err
	}
	return nil
}

// EnvSettingsInit is ..
//Always use this function to initialize a Settings struct
func EnvSettingsInit() {
	ProjectEnvSettings.readEnvironmentVariables()
}
