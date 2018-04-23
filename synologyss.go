package synoss

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

type Client struct {
	URI      string
	SID      string
	Account  string
	Password string
	APILIST  string
}

func New() *Client {
	var s Client

	return &s
}

func (s *Client) Connect(uri string) error {
	fmt.Println("Connecting")
	s.URI = uri

	apiList, err := s.getAPILIST()
	s.APILIST = gjson.Get(apiList, "data").String()

	return err
}

func (s *Client) Raw(api, method string, p map[string]string) (result string, err error) {
	retval, err := s.RawByte(api, method, p)
	result = string(retval)
	if err != nil {
		return result, err
	}

	if gjson.Get(result, "error").String() != "" {
		return "", errors.New(string(result))
	}

	return gjson.Get(result, "data").String(), nil
}

// GET /webapi/entry.cgi? eventId=5753&version="4"&mountId=0&api="SYNO.SurveillanceStation.Event"&analyevent=false&method ="Download"

func (s *Client) RawByte(api, method string, p map[string]string) (result []byte, err error) {
	v := url.Values{}
	v.Set("api", api)
	v.Set("method", method)

	for key, value := range p {
		v.Set(key, value)
	}

	if s.SID != "" {
		v.Set("_sid", s.SID)
	}

	apiPath := gjson.Get(s.APILIST, escapeDots(api)+".path").String()
	fullURI := s.URI + "/webapi/" + apiPath + "?" + v.Encode()

	// fmt.Println(fullURI)

	resp, err := http.Get(fullURI)
	if err != nil {
		return []byte(""), err
	}

	// fmt.Println(resp.Status)

	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes, nil
}

func (s *Client) Login(account, password string) error {
	p := make(map[string]string)

	p["version"] = "2"
	p["account"] = account
	p["passwd"] = password
	p["session"] = "SurveillanceStation"

	res, err := s.Raw("SYNO.API.Auth", "Login", p)
	if err != nil {
		return err
	}

	s.SID = gjson.Get(res, "sid").String()
	if s.SID == "" {
		return errors.New("sid not found in Login response")
	}

	return nil
}

func (s *Client) Logout() error {
	p := make(map[string]string)

	p["version"] = "2"

	_, err := s.Raw("SYNO.API.Auth", "Logout", p)
	if err != nil {
		return err
	}

	return nil
}

func escapeDots(buf string) string {
	return strings.Replace(buf, `.`, `\.`, -1)
}

func (s *Client) getAPILIST() (apiList string, apiErr error) {
	v := url.Values{}
	v.Set("api", "SYNO.API.Info")
	v.Set("version", "1")
	v.Set("query", "All")
	v.Set("method", "Query")

	fullURI := s.URI + "/webapi/query.cgi?" + v.Encode()

	resp, err := http.Get(fullURI)
	if err != nil {
		return "", err
	}

	retval, _ := ioutil.ReadAll(resp.Body)
	apiList = string(retval)

	if gjson.Get(string(apiList), "error").String() != "" {
		return "", errors.New(string(apiList))
	}

	return
}
