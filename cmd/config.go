package cmd

type Config struct {
	Production bool
	Admins     []string
	Spaces     []SpacesConfig
	DB         string
	S3         S3Config
	Addr       string
	OIDC       OIDCConfig
	Session    SessionConfig
	CSRF       CSRFConfig
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

type SessionConfig struct {
	Name   string
	MaxAge int `mapstructure:"max_age"`
	Secret string
}

type CSRFConfig struct {
	Secret string
}
