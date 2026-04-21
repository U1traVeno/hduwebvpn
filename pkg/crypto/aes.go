package crypto

import (
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
)

// EncryptPasswordAES 使用 AES-ECB + PKCS7 padding 对密码进行加密。
// 这是 HDU SSO 系统的标准密码加密方式。
//
// 参数：
//   - cryptoKey: Base64 编码的 AES 密钥（从 SSO 登录页获取）
//   - password: 明文密码
//
// 返回：
//   - Base64 编码的加密密文
//   - 错误（如果加密失败）
func EncryptPasswordAES(cryptoKey, password string) (string, error) {
	// 1. 解码 Base64 密钥
	keyBytes, err := base64.StdEncoding.DecodeString(cryptoKey)
	if err != nil {
		return "", fmt.Errorf("decode crypto key: %w", err)
	}

	// 2. 验证 AES 密钥长度（16/24/32 字节）
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", fmt.Errorf("invalid crypto key length: got %d bytes, want 16/24/32", len(keyBytes))
	}

	// 3. 创建 AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("create aes cipher: %w", err)
	}

	// 4. PKCS7 padding
	plaintext := []byte(password)
	blockSize := block.BlockSize()
	padding := blockSize - len(plaintext)%blockSize
	padtext := make([]byte, len(plaintext)+padding)
	copy(padtext, plaintext)
	for i := len(plaintext); i < len(padtext); i++ {
		padtext[i] = byte(padding)
	}

	// 5. AES-ECB 加密（逐块独立加密）
	ciphertext := make([]byte, len(padtext))
	for i := 0; i < len(padtext); i += blockSize {
		block.Encrypt(ciphertext[i:i+blockSize], padtext[i:i+blockSize])
	}

	// 6. Base64 编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func MD5Hash(text string) string {
	// 计算 MD5 哈希
	hash := md5.Sum([]byte(text))
	// 返回十六进制字符串
	return fmt.Sprintf("%x", hash)
}
