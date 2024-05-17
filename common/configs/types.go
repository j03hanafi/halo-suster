package configs

const (
	bcryptSalt         = "BCRYPT_SALT"
	jwtSecret          = "JWT_SECRET"
	dbUsername         = "DB_USERNAME"
	dbPassword         = "DB_PASSWORD"
	dbHost             = "DB_HOST"
	dbPort             = "DB_PORT"
	dbName             = "DB_NAME"
	dbParams           = "DB_PARAMS"
	awsAccessKeyID     = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKey = "AWS_SECRET_ACCESS_KEY" // #nosec G10
	awsS3BucketName    = "AWS_S3_BUCKET_NAME"
	awsRegion          = "AWS_REGION"
)

type RuntimeConfig struct {
	App appCfg `mapstructure:"APP"`
	API apiCfg `mapstructure:"API"`
	JWT jwtCfg `mapstructure:"JWT"`
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
	DebugMode      bool   `mapstructure:"DEBUG_MODE"`
}

type apiCfg struct {
	BaseURL    string `mapstructure:"BASE_URL"`
	Timeout    int    `mapstructure:"TIMEOUT"`
	BCryptSalt int    `mapstructure:"BCRYPT_SALT"`
}

type jwtCfg struct {
	Expire    int    `mapstructure:"EXPIRE"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
}

type dbCfg struct {
	Username           string  `mapstructure:"DB_USERNAME"`
	Password           string  `mapstructure:"DB_PASSWORD"`
	Host               string  `mapstructure:"DB_HOST"`
	Port               int     `mapstructure:"DB_PORT"`
	Name               string  `mapstructure:"DB_NAME"`
	Params             string  `mapstructure:"DB_PARAMS"`
	MaxConnPool        int     `mapstructure:"MAX_CONN_POOL"`
	MaxConnPoolPercent float64 `mapstructure:"MAX_CONN_POOL_PERCENT"`
}

type s3Cfg struct {
	AccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	BucketName      string `mapstructure:"AWS_S3_BUCKET_NAME"`
	Region          string `mapstructure:"AWS_REGION"`
}
