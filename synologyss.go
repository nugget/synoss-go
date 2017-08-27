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

type SSS struct {
	URI      string
	SID      string
	Account  string
	Password string
	APILIST  string
}

func New() SSS {
	var s SSS

	return s
}

func (s *SSS) Connect(uri string) error {
	fmt.Println("Connecting")
	s.URI = uri

	if s.URI[len(s.URI)-1:] != "/" {
		s.URI = s.URI + "/"
	}
	s.URI = s.URI + "webapi/"

	apiList, err := s.getAPILIST()
	s.APILIST = gjson.Get(apiList, "data").String()

	return err
}

func (s *SSS) Raw(api, method string, p map[string]string) (result string, err error) {
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
	fullURI := s.URI + apiPath + "?" + v.Encode()

	resp, err := http.Get(fullURI)
	if err != nil {
		return "", err
	}

	retval, _ := ioutil.ReadAll(resp.Body)
	result = string(retval)

	if gjson.Get(result, "error").String() != "" {
		return "", errors.New(string(result))
	}

	return gjson.Get(result, "data").String(), nil
}

func (s *SSS) Login(account, password string) error {
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

func (s *SSS) Logout() error {
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

func (s *SSS) getAPILIST() (apiList string, apiErr error) {
	v := url.Values{}
	v.Set("api", "SYNO.API.Info")
	v.Set("version", "1")
	v.Set("query", "All")
	v.Set("method", "Query")

	fullURI := s.URI + "query.cgi?" + v.Encode()

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
