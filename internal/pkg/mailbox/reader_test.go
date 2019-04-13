package mailbox_test

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailchain/mailchain/internal/pkg/mailbox"
	"github.com/mailchain/mailchain/internal/pkg/testutil"
	"github.com/mailchain/mailchain/internal/pkg/testutil/mocks"
	"github.com/stretchr/testify/assert"
)

func TestReadMessage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	assert := assert.New(t)
	cases := []struct {
		name                   string
		txData                 []byte
		expectedID             string
		err                    string
		decrypterLocationCalls int
		decrypterLocationRet   []interface{}
		decrypterContentsCalls int
		decrypterFile          string
		decrypterContentsError []interface{}
	}{
		{"invalid protobuf prefix",
			testutil.MustHexDecodeString("08010f7365637265742d6c6f636174696f6e1a221620aff34d74dcb62c288b1a2f41a4852e82aff6c95e5c40c891299b3488b4340769"),
			"",
			"invalid encoding prefix",
			0,
			[]interface{}{[]byte("test://TestReadMessage/success-2204f3d89e5a"), nil},
			0,
			"",
			nil,
		},
		{"invalid protobuf format",
			testutil.MustHexDecodeString("5008010f7365637265742d6c6f636174696f6e1a221620aff34d74dcb62c288b1a2f41a4852e82aff6c95e5c40c891299b3488b4340769"),
			"",
			"could not unmarshal to data: proto: can't skip unknown wire type 7",
			0,
			[]interface{}{[]byte("test://TestReadMessage/success-2204f3d89e5a"), nil},
			0,
			"",
			nil,
		},
		{"fail decrypted location",
			testutil.MustHexDecodeString("500801120f7365637265742d6c6f636174696f6e1a221620aff34d74dcb62c288b1a2f41a4852e82aff6c95e5c40c891299b3488b4340769"),
			"",
			"could not decrypt location: could not decrypt",
			1,
			[]interface{}{nil, errors.New("could not decrypt")},
			0,
			"",
			nil,
		},
		{"success",
			testutil.MustHexDecodeString("500801120f7365637265742d6c6f636174696f6e1a221620aff34d74dcb62c288b1a2f41a4852e82aff6c95e5c40c891299b3488b4340769"),
			"002c47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
			"",
			1,
			[]interface{}{[]byte("test://TestReadMessage/success-2204f3d89e5a"), nil},
			1,
			"simple.golden.eml",
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			decrypter := mocks.NewMockDecrypter(mockCtrl)
			decrypter.EXPECT().Decrypt(gomock.Any()).Return(tc.decrypterLocationRet...).Times(tc.decrypterLocationCalls)
			decrypted, _ := ioutil.ReadFile("./testdata/" + tc.decrypterFile)
			decrypter.EXPECT().Decrypt(gomock.Any()).Return(decrypted, nil).Times(tc.decrypterContentsCalls)
			actual, err := mailbox.ReadMessage(tc.txData, decrypter)
			_ = actual
			if tc.err == "" {
				assert.NoError(err)
				assert.Equal(tc.expectedID, actual.ID.HexString())
			}
			if tc.err != "" {
				assert.EqualError(err, tc.err)
			}
		})
	}
}