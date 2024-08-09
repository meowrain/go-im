package ini_analyzer

import (
	"fmt"
	"im/model"
	"testing"
)

func TestIniAnalyzer(t *testing.T) {
	conf := &model.Config{}
	err := LoadIni("./conf.ini", conf)
	fmt.Println(conf)
	if err != nil {
		t.Fatal(err)
	}
}
