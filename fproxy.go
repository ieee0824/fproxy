package fproxy

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	secure "github.com/ieee0824/secure-string"
)

const endpoint = "http://pubproxy.com/api/proxy"

// Setting parameter
type Setting struct {
	API        secure.String
	Format     string
	Level      string
	Type       string
	LastCheck  *int
	Speed      *int
	Limit      *int
	Country    string
	NotCountry string
	Google     *bool
	HTTPS      *bool
	Post       *bool
	UserAgent  *bool
	Cookies    *bool
	Referer    *bool
}

func (s *Setting) Query() url.Values {
	ret := url.Values{}

	if s.API != "" {
		ret.Set("api", string(s.API))
	}

	if s.Format != "" {
		ret.Set("format", s.Format)
	}

	if s.Level != "" {
		ret.Set("level", s.Level)
	}

	if s.Type != "" {
		ret.Set("type", s.Type)
	}

	if s.LastCheck != nil {
		ret.Set("last_check", strconv.Itoa(*s.LastCheck))
	}

	if s.Speed != nil {
		ret.Set("speed", strconv.Itoa(*s.Speed))
	}

	if s.Limit != nil {
		ret.Set("limit", strconv.Itoa(*s.Limit))
	}

	if s.Country != "" {
		ret.Set("country", s.Country)
	}

	if s.NotCountry != "" {
		ret.Set("not_country", s.Country)
	}

	if s.Google != nil {
		ret.Set("google", strconv.FormatBool(*s.Google))
	}

	if s.HTTPS != nil {
		ret.Set("https", strconv.FormatBool(*s.HTTPS))
	}

	if s.Post != nil {
		ret.Set("post", strconv.FormatBool(*s.Post))
	}

	if s.UserAgent != nil {
		ret.Set("user_agent", strconv.FormatBool(*s.UserAgent))
	}

	if s.Cookies != nil {
		ret.Set("cookies", strconv.FormatBool(*s.Cookies))
	}

	if s.Referer != nil {
		ret.Set("referer", strconv.FormatBool(*s.Referer))
	}

	return ret
}

func NewSetting() *Setting {
	ret := &Setting{
		Google:    new(bool),
		LastCheck: new(int),
	}
	*ret.Google = true
	*ret.LastCheck = 10
	ret.Country = "US"
	return ret
}

// CLBool is C like boolean
type CLBool int

func (c CLBool) GoBool() bool {
	if c == 1 {
		return true
	}
	return false
}

type ProxyInfo struct {
	Data []struct {
		IPPoort     string `json:"ipPort"`
		IP          string `json:"ip"`
		Port        string `json:"port"`
		Country     string `json:"country`
		LastChecked string `json:"last_checked"`
		ProxyLevel  string `json:"proxy_level"`
		Type        string `json:"type"`
		Speed       string `json:"speed"`
		Support     struct {
			HTTPS     CLBool `json:"https"`
			Get       CLBool `json:"get"`
			Post      CLBool `json:"post"`
			Cookies   CLBool `json:"cookies"`
			Referer   CLBool `json:"referer"`
			UserAgent CLBool `json:"user_agent"`
			Google    CLBool `jsoon:"google"`
		} `json:"support"`
	} `json:"data"`
}

func NewProxy(settings ...*Setting) (*http.Transport, error) {
	var s *Setting
	if settings == nil {
		s = NewSetting()
	} else {
		s = settings[0]
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	u.RawQuery = s.Query().Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	info := new(ProxyInfo)
	if err := json.NewDecoder(resp.Body).Decode(info); err != nil {
		return nil, err
	}

	if len(info.Data) == 0 {
		return nil, errors.New("no proxy error")
	}

	proxyURL, err := url.Parse(info.Data[0].Type + info.Data[0].IPPoort)
	if err != nil {
		return nil, err
	}

	tr := new(http.Transport)
	tr.Proxy = http.ProxyURL(proxyURL)

	return tr, nil
}

func NewClient() (*http.Client, error) {
	proxy, err := NewProxy()
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: proxy}, nil
}
