package app

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name    string // 服务名称
	Version string // 服务版本
}

type Option func(o *App)

// WithConfigPath 设置配置文件路径
func WithConfigPath(path string) Option {
	return func(a *App) {
		a.configPath = path
	}
}

// WithServiceName 设置服务名称
func WithServiceName(name string) Option {
	return func(a *App) {
		a.serviceInfo.Name = name
	}
}

// WithVersion 设置服务版本
func WithVersion(version string) Option {
	return func(a *App) {
		a.serviceInfo.Version = version
	}
}

// WithGormMigrators 设置gorm迁移器
func WithGormMigrators(migrators []interface{}) Option {
	return func(a *App) {
		a.gormMigrators = migrators
	}
}

type DBDriver string

const (
	DBDriverMysql      DBDriver = "mysql"
	DBDriverPostgres   DBDriver = "postgres"
	DBDriverSqlite     DBDriver = "sqlite"
	DBDriverClickHouse DBDriver = "clickhouse"
	DBDriverSqlServer  DBDriver = "sqlserver"
)
