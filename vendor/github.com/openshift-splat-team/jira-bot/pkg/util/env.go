package util

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var bindEnvVars = []string{"personal_access_token", "project", "board"}

func CheckForMissingEnvVars() error {
	for _, envVar := range bindEnvVars {
		if len(viper.GetString(envVar)) == 0 {
			return fmt.Errorf("the environment variable: [%s] must be exported", strings.ToUpper(fmt.Sprintf("jira_%s", envVar)))
		}
	}
	return nil
}

func BindEnvVars() error {
	viper.SetEnvPrefix("jira") // Set a prefix for environment variables
	for _, envVar := range bindEnvVars {
		err := viper.BindEnv(envVar)
		if err != nil {
			return fmt.Errorf("unable to bind env var %s: %v", envVar, err)
		}
	}
	viper.AutomaticEnv() // Automatically read environment variables
	return nil
}
