package defineStruct

// 总配置结构体
type Config struct {
	Server   ServerConfig   `yaml:"server"`   // 服务器相关配置
	Database DatabaseConfig `yaml:"database"` // 数据库相关配置
	Redis    RedisConfig    `yaml:"redis"`    // Redis相关配置
	Core     CoreConfig     `yaml:"core"`     // 核心配置
	Modules  ModulesConfig  `yaml:"modules"`  // 业务模块配置
}

// 服务器相关配置
type ServerConfig struct {
	Address           string   `yaml:"address"`           // 监听地址
	OpenapiPath       string   `yaml:"openapiPath"`       // OpenAPI 路径
	SwaggerPath       string   `yaml:"swaggerPath"`       // Swagger 路径
	ServerRoot        string   `yaml:"serverRoot"`        // 静态资源根目录
	ClientMaxBodySize int      `yaml:"clientMaxBodySize"` // 客户端最大请求体积
	Paths             []string `yaml:"paths"`             // 额外模板路径
	DefaultFile       string   `yaml:"defaultFile"`       // 默认首页文件
	Delimiters        []string `yaml:"delimiters"`        // 模板分隔符
}

// 数据库相关配置
// 支持sqlite和mysql两种类型
// type 字段决定使用哪种数据库
// sqlite只需填写type和name，extra，mysql需填写除extra的其他字段
// link为完整连接字符串，可选
// charset、timezone等为mysql专用
type DatabaseConfig struct {
	Link      string `yaml:"link"`      // 数据库连接字符串（可选）
	Type      string `yaml:"type"`      // 数据库类型（sqlite/mysql）
	Name      string `yaml:"name"`      // 数据库名或sqlite文件名
	Host      string `yaml:"host"`      // 主机地址（mysql专用）
	Port      string `yaml:"port"`      // 端口（mysql专用）
	User      string `yaml:"user"`      // 用户名（mysql专用）
	Pass      string `yaml:"pass"`      // 密码（mysql专用）
	Charset   string `yaml:"charset"`   // 字符集（mysql专用）
	Timezone  string `yaml:"timezone"`  // 时区（mysql专用）
	Extra     string `yaml:"extra"`     // 额外参数（sqlite专用）
	CreatedAt string `yaml:"createdAt"` // 创建时间字段
	UpdatedAt string `yaml:"updatedAt"` // 更新时间字段
	DeletedAt string `yaml:"deletedAt"` // 删除时间字段
	Debug     bool   `yaml:"debug"`     // 是否开启调试
}

// Redis相关配置
type RedisConfig struct {
	Enable  int           `yaml:"enable"`  // 是否启用Redis
	DBRedis DBRedisConfig `yaml:"dbRedis"` // 数据库缓存相关
	Core    RedisCore     `yaml:"core"`    // Redis核心配置
}

// 数据库缓存相关配置
type DBRedisConfig struct {
	Enable int `yaml:"enable"` // 是否启用数据库缓存
	Expire int `yaml:"expire"` // 缓存过期时间（毫秒）
	DB     int `yaml:"db"`     // Redis数据库编号
}

// Redis核心配置
type RedisCore struct {
	Address string `yaml:"address"` // Redis地址
	DB      int    `yaml:"db"`      // Redis数据库编号
	Pass    string `yaml:"pass"`    // Redis密码（如有需要可取消注释）
}

// 核心配置
type CoreConfig struct {
	AppName     string       `yaml:"appName"`     // 应用名称
	IsDesktop   bool         `yaml:"isDesktop"`   // 是否桌面端
	IsProd      bool         `yaml:"isProd"`      // 是否生产模式
	AutoMigrate bool         `yaml:"autoMigrate"` // 是否自动建表
	Eps         bool         `yaml:"eps"`         // 是否生成前端路由
	SQLLogger   LoggerConfig `yaml:"sqlLogger"`   // SQL日志配置
	GFLogger    LoggerConfig `yaml:"gfLogger"`    // GF日志配置
	RunLogger   RunLogger    `yaml:"runLogger"`   // 运行日志配置
	File        FileConfig   `yaml:"file"`        // 文件上传配置
}

