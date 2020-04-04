package PageSnapShot

import (
	"PageSnapShot/src"
	"net"
	"net/http"
	"time"
)

type PageSnapShot struct {
	Client *http.Client
}

func defaultClient() *http.Client {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Timeout:   time.Second * 30,
		Transport: netTransport,
	}
}

func (snapshot *PageSnapShot) NewTaskItem(url string, output string) (*src.TaskItem, error) {
	if snapshot.Client == nil {
		snapshot.Client = defaultClient()
	}
	return src.NewPageTaskItem(url, output, snapshot.Client)
}
