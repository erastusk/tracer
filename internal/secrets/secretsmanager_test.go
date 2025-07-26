package secrets_test

import (
	"testing"

	"github/erastusk/tracer/internal/secrets"
	_types "github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/mocks"

	"go.uber.org/mock/gomock"
)

func TestGetSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tc := mocks.NewMockSecretsManager(ctrl)
	// Set up expected behavior
	name := "/dev/mock/path"
	region := "us-east-1"
	var s *secrets.SecretsSession
	tc.EXPECT().
		GetSession(name, region).
		Return(s, nil)

	// Call the method
	session, err := tc.GetSession(name, region)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if session != s {
		t.Errorf("expected session 'mock-session', got %v", session)
	}
}

func TestGetSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tc := mocks.NewMockSecretsManager(ctrl)
	// Set up expected behavior
	var s _types.Secrets
	tc.EXPECT().GetSecrets().Return(s, nil)

	// Call the method
	result, err := tc.GetSecrets()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != s {
		t.Errorf("expected secret '%s', got %s", s, result)
	}
}
