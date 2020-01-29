package config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"log"
	"strings"
)

type configuration struct {
	Redis struct {
		URL string
	}
	PagerDuty struct {
		Key               string
		Service           string
		Escalation_Policy string
	}
}

var Configuration configuration

// set dummy values to force viper to search for these keys in environment variables
// the AutomaticEnv() only searches for already defined keys in a config file, default values or kvstore struct.
func setDefaults() {
	viper.SetDefault("Redis.URL", "redis://localhost:6379")
	viper.SetDefault("PagerDuty.Key", "")
	viper.SetDefault("PagerDuty.Service", "")
	viper.SetDefault("PagerDuty.Escalation_Policy", "")
}

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	setDefaults()
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(&Configuration); err != nil {
		logger.Error(err, "Error Unmarshal Viper Config File")
	}
	log.Printf("REDIS_URL: %s", Configuration.Redis.URL)
	log.Printf("PAGERDUTY_KEY: %s", Configuration.PagerDuty.Key)
	log.Printf("PAGERDUTY_SERVICE: %s", Configuration.PagerDuty.Service)
	log.Printf("PAGERDUTY_ESCALATION_POLICY: %s", Configuration.PagerDuty.Escalation_Policy)
}
