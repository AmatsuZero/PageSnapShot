package src

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type PageElementItem struct {
	Src    *url.URL
	Output string
	Client *http.Client
	Node   *goquery.Selection
	UA     string
}

func (item *PageElementItem) save() error {
	err := item.createFolder()
	if err != nil { // 创建资源目录
		return err
	}
	_, err = item.download() // 下载
	return err
}

/// 将资源地址换位本地路径
func (item *PageElementItem) rewrite(baseDir string) {
	path, err := filepath.Rel(baseDir, item.Output)
	if err != nil {
		path = item.Output
	}
	item.Node.SetAttr("src", path)
}

/// 下载
func (item *PageElementItem) download() (int64, error) {
	req, err := http.NewRequest("GET", item.Src.String(), nil)
	if err != nil {
		return -1, err
	}
	req.Header.Set("User-Agent", item.UA)
	resp, err := item.Client.Do(req)
	if err != nil {
		return -1, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	out, err := os.Create(item.Output)
	if err != nil {
		return -1, err
	}
	defer func() {
		_ = out.Close()
	}()
	return io.Copy(out, resp.Body)
}

/// 创建目录
func (item *PageElementItem) createFolder() error {
	dirPath, _ := filepath.Split(item.Output)
	dirPath = filepath.Clean(dirPath)
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
	}
	return err
}
