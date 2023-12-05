package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/micro/micro/v3/service/logger"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

type Config struct {
	GO_MICRO_SERVICE_PORT         string `yaml:"GO_MICRO_SERVICE_PORT",envconfig:"GO_MICRO_SERVICE_PORT"`
	GO_MICRO_SERVICE_REGISTRY_URL string `yaml:"GO_MICRO_SERVICE_REGISTRY_URL",envconfig:"GO_MICRO_SERVICE_REGISTRY_URL"`
	GO_MICRO_RENEWALINTERVALINSEC string `yaml:"GO_MICRO_RENEWALINTERVALINSEC",envconfig:"GO_MICRO_RENEWALINTERVALINSEC"`
	GO_MICRO_DURATIONINSECS       string `yaml:"GO_MICRO_DURATIONINSECS",envconfig:"GO_MICRO_DURATIONINSECS"`
	GO_MICRO_KEYCLOAK_URL         string `yaml:"GO_MICRO_KEYCLOAK_URL",envconfig:"GO_MICRO_KEYCLOAK_URL"`
	GO_MICRO_CLIENT_ID            string `yaml:"GO_MICRO_CLIENT_ID",envconfig:"GO_MICRO_CLIENT_ID"`
	GO_MICRO_CLIENT_SECRET        string `yaml:"GO_MICRO_CLIENT_SECRET",envconfig:"GO_MICRO_CLIENT_SECRET"`
	GO_MICRO_REALM_NAME           string `yaml:"GO_MICRO_REALM_NAME",envconfig:"GO_MICRO_REALM_NAME"`
	GO_MICRO_DB_URL               string `yaml:"GO_MICRO_DB_URL",envconfig:"GO_MICRO_DB_URL"`
	GO_MICRO_MESSAGE_BROKER       string `yaml:"GO_MICRO_MESSAGE_BROKER",envconfig:"GO_MICRO_MESSAGE_BROKER"`
	GO_MICRO_MONGODB_URL          string `yaml:"GO_MICRO_MONGODB_URL",envconfig:"GO_MICRO_MONGODB_URL"`
	GO_MICRO_ACTIVE_PROFILE       string `yaml:"GO_MICRO_ACTIVE_PROFILE",envconfig:"GO_MICRO_ACTIVE_PROFILE"`
}

var cfg Config

func Setconfig() {
	profile := "dev"
	profiles := map[string]bool{"dev": true, "prod": true}
	if profiles[os.Getenv("GO_MICRO_ACTIVE_PROFILE")] {
		profile = os.Getenv("GO_MICRO_ACTIVE_PROFILE")
	}
	cfg.GO_MICRO_ACTIVE_PROFILE = profile
	readFile(&cfg, "resources/app.yaml")
	readFile(&cfg, "resources/"+profile+"-config.yaml")
	readEnv(&cfg)
}

func readFile(cfg *Config, path string) {
	f, err := os.Open(path)
	if err != nil {
		logger.Errorf(err.Error())
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		logger.Errorf(err.Error())
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		logger.Errorf(err.Error())
	}
}

func GetVal(key string) string {
	r := reflect.ValueOf(cfg)
	f := reflect.Indirect(r).FieldByName(key)
	return string(f.String())
}

func SetVal(key, value string) {
	r := reflect.ValueOf(&cfg).Elem()
	f := r.FieldByName(key)
	if f.IsValid() && f.CanSet() {
		f.SetString(value)
	} else {
		logger.Errorf("Invalid or unsettable field: %s", key)
	}
}
