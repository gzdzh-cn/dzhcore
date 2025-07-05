// Package local 提供本地文件上传支持
package local

import (
	"path/filepath"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"

	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/corefile"
	"github.com/gzdzh-cn/dzhcore/utility/util"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

var (
	ctx = gctx.GetInitCtx()
)

type Local struct {
}

func init() {
	g.Log().Debug(ctx, "------------ local init start")

	var (
		err         error
		driverObj   = New()
		driverNames = g.SliceStr{"local"}
	)
	for _, driverName := range driverNames {
		if err = corefile.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}

	g.Log().Debug(ctx, "------------ local init end")
}

func NewInit() {
	g.Log().Debug(ctx, "------------ local NewInit start")
	uploadPath := util.NewToolUtil().GetUploadPath(coreconfig.Config.Core.IsProd, coreconfig.Config.Core.AppName, coreconfig.Config.Core.IsDesktop, coreconfig.Config.Core.File.UploadPath)

	if !gfile.Exists(uploadPath) {
		err := gfile.Mkdir(uploadPath)
		if err != nil {
			panic(err)
		}
	}

	s := g.Server()
	s.AddStaticPath(coreconfig.Config.Core.File.UploadPath, uploadPath)

	g.Log().Debugf(ctx, "uploadPath:%v", uploadPath)
	g.Log().Debug(ctx, "------------ local NewInit end")

}

func New() corefile.Driver {
	return &Local{}
}

func (l *Local) New() corefile.Driver {
	return &Local{}
}

func (l *Local) Upload(ctx g.Ctx) (string, error) {
	var (
		err     error
		Request = g.RequestFromCtx(ctx)
	)

	file := Request.GetUploadFile("file")
	if file == nil {
		return "", gerror.New("上传文件为空")
	}
	// 以当前年月日为目录
	dir := gtime.Now().Format("Ymd")
	defaultPath := coreconfig.Config.Server.ServerRoot
	uploadPath := util.NewToolUtil().GetUploadPath(coreconfig.Config.Core.IsProd, coreconfig.Config.Core.AppName, coreconfig.Config.Core.IsDesktop, defaultPath)
	fileName, err := file.Save(filepath.Join(uploadPath, dir), true)
	if err != nil {
		return "", err
	}

	path := filepath.Join(uploadPath, dir) + "/" + fileName
	if coreconfig.Config.Core.IsDesktop && coreconfig.Config.Core.IsProd {
		return path, err
	}

	return coreconfig.Config.Core.File.FilePrefix + "/" + path, err

}

func (l *Local) GetMode() (data interface{}, err error) {
	data = g.MapStrStr{
		"mode": coreconfig.Config.Core.File.Mode,
		"type": "local",
	}
	return
}
