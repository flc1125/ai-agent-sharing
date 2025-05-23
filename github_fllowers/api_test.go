package github_fllowers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	t.Logf("获取 https://flc.io/ 的作者的Github账号所对应的粉丝数量：%d",
		getFlowers(t, "https://flc.io/"))
}

type githubUserResp struct {
	Login     string `json:"login"`
	Followers int    `json:"followers"`
}

func getFlowers(tb testing.TB, url string) int {
	// 提取 URL 中的用户名
	resp, err := http.Get(url)
	require.NoError(tb, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(tb, err)
	assert.NotEmpty(tb, body)
	assert.NoError(tb, resp.Body.Close())

	re := regexp.MustCompile(`.+https://github.com/([^"]+).+关注我.+`)
	matches := re.FindStringSubmatch(string(body))
	require.Len(tb, matches, 2)
	assert.NotEmpty(tb, matches[1])

	tb.Logf("Github username: %s", matches[1])

	// 提取粉丝数量
	resp, err = http.Get("https://api.github.com/users/" + matches[1])
	require.NoError(tb, err)
	defer func() {
		assert.NoError(tb, resp.Body.Close())
	}()

	var response githubUserResp
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(tb, err)
	assert.Equal(tb, matches[1], response.Login)

	return response.Followers
}
