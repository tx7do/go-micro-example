package app

import (
	"context"
	"flag"
	"log"
	"time"

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
	srv micro.Service

	web            web.Service
	grpcGatewayMux *runtime.ServeMux
	ginRouter      *gin.Engine

	logger logger.Logger

	cfg *confV1.Bootstrap

	microClient client.Client
}

func New() *App {
	return &App{}
}

func (a *App) Service() micro.Service {
	return a.srv
}

func (a *App) WebService() web.Service {
	return a.web
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

func (a *App) MicroClient() client.Client {
	if a.srv != nil {
		return a.srv.Client()
	}
	//return client.DefaultClient
	return a.microClient
}

// LoadConfig 加载配置
func (a *App) LoadConfig(path string) (*confV1.Bootstrap, error) {
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

// InitRegistry 初始化注册器
func (a *App) InitRegistry(cfg *confV1.Registry) registry.Registry {
	if cfg == nil {
		return nil
	}

	switch cfg.GetType() {
	case "consul":
		return consul.NewRegistry(
			func(options *registry.Options) {
				options.Addrs = []string{cfg.Consul.GetAddress()}
				//options.Timeout = 5 * time.Second
			},
		)

	case "etcd":
		return etcd.NewRegistry()

	case "zookeeper":
		//return zookeeper.NewRegistry()
		return nil

	case "nacos":
		//return nacos.NewRegistry()
		return nil

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
}

// InitLogger 初始化日志
func (a *App) InitLogger(cfg *confV1.Logger) logger.Logger {
	if cfg == nil {
		return nil
	}

	switch cfg.GetType() {
	case "logrus":
		level, err := logger.GetLevel(cfg.Logrus.GetLevel())
		if err != nil {
			level = logger.InfoLevel
		}
		return logrus.NewLogger(
			logger.WithLevel(level),
		)

	case "zap":
		level, err := logger.GetLevel(cfg.Zap.GetLevel())
		if err != nil {
			level = logger.InfoLevel
		}
		aLog, err := microZap.NewLogger(
			logger.WithLevel(level),
			microZap.WithNamespace("micro"),
		)
		if err != nil {
			panic(err)
			return nil
		}
		return aLog

	default:
		return nil
	}
}

func (a *App) CreateGrpcService(ctx context.Context, serviceName, version string, cfg *confV1.Server_GRPC, aLog logger.Logger, reg registry.Registry) micro.Service {

	return nil
}

func (a *App) CreateMicroService(ctx context.Context, serviceName, version string, cfg *confV1.Server_Micro, aLog logger.Logger, reg registry.Registry) micro.Service {
	var opts = []micro.Option{
		micro.Context(ctx),
		micro.Name(serviceName),
		micro.Version(version),

		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 15),
	}

	if reg != nil {
		opts = append(opts, micro.Registry(reg))
	}
	if aLog != nil {
		opts = append(opts, micro.Logger(aLog))
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

	a.srv = srv

	return srv
}

func (a *App) CreateRestService(ctx context.Context, serviceName, version string, cfg *confV1.Server_REST, aLog logger.Logger, reg registry.Registry) web.Service {

	var opts = []web.Option{
		web.Context(ctx),
		web.Name(serviceName),
		web.Version(version),
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

	if reg != nil {
		opts = append(opts, web.Registry(reg))
	}
	if aLog != nil {
		opts = append(opts, web.Logger(aLog))
	}

	// 创建一个新的微服务实例
	srv := web.NewService(
		opts...,
	//web.Transport(transport),
	//web.WrapHandler(middlewares.RecoverWrapper),
	)

	// 初始化服务
	if err := srv.Init(); err != nil {
		panic(err)
		return nil
	}

	a.web = srv

	a.createMicroClient(reg)

	return srv
}

func (a *App) createMicroClient(reg registry.Registry) {
	a.microClient = client.NewClient(
		client.Registry(reg),
		//client.ContentType("application/protobuf"),
		//client.WithLogger(aLog),
		//client.Codec("application/json", client.NewCodec),
	)
}

func (a *App) Start(ctx context.Context, confPath string, serviceName, version string) error {
	cfg, err := a.LoadConfig(confPath)
	if err != nil {
		//panic(err)
		return err
	}

	aLog := a.InitLogger(cfg.Logger)

	reg := a.InitRegistry(cfg.Registry)

	a.logger = aLog

	if cfg.Server != nil && cfg.Server.Micro != nil && cfg.Server.Micro.GetEnable() {
		a.CreateMicroService(ctx, serviceName, version, cfg.Server.Micro, aLog, reg)
	}
	if cfg.Server != nil && cfg.Server.Grpc != nil && cfg.Server.Grpc.GetEnable() {
		a.CreateGrpcService(ctx, serviceName, version, cfg.Server.Grpc, aLog, reg)
	}
	if cfg.Server != nil && cfg.Server.Rest != nil && cfg.Server.Rest.GetEnable() {
		a.CreateRestService(ctx, serviceName, version, cfg.Server.Rest, aLog, reg)
	}

	return nil
}

func (a *App) Stop() error {
	if a.srv != nil {
		a.srv = nil
	}

	if a.web != nil {
		if err := a.web.Stop(); err != nil {
			return err
		}
		a.web = nil
	}

	return nil
}

func (a *App) Run() error {
	if a.srv != nil {
		return a.runService()
	}

	if a.web != nil {
		return a.runWeb()
	}

	return nil
}

func (a *App) runService() error {
	// 启动服务
	if err := a.srv.Run(); err != nil {
		log.Fatalf("Failed to run service: %v", err)
		return err
	}
	return nil
}

func (a *App) runWeb() error {
	// 启动服务
	if err := a.web.Run(); err != nil {
		log.Fatalf("Failed to run web service: %v", err)
		return err
	}
	return nil
}
