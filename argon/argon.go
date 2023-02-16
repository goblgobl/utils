package argon

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"src.goblgobl.com/utils"
	"src.goblgobl.com/utils/ascii"
)

var (
	b64       = base64.RawStdEncoding
	b64Strict = b64.Strict()

	// make these vars mostly so that we can make this faster in tests
	defaultTimeCost   uint32
	defaultMemoryCost uint32
	defaultParallaism uint8
	defaultKeyLength  uint32 = 16

	encodedPrefix = fmt.Sprintf("$argon2id$v=%d", argon2.Version)
	encodedHeader []byte
)

// so not thread safe
func Config(timeCost uint32, memoryCost uint32, parallelism uint8) {
	defaultTimeCost = timeCost
	defaultMemoryCost = memoryCost
	defaultParallaism = parallelism
	encodedHeader = []byte(fmt.Sprintf("%s$m=%d,t=%d,p=%d$", encodedPrefix, memoryCost, timeCost, parallelism))
}

func init() {
	Config(4, 64*1024, 2)
}

func Hash(plainText string) (string, error) {
	var salt [16]byte
	saltSlice := salt[:]
	_, err := rand.Read(saltSlice)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey(utils.S2B(plainText), saltSlice, defaultTimeCost, defaultMemoryCost, defaultParallaism, defaultKeyLength)

	// we're going to encode the salt into this scrap, write that into our builder
	// then reuse the scrap to encode the key

	encodedKeyLength := b64.EncodedLen(len(key))
	encodedSaltLength := 22 // salt is 16 bytes, which base64 encodes to 22 bytes

	scrapLength := encodedSaltLength
	if encodedKeyLength > encodedSaltLength {
		scrapLength = encodedKeyLength
	}
	scrap := make([]byte, scrapLength)

	var builder strings.Builder
	// +1 for the $ seperator
	builder.Grow(len(encodedHeader) + encodedSaltLength + 1 + encodedKeyLength)
	builder.Write(encodedHeader)

	b64.Encode(scrap, saltSlice)
	builder.Write(scrap[:encodedSaltLength])
	builder.WriteByte('$')

	b64.Encode(scrap, key)
	builder.Write(scrap[:encodedKeyLength])

	return builder.String(), nil
}

func Compare(plainText string, hash string) (bool, error) {
	if len(hash) < len(encodedHeader) {
		return false, errors.New("Invalid hash length")
	}
	if !strings.HasPrefix(hash, encodedPrefix) {
		return false, errors.New("Invalid hash prefix")
	}

	hash = hash[len(encodedPrefix):]
	if hash[:3] != "$m=" {
		return false, errors.New("Invalid hash memory header")
	}

	memoryCost, hash := ascii.Atoi(hash[3:])
	if memoryCost == 0 {
		return false, errors.New("Invalid hash memory parameter")
	}

	if hash[:3] != ",t=" {
		return false, errors.New("Invalid hash time header")
	}

	timeCost, hash := ascii.Atoi(hash[3:])
	if timeCost == 0 {
		return false, errors.New("Invalid hash time parameter")
	}

	if hash[:3] != ",p=" {
		return false, errors.New("Invalid hash parallelism header")
	}

	parallelism, hash := ascii.Atoi(hash[3:])
	if parallelism == 0 {
		return false, errors.New("Invalid hash parallelism parameter")
	}

	if hash[0] != '$' {
		return false, errors.New("Invalid hash header separator")
	}

	hash = hash[1:]
	seperator := strings.Index(hash, "$")
	if seperator == -1 {
		return false, errors.New("Invalid hash data separator")
	}

	salt, err := b64Strict.DecodeString(hash[:seperator])
	if err != nil {
		return false, errors.New("Invalid hash salt")
	}

	key, err := b64Strict.DecodeString(hash[seperator+1:])
	if err != nil {
		return false, errors.New("Invalid hash key")
	}
	keyLen := len(key)
	reKey := argon2.IDKey(utils.S2B(plainText), salt, uint32(timeCost), uint32(memoryCost), uint8(parallelism), uint32(keyLen))
	return subtle.ConstantTimeEq(int32(keyLen), int32(len(reKey))) == 1 && subtle.ConstantTimeCompare(key, reKey) == 1, nil
}
