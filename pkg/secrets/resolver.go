package secrets

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/99minutos/shipments-snapshot-service/pkg/logging"
	"github.com/sethvargo/go-envconfig"
)

const (
	SecretPrefix = "secret://"
	FileSuffix   = "?target=file"
)

func Resolver(sm SecretManager, config *Config) envconfig.MutatorFunc {
	if sm == nil {
		return nil
	}

	resolver := &secretResolver{
		sm:  sm,
		dir: config.SecretsDir,
	}

	return func(ctx context.Context, key, value string) (string, error) {
		s, err := resolver.resolve(ctx, key, value)
		if err != nil {
			return "", err
		}
		return s, nil
	}
}

type secretResolver struct {
	sm  SecretManager
	dir string
}

func (r *secretResolver) resolve(ctx context.Context, envName, secretRef string) (string, error) {
	logger := logging.FromContext(ctx)

	if !strings.HasPrefix(secretRef, SecretPrefix) {
		return secretRef, nil
	}

	if r.sm == nil {
		return "", errors.New("cannot get secret if there is no secret manager configured")
	}

	secretRef = strings.TrimPrefix(secretRef, SecretPrefix)

	toFile := false
	if strings.HasSuffix(secretRef, FileSuffix) {
		toFile = true
		secretRef = strings.TrimSuffix(secretRef, FileSuffix)
	}

	logger.Infof("resolving secret %q (toFile=%t) with ref %q", envName, toFile, secretRef)

	secretVal, err := r.sm.GetSecretValue(ctx, secretRef)
	if err != nil {
		return "", fmt.Errorf("failed to resolve %q: %w", secretRef, err)
	}

	if toFile {
		if err := r.ensureSecureDir(); err != nil {
			return "", err
		}

		secretFileName := filenameForSecret(envName + "." + secretRef)
		secretFilePath := path.Join(r.dir, secretFileName)
		if err := os.WriteFile(secretFilePath, []byte(secretVal), 0o600); err != nil {
			return "", fmt.Errorf("failed to write secret file for %q: %w", envName, err)
		}

		secretVal = secretFilePath
	}

	return secretVal, nil
}

func filenameForSecret(name string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(name))) //nolint:gosec
}

func (r *secretResolver) ensureSecureDir() error {
	stat, err := os.Stat(r.dir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if the secure directory %q exists: %w", r.dir, err)
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(r.dir, 0o700); err != nil {
			return fmt.Errorf("failed to create the secure directory %q: %w", r.dir, err)
		}
	} else if stat.Mode().Perm() != 0o700 {
		return fmt.Errorf("the secure directory %q exists and is not restricted %v", r.dir, stat.Mode())
	}
	return nil
}
