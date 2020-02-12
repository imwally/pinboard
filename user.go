package pinboard

import (
	"encoding/json"
)

// userResponse holds the response from /user/secret and
// /user/api_token
type userResponse struct {
	Result string `json:"result"`
}

// UserSecret returns the user's secret RSS key (for viewing private
// feeds).
func UserSecret() (string, error) {
	resp, err := get("userSecret", nil)
	if err != nil {
		return "", err
	}

	var ur userResponse
	err = json.Unmarshal(resp, &ur)
	if err != nil {
		return "", err
	}

	return ur.Result, nil
}

// UserAPIToken returns the user's API token (for making API calls
// without a password).
func UserAPIToken() (string, error) {
	resp, err := get("userAPIToken", nil)
	if err != nil {
		return "", err
	}

	var ur userResponse
	err = json.Unmarshal(resp, &ur)
	if err != nil {
		return "", err
	}

	return ur.Result, nil
}
