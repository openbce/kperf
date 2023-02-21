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
	"fmt"
	"net"
	"strings"
)

// isPKeyValid check if the pkey is in the valid (15bits long)
func isPKeyValid(pkey int32) bool {
	return pkey == (pkey & 0x7fff)
}

// GuidToString return string guid from HardwareAddr
func guidToString(guidAddr net.HardwareAddr) string {
	return strings.Replace(guidAddr.String(), ":", "", -1)
}

func ParsePkey(pkeyStr string) (int32, error) {
	var pkey int32
	if _, err := fmt.Sscanf(pkeyStr, "0x%x", &pkey); err != nil {
		return -1, err
	}

	if !isPKeyValid(pkey) {
		return -1, fmt.Errorf("invalid pkey")
	}

	return pkey, nil
}

func BuidPKey(pkey int32) (string, error) {
	if !isPKeyValid(pkey) {
		return "", fmt.Errorf("invalid pkey")
	}

	res := fmt.Sprintf("0x%x", pkey)
	return res, nil
}

func buildIBNetwork(pkey int32, param *PKey) *IBNetwork {
	var guids []string
	index0 := false

	for _, id := range param.GUIDs {
		guids = append(guids, id.GUID)
		index0 = id.Index0
	}

	return &IBNetwork{
		Name:         param.Partition,
		PKey:         pkey,
		EnableSharp:  false,
		GUIDs:        guids,
		MTU:          param.Qos.MTU,
		IPOverIB:     param.IPoIB,
		Index0:       index0,
		ServiceLevel: param.Qos.ServiceLevel,
		RateLimit:    param.Qos.RateLimit,
	}
}

func ParseField(f string) Field {
	switch f {
	case string(QoSField):
		return QoSField
	case string(GUIDField):
		return GUIDField
	}

	return UnknownField
}

func ParseStrategy(s string) Strategy {
	switch s {
	case string(AddStrategy):
		return AddStrategy
	case string(DeleteStrategy):
		return DeleteStrategy
	case string(SetStrategy):
		return SetStrategy
	}

	return UnknownStrategy
}
