package tests

import (
	"PageSnapShot"
	"os/user"
	"path"
	"testing"
)

func TestTaskItem(t *testing.T) {
	snapshot := PageSnapShot.PageSnapShot{}
	myself, err := user.Current()
	if err != nil {
		t.Fail()
	}
	homedir := myself.HomeDir
	item, err := snapshot.NewTaskItem("https://www.gcores.com", path.Join(homedir, "Desktop", "output"))
	if err != nil {
		t.Fail()
	}
	if item != nil {
		item.Export()
	}
}
