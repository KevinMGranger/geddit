package geddit

import (
	"net/http"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)



func WithClient(ctx context.Context, client *http.Client) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, client)
}

//TODO: document about non-cancelable transport (won't warn when using timeout)
// Always gets a client.
func ClientFromContext(ctx context.Context) *http.Client {
	if cli, ok := ctx.Value(oauth2.HTTPClient).(*http.Client); ok {
		return cli
	}

	return http.DefaultClient
}

func WithUserAgent(ctx context.Context, useragent string) context.Context {
	// TODO: copy, don't modify client
	cli := ClientFromContext(ctx)
	cli.Transport = requestmod.NewTransport(cli.Transport, func(req *http.Request) error {
		req.Header.Set("User-Agent", useragent)
		return nil
	})
	return WithClient(ctx, cli)
}

type Client http.Client

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type ContextClient interface {
	DoWithContext(ctx context.Context, req *http.Request) (resp *http.Response, err error)
}

// to easily wrap a non-context supporting client
func contextRequest(cli Client, ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	type httpresult struct {
		resp *http.Response
		err error
	}
	reschan := make(chan httpresult, 1)

	go func() { res, err := cli.Do(req); reschan <- httpresult{ res, err } }()

	select {
		case res := <- reschan:
			return res.resp, res.err
		case <-ctx.Done():
			return nil, ctx.Err()
	}
}

func jsonResponse(cli ContextClient, ctx context.Context, req *http.Request, res interface{}) error {
	resp, err := cli.DoWithContext(ctx, req)
	if err != nil {
		return err
	}

	bod, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(bod, res)
	if err != nil {
		return err
	}
}