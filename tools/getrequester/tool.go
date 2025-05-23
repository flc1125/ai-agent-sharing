package getrequester

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type Input struct {
	URL string `json:"url" jsonschema_description:"The URL to make the GET request"`
}

func NewTool() (tool.InvokableTool, error) {
	httpClient := http.DefaultClient

	return utils.InferTool(
		"get-requester",
		`A portal to the internet. Use this when you need to get specific
		content from a website. Input should be a URL (i.e. https://www.google.com).
		The output will be the text response of the GET request.`,
		func(ctx context.Context, input *Input) (output string, err error) {
			if input.URL == "" {
				return "", errors.New("url is required")
			}

			req, err := http.NewRequestWithContext(ctx, "GET", input.URL, nil)
			if err != nil {
				return "", nil
			}

			// configure user agent
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

			// do the request
			resp, err := httpClient.Do(req)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}

			return string(body), nil
		},
	)
}
