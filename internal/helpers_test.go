package internal

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type testCredential struct {
	token azcore.AccessToken
	err   error
}

func (t testCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return t.token, t.err
}

func TestValidateAzureCredentialReturnsAuthError(t *testing.T) {
	verbose = false
	err := ValidateAzureCredential(testCredential{err: errors.New("token fetch failed")})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "unable to authenticate with Azure") {
		t.Fatalf("expected auth error but got: %v", err)
	}
}

func TestValidateAzureCredentialVerboseIncludesOriginalError(t *testing.T) {
	verbose = true
	t.Cleanup(func() {
		verbose = false
	})
	err := ValidateAzureCredential(testCredential{err: errors.New("token fetch failed")})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if !strings.Contains(err.Error(), "token fetch failed") {
		t.Fatalf("expected wrapped source error but got: %v", err)
	}
}
