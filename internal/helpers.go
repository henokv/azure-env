package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"os"
	"strings"
)

type Secret struct {
	EnvRef   string `json:"env_ref"`
	Key      string `json:"key"`
	Value    string `json:"value"`
	ValueRef string `json:"value_ref"`
	Env      string `json:"env"`
}

func GetEnvAsSecret() (secrets []Secret, otherEnv []string, error error) {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		valueRef := parts[1]
		if strings.HasPrefix(valueRef, "azure://") {
			secret, error := GetSecretByRef(valueRef)
			if error != nil {
				return []Secret{}, []string{}, error
			}
			value := *secret.Value
			secrets = append(secrets, Secret{
				EnvRef:   fmt.Sprintf("%s=%s", key, valueRef),
				Key:      key,
				ValueRef: valueRef,
				Value:    value,
				Env:      fmt.Sprintf("%s=%s", key, value),
			})
		} else {
			otherEnv = append(otherEnv, env)
		}
	}
	return secrets, otherEnv, nil
}

func SetSecretsToEnv(secrets []Secret) {
	for _, secret := range secrets {
		os.Setenv(secret.Key, secret.Value)
	}
}

func GetOriginalEnv(secrets []Secret) (env []string) {
	for _, secret := range secrets {
		env = append(env, secret.EnvRef)
	}
	return env
}

func GetRenderedEnv(secrets []Secret) (env []string) {
	for _, secret := range secrets {
		env = append(env, secret.Env)
	}
	return env
}

func GetFullRenderedEnv(secrets []Secret, otherEnv []string) (env []string) {
	for _, secret := range secrets {
		env = append(env, secret.Env)
	}
	env = append(env, otherEnv...)
	return env
}

func DecodeRef(ref string) (vaultUrl, secretName string, error error) {
	if !strings.HasPrefix(ref, "azure://") {
		return "", "", fmt.Errorf("reference requires prefix azure://, but got '%s'", ref)
	}
	ref = strings.TrimPrefix(ref, "azure://")
	refs := strings.Split(ref, "/")
	if len(refs) > 2 {
		return vaultUrl, secretName, fmt.Errorf("reference requires prefix azure://, but got '%s'", ref)
	}
	vaultUrl = fmt.Sprintf("https://%s", refs[0])
	secretName = refs[1]
	return vaultUrl, secretName, nil
}

func GetSecretByRef(ref string) (azsecrets.GetSecretResponse, error) {
	vaultName, secretName, err := DecodeRef(ref)
	if err != nil {
		return azsecrets.GetSecretResponse{}, err
	}
	return GetSecret(vaultName, secretName)
}

func GetSecret(vaultUrl, secretName string) (azsecrets.GetSecretResponse, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		var responseError azidentity.AuthenticationFailedError
		errors.As(err, &responseError)
		return azsecrets.GetSecretResponse{}, fmt.Errorf("authentication error: ", responseError.RawResponse.Status)
	}
	client, err := azsecrets.NewClient(vaultUrl, cred, nil)
	if err != nil {
		return azsecrets.GetSecretResponse{}, fmt.Errorf("client creation error: %s", err)
	}
	ctx := context.Background()
	secret, err := client.GetSecret(ctx, secretName, "", nil)
	if err != nil {
		return azsecrets.GetSecretResponse{}, fmt.Errorf("get secret error: %s", err)
	}
	return secret, nil
}
