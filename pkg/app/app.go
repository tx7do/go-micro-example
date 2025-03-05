package app

import (
	"context"
	"flag"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/logger"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/web"

	"github.com/micro/plugins/v5/logger/logrus"
	microZap "github.com/micro/plugins/v5/logger/zap"

	"github.com/micro/plugins/v5/registry/consul"
	"github.com/micro/plugins/v5/registry/etcd"
	//"github.com/micro/plugins/v5/registry/nacos"
	//"github.com/micro/plugins/v5/registry/zookeeper"

	"github.com/micro/plugins/v5/config/encoder/yaml"
	"go-micro.dev/v5/config"
	"go-micro.dev/v5/config/reader"
	"go-micro.dev/v5/config/reader/json"
	"go-micro.dev/v5/config/source/env"
	"go-micro.dev/v5/config/source/file"

	confV1 "go-micro-example/api/gen/go/common/conf"
)

func init() {
	flag.Parse()
}

type App struct {
	configPath  string
	serviceInfo ServiceInfo
	cfg         *confV1.Bootstrap

	microService micro.Service

	webService     web.Service
	grpcGatewayMux *runtime.ServeMux
	ginRouter      *gin.Engine

	logger logger.Logger

	microClient client.Client

	gormDB        *gorm.DB
	gormMigrators []interface{}

	rdb *redis.Client

	registry registry.Registry
}

func New(opts ...Option) *App {
	app := &App{
		configPath: "./configs/server.yaml",
		serviceInfo: ServiceInfo{
			Version: "1.0.0",
		},
		gormMigrators: make([]interface{}, 0),
	}

	app.init(opts...)

	return app
}

func (a *App) init(opts ...Option) {
	for _, o := range opts {
		o(a)
	}
}

func (a *App) MicroService() micro.Service {
	return a.microService
}

func (a *App) WebService() web.Service {
	return a.webService
}

func (a *App) GrpcGatewayServeMux() *runtime.ServeMux {
	return a.grpcGatewayMux
}

func (a *App) GinRouter() *gin.Engine {
	return a.ginRouter
}

func (a *App) Logger() logger.Logger {
	return a.logger
}

func (a *App) Config() *confV1.Bootstrap {
	return a.cfg
}

func (a *App) Gorm() *gorm.DB {
	return a.gormDB
}

func (a *App) Redis() *redis.Client {
	return a.rdb
}

func (a *App) MicroClient() client.Client {
	if a.microService != nil {
		return a.microService.Client()
	}
	//return client.DefaultClient
	return a.microClient
}

func (a *App) AddGormMigrator(migrator ...interface{}) {
	a.gormMigrators = append(a.gormMigrators, migrator...)
}

