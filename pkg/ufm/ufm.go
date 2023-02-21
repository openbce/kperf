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
	"encoding/json"
	"fmt"
	"strings"

	envv6 "github.com/caarlos0/env/v6"
)

type UFM struct {
	PluginName  string
	SpecVersion string
	conf        UFMConfig
	client      UFMClient
}

const (
	pluginName  = "ufm"
	specVersion = "1.0"
	httpsProto  = "https"
)

type UFMConfig struct {
	Username    string `env:"UFM_USERNAME"`    // Username of ufm
	Password    string `env:"UFM_PASSWORD"`    // Password of ufm
	Address     string `env:"UFM_ADDRESS"`     // IP address or hostname of ufm server
	Port        int    `env:"UFM_PORT"`        // REST API port of ufm
	HTTPSchema  string `env:"UFM_HTTP_SCHEMA"` // http or https
	Certificate string `env:"UFM_CERTIFICATE"` // Certificate of ufm
}

func NewUFM() (*UFM, error) {
	ufmConf := UFMConfig{}
	if err := envv6.Parse(&ufmConf); err != nil {
		return nil, err
	}

	if ufmConf.Username == "" || ufmConf.Password == "" || ufmConf.Address == "" {
		return nil, fmt.Errorf("missing one or more required fileds for ufm [\"username\", \"password\", \"address\"]")
	}

	// set httpSchema and port to ufm default if missing
	ufmConf.HTTPSchema = strings.ToLower(ufmConf.HTTPSchema)
	if ufmConf.HTTPSchema == "" {
		ufmConf.HTTPSchema = httpsProto
	}
	if ufmConf.Port == 0 {
		if ufmConf.HTTPSchema == httpsProto {
			ufmConf.Port = 443
		} else {
			ufmConf.Port = 80
		}
	}

	isSecure := strings.EqualFold(ufmConf.HTTPSchema, httpsProto)
	auth := &BasicAuth{Username: ufmConf.Username, Password: ufmConf.Password}
	client, err := NewClient(isSecure, auth, ufmConf.Certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to create http ufmclient err: %v", err)
	}
	return &UFM{PluginName: pluginName,
		SpecVersion: specVersion,
		conf:        ufmConf,
		client:      client}, nil
}

func (u *UFM) Name() string {
	return u.PluginName
}

func (u *UFM) Spec() string {
	return u.SpecVersion
}

func (u *UFM) Validate() error {
	_, err := u.client.Get(u.buildURL("/ufmRest/app/ufm_version"))

	if err != nil {
		return fmt.Errorf("failed to connect to ufm subnet manager: %v", err)
	}

	return nil
}

func (u *UFM) GetIBNetwork(pkey int32) (*IBNetwork, *UFMError) {
	if !isPKeyValid(pkey) {
		return nil, &UFMError{
			Code:    InvalidPKeyErr,
			Message: fmt.Sprintf("invalid pkey 0x%04X, out of range 0x0001 - 0xFFFE", pkey),
		}
	}

	res := &PKey{}
	path := fmt.Sprintf("/ufmRest/resources/pkeys/0x%x?guids_data=true&qos_conf=true", pkey)
	if data, err := u.client.Get(u.buildURL(path)); err != nil {
		return nil, &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to get pkey 0x%04X with error: %v", pkey, err),
		}
	} else {
		if string(data) == "{}" {
			return nil, &UFMError{
				Code:    NotFoundErr,
				Message: "NotFound",
			}
		}

		if err := json.Unmarshal(data, res); err != nil {
			return nil, &UFMError{
				Code:    UnknownErr,
				Message: fmt.Sprintf("failed to unmarshal pkey 0x%04X with error: %v", pkey, err),
			}
		}
	}

	return buildIBNetwork(pkey, res), nil
}

func (u *UFM) CreateIBNetwork(ib *IBNetwork) *UFMError {
	return u.addGuids(ib)
}

func (u *UFM) patchQoS(ib *IBNetwork, _ Strategy) *UFMError {
	pkey, _ := BuidPKey(ib.PKey)

	qos := struct {
		PKey         string  `json:"pkey"`
		ServiceLevel int32   `json:"service_level"`
		MTU          int32   `json:"mtu_limit"`
		RateLimit    float64 `json:"rate_limit"`
	}{
		PKey:         pkey,
		RateLimit:    ib.RateLimit,
		MTU:          ib.MTU,
		ServiceLevel: ib.ServiceLevel,
	}

	qosData, err := json.Marshal(qos)
	if err != nil {
		return &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to marshal IB with error: %v", err),
		}
	}

	if _, err := u.client.Post(u.buildURL("/ufmRest/resources/pkeys/qos_conf"), qosData); err != nil {
		return &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to update PKey 0x%04X with error: %v", ib.PKey, err),
		}
	}

	return nil
}

