package argon

import (
	"strings"
	"testing"

	"src.goblgobl.com/tests"
	"src.goblgobl.com/tests/assert"
)

func Test_Hash_And_Compare(t *testing.T) {
	Insecure()
	for i := 0; i < 40; i++ {
		plainText := tests.Generator.String()
		h, err := Hash(plainText)
		assert.Nil(t, err)

		ok, err := Compare(plainText, h)
		assert.Nil(t, err)
		assert.True(t, ok)

		ok, err = Compare(plainText+"!", h)
		assert.Nil(t, err)
		assert.False(t, ok)
	}
}

func Test_Compare_Invalid_Hash(t *testing.T) {
	_, err := Compare("plain", "")
	assert.Equal(t, err.Error(), "Invalid hash length")

	_, err = Compare("plain", strings.Repeat("a", len(encodedHeader)))
	assert.Equal(t, err.Error(), "Invalid hash prefix")

	// wrong version
	_, err = Compare("plain", "$argon2id$v=18$m=1024,t=1,p=1$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash prefix")

	// wrong m= parameters
	_, err = Compare("plain", "$argon2id$v=19$X=1024,t=1,p=1$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash memory header")
	_, err = Compare("plain", "$argon2id$v=19$m=,t=1,p=1$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash memory parameter")
	_, err = Compare("plain", "$argon2id$v=19$m=a,t=1,p=1$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash memory parameter")

	// wrong t= parameters
	_, err = Compare("plain", "$argon2id$v=19$m=1024,X=1,p=1$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash time header")
	_, err = Compare("plain", "$argon2id$v=19$m=1024,t=,p=1$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash time parameter")
	_, err = Compare("plain", "$argon2id$v=19$m=1024,t=a,p=1$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash time parameter")

	// wrong p= parameters
	_, err = Compare("plain", "$argon2id$v=19$m=1024,t=1,X=1$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash parallelism header")
	_, err = Compare("plain", "$argon2id$v=19$m=1024,t=1,p=$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash parallelism parameter")
	_, err = Compare("plain", "$argon2id$v=19$m=1024,t=1,p=a$OgpTGmRelTa9Iuea6yy/4w$Q65Wg7OCGs+v4Uf+Zgs9SQ")
	assert.Equal(t, err.Error(), "Invalid hash parallelism parameter")

	_, err = Compare("plain", "$argon2id$v=19$m=1024,t=1,p=1!")
	assert.Equal(t, err.Error(), "Invalid hash header separator")

	_, err = Compare("plain", string(encodedHeader))
	assert.Equal(t, err.Error(), "Invalid hash data separator")

	_, err = Compare("plain", "$argon2id$v=19$m=1024,t=1,p=1$_$_")
	assert.Equal(t, err.Error(), "Invalid hash salt")

	ok, err := Compare("plain", "$argon2id$v=19$m=1024,t=1,p=1$YQ$YQ")
	assert.Nil(t, err)
	assert.False(t, ok)
}
