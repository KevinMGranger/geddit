package geddit

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// A ScriptOAuthConfig represents a token source and configuration for a script-type app.
// For more information on script-type apps, see https://github.com/reddit/reddit/wiki/OAuth2-App-Types#script
// and https://github.com/reddit/reddit/wiki/OAuth2-Quick-Start-Example
type ScriptOAuthConfig struct {
	// A standard oauth2 Config.
	// Scopes are ignored for script apps.
	// The Endpoint should be RedditEndpoint.
	Config oauth2.Config

	Username string

	Password string

	// The context which will be passed on to oauth2 code.
	// You may specify an HTTP client here.
	// Do not leave as nil, instead use context.TODO().
	Context context.Context
}

// Token gets a proper token from the given configuration, or an error if one could not be obtained.
// The request to get the token will use the HTTP client from the Context.
// It is not recommended to use a ScriptOAuthConfig directly as a TokenSource, as it will not cache.
// Use .TokenSource() instead.
func (src *ScriptOAuthConfig) Token() (*oauth2.Token, error) {
	return src.Config.PasswordCredentialsToken(src.Context, src.Username, src.Password)
}

func (src *ScriptOAuthConfig) TokenSource() oauth2.TokenSource {
	return oauth2.ReuseTokenSource(nil, src)
}

func (src *ScriptOAuthConfig) Client() *Client {
	return oauth2.NewClient(src.Context, src.TokenSource())
}