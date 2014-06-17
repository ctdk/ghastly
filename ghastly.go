//Package ghastly provides a golang interface for using the Fastly
//(http://www.fastly.com) CDN's API.
package ghastly

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"
)

type Client struct {
	ApiKey   string
	User     string
	Password string
	BaseUrl  string
	Http     *http.Client
}

type Ghastly struct {
	*Client
}

// Initialize a new ghastly object, create the HTTP client, and log in.
func New(opts map[string]string) (*Ghastly, error) {
	g := &Ghastly{}
	client, err := login(opts["user"], opts["password"], opts["base_url"])
	if err != nil {
		return nil, err
	}
	g.Client = client
	return g, nil
}

func login(username, password string, base_url string) (*Client, error) {
	values := make(url.Values)
	values.Set("user", username)
	values.Set("password", password)
	client := new(Client)
	if base_url != "" {
		client.BaseUrl = base_url
	} else {
		client.BaseUrl = "https://api.fastly.com"
	}
	jar, jerr := cookiejar.New(nil)
	if jerr != nil {
		return nil, jerr
	}
	client.Http = &http.Client{Jar: jar}
	resp, err := client.PostForm("/login", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error logging in: %s", resp.Status)
		return nil, err
	}
	return client, err
}

// Convenience wrapper around http.Client.Get.
func (c *Client) Get(url string) (*http.Response, error) {
	resp, err := c.Http.Get(c.makeURL(url))
	if err != nil {
		return nil, err
	}
	if err = c.checkRespErr(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Convenience wrapper for GET requests with query parameters.
func (c *Client) GetParams(baseURL string, queryParams map[string]string) (*http.Response, error) {
	u, err := url.Parse(c.makeURL(baseURL))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for k, v := range queryParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	resp, err := c.Http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if err = c.checkRespErr(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Convenience wrapper around http.Client.PostForm.
func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	resp, err := c.Http.PostForm(c.makeURL(url), data)
	if err != nil {
		return nil, err
	}
	if err = c.checkRespErr(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Post a form to the server with a map[string]string of paramaters.
func (c *Client) PostFormParams(url string, params map[string]string) (*http.Response, error) {
	values := c.makeValues(params)
	resp, err := c.Http.PostForm(c.makeURL(url), values)
	if err != nil {
		return nil, err
	}
	if err = c.checkRespErr(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Convenience wrapper around http.Client.Post.
func (c *Client) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	resp, err := c.Http.Post(c.makeURL(url), bodyType, body)
	if err != nil {
		return nil, err
	}
	if err = c.checkRespErr(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Convenience wrapper for DELETE requests.
func (c *Client) Delete(url string) (*http.Response, error) {
	request, err := http.NewRequest("DELETE", c.makeURL(url), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Http.Do(request)
	if err != nil {
		return nil, err
	}
	if err = c.checkRespErr(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Convenience wrapper for PUT requests.
func (c *Client) Put(url string, data url.Values, contentType ...string) (*http.Response, error) {
	bodyStr := data.Encode()
	request, err := http.NewRequest("PUT", c.makeURL(url), strings.NewReader(bodyStr))
	request.Header.Set("content-type", setContentType(cType))
	if err != nil {
		return nil, err
	}
	resp, err := c.Http.Do(request)
	if err != nil {
		return nil, err
	}
	if err = c.checkRespErr(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Convenience wrapper for PUT requests, taking a map of strings to create the
// url.Values to send.
func (c *Client) PutParams(url string, params map[string]string) (*http.Response, error) {
	values := c.makeValues(params)
	return c.Put(url, values, "application/x-www-form-urlencoded")
}

func setContentType(contentType []string) string {
	var cType string
	if contentType != nil {
		cType = contentType[0]
	} else {
		cType = "application/json"
	}
	return cType
}

func (c *Client) makeURL(url string) string {
	url = path.Clean(url)
	if !path.IsAbs(url) {
		url = fmt.Sprintf("/%s", url)
	}
	return fmt.Sprintf("%s%s", c.BaseUrl, url)
}

func (c *Client) makeValues(params map[string]string) url.Values {
	values := make(url.Values)
	for k, v := range params {
		values.Set(k, v)
	}
	return values
}

func (c *Client) checkRespErr(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		rdata, err := ParseJson(resp.Body)
		if err != nil {
			return err
		}
		detail, _ := rdata["detail"].(string)
		err = fmt.Errorf("%s :: %s %s", resp.Status, rdata["msg"].(string), detail)
		return err
	}
	return nil
}

func ParseJson(data io.ReadCloser) (map[string]interface{}, error) {
	respData := make(map[string]interface{})
	dec := json.NewDecoder(data)
	if err := dec.Decode(&respData); err != nil {
		return nil, err
	}
	return respData, nil
}
