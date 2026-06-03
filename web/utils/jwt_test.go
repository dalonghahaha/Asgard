 package utils

 import (
 	"testing"
 	"time"
 )

 // T-211 JWT HS256 单元测试：签发/解析/过期
 func TestIssueAndParseToken(t *testing.T) {
 	secret := "test-secret"
 	token, exp, err := IssueToken(42, secret, 60)
 	if err != nil {
 	t.Fatalf("IssueToken err: %v", err)
 	}
 	if exp <= time.Now().Unix() {
 	t.Fatalf("exp not in future: %d", exp)
 	}
 	uid, err := ParseToken(token, secret)
 	if err != nil {
 	t.Fatalf("ParseToken err: %v", err)
 	}
 	if uid != 42 {
 	t.Fatalf("uid = %d, want 42", uid)
 	}
 }

 func TestParseTokenWrongSecret(t *testing.T) {
 	token, _, err := IssueToken(1, "secret-A", 60)
 	if err != nil {
 	t.Fatal(err)
 	}
 	if _, err := ParseToken(token, "secret-B"); err == nil {
 	t.Fatal("expected error on wrong secret, got nil")
 	}
 }

 func TestParseTokenExpired(t *testing.T) {
 	token, _, err := IssueToken(1, "s", 0)
 	if err != nil {
 	t.Fatal(err)
 	}
 	// exp = now；再 sleep 1s 强制过期
 	time.Sleep(1100 * time.Millisecond)
 	if _, err := ParseToken(token, "s"); err == nil {
 	t.Fatal("expected error on expired token, got nil")
 	}
 }
