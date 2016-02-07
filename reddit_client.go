package geddit

import (
	"net/http"
	"golang.org/x/net/context"
)

type RedditClient interface {
	Me() (*Redditor, error)
}

type ContextRedditClient interface {
	Me(ctx context.Context) (*Redditor, error)
}

const (
	meURL string = "https://oauth.reddit.com/api/v1/me"
)


func me(cli *Client, ctx context.Context) (r *Redditor, err error) {
	req , err := http.NewRequest("GET", meURL, nil)

	if err != nil {
		return nil, err
	}

	r = &Redditor{}
	err = jsonResponse(cli, ctx, req, r)
	return
}

func (cli *Client) Me(ctx context.Context) (r *Redditor, err error) {
	return me(cli, ctx)
}

func (cli *NoCtxClient) Me() (r *Redditor, err error) {
	return me((*Client)(cli), context.TODO())
}
	