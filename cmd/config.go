package cmd

type Config struct {
	Production   bool
	Admins       []string
	DB           string
	S3           S3Config
	Host         string
	Port         int
	OIDC         OIDCConfig
	MaxFileSize  int64  `mapstructure:"max_file_size"`
	CookieSecret string `mapstructure:"cookie_secret"`
}

type SpacesConfig struct {
	ID     string
	Admins []string
}

type S3Config struct {
	URL    string
	Region string
	ID     string
	Secret string
	Bucket string
}

type OIDCConfig struct {
	URL         string
	ID          string
	Secret      string
	RedirectURL string `mapstructure:"redirect_url"`
}

type CSRFConfig struct {
	Secret string
}
