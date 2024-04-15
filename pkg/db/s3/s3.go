package s3

type S3Config struct {
	Endpoint  string `yaml:"endpoint" env:"S3_ENDPOINT"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
}
