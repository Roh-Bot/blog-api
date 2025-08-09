package config

import (
	"context"
	"errors"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type AtomicConfig struct {
	value atomic.Value
}

func (a *AtomicConfig) Set(cfg *Config) {
	a.value.Store(cfg)
}

func (a *AtomicConfig) Get() *Config {
	return a.value.Load().(*Config)
}

type Config struct {
	Server   Server   `koanf:"Server"`
	Auth     Auth     `koanf:"Auth"`
	Database Database `koanf:"Database"`
	Cache    Cache    `koanf:"Cache"`
	Logger   Logger   `koanf:"Logger"`
}

type Server struct {
	Address        string `koanf:"address"`
	ContextTimeout int32  `koanf:"context_timeout"`
}

type Auth struct {
	Secret          string   `koanf:"secret"` //TODO: store this in the env
	Issuer          string   `koanf:"issuer"`
	Audience        string   `koanf:"audience"`
	TokenTTL        int32    `koanf:"token_ttl"`
	RefreshTokenTTL int32    `koanf:"refresh_token_ttl"`
	ValidUsers      []string `koanf:"valid_users"`
	EncryptionKey   string   `koanf:"encryption_key"`
}

type Database struct {
	Host                  string        `koanf:"host"`
	Port                  string        `koanf:"port"`
	User                  string        `koanf:"user"`
	Password              string        `koanf:"password"`
	Database              string        `koanf:"database"`
	SSLMode               string        `koanf:"ssl_mode"`
	MaxConnectionIdleTime time.Duration `koanf:"max_connection_idle_time"`
	MaxConnectionLifetime time.Duration `koanf:"max_connection_lifetime"`
}
type Cache struct {
	Host            string        `koanf:"host"`
	Protocol        int           `koanf:"protocol"`
	Port            int           `koanf:"port"`
	User            string        `koanf:"user"`
	Password        string        `koanf:"password"`
	Db              int           `koanf:"db"`
	MaxRetries      int           `koanf:"max_retries"`
	MinRetryBackoff time.Duration `koanf:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `koanf:"max_retry_backoff"`
	DialTimeout     time.Duration `koanf:"dial_timeout"`
	ReadTimeout     time.Duration `koanf:"read_timeout"`
	WriteTimeout    time.Duration `koanf:"write_timeout"`
	PoolFifo        bool          `koanf:"pool_fifo"`
	PoolSize        int           `koanf:"pool_size"`
	PoolTimeout     time.Duration `koanf:"pool_timeout"`
	MinIdleConns    int           `koanf:"min_idle_conns"`
	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
}

type Logger struct {
	Level         string `koanf:"level"`
	FilePath      string `koanf:"filepath"`
	IsDevelopment bool   `koanf:"is_development"`
}

func LoadConfiguration(ctx context.Context) (*AtomicConfig, error) {
	// Load YAML config
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.yaml")

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		configPath = os.Getenv("CONFIG_PATH")
	}
	fileProvider := file.Provider(configPath)

	k := koanf.New(".")
	if err := k.Load(fileProvider, yaml.Parser()); err != nil {
		return nil, err
	}

	c, err := unmarshalIntoStruct(k)
	if err != nil {
		return nil, err
	}

	atomicConfig := &AtomicConfig{}
	atomicConfig.Set(c)

	//go reloader(ctx, fileProvider, atomicConfig)
	// Get secrets from env
	//c.Server.auth.Secret = os.Getenv("JWT_SECRET")

	return atomicConfig, nil
}

func reloader(ctx context.Context, f *file.File, ac *AtomicConfig) {
	err := f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Println("watch error: ", err.Error())
			return
		}

		// Throw away the old config and load a fresh copy.
		log.Println("config changed. Reloading ...")
		k := koanf.New(".")
		if err := k.Load(f, yaml.Parser()); err != nil {
			return
		}
		c, err := unmarshalIntoStruct(k)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ac.Set(c)

		log.Println("config reload complete.")
	})
	if err != nil {
		return
	}

	log.Println("waiting for config changes...")
	<-ctx.Done()
	if err := f.Unwatch(); err != nil {
		log.Println("failed to unwatch: ", err.Error())
	}
}

func unmarshalIntoStruct(k *koanf.Koanf) (*Config, error) {
	c := &Config{}
	if err := k.Unmarshal("", c); err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
		return nil, err
	}
	return c, nil
}
