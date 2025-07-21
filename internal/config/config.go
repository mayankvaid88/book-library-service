package config

import (
	"encoding/json"
	"os"

	"github.com/go-playground/validator/v10"
)

type Config interface {
	GetUser() string
	GetPassword() string
	GetHost() string
	GetPort() string
	GetName() string
}

type DBConfig struct {
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     string `json:"port" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type config struct {
	DB DBConfig `json:"db" validate:"required"`
}

func (c config) GetUser() string {
	return c.DB.User
}
func (c config) GetPassword() string {
	return c.DB.Password
}
func (c config) GetHost() string {
	return c.DB.Host
}
func (c config) GetPort() string {
	return c.DB.Port
}
func (c config) GetName() string {
	return c.DB.Name
}

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	val := validator.New()
	if err := val.Struct(&cfg); err!=nil{
		return nil,err
	}
	return &cfg, nil
}
