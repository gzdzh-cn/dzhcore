package dzhcore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadAddonManifestsAndBuiltinTables(t *testing.T) {
	addonDirectory := t.TempDir()
	for directory, content := range map[string]string{
		"task": `{"name":"task","builtin":true,"tables":["addons_task_info","addons_task_log"]}`,
		"crm":  `{"name":"crm","builtin":false,"tables":["addons_crm_customer"]}`,
	} {
		pluginDirectory := filepath.Join(addonDirectory, directory)
		if err := os.MkdirAll(pluginDirectory, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(pluginDirectory, addonManifestFileName), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	manifests, err := LoadAddonManifests(addonDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifests) != 2 || manifests[0].Name != "crm" || manifests[1].Name != "task" {
		t.Fatalf("unexpected manifests: %#v", manifests)
	}

	tables, err := builtinAddonTables(addonDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := tables["addons_task_info"]; !ok {
		t.Fatal("builtin task table was not loaded")
	}
	if _, ok := tables["addons_crm_customer"]; ok {
		t.Fatal("optional plugin table must not be marked as builtin")
	}
}

type manifestTestModel struct{ table string }

func (m *manifestTestModel) TableName() string { return m.table }
func (m *manifestTestModel) GroupName() string { return "default" }

type manifestTestAddon struct{ Addon }

func (*manifestTestAddon) NewInit() {}

func TestValidateAddonModelsOptionalImports(t *testing.T) {
	optional := AddonManifest{Name: "optional", Tables: []string{"addons_optional_one", "addons_optional_two"}}
	for _, test := range []struct {
		name      string
		manifest  AddonManifest
		tables    []string
		compiled  bool
		wantError string
	}{
		{name: "disabled external manifest may remain", manifest: optional},
		{name: "compiled external must register its models", manifest: optional, compiled: true, wantError: "声明了未注册模型表"},
		{name: "partial model import remains invalid", manifest: optional, tables: []string{"addons_optional_one"}, wantError: "addons_optional_two"},
		{name: "builtin cannot omit all models", manifest: AddonManifest{Name: "optional", Builtin: true, Tables: optional.Tables}, wantError: "声明了未注册模型表"},
		{name: "complete compiled external", manifest: optional, tables: optional.Tables, compiled: true},
		{name: "model import without addon is still checked", manifest: optional, tables: optional.Tables},
		{name: "undeclared registered model remains invalid", manifest: optional, tables: []string{"addons_undeclared"}, wantError: "未在 plugin.json 中声明"},
	} {
		t.Run(test.name, func(t *testing.T) {
			modelsMu.Lock()
			previousModels, previousDeferred := Models, deferredAutoMigrateModels
			Models, deferredAutoMigrateModels = nil, nil
			modelsMu.Unlock()
			previousAddons := Addons
			Addons = nil
			t.Cleanup(func() {
				modelsMu.Lock()
				Models, deferredAutoMigrateModels = previousModels, previousDeferred
				modelsMu.Unlock()
				Addons = previousAddons
			})
			for _, table := range test.tables {
				AddModel(&manifestTestModel{table: table})
			}
			if test.compiled {
				AddAddon(&manifestTestAddon{Addon: Addon{Name: "optional"}})
			}
			err := validateAddonModels([]AddonManifest{test.manifest})
			if test.wantError == "" {
				if err != nil {
					t.Fatalf("unexpected manifest error: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("expected error containing %q, got %v", test.wantError, err)
			}
		})
	}
}

func TestValidateAddonManifestsRejectsDuplicateOwnership(t *testing.T) {
	_, err := validateAddonManifests([]AddonManifest{
		{Name: "one", MenuID: "menu-1", Tables: []string{"addons_one"}},
		{Name: "two", MenuID: "menu-1", Tables: []string{"addons_two"}},
	})
	if err == nil {
		t.Fatal("expected duplicate menuId to be rejected")
	}

	_, err = validateAddonManifests([]AddonManifest{
		{Name: "one", Tables: []string{"addons_shared"}},
		{Name: "two", Tables: []string{"addons_shared"}},
	})
	if err == nil {
		t.Fatal("expected duplicate table ownership to be rejected")
	}
}
