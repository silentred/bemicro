package gateway

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const hmacSecret = "fo3NvF8nWpOnXm2n1PnwkhBGTtyDx9PK"
const UserClaimKey = "user"

var hmacSecretByte []byte

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	hmacSecretByte = []byte(hmacSecret)
	privateKeyPemByte := []byte(prvKeyPem)
	publicKeyPemByte := []byte(pubKeyPem)

	var err error
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyPemByte)
	if err != nil {
		panic(err)
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyPemByte)
	if err != nil {
		panic(err)
	}
}

type User struct {
	UID   uint64 `json:"uid"`
	Name  string `json:"nam"`
	Email string `json:"eml"`
}

type UserClaims struct {
	jwt.StandardClaims
	User
}

func GetUserClaimPair(token string) []string {
	return []string{UserClaimKey, token}
}

// GenerateUserToken returns jwt with UserClaims
func GenerateUserToken(uid uint64, name, email string) (string, error) {
	claims := &UserClaims{
		jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24).Unix()},
		User{uid, name, email},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSecretByte)

	// token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// return token.SignedString(privateKey)
}

// Verify the token string
func VerifyUserToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		switch token.Method.(type) {
		case *jwt.SigningMethodHMAC:
			return hmacSecretByte, nil
		case *jwt.SigningMethodRSA:
			return publicKey, nil
		default:
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	})

	if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return userClaims, nil
	}

	return nil, err
}

const prvKeyPem = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAxW+Z9YYAN6m1kshEgVmsdedxlP7F1iyM2lyqIK7RuUcdWXao
nV8ghfMpLAoGfEOX+jGSk+jtXSBZ26vUmGi9s28kcLIdWuHx6WERM92nRaW7uKxz
YKTADNz33DoAvqHeyNyP/hELY0HukFjwTSwhTw3WoRE0WfRORga4T5e0YhwYnVjn
IAH7FlCG8Ltaiq9/rVKdIUnhPMFSf+wX0NpEwjLdCA6eeQfGeHQ7YCX820y7xpbF
5fwnQKynBzcQPjIO/rdTKj7kPxqqiMXu5EHCspYYFbE+9kuMmck+1pV4KwasaYT2
ivMkF5Tus9xzf3bxNltnjGR9cdsKS8dBxCB3mwIDAQABAoIBAQCBPKtH5x4vUXyU
h1koXp2gVA6qXBb+Og09RpjqaeTIZf+VNzHqSYGNjPzvYeSa5NgPovFytm7hnbKU
M6cm2LEMSn1M85p5ihsDDFHpZHcBBRqbKO8hXNaF1QK9+o3QOz8MtivfQCL3JwpV
HJK3wWJQUBulNRDSrTOrbOyq1P/zk5As9aNA8UNS4hPAJUvtAiGIuhhN+KYlJyaZ
BYIjpAjnqZV5OfS3tf2kgocSNNbPsVSIIZBA1ZZCRIYLitYHNk6cZ40OhM5JqYpM
GSEovawR+748HKmhzN/z42oIbdQcBb7WtWXUS5epzAiB83gxoSoPQSGDRb51vEYs
4z+DgSoZAoGBAOcWuyvs/l+8LPMkp5y6BptxPxk+ad6SnqXY2pvaeC4T2CL1Pu3w
+ntvgnJh5ULr0l3uxoqGqE8Y6pxKmavnvzZLM6yHFEXVteByfnKq/SVsvUclsT+Y
2dQ7T2z5oi8PfbxN2t19y3230+qtDV1JGOjpwWTYAWzxSDIWdq72H0itAoGBANq4
Klqu8CzQrNN/vArmx+tZShT0tMsEz1b21E90tU7kDA5xdZYoNaY0GDBDLxwxd0AM
7HMwppC7KaFHVIZ6PV3kRhZfXVBcIotb2Iq1JgIVfky7+E0MF5gs/0Dv6/78buJ2
sPyznCXU52QUZB/U1pb94IxA0x+15pxoLFqxmGJnAoGAdPs01Q+r1ZrUxmEP2G7z
WU0CvCy0O0/Nr/cO80as/+Zby5aKvLj4k/Pm/TBBdpcabyKorwdrvF7IpUW+dR9j
1IBNMFFRGekNoQlUqYeVjpR1XMbf62ndG2rK0kesqlYVOHXRDb7YfFPKm0nvMgIG
8iEjHYGbdyLNgU2N1xQQ0iECgYAg7RyjLjbF6Fw8MrySP4/VJEn8waH99ilohBwO
IhmxWK9f9UCobEE3VhxWF6cd7WxwXgGyjZ5lp2dq+hwFap2WZukOMSkREe25YQhG
SWMBaU7sKlgE8U8T/6IlmnjCmCnxOcEHKdrV7ykubcts51Ouw2Vsd83QtkeTQDN9
K8Mu/QKBgC0S9UmZDgHI6CE8jlN7nc3s7xBTMoDN+JaE7MIqQ5onf9qYfiS2EuXG
zrBa7jxzUGHtUcaW2UinQ2+BojyB4QawtL5qaVFgZOTtxpWxridUtEKrhhtaOiY1
FHaY1J8WwPjNLYFerrfEhGWRpA+XCsvqYWAxGkcvveiefBo4KNwD
-----END RSA PRIVATE KEY-----`

const pubKeyPem = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxW+Z9YYAN6m1kshEgVms
dedxlP7F1iyM2lyqIK7RuUcdWXaonV8ghfMpLAoGfEOX+jGSk+jtXSBZ26vUmGi9
s28kcLIdWuHx6WERM92nRaW7uKxzYKTADNz33DoAvqHeyNyP/hELY0HukFjwTSwh
Tw3WoRE0WfRORga4T5e0YhwYnVjnIAH7FlCG8Ltaiq9/rVKdIUnhPMFSf+wX0NpE
wjLdCA6eeQfGeHQ7YCX820y7xpbF5fwnQKynBzcQPjIO/rdTKj7kPxqqiMXu5EHC
spYYFbE+9kuMmck+1pV4KwasaYT2ivMkF5Tus9xzf3bxNltnjGR9cdsKS8dBxCB3
mwIDAQAB
-----END PUBLIC KEY-----`
