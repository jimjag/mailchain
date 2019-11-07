// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pubkey

import (
	"testing"

	"github.com/mailchain/mailchain/internal/encoding"
	"github.com/mailchain/mailchain/internal/testutil"
)

func TestEncodeByProtocol(t *testing.T) {
	type args struct {
		in       []byte
		protocol string
	}
	tests := []struct {
		name             string
		args             args
		wantEncoded      string
		wantEncodingType string
		wantErr          bool
	}{
		{
			"ethereum",
			args{
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				"ethereum",
			},
			"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
			encoding.TypeHex0XPrefix,
			false,
		},
		{
			"substrate",
			args{
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				"substrate",
			},
			"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
			encoding.TypeHex0XPrefix,
			false,
		},
		{
			"err",
			args{
				testutil.MustHexDecodeString("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
				"invalid",
			},
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncoded, gotEncodingType, err := EncodeByProtocol(tt.args.in, tt.args.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeByProtocol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEncoded != tt.wantEncoded {
				t.Errorf("EncodeByProtocol() gotEncoded = %v, want %v", gotEncoded, tt.wantEncoded)
			}
			if gotEncodingType != tt.wantEncodingType {
				t.Errorf("EncodeByProtocol() gotEncodingType = %v, want %v", gotEncodingType, tt.wantEncodingType)
			}
		})
	}
}