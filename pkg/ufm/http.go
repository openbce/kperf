/*
Copyright 2023 The openBCE Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ufm

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

type UFMClient interface {
	Get(url string) ([]byte, *UFMError)
	Post(url string, body []byte) ([]byte, *UFMError)
	Put(url string, body []byte) ([]byte, *UFMError)
	Delete(url string) ([]byte, *UFMError)
}

type BasicAuth struct {
	Username string
	Password string
}

type ufmclient struct {
	basicAuth  *BasicAuth
	httpClient *http.Client
}

func NewClient(isSecure bool, basicAuth *BasicAuth, cert string) (UFMClient, *UFMError) {
	log.Debug().Msgf("creating http ufmclient, isSecure %v, basicAuth %+v, cert %s", isSecure, basicAuth, cert)
	if basicAuth == nil {
		return nil, &UFMError{
			Code:    AuthErr,
			Message: fmt.Sprintf("invalid basicAuth value %v", basicAuth),
		}
	}
	httpClient := &http.Client{Transport: http.DefaultTransport}
	if isSecure {
		if cert == "" {
			/* #nosec */
			httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM([]byte(cert))
			httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{RootCAs: caCertPool}
		}
	}

	return &ufmclient{basicAuth: basicAuth, httpClient: httpClient}, nil
}

func (c *ufmclient) Get(url string) ([]byte, *UFMError) {
	log.Debug().Msgf("Http ufmclient GET: url %s", url)
	return c.executeRequest(http.MethodGet, url, nil)
}

func (c *ufmclient) Post(url string, body []byte) ([]byte, *UFMError) {
	log.Debug().Msgf("Http ufmclient POST: url %s,  body %s", url, string(body))
	return c.executeRequest(http.MethodPost, url, body)
}

func (c *ufmclient) Put(url string, body []byte) ([]byte, *UFMError) {
	log.Debug().Msgf("Http ufmclient PUT: url %s,  body %s", url, string(body))
	return c.executeRequest(http.MethodPut, url, body)
}

func (c *ufmclient) Delete(url string) ([]byte, *UFMError) {
	log.Debug().Msgf("Http ufmclient DELETE: url %s", url)
	return c.executeRequest(http.MethodDelete, url, nil)
}

func (c *ufmclient) createRequest(method, url string, body io.Reader) (*http.Request, *UFMError) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to create request object %v", err),
		}
	}

	req.SetBasicAuth(c.basicAuth.Username, c.basicAuth.Password)

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return req, nil
}

func (c *ufmclient) executeRequest(method, url string, body []byte) ([]byte, *UFMError) {
	req, ufmErr := c.createRequest(method, url, bytes.NewBuffer(body))
	if ufmErr != nil {
		return nil, ufmErr
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &UFMError{
			Code:    UnknownErr,
			Message: err.Error(),
		}
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case http.StatusOK:
		return responseBody, nil
	case http.StatusNotFound:
		return nil, &UFMError{
			Code:    NotFoundErr,
			Message: http.StatusText(http.StatusNotFound),
		}
	}

	return nil, &UFMError{
		Code: UnknownErr,
		Message: fmt.Sprintf("http status (%d): %s",
			resp.StatusCode, http.StatusText(resp.StatusCode)),
	}
}
