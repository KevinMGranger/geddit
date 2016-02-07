package geddit

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/kevinmgranger/requestmod"
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

func ClientSetUserAgent(cli *http.Client, useragent string) {
	cli.Transport = requestmod.NewTransport(cli.Transport, func(req *http.Request) error {
		req.Header.Set("User-Agent", useragent)
		return nil
	})
}

type Client struct {
	*http.Client
}

type NoCtxClient Client

// to easily wrap a non-context supporting client
func contextRequest(cli *Client, ctx context.Context, req *http.Request) (resp *http.Response, err error) {
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
			// TODO: cancel request here?
			return nil, ctx.Err()
	}
}

func jsonResponse(cli *Client, ctx context.Context, req *http.Request, res interface{}) (err error) {
	resp, err := contextRequest(cli, ctx, req)
	if err != nil {
		return
	}

	bod, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(bod, res)
	return
}