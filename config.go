package main

import "github.com/joeshaw/envdecode"

type Config struct {
	ApiUrl      string `env:"MSPL_API_URL,default=http://mcflyapi.ngrok.io"`
	RabbitMQUrl string `env:"RABBITMQ_URL,default=amqp://guest:guest@localhost:5672/"`
}

func NewConfigFromEnvironment() (*Config, error) {
	var cfg Config
	err := envdecode.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
