package src

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/reactivex/rxgo/v2"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type TaskItem struct {
	EntryURL  *url.URL
	Client    *http.Client
	OutputDir string
	UA        string
	document  *goquery.Document
}

func NewPageTaskItem(src string, outputDir string, client *http.Client, ua string) (*TaskItem, error) {
	urlItem, err := url.Parse(src)
	if err != nil {
		return nil, err
	}
	item := &TaskItem{
		EntryURL:  urlItem,
		Client:    client,
		OutputDir: outputDir,
		UA:        ua,
	}
	return item, nil
}

func (item *TaskItem) prepareDocument() rxgo.Observable {
	return rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		req, err := http.NewRequest("GET", item.EntryURL.String(), nil)
		if err != nil {
			ch <- rxgo.Item{
				V: nil,
				E: err,
			}
			return
		}
		req.Header.Set("User-Agent", item.UA)
		resp, err := item.Client.Do(req)
		if err != nil {
			ch <- rxgo.Item{
				V: nil,
				E: err,
			}
			return
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			ch <- rxgo.Item{
				V: nil,
				E: err,
			}
		} else {
			ch <- rxgo.Item{
				V: doc,
				E: nil,
			}
		}
	}})
}

func (item *TaskItem) analyze() (rxgo.Observable, error) {
	var document *goquery.Document
	for item := range item.prepareDocument().Observe() {
		if item.E != nil {
			return nil, item.E
		} else {
			document = item.V.(*goquery.Document)
			break
		}
	}
	if document == nil {
		return nil, fmt.Errorf("%v未能成功解析", item.EntryURL)
	}
	item.document = document
	observable := rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		document.Find("*").Each(func(i int, selection *goquery.Selection) {
			if selection.Size() == 0 {
				return
			}
			for _, node := range selection.Nodes {
				switch node.Data {
				case "script":
					fallthrough
				case "img":
					fallthrough
				case "style":
					ch <- item.walker(node, selection)
					break
				default:
					break
				}
			}
		})
	}})
	return observable, nil
}

func (item *TaskItem) Export() error {
	result, err := item.analyze()
	if err != nil {
		return err
	}
	for value := range result.Observe() {
		if value.E != nil || value.V == nil {
			continue
		}
		pageItem := value.V.(*PageElementItem)
		if err := pageItem.save(); err == nil { // 下载并保存成功
			pageItem.rewrite(item.OutputDir)
		}
	}
	return item.exportHTML()
}

func (item *TaskItem) exportHTML() error {
	htmlStr, err := item.document.Html()
	if err != nil {
		return err
	}
	_, file := filepath.Split(item.EntryURL.String())
	ext := filepath.Ext(file)
	// 检查文件夹是否已经创建
	_, err = os.Stat(item.OutputDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(item.OutputDir, os.ModePerm)
	}
	if err != nil {
		return err
	}
	if ext != "html" && ext != "htm" {
		file = "index.html"
	}
	return ioutil.WriteFile(filepath.Join(item.OutputDir, file), []byte(htmlStr), os.ModePerm)
}

func (item *TaskItem) walker(node *html.Node, selection *goquery.Selection) rxgo.Item {
	link := ""
	for _, attr := range node.Attr {
		if attr.Key == "src" {
			link = attr.Val
		}
	}
	addr, err := url.Parse(link)
	if err != nil || len(link) == 0 {
		return rxgo.Item{
			V: nil,
			E: err,
		}
	}
	output := addr.String()
	if !addr.IsAbs() { // 是否是相对路径
		addr = item.EntryURL.ResolveReference(addr)
	} else {
		output = addr.Host + addr.Path
	}
	output = filepath.FromSlash(output)
	output = filepath.Join(item.OutputDir, output)
	element := &PageElementItem{
		Src:    addr,
		Output: output,
		Client: item.Client,
		Node:   selection,
	}
	return rxgo.Item{
		V: element,
		E: err,
	}
}
