package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Encrypt(password string) (string, string) {
	// 生成盐
	salt := GetRandomString2(22)
	// 对密码进行加密
	password = fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
	return password, salt
}

func Analysis(password, salt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
}

func GetRandomString2(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func Encrypt(password string) (string, string) {
	// 生成盐
	salt := GetRandomString2(22)
	// 对密码进行加密
	password = fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
	return password, salt
}

func Analysis(password, salt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
}

func GetRandomString2(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
