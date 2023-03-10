[![goreleaser](https://github.com/henokv/azure-env/actions/workflows/release.yml/badge.svg)](https://github.com/henokv/azure-env/actions/workflows/release.yml)

# azure-env

Azure env is a CLI tool built to help solve a specific problem.
Sometimes while testing certain CLI tools (e.g.: terraform import) locally I can't automatically use all credentials in the shell.

This tool allows you to use env vars and still use azure key vault to store the secrets.
Later on support for azure app configuration might get added.

## Installation
To install download the latest version from the [releases](https://github.com/henokv/azure-env/releases) page or if you have go installed run the command
```shell
go install github.com/henokv/azure-env@latest
```

## Usage 
### Linux config file
Contents of .env
```
AZURE_CLIENT_SECRET=azure://knox.vault.azure.net/secret
PUBLIC_ENV_VAR=notasecret
```
```shell
# Run the command which will consume the secrets from the env file
azure-env run -f .env terraform plan
```
### Linux env vars
```shell
# Set a secret supported by azure cli & terraform
export AZURE_CLIENT_SECRET=azure://knox.vault.azure.net/secret
# Run the command which will consume the secret
azure-env run terraform plan
```

### Windows env vars
```shell
# Set a secret supported by azure cli & terraform
$Env:AZURE_CLIENT_SECRET=azure://knox.vault.azure.net/secret
# Run the command which will consume the secret
azure-env run terraform plan
```
```