package configs

const (
	dbName     = "DB_NAME"
	dbPort     = "DB_PORT"
	dbHost     = "DB_HOST"
	dbUsername = "DB_USERNAME"
	dbPassword = "DB_PASSWORD"
	dbParams   = "DB_PARAM"
	jwtSecret  = "JWT_SECRET"
	bcryptSalt = "BCRYPT_SALT"
	s3ID       = "S3_ID"
	s3Secret   = "S3_SECRET_KEY" // #nosec G101
	s3Bucket   = "S3_BUCKET_NAME"
	s3Region   = "S3_REGION"

	defaultDBName     = "halo-suster"
	defaultDBPort     = 5432
	defaultDBHost     = "localhost"
	defaultDBUsername = "postgres"
	defaultDBPassword = "password"
	defaultDBParam    = "sslmode=disable"
	defaultJWTSecret  = "secret"
	defaultBCryptSalt = 8
	defaultS3ID       = ""
	defaultS3Secret   = ""
	defaultS3Bucket   = ""
	defaultS3Region   = ""
)

type RuntimeConfig struct {
	App appCfg `mapstructure:"APP"`
	API apiCfg `mapstructure:"API"`
	DB  dbCfg  `mapstructure:"DB"`
	S3  s3Cfg  `mapstructure:"S3"`
}

type appCfg struct {
	Name           string `mapstructure:"NAME"`
	Host           string `mapstructure:"HOST"`
	Port           int    `mapstructure:"PORT"`
	Version        string `mapstructure:"VERSION"`
	PreFork        bool   `mapstructure:"PREFORK"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`
}

type apiCfg struct {
	BaseURL    string `mapstructure:"BASE_URL"`
	Timeout    int    `mapstructure:"TIMEOUT"`
	DebugMode  bool   `mapstructure:"DEBUG_MODE"`
	BCryptSalt int    `mapstructure:"BCRYPT_SALT"`
	JWT        jwt    `mapstructure:"JWT"`
}

type jwt struct {
	Expire    int    `mapstructure:"EXPIRE"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

type dbCfg struct {
	Name               string  `mapstructure:"DB_NAME"`
	Port               int     `mapstructure:"DB_PORT"`
	Host               string  `mapstructure:"DB_HOST"`
	Username           string  `mapstructure:"DB_USERNAME"`
	Password           string  `mapstructure:"DB_PASSWORD"`
	Param              string  `mapstructure:"DB_PARAM"`
	MaxConnPool        int     `mapstructure:"MAX_CONN_POOL"`
	MaxConnPoolPercent float64 `mapstructure:"MAX_CONN_POOL_PERCENT"`
}

type s3Cfg struct {
	s3ID         string `mapstructure:"S3_ID"`
	s3SecretKey  string `mapstructure:"S3_SECRET_KEY"`
	s3BucketName string `mapstructure:"S3_BUCKET_NAME"`
	s3Region     string `mapstructure:"S3_REGION"`
}
