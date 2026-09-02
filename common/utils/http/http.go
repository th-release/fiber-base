package httpclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func Get[T any](client *resty.Client, url string, query, headers map[string]string) (*resty.Response, T, error) {
	var result T
	req := client.R().SetResult(&result)
	applyQuery(req, query)
	applyHeaders(req, headers)
	resp, err := req.Get(url)
	if err != nil {
		return nil, result, err
	}
	return resp, result, nil
}

func Post[T any](client *resty.Client, url string, query, headers map[string]string, body interface{}) (*resty.Response, T, error) {
	return send[T](client, query, headers, body, func(req *resty.Request) (*resty.Response, error) {
		return req.Post(url)
	})
}

func Put[T any](client *resty.Client, url string, query, headers map[string]string, body interface{}) (*resty.Response, T, error) {
	return send[T](client, query, headers, body, func(req *resty.Request) (*resty.Response, error) {
		return req.Put(url)
	})
}

func Patch[T any](client *resty.Client, url string, query, headers map[string]string, body interface{}) (*resty.Response, T, error) {
	return send[T](client, query, headers, body, func(req *resty.Request) (*resty.Response, error) {
		return req.Patch(url)
	})
}

func Delete[T any](client *resty.Client, url string, query, headers map[string]string, body interface{}) (*resty.Response, T, error) {
	return send[T](client, query, headers, body, func(req *resty.Request) (*resty.Response, error) {
		return req.Delete(url)
	})
}

func CheckStatus(resp *resty.Response) error {
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode(), resp.String())
	}
	return nil
}

func send[T any](client *resty.Client, query, headers map[string]string, body interface{}, do func(*resty.Request) (*resty.Response, error)) (*resty.Response, T, error) {
	var result T
	req := client.R().SetResult(&result)
	applyQuery(req, query)
	applyHeaders(req, headers)
	if err := applyBody(req, headers, body); err != nil {
		return nil, result, err
	}
	resp, err := do(req)
	if err != nil {
		return nil, result, err
	}
	return resp, result, nil
}

func applyQuery(req *resty.Request, query map[string]string) {
	if len(query) > 0 {
		req.SetQueryParams(query)
	}
}

func applyHeaders(req *resty.Request, headers map[string]string) {
	if len(headers) > 0 {
		req.SetHeaders(headers)
	}
}

func applyBody(req *resty.Request, headers map[string]string, body interface{}) error {
	if body == nil {
		return nil
	}
	ct := headers["Content-Type"]
	if ct == "application/x-www-form-urlencoded" || ct == "application/form-data" {
		m, ok := body.(map[string]interface{})
		if !ok {
			return fmt.Errorf("form-urlencoded body must be map[string]interface{}")
		}
		params := make(map[string]string, len(m))
		for k, v := range m {
			params[k] = fmt.Sprintf("%v", v)
		}
		req.SetFormData(params)
		return nil
	}
	req.SetBody(body)
	return nil
}
