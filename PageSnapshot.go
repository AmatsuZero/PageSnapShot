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

func (snapshot *PageSnapShot) defaultClient() *http.Client {
	netTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
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
		snapshot.Client = snapshot.defaultClient()
	}
	return src.NewPageTaskItem(url, output, snapshot.Client)
}
