package middleware

import (
	"encoding/json"
	"time"
)

type (
	UserRoles struct {
		IdentityProvider string    `json:"identity_provider"`
		InstitutionId    string    `json:"institution_id"`
		SessionId        string    `json:"session_id"`
		UserId           string    `json:"user_id"`
		ExternalSubject  string    `json:"external_subject"`
		Roles            []Roles   `json:"roles"`
		AccessToken      string    `json:"access_token"`
		ExpiresAt        time.Time `json:"expires_at"`
	}
	Roles struct {
		Groups      []string `json:"groups"`
		RoleId      string   `json:"role_id"`
		RoleName    string   `json:"role_name"`
		Permissions []string `json:"permissions"`
	}
)

func (r *UserRoles) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Roles) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *UserRoles) Groups() []string {
	var groups []string
	for _, role := range r.Roles {
		groups = append(groups, role.Groups...)
	}

	return groups
}

func (r *UserRoles) Permissions() []string {
	var permissions []string
	for _, role := range r.Roles {
		permissions = append(permissions, role.Permissions...)
	}

	return permissions
}
