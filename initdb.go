package dzhcore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"gorm.io/gorm"

	"github.com/gogf/gf/v2/database/gdb"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/coredb"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gres"
)

var (
	// Models contains the models participating in startup AutoMigrate. It remains
	// exported for backward compatibility; use RegisteredModels when a caller
	// needs every model registered by the application.
	Models []IModel

	deferredAutoMigrateModels []IModel
	modelsMu                  sync.RWMutex
	manifestMu                sync.RWMutex
	manifestMode              = ManifestModeFilesystem
)

type ManifestMode string

const (
	ManifestModeFilesystem ManifestMode = "filesystem"
	ManifestModePacked     ManifestMode = "packed"
)

const (
	// DefaultAddonManifestDirectory is the physical directory containing each
	// plugin's plugin.json manifest for development mode.
	DefaultAddonManifestDirectory = "addons"
	addonManifestFileName         = "plugin.json"
	packedAddonManifestDirectory  = "addons"
)

// AddonManifest describes the startup behaviour and table ownership of a
// plugin. Development reads addons/<name>/plugin.json. Production reads the
// same file after it has been packed into addons/<name>/resource/plugin.json.
type AddonManifest struct {
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	MenuID  string   `json:"menuId"`
	Builtin bool     `json:"builtin"`
	Tables  []string `json:"tables"`
}

// AddModel registers a model for startup AutoMigrate by default.
func AddModel(model IModel) {
	modelsMu.Lock()
	defer modelsMu.Unlock()
	Models = append(Models, model)
}

// ExcludeModelsFromAutoMigrate keeps matching registered models available to
// callers while excluding them from InitModels. It is useful for models whose
// schema must be created explicitly, for example when a plugin is installed.
func ExcludeModelsFromAutoMigrate(match func(IModel) bool) {
	if match == nil {
		return
	}

	modelsMu.Lock()
	defer modelsMu.Unlock()

	autoMigrateModels := make([]IModel, 0, len(Models))
	for _, item := range Models {
		if match(item) {
			deferredAutoMigrateModels = append(deferredAutoMigrateModels, item)
			continue
		}
		autoMigrateModels = append(autoMigrateModels, item)
	}
	Models = autoMigrateModels
}

// SetAddonManifestMode selects the manifest source used by startup and
// installation. Development uses the filesystem; production uses GoFrame's
// packed resource manager.
func SetAddonManifestMode(mode ManifestMode) {
	manifestMu.Lock()
	defer manifestMu.Unlock()
	manifestMode = mode
}

func LoadConfiguredAddonManifests() ([]AddonManifest, error) {
	manifestMu.RLock()
	mode := manifestMode
	manifestMu.RUnlock()

	if mode == ManifestModePacked {
		return LoadPackedAddonManifests()
	}
	return LoadAddonManifests(resolveAddonManifestDirectory(DefaultAddonManifestDirectory))
}

// LoadPackedAddonManifests loads plugin manifests from GoFrame's resource
// manager. Each plugin manifest must be packed at
// addons/<plugin>/resource/plugin.json.
func LoadPackedAddonManifests() ([]AddonManifest, error) {
	files := gres.ScanDirFile(packedAddonManifestDirectory, addonManifestFileName, true)
	if len(files) == 0 {
		return nil, fmt.Errorf("生产环境插件清单未打包，请执行 gf pack")
	}

	manifests := make([]AddonManifest, 0, len(files))
	for _, file := range files {
		fileName := strings.Trim(file.Name(), "/")
		parts := strings.Split(fileName, "/")
		if len(parts) != 4 || parts[0] != "addons" || parts[2] != "resource" || parts[3] != addonManifestFileName {
			return nil, fmt.Errorf("生产环境插件清单路径不正确: %s，必须为 addons/<插件名>/resource/plugin.json", file.Name())
		}

		addonName := parts[1]
		manifest := AddonManifest{}
		if err := json.Unmarshal(file.Content(), &manifest); err != nil {
			return nil, fmt.Errorf("解析打包插件清单 %q 失败: %w", file.Name(), err)
		}
		if manifest.Name == "" {
			return nil, fmt.Errorf("插件清单 %q 缺少 name", file.Name())
		}
		if manifest.Name != addonName {
			return nil, fmt.Errorf("插件清单 %q 的 name 必须与目录名 %q 一致", file.Name(), addonName)
		}
		manifests = append(manifests, manifest)
	}

	return validateAddonManifests(manifests)
}

