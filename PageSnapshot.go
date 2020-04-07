package PageSnapShot

import (
	"PageSnapShot/src"
	"net"
	"net/http"
	"time"
)

type PageSnapShot struct {
	Client *http.Client
	UA     string
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
	if len(snapshot.UA) == 0 {
		snapshot.UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36 OPR/67.0.3575.115"
	}
	return src.NewPageTaskItem(url, output, snapshot.Client, snapshot.UA)
}
