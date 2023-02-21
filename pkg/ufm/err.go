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

type ErrCode int32

const (
	UnknownErr     ErrCode = -1
	NotFoundErr    ErrCode = 1
	InvalidPKeyErr ErrCode = 2
	AuthErr        ErrCode = 3
)

type UFMError struct {
	Code    ErrCode
	Message string
}

func (u *UFMError) Error() string {
	return u.Message
}

func (u *UFMError) IsNotFound() bool {
	return u.Code == NotFoundErr
}
