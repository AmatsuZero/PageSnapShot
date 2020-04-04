package src

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
	"path"
)

type PageElementItem struct {
	Src    *url.URL
	Output string
}

func NewPageElementItem(src string, baseURL *url.URL, mainFolder string) (*PageElementItem, error) {
	addr, err := url.Parse(src)
	if err != nil {
		return nil, err
	}
	output := addr.String()
	return &PageElementItem{
		Src:    addr,
		Output: path.Join(output, mainFolder),
	}, nil
}

func (item *PageElementItem) save() {
	fmt.Println(item.Src)
}

/// 将资源地址换位本地路径
func (item *PageElementItem) rewrite(node *goquery.Selection) {

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
