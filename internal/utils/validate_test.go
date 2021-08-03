/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestObj struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func ReadFile(t *testing.T, file string) []byte {
	wd, _ := os.Getwd()
	dir := wd[:strings.Index(wd, "internal")]
	path := filepath.Join(dir, "test/testdata/", file)
	bs, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	return bs
}

func TestJsonSchemaValidator_Validate(t *testing.T) {
	tests := []struct {
		givePath        string
		giveObj         interface{}
		wantNewErr      error
		wantValidateErr []error
	}{
		{
			givePath: "validate_test.json",
			giveObj: TestObj{
				Name:  "lessName",
				Email: "too long name greater than 10",
				Age:   12,
			},
			wantValidateErr: []error{
				fmt.Errorf("name: String length must be greater than or equal to 10\nemail: String length must be less than or equal to 10"),
				fmt.Errorf("email: String length must be less than or equal to 10\nname: String length must be greater than or equal to 10"),
			},
		},
	}

	for _, tc := range tests {
		bs := ReadFile(t, tc.givePath)
		v, err := NewJsonSchemaValidator(string(bs))
		if err != nil {
			assert.Equal(t, tc.wantNewErr, err)
			continue
		}
		err = v.Validate(tc.giveObj)
		ret := false
		for _, wantErr := range tc.wantValidateErr {
			if wantErr.Error() == err.Error() {
				ret = true
			}
		}
		assert.True(t, ret)
	}
}
