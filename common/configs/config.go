package configs

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	once    sync.Once
	runtime *RuntimeConfig
)

func Get() *RuntimeConfig {
	once.Do(func() {
		callerInfo := "[configs.Get]"

		runtimeViper := viper.New()

		// ListCats config from file
		err := readConfigFile(runtimeViper, "config", "toml", "configs")
		if err != nil {
			panic(fmt.Errorf("%s failed to read config file: %v\n", callerInfo, err))
		}

		log.Printf("\n\n%#v\n\n", runtimeViper.AllSettings())

		// ListCats config from env vars
		err = readEnv(runtimeViper)
		if err != nil {
			panic(fmt.Errorf("%s failed to read env vars: %v\n", callerInfo, err))
		}

		log.Printf("\n\n%#v\n\n", runtimeViper.AllSettings())

		// Load config into runtimeConfig
		runtime, err = loadConfig(runtimeViper)
		if err != nil {
			panic(fmt.Errorf("%s failed to load config: %v\n", callerInfo, err))
		}
	})

	if runtime == nil {
		panic(fmt.Errorf("[configs.Get] runtime is nil"))
	}

	return runtime
}

func readEnv(runtimeViper *viper.Viper) error {
	// Load env vars
	runtimeViper.AllowEmptyEnv(false)

	if err := runtimeViper.BindEnv("API."+bcryptSalt, bcryptSalt); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("JWT."+jwtSecret, jwtSecret); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("DB."+dbUsername, dbUsername); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("DB."+dbPassword, dbPassword); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("DB."+dbHost, dbHost); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("DB."+dbPort, dbPort); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("DB."+dbName, dbName); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("DB."+dbParams, dbParams); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("S3."+awsAccessKeyID, awsAccessKeyID); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("S3."+awsSecretAccessKey, awsSecretAccessKey); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("S3."+awsS3BucketName, awsS3BucketName); err != nil {
		return err
	}
	if err := runtimeViper.BindEnv("S3."+awsRegion, awsRegion); err != nil {
		return err
	}

	return nil
}

func readConfigFile(runtimeViper *viper.Viper, fileName, fileType string, filePath ...string) error {
	callerInfo := "[configs.readConfigFile]"

	runtimeViper.SetConfigName(fileName)
	runtimeViper.SetConfigType(fileType)
	for _, path := range filePath {
		runtimeViper.AddConfigPath(path)
	}

	err := runtimeViper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("%s failed to read config file: %v\n", callerInfo, err)
	}
	return nil
}

func loadConfig(runtimeViper *viper.Viper) (*RuntimeConfig, error) {
	callerInfo := "[configs.loadConfig]"

	// load env vars to dbCfg and apiCfg
	dbConfig, apiConfig, jwtConfig, s3Config := &dbCfg{}, &apiCfg{}, &jwtCfg{}, &s3Cfg{}
	err := runtimeViper.Unmarshal(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode dbConfig: %v\n", callerInfo, err)
	}
	err = runtimeViper.Unmarshal(apiConfig)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode apiConfig: %v\n", callerInfo, err)
	}
	err = runtimeViper.Unmarshal(jwtConfig)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode jwtConfig: %v\n", callerInfo, err)
	}
	err = runtimeViper.Unmarshal(s3Config)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode s3Config: %v\n", callerInfo, err)
	}

	// set env vars to runtimeConfig before decode from config file
	runtimeConfig := &RuntimeConfig{
		API: *apiConfig,
		DB:  *dbConfig,
		JWT: *jwtConfig,
		S3:  *s3Config,
	}
	err = runtimeViper.Unmarshal(runtimeConfig)
	if err != nil {
		return runtimeConfig, fmt.Errorf("%s failed to decode runtimeConfig: %v\n", callerInfo, err)
	}

	return runtimeConfig, nil
}