// 日志配置
type LoggerConfig struct {
	Path     string `yaml:"path"`     // 日志路径
	File     string `yaml:"file"`     // 日志文件名
	Level    string `yaml:"level"`    // 日志级别
	Stdout   bool   `yaml:"stdout"`   // 是否输出到控制台
	Flags    int    `yaml:"flags"`    // 日志标志位
	StStatus int    `yaml:"stStatus"` // 日志状态
	StSkip   int    `yaml:"stSkip"`   // 日志跳过
}

// SQL日志配置
type SQLLogger struct {
	Path     string `yaml:"path"`     // 日志路径
	File     string `yaml:"file"`     // 日志文件名
	Level    string `yaml:"level"`    // 日志级别
	Stdout   bool   `yaml:"stdout"`   // 是否输出到控制台
	Flags    int    `yaml:"flags"`    // 日志标志位
	StStatus int    `yaml:"stStatus"` // 日志状态
	StSkip   int    `yaml:"stSkip"`   // 日志跳过
}

// GF日志配置
type GFLogger struct {
	Path   string `yaml:"path"`   // 日志路径
	File   string `yaml:"file"`   // 日志文件名
	Level  string `yaml:"level"`  // 日志级别
	Stdout bool   `yaml:"stdout"` // 是否输出到控制台
	Flags  int    `yaml:"flags"`  // 日志标志位
}

// 运行日志配置
type RunLogger struct {
	LoggerConfig `yaml:"logger"`
	Enable       bool   `yaml:"enable"`     // 是否启用
	RotateSize   string `yaml:"rotateSize"` // 日志切割大小
}

// 文件上传配置
type FileConfig struct {
	Mode   string    `yaml:"mode"`   // 上传模式（local/oss）
	Domain string    `yaml:"domain"` // 域名或目录映射
	Oss    OssConfig `yaml:"oss"`    // OSS配置
}

// OSS配置
type OssConfig struct {
	Endpoint        string `yaml:"endpoint"`        // OSS服务地址
	AccessKeyID     string `yaml:"accessKeyID"`     // 访问Key
	SecretAccessKey string `yaml:"secretAccessKey"` // 访问密钥
	BucketName      string `yaml:"bucketName"`      // 存储桶名称
	UseSSL          bool   `yaml:"useSSL"`          // 是否使用SSL
	Location        string `yaml:"location"`        // 区域
}

// 业务模块配置
type ModulesConfig struct {
	Base BaseModuleConfig `yaml:"base"` // 基础模块
}

// 基础模块配置
type BaseModuleConfig struct {
	JWT        JWTConfig        `yaml:"jwt"`        // JWT配置
	Middleware MiddlewareConfig `yaml:"middleware"` // 中间件配置
	HTTP       HTTPConfig       `yaml:"http"`       // HTTP代理配置
	Img        ImgConfig        `yaml:"img"`        // 图片相关配置
}

// JWT配置
type JWTConfig struct {
	SSO    bool        `yaml:"sso"`    // 是否单点登录
	Secret string      `yaml:"secret"` // JWT密钥
	Token  TokenConfig `yaml:"token"`  // Token相关配置
}

// Token相关配置
type TokenConfig struct {
	Expire        int `yaml:"expire"`        // Token过期时间（秒）
	RefreshExpire int `yaml:"refreshExpire"` // 刷新Token过期时间（秒）
}

// 中间件配置
type MiddlewareConfig struct {
	CORS      bool          `yaml:"cors"`      // 是否启用CORS
	Authority AuthorityConf `yaml:"authority"` // 权限中间件配置
	Log       LogConf       `yaml:"log"`       // 日志中间件配置
}

// 权限中间件配置
type AuthorityConf struct {
	Enable bool `yaml:"enable"` // 是否启用权限中间件
}

// 日志中间件配置
type LogConf struct {
	Enable     bool   `yaml:"enable"`     // 是否启用日志
	IgnorePath string `yaml:"ignorePath"` // 忽略日志的路径
	IgnoreReg  string `yaml:"ignoreReg"`  // 忽略日志的正则
}

// HTTP代理配置
type HTTPConfig struct {
	ProxyOpen bool   `yaml:"proxy_open"` // 是否开启代理
	ProxyURL  string `yaml:"proxy_url"`  // 代理地址
}

// 图片相关配置
type ImgConfig struct {
	CDNUrl string `yaml:"cdn_url"` // CDN地址
}
