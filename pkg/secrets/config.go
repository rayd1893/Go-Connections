package secrets

type Config struct {
	Type       string `env:"SECRET_MANAGER, default=GOOGLE_SECRET_MANAGER"`
	SecretsDir string `env:"SECRETS_DIR, default=/var/run/secrets"`

	FilesystemRoot string `env:"SECRET_FILESYSTEM_ROOT"`
}