// loadConfig 加载配置
func (a *App) loadConfig(path string) (*confV1.Bootstrap, error) {
	enc := yaml.NewEncoder()

	c, err := config.NewConfig(
		config.WithReader(
			json.NewReader( // json reader for internal config merge
				reader.WithEncoder(enc),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	err = c.Load(
		// base config from env
		env.NewSource(),

		// override env with flags
		//flag.NewSource(),

		file.NewSource(
			file.WithPath(path),
		),
	)
	if err != nil {
		return nil, err
	}

	var cfg = &confV1.Bootstrap{}
	err = c.Scan(cfg)
	if err != nil {
		return nil, err
	}

	a.cfg = cfg

	return cfg, nil
}

// initRegistry 初始化注册器
func (a *App) initRegistry(cfg *confV1.Registry) registry.Registry {
	if cfg == nil {
		return nil
	}

	var reg registry.Registry

	switch cfg.GetType() {
	case "consul":
		reg = consul.NewRegistry(
			func(options *registry.Options) {
				options.Addrs = []string{cfg.Consul.GetAddress()}
				//options.Timeout = 5 * time.Second
			},
		)

	case "etcd":
		reg = etcd.NewRegistry()

	case "zookeeper":
		//reg = zookeeper.NewRegistry()

	case "nacos":
		//reg = nacos.NewRegistry()

	case "kubernetes":
		return nil

	case "eureka":
		//return eureka.NewRegistry()
		return nil

	case "polaris":
		//return polaris.NewRegistry()
		return nil

	case "servicecomb":
		return nil

	default:
		return nil
	}

	a.registry = reg

	return reg
}

// initLogger 初始化日志
func (a *App) initLogger(cfg *confV1.Logger) logger.Logger {
	if cfg == nil {
		a.logger = logger.DefaultLogger
		return nil
	}

	var l logger.Logger

	switch cfg.GetType() {
	case "logrus":
		level, err := logger.GetLevel(cfg.Logrus.GetLevel())
		if err != nil {
			level = logger.InfoLevel
		}

		l = logrus.NewLogger(
			logger.WithLevel(level),
		)

	case "zap":
		level, err := logger.GetLevel(cfg.Zap.GetLevel())
		if err != nil {
			level = logger.InfoLevel
		}

		l, err = microZap.NewLogger(
			logger.WithLevel(level),
			microZap.WithNamespace("micro"),
		)
		if err != nil {
			panic(err)
		}
	}

	a.logger = l
	if l == nil {
		a.logger = logger.DefaultLogger
	}

	return l
}

func (a *App) createGrpcService(ctx context.Context, cfg *confV1.Server_GRPC) micro.Service {
	return nil
}

func (a *App) createMicroService(ctx context.Context, cfg *confV1.Server_Micro) micro.Service {
	var opts = []micro.Option{
		micro.Context(ctx),
		micro.Name(a.serviceInfo.Name),
		micro.Version(a.serviceInfo.Version),

		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 15),
	}

	if a.registry != nil {
		opts = append(opts, micro.Registry(a.registry))
	}
	if a.logger != nil {
		opts = append(opts, micro.Logger(a.logger))
	}

	if cfg != nil {
		opts = append(opts, micro.Address(cfg.Addr))
		//rpcServer := server.NewServer(
		//	//server.Name(serviceName),
		//	server.Address(cfg.Addr),
		//	opts...,
		//)
		//opts = append(opts, micro.Server(rpcServer))
	}

	// 创建一个新的微服务实例
	srv := micro.NewService(
		opts...,
	//micro.Transport(grpc.NewTransport()),
	)

	// 初始化服务
	srv.Init()

	a.microService = srv

	return srv
}

func (a *App) createRestService(ctx context.Context, cfg *confV1.Server_REST) web.Service {

	var opts = []web.Option{
		web.Context(ctx),
		web.Name(a.serviceInfo.Name),
		web.Version(a.serviceInfo.Version),
	}

	if cfg.GetEnableGrpcGateway() {
		gatewayMux := runtime.NewServeMux()
		a.grpcGatewayMux = gatewayMux

		if gatewayMux != nil {
			opts = append(opts, web.Handler(gatewayMux))
		}
	} else {
		gin.SetMode(gin.DebugMode)
		router := gin.New()
		router.Use(gin.Recovery())
		a.ginRouter = router

		if router != nil {
			opts = append(opts, web.Handler(router))
		}
	}

	if cfg != nil {
		opts = append(opts, web.Address(cfg.Addr))
	}

	if a.registry != nil {
		opts = append(opts, web.Registry(a.registry))
	}
	if a.logger != nil {
		opts = append(opts, web.Logger(a.logger))
	}

	// 创建一个新的微服务实例
	srv := web.NewService(
		opts...,
	//webService.Transport(transport),
	//webService.WrapHandler(middlewares.RecoverWrapper),
	)

	// 初始化服务
	if err := srv.Init(); err != nil {
		panic(err)
		return nil
	}

	a.webService = srv

	a.createMicroClient()

	return srv
}

func (a *App) createMicroClient() {
	a.microClient = client.NewClient(
		client.Registry(a.registry),
		client.WithLogger(a.logger),
	)
}

func (a *App) Start(ctx context.Context) error {
	cfg, err := a.loadConfig(a.configPath)
	if err != nil {
		return err
	}

	// 初始化日志记录器
	l := a.initLogger(cfg.Logger)

	// 初始化注册器
	a.initRegistry(cfg.Registry)

	// 初始化数据库
	a.createGormClient(cfg.Data, l)

	// 初始化Redis
	a.createRedisClient(cfg.Data, l)

	// 初始化RPC服务
	a.initService(ctx, cfg.Server)

	return nil
}

func (a *App) initService(ctx context.Context, cfg *confV1.Server) {
	if cfg == nil {
		return
	}

	if cfg.Micro != nil && cfg.Micro.GetEnable() {
		a.createMicroService(ctx, cfg.Micro)
	}
	if cfg.Grpc != nil && cfg.Grpc.GetEnable() {
		a.createGrpcService(ctx, cfg.Grpc)
	}
	if cfg.Rest != nil && cfg.Rest.GetEnable() {
		a.createRestService(ctx, cfg.Rest)
	}
}

func (a *App) Stop() {
	if a.microService != nil {
		a.microService = nil
	}

	if a.webService != nil {
		//if err := a.webService.Stop(); err != nil {
		//	panic(err)
		//	return
		//}
		a.webService = nil
	}
}

func (a *App) Run() error {
	if a.microService != nil {
		return a.runMicroService()
	}

	if a.webService != nil {
		return a.runWebService()
	}

	return nil
}

func (a *App) runMicroService() error {
	// 启动服务
	if err := a.microService.Run(); err != nil {
		log.Fatalf("Failed to run micro service: %v", err)
		return err
	}
	return nil
}

func (a *App) runWebService() error {
	// 启动服务
	if err := a.webService.Run(); err != nil {
		log.Fatalf("Failed to run webService service: %v", err)
		return err
	}
	return nil
}

// createGormClient 创建数据库gorm客户端
func (a *App) createGormClient(cfg *confV1.Data, l logger.Logger) *gorm.DB {
	if cfg == nil || cfg.Database == nil {
		return nil
	}

	var driver gorm.Dialector
	switch cfg.Database.Driver {
	default:
		fallthrough
	case "mysql":
		driver = mysql.Open(cfg.Database.Source)
		break
	case "postgres":
		driver = postgres.Open(cfg.Database.Source)
		break
		//case DBDriverClickHouse:
		//	driver = clickhouse.Open(cfg.Database.Source)
		//	break
		//case DBDriverSqlite:
		//	driver = sqlite.Open(cfg.Database.Source)
		//	break
		//case DBDriverSqlServer:
		//	driver = sqlserver.Open(cfg.Database.Source)
		//break
	}

	cli, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		l.Logf(logger.FatalLevel, "[gorm] failed opening connection to db: %v", err)
		return nil
	}

	// 运行数据库迁移工具
	if cfg.Database.Migrate {
		if err = cli.AutoMigrate(
			a.gormMigrators...,
		); err != nil {
			l.Logf(logger.FatalLevel, "[gorm] failed creating schema resources: %v", err)
			return nil
		}
	}

	a.gormDB = cli

	return cli
}

func (a *App) createRedisClient(cfg *confV1.Data, l logger.Logger) (rdb *redis.Client) {
	if cfg == nil || cfg.Redis == nil {
		return nil
	}

	if rdb = redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedis().GetAddr(),
		Password:     cfg.GetRedis().GetPassword(),
		DB:           int(cfg.GetRedis().GetDb()),
		DialTimeout:  cfg.GetRedis().GetDialTimeout().AsDuration(),
		WriteTimeout: cfg.GetRedis().GetWriteTimeout().AsDuration(),
		ReadTimeout:  cfg.GetRedis().GetReadTimeout().AsDuration(),
	}); rdb == nil {
		log.Fatalf("[redis] failed opening connection to redis")
		return nil
	}

	// open tracing instrumentation.
	if cfg.GetRedis().GetEnableTracing() {
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			l.Logf(logger.FatalLevel, "[redis] failed open tracing: %s", err.Error())
			return nil
		}
	}

	// open metrics instrumentation.
	if cfg.GetRedis().GetEnableMetrics() {
		if err := redisotel.InstrumentMetrics(rdb); err != nil {
			l.Logf(logger.FatalLevel, "[redis] failed open metrics: %s", err.Error())
			return nil
		}
	}

	a.rdb = rdb

	return rdb
}
