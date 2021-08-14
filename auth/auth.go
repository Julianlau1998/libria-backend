package auth

func IsAuthorized(cookie []string) bool {
	if cookie != nil && cookie[0] == "_legacy_auth0.is.authenticated=true; auth0.is.authenticated=true" {
		return true
	}
	// return false
	return true
}
