package cmd

type Config struct {
	Production bool
	Admin      []string
	Spaces     []SpacesConfig
	DB         string
	S3         S3Config
	Addr       string
	Oidc       OidcConfig
	Session    SessionConfig
	Csrf       CsrfConfig
}

type SpacesConfig struct {
	ID    string
	Admin []string
}

type S3Config struct {
	URL    string
	Region string
	ID     string
	Secret string
	Bucket string
}

type OidcConfig struct {
	URL         string
	ID          string
	Secret      string
	RedirectURL string
}

type SessionConfig struct {
	Name   string
	MaxAge int
	Secret string
}

type CsrfConfig struct {
	Secret string
}
