package pinboard

import (
	"testing"
)

func TestUserSecret(t *testing.T) {
	secret, err := UserSecret()
	if err != nil {
		t.Errorf("error: UserSecret: %s", err)
	}

	if secret == "" {
		t.Error("error: UserSecret: expected secret string")
	}
}

func TestUserAPIToken(t *testing.T) {
	token, err := UserAPIToken()
	if err != nil {
		t.Errorf("error: UserAPIToken: %s", err)
	}

	if token == "" {
		t.Error("error: UserAPIToken: expected token string")
	}
}
