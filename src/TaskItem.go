package src

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/reactivex/rxgo/v2"
	"net/http"
	"net/url"
	"path/filepath"
)

type TaskItem struct {
	EntryURL  *url.URL
	Client    *http.Client
	OutputDir string
	UA        string
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
	observable := rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
		document.Find("*").Each(func(i int, selection *goquery.Selection) {
			//ch <- item.imageWalker(i, selection)
			//ch <- item.styleWalker(i, selection)
			ch <- item.scriptWalker(i, selection)
		})
	}})
	return observable, nil
}

func (item *TaskItem) Export() {
	result, err := item.analyze()
	if err != nil {
		fmt.Println(err)
		return
	}

	for value := range result.Observe() {
		if value.E != nil || value.V == nil {
			continue
		}
		pageItem := value.V.(*PageElementItem)
		if err := pageItem.save(); err == nil { // 下载并保存成功
			pageItem.rewrite()
		}
	}
}

func (item *TaskItem) walker(index int, node *goquery.Selection, name string) (*PageElementItem, error) {
	linkTag := node.Find(name)
	link, isExist := linkTag.Attr("src")
	if !isExist {
		return nil, fmt.Errorf("%v: %v 标签不存在", index, name)
	}
	addr, err := url.Parse(link)
	if err != nil {
		return nil, err
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
		Node:   node,
	}
	return element, err
}

func (item *TaskItem) imageWalker(index int, node *goquery.Selection) rxgo.Item {
	element, err := item.walker(index, node, "img")
	return rxgo.Item{
		V: element,
		E: err,
	}
}

func (item *TaskItem) styleWalker(index int, node *goquery.Selection) rxgo.Item {
	element, err := item.walker(index, node, "style")
	return rxgo.Item{
		V: element,
		E: err,
	}
}

func (item *TaskItem) scriptWalker(index int, node *goquery.Selection) rxgo.Item {
	element, err := item.walker(index, node, "script")
	return rxgo.Item{
		V: element,
		E: err,
	}
}