// GetAddonResourceContent reads an addon resource from the packed resource
// manager. Development mode may fall back to the physical resource file;
// production mode deliberately does not depend on /app/addons.
func GetAddonResourceContent(path string) ([]byte, error) {
	if content := gres.GetContent(path); len(content) > 0 {
		return content, nil
	}
	manifestMu.RLock()
	mode := manifestMode
	manifestMu.RUnlock()
	if mode == ManifestModePacked {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(path)
}

func resolveAddonManifestDirectory(directory string) string {
	if _, err := os.Stat(directory); err == nil {
		return directory
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		return directory
	}
	for {
		candidate := filepath.Join(workingDirectory, directory)
		if _, err = os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(workingDirectory)
		if parent == workingDirectory {
			return directory
		}
		workingDirectory = parent
	}
}

func FindAddonManifest(addonName string) (AddonManifest, bool, error) {
	manifests, err := LoadConfiguredAddonManifests()
	if err != nil {
		return AddonManifest{}, false, err
	}
	for _, manifest := range manifests {
		if manifest.Name == addonName {
			return manifest, true, nil
		}
	}
	return AddonManifest{}, false, nil
}

// PrepareAddonModelsForStartup excludes every addons_* model from startup
// AutoMigrate. Builtin tables are initialized by the builtin lifecycle, while
// external tables are created only by the explicit install operation. All
// model definitions remain available through RegisteredModels.
func PrepareAddonModelsForStartup() error {
	manifests, err := LoadConfiguredAddonManifests()
	if err != nil {
		return err
	}
	if err = validateAddonModels(manifests); err != nil {
		return err
	}

	ExcludeModelsFromAutoMigrate(func(model IModel) bool {
		tableName := strings.ToLower(model.TableName())
		return strings.HasPrefix(tableName, "addons_")
	})
	return nil
}

// ValidateConfiguredAddonModels checks that every registered addon model is
// declared by exactly one manifest. It is exported for build/tests that need
// to validate a generated production manifest without opening a database.
func ValidateConfiguredAddonModels(manifests []AddonManifest) error {
	return validateAddonModels(manifests)
}

func builtinAddonTables(addonDirectory string) (map[string]struct{}, error) {
	manifests, err := LoadAddonManifests(addonDirectory)
	if err != nil {
		return nil, err
	}

	tables := make(map[string]struct{})
	for _, manifest := range manifests {
		if !manifest.Builtin {
			continue
		}
		for _, tableName := range manifest.Tables {
			tableName = strings.ToLower(strings.TrimSpace(tableName))
			if !strings.HasPrefix(tableName, "addons_") {
				return nil, fmt.Errorf("内置插件 %q 的表名必须以 addons_ 开头: %q", manifest.Name, tableName)
			}
			tables[tableName] = struct{}{}
		}
	}
	return tables, nil
}

// LoadAddonManifests loads physical plugin.json files from direct child plugin
// directories. Missing manifests are allowed and mean the plugin is optional.
func LoadAddonManifests(addonDirectory string) ([]AddonManifest, error) {
	entries, err := os.ReadDir(addonDirectory)
	if err != nil {
		return nil, fmt.Errorf("读取插件目录 %q 失败: %w", addonDirectory, err)
	}

	manifests := make([]AddonManifest, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(addonDirectory, entry.Name(), addonManifestFileName)
		content, readErr := os.ReadFile(manifestPath)
		if os.IsNotExist(readErr) {
			continue
		}
		if readErr != nil {
			return nil, fmt.Errorf("读取插件清单 %q 失败: %w", manifestPath, readErr)
		}

		manifest := AddonManifest{}
		if err = json.Unmarshal(content, &manifest); err != nil {
			return nil, fmt.Errorf("解析插件清单 %q 失败: %w", manifestPath, err)
		}
		if manifest.Name == "" {
			return nil, fmt.Errorf("插件清单 %q 缺少 name", manifestPath)
		}
		if manifest.Name != entry.Name() {
			return nil, fmt.Errorf("插件清单 %q 的 name 必须与目录名 %q 一致", manifestPath, entry.Name())
		}
		manifests = append(manifests, manifest)
	}

	return validateAddonManifests(manifests)
}

func validateAddonManifests(manifests []AddonManifest) ([]AddonManifest, error) {
	seenNames := make(map[string]struct{}, len(manifests))
	seenMenuIDs := make(map[string]string, len(manifests))
	seenTables := make(map[string]string)
	for _, manifest := range manifests {
		if manifest.Name == "" {
			return nil, fmt.Errorf("插件清单缺少 name")
		}
		if _, ok := seenNames[manifest.Name]; ok {
			return nil, fmt.Errorf("插件清单 name 重复: %s", manifest.Name)
		}
		seenNames[manifest.Name] = struct{}{}
		if manifest.MenuID != "" {
			if owner, ok := seenMenuIDs[manifest.MenuID]; ok && owner != manifest.Name {
				return nil, fmt.Errorf("插件 menuId 重复: %s (%s、%s)", manifest.MenuID, owner, manifest.Name)
			}
			seenMenuIDs[manifest.MenuID] = manifest.Name
		}
		for _, tableName := range manifest.Tables {
			tableName = strings.ToLower(strings.TrimSpace(tableName))
			if tableName == "" || !strings.HasPrefix(tableName, "addons_") {
				return nil, fmt.Errorf("插件 %q 的表名必须以 addons_ 开头: %q", manifest.Name, tableName)
			}
			if owner, ok := seenTables[tableName]; ok && owner != manifest.Name {
				return nil, fmt.Errorf("插件表重复归属: %s (%s、%s)", tableName, owner, manifest.Name)
			}
			seenTables[tableName] = manifest.Name
		}
	}
	sort.Slice(manifests, func(i, j int) bool { return manifests[i].Name < manifests[j].Name })
	return manifests, nil
}

func validateAddonModels(manifests []AddonManifest) error {
	manifestTables := make(map[string]string)
	// A manifest may remain on disk (or in packed resources) after an optional
	// addon import is disabled. Only compiled addons, builtins and addons with
	// registered models require the complete set of declared model tables.
	requiredAddons := make(map[string]bool)
	for _, addon := range Addons {
		requiredAddons[addon.GetName()] = true
	}
	for _, manifest := range manifests {
		if manifest.Builtin {
			requiredAddons[manifest.Name] = true
		}
		for _, tableName := range manifest.Tables {
			manifestTables[strings.ToLower(strings.TrimSpace(tableName))] = manifest.Name
		}
	}
	registeredTables := make(map[string]struct{})
	for _, item := range RegisteredModels() {
		tableName := strings.ToLower(item.TableName())
		if !strings.HasPrefix(tableName, "addons_") {
			continue
		}
		registeredTables[tableName] = struct{}{}
		owner, ok := manifestTables[tableName]
		if !ok {
			return fmt.Errorf("插件模型 %q 未在 plugin.json 中声明", item.TableName())
		}
		requiredAddons[owner] = true
	}
	for tableName, addonName := range manifestTables {
		if !requiredAddons[addonName] {
			continue
		}
		if _, ok := registeredTables[tableName]; !ok {
			return fmt.Errorf("插件 %q 在 plugin.json 声明了未注册模型表 %q", addonName, tableName)
		}
	}
	return nil
}

// RegisteredModels returns every registered model, including models excluded
// from startup AutoMigrate.
func RegisteredModels() []IModel {
	modelsMu.RLock()
	defer modelsMu.RUnlock()

	models := make([]IModel, 0, len(Models)+len(deferredAutoMigrateModels))
	models = append(models, Models...)
	return append(models, deferredAutoMigrateModels...)
}

// IsModelExcludedFromAutoMigrate reports whether a model was explicitly
// excluded from startup migration. Optional plugin controllers use this to
// avoid registering routes for plugins whose tables have not been installed.
func IsModelExcludedFromAutoMigrate(model IModel) bool {
	if model == nil {
		return false
	}

	modelsMu.RLock()
	defer modelsMu.RUnlock()
	for _, excludedModel := range deferredAutoMigrateModels {
		if excludedModel.GroupName() == model.GroupName() && excludedModel.TableName() == model.TableName() {
			return true
		}
	}
	return false
}

func autoMigrateModels() []IModel {
	modelsMu.RLock()
	defer modelsMu.RUnlock()

	return append([]IModel(nil), Models...)
}

func InitModels() {
	if coreconfig.Config.Core.AutoMigrate {
		models := autoMigrateModels()
		g.Log().Debugf(ctx, "InitModels,数量： %v", len(models))
		for _, model := range models {
			g.Log().Debugf(ctx, "model: %v", model.TableName())
			CreateTable(model)
		}
	}
}

// InitDB 初始化数据库连接供gorm使用
func InitDB(group string) (*gorm.DB, error) {
	// 如果group为空，则使用默认的group，否则使用group参数
	if group == "" {
		group = "default"
	}
	config := g.DB(group).GetConfig()
	db, err := coredb.GetConn(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database (group=%s): %w", group, err)
	}

	GormDBS[group] = db
	return db, nil
}

// 根据entity结构体获取 *gorm.DB
func getDBbyModel(model IModel) *gorm.DB {

	group := model.GroupName()
	// 判断是否存在 GormDBS[group] 字段，如果存在，则使用该字段的值作为DB，否则初始化DB
	if _, ok := GormDBS[group]; ok {
		return GormDBS[group]
	} else {

		db, err := InitDB(group)
		if err != nil {
			panic(err)
		}
		// 把重新初始化的GormDBS存入全局变量中
		GormDBS[group] = db
		return db
	}
}

// CreateTable 根据entity结构体创建表
func CreateTable(model IModel) error {

	db := getDBbyModel(model)
	return db.AutoMigrate(model)

}

// FillInitData 数据库填充初始数据
func FillInitData(ctx g.Ctx, moduleName string, model IModel) error {

	mInit := g.DB("default").Model("base_sys_init")
	value, err := mInit.Clone().Where("group", model.GroupName()).Where("module", moduleName).Value("tables")
	if err != nil {
		g.Log().Errorf(ctx, "读取表 base_sys_init 失败: %s", err.Error())
		return err
	}

	// 模块第一次写入数据
	if value.IsEmpty() {
		id := NodeSnowflake.Generate().String()
		if err = updateData(ctx, mInit, moduleName, model); err == nil {
			g.Log().Debugf(ctx, "分组 %v,模块 %v 中的表 %v，第一次写入", model.GroupName(), moduleName, model.TableName())
			_, err = mInit.Insert(g.Map{"id": id, "group": model.GroupName(), "module": moduleName, "tables": model.TableName()})
			if err != nil {
				g.Log().Error(ctx, err.Error())
				return err
			}
		}
		return err
	}

	tableArr := strings.Split(value.String(), ",")
	tableGarry := garray.NewStrArrayFrom(tableArr)
	//写入过了，跳过
	if tableGarry.Contains(model.TableName()) {
		g.Log().Debugf(ctx, "分组 %v, 模块 %v 中的表 %v, 已经初始化过,跳过本次初始化.", model.GroupName(), moduleName, model.TableName())
		return nil
	}

	//更新写入
	if err = updateData(ctx, mInit, moduleName, model); err == nil {
		g.Log().Debugf(ctx, "分组 %v, 模块 %v 中的表 %v, 写入 ", model.GroupName(), moduleName, model.TableName())
		tableGarry.Append(model.TableName())
		str := strings.Join(tableGarry.Slice(), ",")
		_, err := mInit.Where("group", model.GroupName()).Where("module", moduleName).Data(g.Map{"tables": str}).Update()
		if err != nil {
			return err
		}
	}

	g.Log().Debugf(ctx, "分组 %v, 模块 %v 中的表 %v, 初始化完成 ", model.GroupName(), moduleName, model.TableName())
	return err
}

// InitTables returns the table names recorded as initialized for one module.
// The record is deliberately kept separate from the physical table check so
// callers can distinguish a legacy table from a table initialized by the
// current plugin lifecycle.
func InitTables(ctx g.Ctx, moduleName, groupName string) ([]string, error) {
	if groupName == "" {
		groupName = "default"
	}
	value, err := g.DB(groupName).Model("base_sys_init").Ctx(ctx).
		Where("group", groupName).
		Where("module", moduleName).
		Value("tables")
	if err != nil {
		return nil, err
	}
	if value.IsEmpty() {
		return nil, nil
	}
	result := make([]string, 0)
	for _, tableName := range strings.Split(value.String(), ",") {
		tableName = strings.TrimSpace(tableName)
		if tableName != "" {
			result = append(result, tableName)
		}
	}
	return result, nil
}

// RecordInitTable records a successfully prepared plugin table without
// importing seed data. This is used for legacy installations where the table
// already exists but base_sys_init has no record; importing seeds there could
// duplicate business data.
func RecordInitTable(ctx g.Ctx, moduleName, groupName, tableName string) error {
	if groupName == "" {
		groupName = "default"
	}
	tableName = strings.TrimSpace(tableName)
	if tableName == "" {
		return gerror.New("初始化表名不能为空")
	}
	mInit := g.DB(groupName).Model("base_sys_init").Ctx(ctx)
	value, err := mInit.Clone().Where("group", groupName).Where("module", moduleName).Value("tables")
	if err != nil {
		return err
	}
	if value.IsEmpty() {
		_, err = mInit.Insert(g.Map{
			"id":     NodeSnowflake.Generate().String(),
			"group":  groupName,
			"module": moduleName,
			"tables": tableName,
		})
		return err
	}
	seen := make(map[string]struct{})
	ordered := make([]string, 0)
	for _, item := range strings.Split(value.String(), ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		ordered = append(ordered, item)
	}
	if _, ok := seen[tableName]; ok {
		return nil
	}
	ordered = append(ordered, tableName)
	_, err = mInit.Clone().Where("group", groupName).Where("module", moduleName).
		Data("tables", strings.Join(ordered, ",")).Update()
	return err
}

// 写入文件
func updateData(ctx g.Ctx, mInit *gdb.Model, moduleName string, model IModel) error {

	m := g.DB(model.GroupName()).Model(model.TableName())
	pathName := "addons/" + moduleName
	if moduleName == "base" {
		pathName = "internal"
	}
	path := pathName + "/resource/initjson/" + model.TableName() + ".json"
	jsonData, _ := gjson.LoadContent(gres.GetContent(path))

	g.Log().Debugf(ctx, "model.TableName(): %v,path:%v", model.TableName(), path)

	if jsonData.Var().Clone().IsEmpty() {
		g.Log().Debugf(ctx, "分组%s中的表%s无可用的初始化数据,跳过本次初始化. jsonData:%v", model.GroupName(), model.TableName(), jsonData)
		return gerror.New("无可用的初始化数据,跳过本次初始化")
	}
	_, err := m.Data(jsonData).Insert()
	if err != nil {
		g.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
