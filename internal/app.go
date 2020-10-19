package internal

import (
	"fmt"

	"github.com/olivere/elastic/v7"

	fasthttp "gitlab.bcowtech.de/bcow-go/host-fasthttp"
)

type (
	App struct {
		Host            *Host
		Config          *Config
		ServiceProvider *ServiceProvider
	}

	Host fasthttp.Host

	Config struct {
		// fasthttp server
		ListenAddress  string `arg:"address;the combination of IP address and listen port"`
		EnableCompress bool   `arg:"compress;indicates the response enable compress or not"`
		ServerName     string `arg:"hostname;the server name"`

		/* you can put the configuration settings here...
		 *
		 * example:
		 *
		 *   RedisHost      string `env:"*REDIS_HOST"        yaml:"redisHost"`
		 *   RedisPassword  string `env:"REDIS_PASSWORD"     yaml:"redisPassword"`
		 *   RedisDB        int    `env:"REDIS_DB"           yaml:"redisDB"`
		 *   RedisPoolSize  int    `env:"REDIS_POOL_SIZE"    yaml:"redisPoolSize"`
		 *   RedisWorkspace string `env:"REDIS_WORKSPACE"    yaml:"workspace"`
		 */
		CacheHost     string `env:"*CACHE_HOST"        yaml:"cacheHost"`
		CachePassword string `env:"CACHE_PASSWORD"     yaml:"cachePassword"`
		CacheDB       int    `env:"CACHE_DB"           yaml:"cacheDB"`
		CachePoolSize int    `env:"CACHE_POOL_SIZE"    yaml:"cachePoolSize"`
		Workspace     string `env:"WORKSPACE"          yaml:"workspace"`
		ESHost        string `env:"*ES_HOST"           yaml:"ESHost"`
	}

	ServiceProvider struct {
		/* you can put the service here...
		 *
		 * example:
		 *
		 *   RedisClient      *redis.Client
		 *   OrderRepository  *db.OrderRepository
		 */
		CacheServer *MockCacheServer
		Workspace   string
		ESClient    *elastic.Client
	}
)

func (provider *ServiceProvider) Init(conf *Config) {
	provider.CacheServer = &MockCacheServer{
		Host:     conf.CacheHost,
		Password: conf.CachePassword,
		DB:       conf.CacheDB,
		PoolSize: conf.CachePoolSize,
	}
	provider.Workspace = conf.Workspace
	fmt.Println("TEST", conf.ESHost)
	provider.ESClient, _ = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(conf.ESHost))
}

func (h *Host) Init(conf *Config) {
	h.Server = &fasthttp.Server{
		Name:                          conf.ServerName,
		DisableKeepalive:              true,
		DisableHeaderNamesNormalizing: true,
	}
	h.ListenAddress = conf.ListenAddress
	h.EnableCompress = conf.EnableCompress
}
