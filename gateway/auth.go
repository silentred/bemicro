package gateway

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	AuthSecret = "aJ8yCtDo9Ui"
	AuthKey    = "auth"
)

var AuthSecretByte []byte

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))

	AuthSecretByte = []byte(AuthSecret)
}

// GetAuthInfoPair returns key value pair of auth info
func GetAuthInfoPair() []string {
	sign := makeSign()
	return []string{AuthKey, sign}
}

// IsValidAuth validates the auth info
func IsValidAuth(ctx context.Context) bool {
	if md, ok := metadata.FromContext(ctx); ok {
		values := md[AuthKey]
		if len(values) > 0 {
			return isValid(values[0])
		}
	}
	return false
}

func makeSign() (sign string) {
	num := rand.Uint32()

	numBytes := make([]byte, 8, 8+len(AuthSecretByte))
	binary.PutUvarint(numBytes, uint64(num))
	encodedNum := base64.StdEncoding.EncodeToString(numBytes[:8])

	lastByte := append(numBytes, AuthSecretByte...)
	signPart := md5.Sum(lastByte)

	sign = fmt.Sprintf("%s:%x", encodedNum, signPart)
	return
}

func isValid(sign string) bool {
	parts := strings.Split(sign, ":")
	if len(parts) != 2 {
		return false
	}

	numBytes, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		fmt.Println(err)
	}

	lastByte := append(numBytes, AuthSecretByte...)
	correctSignPart := md5.Sum(lastByte)

	signPart := parts[1]
	if fmt.Sprintf("%x", correctSignPart) != signPart {
		return false
	}

	return true
}
