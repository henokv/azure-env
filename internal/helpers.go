package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"os"
	"strings"
	"sync"
)

// Secret represents an environment variable that may reference an Azure Key Vault secret
type Secret struct {
	Env      string `json:"env"`      // The final rendered environment variable (KEY=value)
	EnvRef   string `json:"env_ref"`  // The original reference (KEY=azure://vault/secret)
	Key      string `json:"key"`      // The environment variable name
	Value    string `json:"value"`    // The resolved secret value
	ValueRef string `json:"value_ref"` // The Azure reference (azure://vault/secret)
}

// NewSecret creates a new Secret from an Azure reference
func NewSecret(envRef string) *Secret {
	return &Secret{
		EnvRef: envRef,
	}
}

// Config holds the global configuration state for Azure authentication
type Config struct {
	credential *azidentity.DefaultAzureCredential
	mu         *sync.Mutex
	verbose    bool
}

// DefaultConfig is the singleton configuration instance
var defaultConfig = &Config{
	mu: &sync.Mutex{},
}

// GetEnvAsSecret separates environment variables into Azure Key Vault references and regular variables
func GetEnvAsSecret() ([]Secret, []string, error) {
	environ := os.Environ()
	secrets := make([]Secret, 0)
	otherEnv := make([]string, 0, len(environ))

	for _, env := range environ {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		valueRef := parts[1]

		if strings.HasPrefix(valueRef, "azure://") {
			secret, err := GetSecretByRef(valueRef)
			if err != nil {
				return nil, nil, err
			}
			if secret.Value == nil {
				return nil, nil, fmt.Errorf("empty secret value for %s", valueRef)
			}

			secrets = append(secrets, Secret{
				EnvRef:   fmt.Sprintf("%s=%s", key, valueRef),
				Key:      key,
				ValueRef: valueRef,
				Value:    *secret.Value,
				Env:      fmt.Sprintf("%s=%s", key, *secret.Value),
			})
		} else {
			otherEnv = append(otherEnv, env)
		}
	}
	return secrets, otherEnv, nil
}

// SetSecretsToEnv sets resolved secret values to the environment
func SetSecretsToEnv(secrets []Secret) error {
	for _, secret := range secrets {
		if err := os.Setenv(secret.Key, secret.Value); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", secret.Key, err)
		}
	}
	return nil
}

// GetOriginalEnv returns the original environment variable references (before resolution)
func GetOriginalEnv(secrets []Secret) []string {
	env := make([]string, len(secrets))
	for i, secret := range secrets {
		env[i] = secret.EnvRef
	}
	return env
}

// GetRenderedEnv returns the resolved environment variables (after secret resolution)
func GetRenderedEnv(secrets []Secret) []string {
	env := make([]string, len(secrets))
	for i, secret := range secrets {
		env[i] = secret.Env
	}
	return env
}

// GetFullRenderedEnv returns all environment variables (secrets + non-secrets)
func GetFullRenderedEnv(secrets []Secret, otherEnv []string) []string {
	env := make([]string, 0, len(secrets)+len(otherEnv))
	for _, secret := range secrets {
		env = append(env, secret.Env)
	}
	env = append(env, otherEnv...)
	return env
}

// DecodeRef parses an Azure reference (azure://vaultname/secretname) into vault URL and secret name
func DecodeRef(ref string) (string, string, error) {
	const azurePrefix = "azure://"
	if !strings.HasPrefix(ref, azurePrefix) {
		return "", "", fmt.Errorf("reference requires prefix %s, but got '%s'", azurePrefix, ref)
	}

	ref = strings.TrimPrefix(ref, azurePrefix)
	parts := strings.Split(ref, "/")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("reference should contain 2 parts separated by '/', but got '%s'", ref)
	}

	vaultURL := fmt.Sprintf("https://%s.vault.azure.net", parts[0])
	secretName := parts[1]

	return vaultURL, secretName, nil
}

// GetSecretByRef retrieves a secret from Azure Key Vault using a reference string
func GetSecretByRef(ref string) (azsecrets.GetSecretResponse, error) {
	vaultURL, secretName, err := DecodeRef(ref)
	if err != nil {
		return azsecrets.GetSecretResponse{}, err
	}
	return GetSecret(vaultURL, secretName)
}

// formatError returns a user-friendly error message, with detailed info if verbose mode is enabled
func (c *Config) formatError(verbose bool, detailedErr error, userMsg string) error {
	if verbose {
		return fmt.Errorf("%s: %w", userMsg, detailedErr)
	}
	return errors.New(userMsg)
}

// ensureAuthenticated initializes the Azure credential if not already done
func (c *Config) ensureAuthenticated() error {
	if c.credential != nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check again after acquiring lock (double-check locking)
	if c.credential != nil {
		return nil
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		var authErr azidentity.AuthenticationFailedError
		if errors.As(err, &authErr) {
			return c.formatError(c.verbose, err, "not authenticated: check Azure authentication (env vars, CLI, or managed identity), or add -v/--verbosity for details")
		}
		return c.formatError(c.verbose, err, "authentication failed: check Azure authentication (env vars, CLI, or managed identity)")
	}

	c.credential = cred
	return nil
}

// GetSecret retrieves a secret from Azure Key Vault
func GetSecret(vaultURL, secretName string) (azsecrets.GetSecretResponse, error) {
	if err := defaultConfig.ensureAuthenticated(); err != nil {
		return azsecrets.GetSecretResponse{}, err
	}

	client, err := azsecrets.NewClient(vaultURL, defaultConfig.credential, nil)
	if err != nil {
		return azsecrets.GetSecretResponse{}, defaultConfig.formatError(
			defaultConfig.verbose,
			err,
			fmt.Sprintf("unable to create client for vault %s: check Azure authentication", vaultURL),
		)
	}

	ctx := context.Background()
	secret, err := client.GetSecret(ctx, secretName, "", nil)
	if err != nil {
		return azsecrets.GetSecretResponse{}, defaultConfig.formatError(
			defaultConfig.verbose,
			err,
			fmt.Sprintf("unable to retrieve secret %s from vault %s: check secret exists and you have access", secretName, vaultURL),
		)
	}

	return secret, nil
}

// SetVerbosity enables or disables verbose logging
func SetVerbosity(verbose bool) {
	defaultConfig.mu.Lock()
	defer defaultConfig.mu.Unlock()
	defaultConfig.verbose = verbose
}
