package oauth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type Provider struct {
	RequestTokenUrl  string
	AuthorizationUrl string
	AccessTokenUrl   string
}

type Client struct {
	AuthProvider   Provider
	CallbackUrl    string
	ConsumerKey    string
	ConsumerSecret string
	RequestToken   string
	Signature      string
	Method         string
	Timestamp      string
	Version        string
	Verifier       string
}

type RequestSignature struct {
	Method        string
	BaseUrl       string
	RequestParams string
	OauthParams   string
	SignatureBase string
}

func NewClient(provider Provider, callbackUrl, consumerKey, consumerSecret string) *Client {
	//Instantiate client struct containing initial data
	return &Client{provider, callbackUrl, consumerKey, consumerSecret, "", "", "", "", "", ""}
}

func (c *Client) GetRequestToken() string {
	//make call for request token
	client := &http.Client{}
	httpMethod := "POST"
	req, err := http.NewRequest(httpMethod, c.AuthProvider.RequestTokenUrl, nil)
	if err != nil {
	}
	header, parameters := c.GenerateRequest()
	signedHeader := header + ",oauth_signature=\"" + url.QueryEscape(c.GenerateSignature(httpMethod, c.AuthProvider.RequestTokenUrl, parameters)) + "\""
	req.Header.Add("Authorization", signedHeader)
	resp, err := client.Do(req)
	if err != nil {
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	tr := string(bs)
	return tr
}

func (c *Client) GetAccessToken(verifier, token string) string {
	c.Verifier = verifier
	c.RequestToken = token
	client := &http.Client{}
	httpMethod := "POST"
	req, err := http.NewRequest(httpMethod, c.AuthProvider.AccessTokenUrl, nil)
	if err != nil {
	}
	header, parameters := c.GenerateRequest()
	signedHeader := header + ",oauth_signature=\"" + url.QueryEscape(c.GenerateSignature(httpMethod, c.AuthProvider.AccessTokenUrl, parameters)) + "\""
	req.Header.Add("Authorization", signedHeader)
	resp, err := client.Do(req)
	if err != nil {
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}

	fmt.Println("===============\n" + string(signedHeader) + "\n==================")

	fmt.Println("===============\n" + string(parameters) + "\n==================")
	fmt.Println("===============\n" + string(bs) + "\n==================")
	tr := string(bs)
	return tr
}

func (c *Client) GenerateSignature(method, baseUrl, headerContent string) string {
	parsedUrl, err := url.Parse(baseUrl)
	base := parsedUrl.Query()
	head, err := url.ParseQuery(headerContent)

	if err != nil {
	}
	//SORT STRINGS
	combined := make(map[string]string)
	for k, v := range base {
		combined[k] = v[0]
	}
	for k, v := range head {
		combined[k] = v[0]
	}
	var keys []string
	for k := range combined {
		keys = append(keys, k)
	}
	var params string
	sort.Strings(keys)
	for i, k := range keys {
		params += k + "=" + combined[k]
		if i != len(keys)-1 {
			params += "&"
		}
	}

	baseRequest := parsedUrl.Scheme + "://" + parsedUrl.Host + parsedUrl.Path
	sig := method + "&" + url.QueryEscape(baseRequest) + "&" + url.QueryEscape(params)
	fmt.Println(sig)
	sigKey := url.QueryEscape(c.ConsumerSecret) + "&"
	hashfun := hmac.New(sha1.New, []byte(sigKey))
	hashfun.Write([]byte(sig))
	rawsignature := hashfun.Sum(nil)
	base64signature := make([]byte, base64.StdEncoding.EncodedLen(len(rawsignature)))
	base64.StdEncoding.Encode(base64signature, rawsignature)
	return string(base64signature)
}

func (c *Client) GenerateRequest() (string, string) {
	time := generateTimestamp()
	nonce := url.QueryEscape(generateNonce())
	key := c.ConsumerKey
	callback := url.QueryEscape(c.CallbackUrl)
	header := "OAuth "
	header += "oauth_consumer_key=\"" + key + "\",oauth_nonce=\"" + nonce + "\",oauth_signature_method=\"HMAC-SHA1\",oauth_timestamp=\"" + time + "\",oauth_version=\"1.0\""

	parameters := "oauth_consumer_key=" + key + "&oauth_nonce=" + nonce + "&oauth_signature_method=HMAC-SHA1&oauth_timestamp=" + time + "&oauth_version=1.0"
	if c.Verifier == "" {
		header += ",oauth_callback=\"" + callback + "\""
		parameters += "&oauth_callback=" + url.QueryEscape(callback)
	}
	if c.Verifier != "" {
		header += ",oauth_verifier=\"" + c.Verifier + "\"" + ",oauth_token=\"" + c.RequestToken + "\""
		parameters += "&oauth_verifier=" + c.Verifier + "&oauth_token=" + c.RequestToken
	}
	return header, parameters
}

func generateNonce() string {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	uuid[8] = 0x80
	uuid[4] = 0x40
	return hex.EncodeToString(uuid)
}

func generateTimestamp() string {
	ts := time.Now().Unix()
	return strconv.Itoa(int(ts))
}
