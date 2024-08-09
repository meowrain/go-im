package ini_analyzer

import (
	"fmt"
	"im/conf"
	"testing"
)

func TestIniAnalyzer(t *testing.T) {
	conf := &conf.Config{}
	err := LoadIni("./conf.ini", conf)
	fmt.Println(conf)
	if err != nil {
		t.Fatal(err)
	}
}