func (u *UFM) listQoS() (map[string]PKey, *UFMError) {
	if data, err := u.client.Get(u.buildURL("/ufmRest/resources/pkeys?qos_conf=true")); err != nil {
		return nil, &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to list pkey with error: %v", err),
		}
	} else {
		pkeys := map[string]PKey{}
		if err := json.Unmarshal(data, &pkeys); err != nil {
			return nil, &UFMError{
				Code:    UnknownErr,
				Message: fmt.Sprintf("failed to unmarshal pkey with error: %v", err),
			}
		}

		return pkeys, nil
	}
}

func (u *UFM) listGUID() (map[string]PKey, *UFMError) {
	if data, err := u.client.Get(u.buildURL("/ufmRest/resources/pkeys?guids_data=true")); err != nil {
		return nil, &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to list pkey with error: %v", err),
		}
	} else {
		pkeys := map[string]PKey{}
		if err := json.Unmarshal(data, &pkeys); err != nil {
			return nil, &UFMError{
				Code:    UnknownErr,
				Message: fmt.Sprintf("failed to unmarshal pkey with error: %v", err),
			}
		}

		return pkeys, nil
	}
}

func (u *UFM) ListIBNetwork() ([]*IBNetwork, *UFMError) {
	var res []*IBNetwork
	qos, ufmErr := u.listQoS()
	if ufmErr != nil {
		return nil, ufmErr
	}
	guids, ufmErr := u.listGUID()
	if ufmErr != nil {
		return nil, ufmErr
	}

	for pkeyStr, param := range qos {
		if ids, found := guids[pkeyStr]; found {
			param.GUIDs = ids.GUIDs
		}

		pkey, err := ParsePkey(pkeyStr)
		if err != nil {
			continue
		}
		res = append(res, buildIBNetwork(pkey, &param))
	}

	return res, nil
}

func (u *UFM) DeleteIBNetwork(pkey int32) *UFMError {
	path := fmt.Sprintf("/ufmRest/resources/pkeys/0x%x", pkey)
	if _, err := u.client.Delete(u.buildURL(path)); err != nil {
		return &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to delete PKey 0x%04X with error: %v", pkey, err),
		}
	}

	return nil
}

func (u *UFM) Patch(ib *IBNetwork, field Field, op Strategy) *UFMError {
	switch field {
	case GUIDField:
		return u.patchGUIDs(ib, op)
	case QoSField:
		return u.patchQoS(ib, op)
	}

	return &UFMError{
		Code:    UnknownErr,
		Message: "Invalid field",
	}
}

func (u *UFM) buildURL(path string) string {
	return fmt.Sprintf("%s://%s:%d%s", u.conf.HTTPSchema, u.conf.Address, u.conf.Port, path)
}

func (u *UFM) patchGUIDs(ib *IBNetwork, op Strategy) *UFMError {
	switch op {
	case AddStrategy:
		return u.addGuids(ib)
	case DeleteStrategy:
		return u.deleteGuids(ib)
	case SetStrategy:
		return u.addGuids(ib)
	}

	return nil
}

func (u *UFM) deleteGuids(ib *IBNetwork) *UFMError {
	pkey, _ := BuidPKey(ib.PKey)

	guidList := struct {
		PKey  string   `json:"pkey"`
		GUIDs []string `json:"guids"`
	}{
		PKey:  pkey,
		GUIDs: ib.GUIDs,
	}

	data, err := json.Marshal(guidList)
	if err != nil {
		return &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to marshal IB with error: %v", err),
		}
	}

	if _, err := u.client.Post(u.buildURL("/ufmRest/actions/remove_guids_from_pkey"), data); err != nil {
		return &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to create PKey 0x%04X with error: %v", ib.PKey, err),
		}
	}

	return nil
}

func (u *UFM) addGuids(ib *IBNetwork) *UFMError {
	pkey, _ := BuidPKey(ib.PKey)

	guidList := struct {
		PKey       string   `json:"pkey"`
		IPoIB      bool     `json:"ip_over_ib"`
		Index0     bool     `json:"index0"`
		GUIDs      []string `json:"guids"`
		Membership string   `json:"membership"`
	}{
		PKey:       pkey,
		IPoIB:      ib.IPOverIB,
		Membership: "full",
		Index0:     ib.Index0,
		GUIDs:      ib.GUIDs,
	}

	data, err := json.Marshal(guidList)
	if err != nil {
		return &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to marshal IB with error: %v", err),
		}
	}

	if _, err := u.client.Post(u.buildURL("/ufmRest/resources/pkeys"), data); err != nil {
		return &UFMError{
			Code:    UnknownErr,
			Message: fmt.Sprintf("failed to create PKey 0x%04X with error: %v", ib.PKey, err),
		}
	}

	return nil
}
