package router

import (
	"fmt"
	"net/url"
)

// ChangePassword changes the router admin/user password.
func ChangePassword(client *Client, oldPass, newPass, confirmPass string) (bool, error) {
	_, err := client.GetPage(PageUserMgmt)
	if err != nil {
		return false, fmt.Errorf("failed to load user management page: %w", err)
	}

	if client.SessionToken == "" {
		return false, fmt.Errorf("no session token available")
	}

	if newPass != confirmPass {
		return false, fmt.Errorf("new password and confirmation do not match")
	}

	if len(newPass) < 4 {
		return false, fmt.Errorf("password must be at least 4 characters")
	}

	formData := url.Values{}
	formData.Set("IF_ACTION", "apply")
	formData.Set("UserName", "admin")
	formData.Set("OldPassword", oldPass)
	formData.Set("NewPassword", newPass)
	formData.Set("ConfirmedPassword", confirmPass)
	formData.Set("_SESSION_TOKEN", client.SessionToken)

	_, err = client.PostAction(PageUserMgmt, formData.Encode())
	if err != nil {
		return false, fmt.Errorf("failed to change password: %w", err)
	}

	return true, nil
}
