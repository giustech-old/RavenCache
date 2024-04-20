package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/giustech/RavenCache/pkg/cache"
	"github.com/giustech/RavenCache/pkg/cache/redis"
	"github.com/giustech/RavenCache/pkg/configuration"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var mapHeaders = make(map[string][]string)

var mapCache = make(map[string]cache.Repository)

var config configuration.Config

func init() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, cache := range config.HTTP.HeaderCache {
		if cache.Values == nil {
			continue
		}

		if _, ok := mapHeaders[cache.Name]; !ok {
			mapHeaders[cache.Name] = make([]string, 0)
		}
		mapHeaders[cache.Name] = append(mapHeaders[cache.Name], cache.Values...)
	}

	for i := range config.HTTP.Routes {
		config.HTTP.Routes[i].Regexp = regexp.MustCompile(config.HTTP.Routes[i].Regex)
	}

	for _, cacheConfig := range config.Cache {
		switch cacheConfig.Type {
		case "redis":
			result := redis.GetInstance(fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port), cacheConfig.Credentials.Password, cacheConfig.DatabaseNumber)
			mapCache[cacheConfig.Name] = result
			break
		}
	}

}

func main() {
	targetURL, _ := url.Parse("http://localhost:3000")

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		for _, regex := range config.HTTP.Routes {
			if regex.Regexp.MatchString(r.RequestURI) {
				if mapHeaders[regex.HeaderCache] == nil {
					fmt.Println("Erro ao obter chaves de cache")
					continue
				}

				if mapCache[regex.CacheKey] == nil {
					fmt.Println("Erro ao obter chaves de cache")
					continue
				}
				fmt.Println("A URI corresponde à expressão regular:", regex.Regexp.String())
				var builder strings.Builder
				for _, header := range mapHeaders[regex.HeaderCache] {
					if value := r.Header.Get(header); value != "" {
						value = strings.ReplaceAll(value, " ", "")
						builder.WriteString(value)
					}
				}
				builder.WriteString(r.Method)
				builder.WriteString(r.RequestURI)
				hasher := md5.New()
				hasher.Write([]byte(builder.String()))
				md5Hash := hex.EncodeToString(hasher.Sum(nil))

				content := mapCache[regex.CacheKey].Get(md5Hash, regex.MillisTTL)

				if content != "" {
					w.Write([]byte(content))
					return
				} else {
					rec := httptest.NewRecorder()
					proxy.ServeHTTP(rec, r)
					resp := rec.Result()
					body, _ := io.ReadAll(resp.Body)
					if resp.StatusCode >= 200 && resp.StatusCode <= 300 {
						mapCache[regex.CacheKey].Put(md5Hash, body, regex.MillisTTL)
					}
					for k, v := range resp.Header {
						w.Header()[k] = v
					}
					w.WriteHeader(resp.StatusCode)
					w.Write(body)
					return
				}
			}
		}
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
