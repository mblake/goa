package examples

import (
	"fmt"
	"goa/oauth"
	"net/http"
)

func main() {
	http.HandleFunc("/twitter", TwitterAuthorization)
	http.ListenAndServe(":8081", nil)
}

func TwitterAuthorization(w http.ResponseWriter, r *http.Request) {
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
	res := client.GetRequestToken()
	http.Redirect(w, r, "https://api.twitter.com/oauth/authenticate?"+res, http.StatusFound)
}
