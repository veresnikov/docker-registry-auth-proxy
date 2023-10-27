package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/veresnikov/docker-registry-auth-proxy/pkg/application/auth"
)

type UserName string
type PasswordHash []byte

func Load(configPath string, passwordHasher auth.PasswordHasher) (map[UserName]PasswordHash, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	configBody, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var configData []userData
	err = json.Unmarshal(configBody, &configData)
	if err != nil {
		return nil, err
	}

	userMap := make(map[UserName]PasswordHash)
	for _, user := range configData {
		passwordHash, err := passwordHasher.Hash(user.Password)
		if err != nil {
			return nil, err
		}
		userMap[UserName(user.UserName)] = passwordHash
	}
	return userMap, nil
}

type userData struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
