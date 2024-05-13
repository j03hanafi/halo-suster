package configs

import (
	"fmt"
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

		// ListCats config from env vars
		readEnv(runtimeViper)

		// ListCats config from file
		err := readConfigFile(runtimeViper, "config", "toml", "configs")
		if err != nil {
			panic(fmt.Errorf("%s failed to read config file: %v\n", callerInfo, err))
		}

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

func readEnv(runtimeViper *viper.Viper) {
	// Set defaults for env vars
	runtimeViper.SetDefault(dbName, defaultDBName)
	runtimeViper.SetDefault(dbPort, defaultDBPort)
	runtimeViper.SetDefault(dbHost, defaultDBHost)
	runtimeViper.SetDefault(dbUsername, defaultDBUsername)
	runtimeViper.SetDefault(dbPassword, defaultDBPassword)
	runtimeViper.SetDefault(dbParams, defaultDBParam)
	runtimeViper.SetDefault(jwtSecret, defaultJWTSecret)
	runtimeViper.SetDefault(bcryptSalt, defaultBCryptSalt)
	runtimeViper.SetDefault(s3ID, defaultS3ID)
	runtimeViper.SetDefault(s3Secret, defaultS3Secret)
	runtimeViper.SetDefault(s3Bucket, defaultS3Bucket)
	runtimeViper.SetDefault(s3Region, defaultS3Region)

	// Load env vars
	runtimeViper.AllowEmptyEnv(false)
	runtimeViper.AutomaticEnv()
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
	dbConfig, apiConfig, s3Config := &dbCfg{}, &apiCfg{}, &s3Cfg{}
	err := runtimeViper.Unmarshal(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode dbConfig: %v\n", callerInfo, err)
	}
	err = runtimeViper.Unmarshal(apiConfig)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode apiConfig: %v\n", callerInfo, err)
	}
	err = runtimeViper.Unmarshal(s3Config)
	if err != nil {
		return nil, fmt.Errorf("%s failed to decode s3Config: %v\n", callerInfo, err)
	}

	// because we get jwtSecret from env vars but in config runtime it's nested, we need to manually set it
	apiConfig.JWT.JWTSecret = runtimeViper.GetString(jwtSecret)

	// set env vars to runtimeConfig before decode from config file
	runtimeConfig := &RuntimeConfig{
		API: *apiConfig,
		DB:  *dbConfig,
	}
	err = runtimeViper.Unmarshal(runtimeConfig)
	if err != nil {
		return runtimeConfig, fmt.Errorf("%s failed to decode runtimeConfig: %v\n", callerInfo, err)
	}

	return runtimeConfig, nil
}
