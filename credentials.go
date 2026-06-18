package shopware

// Credentials produces the OAuth token request for a Client. The Shopware Admin
// API supports two grant types, each modeled by an implementation here:
//
//   - IntegrationCredentials — client_credentials grant (integration / app).
//   - PasswordCredentials     — password grant (admin user login).
//
// Custom implementations may supply any other body the /api/oauth/token
// endpoint accepts.
type Credentials interface {
	// tokenRequest returns the JSON body POSTed to /api/oauth/token.
	tokenRequest() map[string]string

	// identity returns a stable string that uniquely identifies the principal
	// these credentials authenticate as. It is used as the default token
	// storage key so distinct principals never share a cached token.
	identity() string
}

// IntegrationCredentials authenticates as a Shopware integration using the
// client_credentials grant.
type IntegrationCredentials struct {
	ClientID     string
	ClientSecret string
}

// NewIntegrationCredentials creates IntegrationCredentials.
func NewIntegrationCredentials(clientID, clientSecret string) IntegrationCredentials {
	return IntegrationCredentials{ClientID: clientID, ClientSecret: clientSecret}
}

func (c IntegrationCredentials) tokenRequest() map[string]string {
	return map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
	}
}

func (c IntegrationCredentials) identity() string {
	return "integration:" + c.ClientID
}

// PasswordCredentials authenticates as an admin user using the password grant.
// The client_id is fixed to "administration", matching the Shopware admin SPA.
type PasswordCredentials struct {
	Username string
	Password string
}

// NewPasswordCredentials creates PasswordCredentials.
func NewPasswordCredentials(username, password string) PasswordCredentials {
	return PasswordCredentials{Username: username, Password: password}
}

func (c PasswordCredentials) tokenRequest() map[string]string {
	return map[string]string{
		"grant_type": "password",
		"client_id":  "administration",
		"scopes":     "write",
		"username":   c.Username,
		"password":   c.Password,
	}
}

func (c PasswordCredentials) identity() string {
	return "password:" + c.Username
}
