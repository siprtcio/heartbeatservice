package utils

import (
	"github.com/go-resty/resty/v2"
	"time"
)

var restClient = resty.New()

func Get(data map[string]string, url string) (int,[]byte,error) {
	restClient.SetTimeout(2 * time.Second)
	if data == nil{
		resp, err := restClient.R().
			SetHeader("Accept", "application/json").
			Get(url)
		if resp == nil{
			return 400,nil,err
		}
		return resp.StatusCode(),resp.Body(),err
	}
	resp, err := restClient.R().
		SetQueryParams(data).
		SetHeader("Accept", "application/json").
		Get(url)
	if resp == nil{
		return 400,nil,err
	}
	return resp.StatusCode(),resp.Body(),err
}

func Post(data map[string]string, url string) (int,[]byte,error) {
	restClient.SetTimeout(2 * time.Second)
	if data == nil{
		resp, err := restClient.R().
			SetHeader("Accept", "application/json").
			Post(url)
		if resp == nil{
			return 400,nil,err
		}
		return resp.StatusCode(),resp.Body(),err
	}
	resp, err := restClient.R().
		SetBody(data).
		SetHeader("Accept", "application/json").
		Get(url)
	if resp == nil{
		return 400,nil,err
	}
	return resp.StatusCode(),resp.Body(),err
}