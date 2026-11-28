// Package session defines the server side session state of a signed in user.
//
// The auth-service is the only writer of the session state. It stores one
// Redis entry per issued token pair and removes the entries when the session
// is refreshed or revoked.
//
// A JWT stays valid until its expiry claim passes, so the signature alone
// cannot tell a live token apart from a revoked one. The API gateway reads
// the same Redis entries to reject a token that the auth-service has already
// invalidated, which is what makes a sign out effective right away instead of
// at the end of the token lifetime.
//
// The key format lives in this package because it is a contract between the
// writer and the reader. Keeping it in one place prevents the two from
// drifting apart.
package session

const (
	// AccessKeyPrefix prefixes the Redis entry of an issued access token.
	AccessKeyPrefix = "access:"
	// RefreshKeyPrefix prefixes the Redis entry of an issued refresh token.
	RefreshKeyPrefix = "refresh:"
)

// AccessKey returns the Redis key that holds the state of an access token.
//
// The key exists while the access token is usable. It is removed when the
// session is refreshed, when the user signs out and when it expires.
func AccessKey(jti string) string {
	return AccessKeyPrefix + jti
}

// RefreshKey returns the Redis key that holds the state of a refresh token.
//
// The key exists while the refresh token can still be exchanged. A refresh
// consumes the entry and writes a new one, so a replay of the old token is
// rejected.
func RefreshKey(jti string) string {
	return RefreshKeyPrefix + jti
}
