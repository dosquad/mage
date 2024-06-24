package helper

import (
	"github.com/go-resty/resty/v2"
)

type RestyOpt func(*resty.Client)

func AddAuthToken(token string) RestyOpt {
	return func(c *resty.Client) {
		c.SetHeader(
			"Authorization",
			"token "+token,
		)
	}
}

func AddHeader(key, value string) RestyOpt {
	return func(c *resty.Client) {
		c.SetHeader(key, value)
	}
}
