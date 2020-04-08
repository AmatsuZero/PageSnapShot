package tests

import (
	"PageSnapShot"
	"os/user"
	"path/filepath"
	"testing"
)

func TestTaskItem(t *testing.T) {
	snapshot := PageSnapShot.PageSnapShot{}
	myself, err := user.Current()
	if err != nil {
		t.Fail()
	}
	homedir := myself.HomeDir
	output := filepath.FromSlash("Desktop/output")
	item, err := snapshot.NewTaskItem("https://www.taobao.com", filepath.Join(homedir, output))
	if err != nil {
		t.Fail()
	}
	if item != nil {
		err = item.Export()
		if err != nil {
			t.Fail()
		}
	}
}
