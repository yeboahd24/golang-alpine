package service

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "strings"

    "golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
    time    uint32
    memory  uint32
    threads uint8
    keyLen  uint32
}

func NewPasswordConfig() *PasswordConfig {
    return &PasswordConfig{
        time:    1,
        memory:  64 * 1024,
        threads: 4,
        keyLen:  32,
    }
}

func (c *PasswordConfig) HashPassword(password string) (string, error) {
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    hash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

    // Format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    encodedHash := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
        c.memory, c.time, c.threads, b64Salt, b64Hash)

    return encodedHash, nil
}

func (c *PasswordConfig) VerifyPassword(password, encodedHash string) (bool, error) {
    parts := strings.Split(encodedHash, "$")
    if len(parts) != 6 {
        return false, fmt.Errorf("invalid hash format")
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[4])
    if err != nil {
        return false, err
    }

    hash, err := base64.RawStdEncoding.DecodeString(parts[5])
    if err != nil {
        return false, err
    }

    newHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, uint32(len(hash)))
    return string(hash) == string(newHash), nil
}