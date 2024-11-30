package met

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Cfg    *Config
	Cli    *mongo.Client
	filter bson.M
	update bson.M
	limit  int64
}

var (
	DeFaultMaxPoolSize     uint64        = 100
	DeFaultMinPoolSize     uint64        = 10
	DeFaultMaxConnIdleTime time.Duration = 10 * time.Second
)

func Init(url string, config *Config) (*DB, error) {
	dboption := options.Client().ApplyURI(url)

	if config == nil {
		panic("db config is empty")
	}

	var cfg *Config

	cfg = config

	// dbname is can't be empty
	if cfg.DbName == "" {
		panic("dbname is empty")
	}
	// if maxpoolsize is empty, set default value
	if cfg.MaxPoolSize == 0 {
		cfg.MaxPoolSize = DeFaultMaxPoolSize

	}

	if cfg.MinPoolSize == 0 {
		cfg.MinPoolSize = DeFaultMinPoolSize

	}

	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = DeFaultMaxConnIdleTime

	}

	// apply config
	dboption.SetMaxPoolSize(cfg.MaxPoolSize)
	dboption.SetMinPoolSize(cfg.MinPoolSize)
	dboption.SetMaxConnIdleTime(cfg.MaxConnIdleTime)

	// connect to mongodb
	cli, err := mongo.Connect(context.Background(), dboption)

	if err != nil {
		panic(err)
	}

	// 检查连接
	err = cli.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return &DB{
		Cfg:    cfg,
		Cli:    cli,
		filter: bson.M{},
		update: bson.M{},
	}, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		MaxPoolSize:     DeFaultMaxPoolSize,
		MinPoolSize:     DeFaultMinPoolSize,
		MaxConnIdleTime: DeFaultMaxConnIdleTime,
	}
}

func (d *DB) Close() {
	d.Cli.Disconnect(context.TODO())
}
