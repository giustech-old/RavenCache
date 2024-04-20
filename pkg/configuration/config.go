package configuration

import "regexp"

type (
	Config struct {
		Cache []CacheConfig `yaml:"cache"`
		HTTP  HttpConfig    `yaml:"http"`
	}

	CacheConfig struct {
		Name           string                 `yaml:"name"`
		Type           string                 `yaml:"type"`
		Host           string                 `yaml:"host"`
		Port           int                    `yaml:"port"`
		DatabaseNumber int                    `yaml:"database_number"`
		Credentials    CacheConfigCredentials `yaml:"credentials"`
	}

	CacheConfigCredentials struct {
		Password string `yaml:"password"`
	}

	HttpConfig struct {
		HeaderCache []HeaderCacheConfig `yaml:"header_cache"`
		Routes      []HttpConfigRoute   `yaml:"routes"`
	}

	HeaderCacheConfig struct {
		Name   string   `yaml:"name"`
		Values []string `yaml:"values"`
	}

	HttpConfigRoute struct {
		Regex       string   `yaml:"regex"`
		Method      []string `yaml:"method"`
		HeaderCache string   `yaml:"header_cache"`
		CacheKey    string   `yaml:"cache_key"`
		MillisTTL   int64    `yaml:"millis_ttl"`
		Regexp      *regexp.Regexp
	}
)
