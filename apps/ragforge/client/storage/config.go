package storage

type EngineConfig struct {
	DefaultProvider string        `yaml:"default_provider" json:"default_provider"`
	Local           *LocalConfig  `yaml:"local,omitempty" json:"local,omitempty"`
	MinIO           *MinIOConfig  `yaml:"minio,omitempty" json:"minio,omitempty"`
}

type LocalConfig struct {
	PathPrefix string `yaml:"path_prefix" json:"path_prefix"`
}

type MinIOConfig struct {
	Mode            string `yaml:"mode" json:"mode"`
	Endpoint        string `yaml:"endpoint" json:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id" json:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key" json:"secret_access_key"`
	BucketName      string `yaml:"bucket_name" json:"bucket_name"`
	UseSSL          bool   `yaml:"use_ssl" json:"use_ssl"`
	PathPrefix      string `yaml:"path_prefix" json:"path_prefix"`
}
