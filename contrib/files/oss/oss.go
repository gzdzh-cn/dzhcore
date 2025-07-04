package oss

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gzdzh-cn/dzhcore/coreconfig"
	"github.com/gzdzh-cn/dzhcore/corefile"
)

var (
	ctx = gctx.GetInitCtx()
)

type Oss struct {
	Client *oss.Client
	Bucket *oss.Bucket
}

func NewInit() {
	g.Log().Debug(ctx, "------------ oss NewInit start")
	var (
		err          error
		driverNames  = g.SliceStr{"oss"}
		ossDriverObj = New()
	)

	if err != nil {
		panic(err)
	}

	for _, driverName := range driverNames {
		if err = corefile.Register(driverName, ossDriverObj); err != nil {
			panic(err)
		}
	}

	g.Log().Debug(ctx, "------------ oss NewInit end")
}

func New() corefile.Driver {

	if coreconfig.Config.Core.File.Mode != "oss" {
		return nil
	}
	endpoint := coreconfig.Config.Core.File.Oss.Endpoint
	accessKeyID := coreconfig.Config.Core.File.Oss.AccessKeyID
	secretAccessKey := coreconfig.Config.Core.File.Oss.SecretAccessKey
	bucketName := coreconfig.Config.Core.File.Oss.BucketName
	// Initialize oss client object.
	client, err := oss.New(endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		g.Log().Fatal(ctx, err)
		return nil
	}

	exist, err := client.IsBucketExist(bucketName)

	if err != nil {
		g.Log().Fatal(ctx, err)
		return nil
	}

	if exist {
		g.Log().Debug(ctx, fmt.Sprintf("存储桶%s已存在", bucketName))
	} else {
		// 创建存储桶
		err = client.CreateBucket(bucketName)
		if err != nil {
			g.Log().Fatal(ctx, err)
			return nil
		}
		g.Log().Debug(ctx, fmt.Sprintf("存储桶%s创建成功", bucketName))
	}

	bucket, _ := client.Bucket(bucketName)
	return &Oss{Client: client, Bucket: bucket}
}

func (m *Oss) GetMode() (data interface{}, err error) {
	data = g.MapStrStr{
		"mode": "local",
		"type": "oss",
	}
	return
}

func (m *Oss) Upload(ctx g.Ctx) (string, error) {
	var (
		err     error
		Request = g.RequestFromCtx(ctx)
	)

	file := Request.GetUploadFile("file")
	if file == nil {
		return "", gerror.New("上传文件为空")
	}

	src, err := file.Open()
	if err != nil {
		g.Log().Error(ctx, "文件打开失败")
	}
	defer src.Close()

	// 以当前年月日为目录
	dir := gtime.Now().Format("Ymd")
	fileName := Request.Get("key", grand.S(16, false)).String()
	fullPath := fmt.Sprintf("uploads/%s/%s", dir, fileName)

	// 创建目录
	err = m.Bucket.PutObject(fullPath, src)

	if err != nil {
		return "上传失败", err
	}

	url := fmt.Sprintf("https://%s.%s/%s", m.Bucket.BucketName, coreconfig.Config.Core.File.Oss.Endpoint, fullPath)

	return url, nil
}

// 上传文件
func (m *Oss) UploadFile(ctx g.Ctx, filePath string) (string, error) {

	var (
		err       error
		isWebPath bool
	)

	// 以当前年月日为目录
	dir := gtime.Now().Format("Ymd")

	isWebPath = gstr.HasSuffix(filePath, "https://") || gstr.HasSuffix(filePath, "http://")
	// if !isWebPath {
	// 	isWebPath = gstr.HasSuffix(filePath, "http://")
	// }

	// 如果是网络图片，先下载到系统临时文件夹
	if isWebPath {
		g.Log().Debugf(ctx, "web pic : %v", filePath)
		filePath, _ = downLoadToLocal(ctx, filePath)

	}

	fileName := grand.S(16, false) + ".png"
	fullPath := fmt.Sprintf("uploads/%s/%s", dir, fileName)

	// 创建目录
	err = m.Bucket.PutObjectFromFile(fullPath, filePath)

	if err != nil {
		g.Log().Errorf(ctx, "上传失败 err : %v", err)
		return "上传失败", err
	}
	if isWebPath {
		// 删除临时文件
		gfile.Remove(filePath)
	}

	url := fmt.Sprintf("https://%s.%s/%s", m.Bucket.BucketName, coreconfig.Config.Core.File.Oss.Endpoint, fullPath)

	return url, nil
}

func (m *Oss) New() corefile.Driver {
	return m
}

// 下载网络图片到系统临时本地文件夹
func downLoadToLocal(ctx g.Ctx, filePath string) (string, error) {

	// Make an HTTP GET request
	response, err := http.Get(filePath)
	if err != nil {
		g.Log().Error(ctx, "Make an HTTP GET request err:", err)
		return "", gerror.New("Make an HTTP GET request err:")
	}
	defer response.Body.Close()

	// 检查响应状态码
	if response.StatusCode != http.StatusOK {
		g.Log().Error(ctx, "HTTP response status code error:", response.Status)
		return "", gerror.New("HTTP response status code error")
	}

	// 以当前年月日为目录
	// dir := gtime.Now().Format("Ymd")
	fileName := grand.S(16, false) + ".png"
	// fullPath_ := fmt.Sprintf(gfile.MainPkgPath()+"/public/uploads/%s", dir)
	TempPath := gfile.Temp("gfile_example_basic_dir")
	isExist := gfile.Exists(TempPath)
	if !isExist {
		// 创建文件
		gfile.Mkdir(TempPath)
	}

	fullPath := fmt.Sprintf("%s/%s", TempPath, fileName)

	// 创建本地文件用于保存图片
	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("创建本地文件失败:", err)
		return "", nil
	}
	defer file.Close()

	// 将HTTP响应的内容复制到本地文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("保存文件失败:", err)
		return "", nil
	}

	return fullPath, nil

}
