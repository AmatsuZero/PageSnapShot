package src

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

type PageElementItem struct {
	Src    *url.URL
	Output string
	Client *http.Client
}

func NewPageElementItem(src string, baseURL *url.URL, mainFolder string, client *http.Client) (*PageElementItem, error) {
	addr, err := url.Parse(src)
	if err != nil {
		return nil, err
	}
	if !addr.IsAbs() { // 是否是相对路径
		addr = baseURL.ResolveReference(addr)
	}
	output := addr.String()
	return &PageElementItem{
		Src:    addr,
		Output: path.Join(output, mainFolder),
		Client: client,
	}, nil
}

func (item *PageElementItem) save(node *goquery.Selection) {
	err := item.createFolder()
	if err != nil { // 创建资源目录
		return
	}
	size, err := item.download() // 下载
	if size <= 0 || err != nil {
		return
	}
	item.rewrite(node) // 替换节点里面的路径
}

/// 将资源地址换位本地路径
func (item *PageElementItem) rewrite(node *goquery.Selection) {
	node.SetAttr("src", item.Output)
}

/// 下载
func (item *PageElementItem) download() (int64, error) {
	resp, err := item.Client.Get(item.Src.String())
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
	savePath := path.Join(item.Output)

	_, err := os.Stat(savePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(savePath, 0755)
	}
	return err
}
