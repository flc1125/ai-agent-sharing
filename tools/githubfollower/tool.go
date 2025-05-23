package githubfollower

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type Input struct {
	Username string `json:"username" jsonschema_description:"The username to get the followers count"`
}

type Output struct {
	Followers int `json:"followers" jsonschema_description:"The followers count"`
}

func NewTool() (tool.InvokableTool, error) {
	httpClient := http.DefaultClient
	return utils.InferTool(
		"get-github-followers",
		`A tool to get the followers count of a GitHub user. Input should be a username. Output will be the followers count.`,
		func(ctx context.Context, input *Input) (output *Output, err error) {
			if input.Username == "" {
				return nil, fmt.Errorf("username is required")
			}

			url := fmt.Sprintf("https://api.github.com/users/%s", input.Username)

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return nil, nil
			}

			// configure user agent
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

			// do the request
			resp, err := httpClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			output = new(Output)
			if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
				return nil, err
			}

			return output, nil
		},
	)
}
