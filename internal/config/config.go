package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jessevdk/go-flags"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/viper"

	pb_core "github.com/pocketbase/pocketbase/core"

	pb_models "github.com/pocketbase/pocketbase/models"

	pb_security "github.com/pocketbase/pocketbase/tools/security"
)

type Config struct {
	SMTP struct {
		ServerHost string `mapstructure:"server_host"`
		Port       int    `mapstructure:"port"`
		Username   string `mapstructure:"user_name"`
		Password   string `mapstructure:"password"`
	} `mapstructure:"smtp"`
	Redis struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"redis"`
	ReCAPTCHASecretKey string `mapstructure:"recaptcha_secret_key"`
	S3                 struct {
		AccessKey string `mapstructure:"access_key"`
		Secret    string `mapstructure:"secret"`
		Bucket    string `mapstructure:"bucket"`
		Region    string `mapstructure:"region"`
		Endpoint  string `mapstructure:"endpoint"`
	} `mapstructure:"s3"`
	Environment string
}

func (c Config) IsDev() bool {
	return c.Environment == "development"
}

func (c Config) IsProd() bool {
	return c.Environment == "production"
}

func (c Config) IsStaging() bool {
	return c.Environment == "staging"
}

var AppConfig *Config

func init() {
	var cli struct {
		Config      string `short:"c" long:"config" description:"config file"`
		Environment string `short:"e" long:"environment" description:"app environment" default:"development"`
	}

	parser := flags.NewParser(&cli, flags.Default)
	if _, err := parser.Parse(); err != nil {
		log.Panicln(err)
	}

	if cli.Config == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Panicln(err)
		}

		cli.Config = path.Join(dir, "..", "internal/config/config.yaml")
	}

	viper.SetConfigFile(cli.Config)
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	var configs map[string]Config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error reading config file %s", err)
	}

	err := viper.Unmarshal(&configs)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	config := configs[cli.Environment]
	config.Environment = cli.Environment
	AppConfig = &config
}

func overwriteSMTPConfig(setting *pb_core.SmtpConfig) {
	setting.Enabled = true
	setting.Host = AppConfig.SMTP.ServerHost
	setting.Password = AppConfig.SMTP.Password
	setting.Port = AppConfig.SMTP.Port
	setting.Username = AppConfig.SMTP.Username
	setting.Tls = false
}

func overwriteS3Config(setting *pb_core.S3Config) {
	setting.Enabled = true
	setting.AccessKey = AppConfig.S3.AccessKey
	setting.Secret = AppConfig.S3.Secret
	setting.Bucket = AppConfig.S3.Bucket
	setting.Region = AppConfig.S3.Region
	setting.Endpoint = AppConfig.S3.Endpoint
}

// OverwritePocketBaseConfig load setting values from yaml file and save to _params table
func OverwritePocketBaseConfig(app *pocketbase.PocketBase) {
	paramSettings, err := app.Dao().FindParamByKey(pb_models.ParamAppSettings)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error retrieving settings %s\n", err)
	}

	encryptionKey := os.Getenv(app.EncryptionEnv())

	if paramSettings == nil {
		app.Dao().SaveParam(pb_models.ParamAppSettings, app.Settings(), encryptionKey)
		paramSettings, err = app.Dao().FindParamByKey(pb_models.ParamAppSettings)
	}

	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error retrieving settings %s\n", err)
	}

	newSettings := pb_core.NewSettings()

	// try first without decryption
	plainDecodeErr := json.Unmarshal(paramSettings.Value, newSettings)

	if plainDecodeErr != nil {
		// load without decrypt has failed and there is no encryption key to use for decrypt
		if encryptionKey == "" {
			log.Fatalln(errors.New("failed to load the stored app settings (missing or invalid encryption key)"))
		}

		// decrypt
		decrypted, decryptErr := pb_security.Decrypt(string(paramSettings.Value), encryptionKey)
		if decryptErr != nil {
			log.Fatalf("%v: %s\n", errors.New("failed to decrypt"), decryptErr)
		}

		// decode again
		decryptedDecodeErr := json.Unmarshal(decrypted, newSettings)
		if decryptedDecodeErr != nil {
			log.Fatalf("%v: %s\n", errors.New("failed to decrypt"), decryptedDecodeErr)
		}
	}

	if err := app.RefreshSettings(); err != nil {
		log.Fatalf("failed to refresh setting, %v\n", err)
	}

	overwriteSMTPConfig(&newSettings.Smtp)
	overwriteS3Config(&newSettings.S3)

	if err := app.Settings().Merge(newSettings); err != nil {
		log.Fatalf("failed to merge setting, %v\n", err)
	}

	if plainDecodeErr == nil {
		// save because previously the settings weren't stored encrypted
		var saveErr error
		if encryptionKey != "" {
			saveErr = app.Dao().SaveParam(pb_models.ParamAppSettings, app.Settings(), encryptionKey)
		} else {
			saveErr = app.Dao().SaveParam(pb_models.ParamAppSettings, app.Settings())
		}

		if saveErr != nil {
			log.Fatalf("failed to save settings, %v\n", saveErr)
		}
	}
}
