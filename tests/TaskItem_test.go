package tests

import (
	"PageSnapShot"
	"testing"
)

func TestTaskItem(t *testing.T) {
	snapshot := PageSnapShot.PageSnapShot{}
	item, err := snapshot.NewTaskItem("https://www.gcores.com", "C:\\Users\\jzh16\\Desktop\\output")
	if err != nil {
		t.Fail()
	}
	if item != nil {
		item.Export()
	}
}
