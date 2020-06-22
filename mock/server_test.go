package mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestServer(t *testing.T){
	api_url:="http://0.0.0.0:8080/command"
	data:=url.Values{"target": {"Value"}, "content": {"123"}}
	resp, err := http.PostForm(api_url,data)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	fmt.Println(string(body))
}
