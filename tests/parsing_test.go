package tests

import(
	"github.com/jmespath/go-jmespath"


	"testing"
	"encoding/json"
	"io/ioutil"
	"github.com/beard1ess/gauss/parsing"
	"os"
	"fmt"
)

func ExampleParse() {

	var JsonInput interface{}

	read,_ := ioutil.ReadFile("./one.json")


	_ = json.Unmarshal(read, &JsonInput)
	searched,err := jmespath.Search("Resources.ALBListener.Properties.DefaultActions[0]", JsonInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	out,_ := json.Marshal(searched)
	os.Stdout.Write(out)
	fmt.Println()

	t := []string{"g", "h", "i"}
	parsing.PathFormatter(t)
}

func TestMain(*testing.M) {
	ExampleParse()
}
