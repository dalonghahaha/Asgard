 package utils

 import (
 	"crypto/hmac"
 	"crypto/sha256"
 	"encoding/base64"
 	"encoding/json"
 	"errors"
 	"strconv"
 	"strings"
 	"time"
 )

 // 极简 JWT（HS256）实现：只依赖标准库，避免引入第三方依赖。
 // token 形如 base64url(header).base64url(payload).base64url(hmacSha256)。

 type jwtHeader struct {
 	Alg string `json:"alg"`
 	Typ string `json:"typ"`
 }

 type jwtClaims struct {
 	Sub string `json:"sub"`
 	Exp int64  `json:"exp"`
 	Iat int64  `json:"iat"`
 }

 func base64URLEncode(b []byte) string {
 	return base64.RawURLEncoding.EncodeToString(b)
 }

 func base64URLDecode(s string) ([]byte, error) {
 	return base64.RawURLEncoding.DecodeString(s)
 }

 func jwtSign(secret string, header, payload []byte) string {
 	hh := base64URLEncode(header)
 	pp := base64URLEncode(payload)
 	mac := hmac.New(sha256.New, []byte(secret))
 	mac.Write([]byte(hh + "." + pp))
 	sig := base64URLEncode(mac.Sum(nil))
 	return hh + "." + pp + "." + sig
 }

 // IssueToken 签发 userID 的 JWT，过期时间由 ttlSeconds 控制（秒）。
 func IssueToken(userID int64, secret string, ttlSeconds int) (string, int64, error) {
 	now := time.Now().Unix()
 	exp := now + int64(ttlSeconds)
 	header, _ := json.Marshal(jwtHeader{Alg: "HS256", Typ: "JWT"})
 	claims, _ := json.Marshal(jwtClaims{Sub: strconv.FormatInt(userID, 10), Exp: exp, Iat: now})
 	return jwtSign(secret, header, claims), exp, nil
 }

 // ParseToken 解析并校验签名/过期；返回 userID。
 func ParseToken(token, secret string) (int64, error) {
 	parts := strings.Split(token, ".")
 	if len(parts) != 3 {
 		return 0, errors.New("token 格式错误")
 	}
 	mac := hmac.New(sha256.New, []byte(secret))
 	mac.Write([]byte(parts[0] + "." + parts[1]))
 	expected := base64URLEncode(mac.Sum(nil))
 	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
 		return 0, errors.New("token 签名错误")
 	}
 	payload, err := base64URLDecode(parts[1])
 	if err != nil {
 		return 0, err
 	}
 	var claims jwtClaims
 	if err := json.Unmarshal(payload, &claims); err != nil {
 		return 0, err
 	}
 	if claims.Exp > 0 && time.Now().Unix() >= claims.Exp {
 		return 0, errors.New("token 已过期")
 	}
 	userID, err := strconv.ParseInt(claims.Sub, 10, 64)
 	if err != nil {
 		return 0, err
 	}
 	return userID, nil
 }
