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

type IBNetwork struct {
	// The name of IB network.
	Name string `json:"name"`
	// The pkeys for IB network. If not provided, it'll be generated automatically; the generated pkeys is only used by the service.
	PKey int32 `json:"pkey"`
	// Default false; create sharp allocation accordingly.
	EnableSharp bool `json:"-"`
	// The GUID list of the IB network.
	GUIDs []string `json:"guids"`
	// Default 2k; one of 2k or 4k; the MTU of the services.
	MTU int32 `json:"mtu"`
	// Default false
	IPOverIB bool `json:"ip_over_ib"`
	// Default false; store the PKey at index 0 of the PKey table of the GUID.
	Index0 bool `json:"index0"`
	// Default is None, value can be range from 0-15
	ServiceLevel int32 `json:"service_level"`
	// Default is None, can be one of the following: 2.5, 10, 30, 5, 20, 40, 60, 80, 120, 14, 56, 112, 168, 25, 100, 200, or 300
	RateLimit float64 `json:"rate_limit"`
}

type PKey struct {
	Partition string `json:"partition"`
	IpOverIb  bool   `json:"ip_over_ib"`
	Qos       struct {
		ServiceLevel int32   `json:"service_level"`
		MTU          int32   `json:"mtu_limit"`
		RateLimit    float64 `json:"rate_limit"`
	} `json:"qos_conf"`
	GUIDs []struct {
		GUID       string `json:"guid"`
		Index0     bool   `json:"index0"`
		Membership string `json:"membership"`
	}
}

type Field string

const (
	UnknownField Field = "unknown"
	QoSField     Field = "qos"
	GUIDField    Field = "guid"
)

type Strategy string

const (
	UnknownStrategy Strategy = "unknown"
	AddStrategy     Strategy = "add"
	DeleteStrategy  Strategy = "delete"
	SetStrategy     Strategy = "set"
)
