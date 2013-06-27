package examples

import (
	"fmt"
	"goa/oauth"
)

func main() {
	provider := oauth.Provider{
		RequestTokenUrl:  "https://api.twitter.com/oauth/request_token",
		AuthorizationUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:   "https://api.twitter.com/oauth/access_token"}

	client := oauth.NewClient(provider,
		"http://localhost:8080",
		"BNdWbiUiSfiPj1MyWXhpQA",
		"iNhtpL4Kfktxj5k3WuDBMwkUBc3n3Qe8hUoJNDVookg")
	if client != nil {
		fmt.Println(client.ConsumerKey)
	}
	client.GetRequestToken()
}
