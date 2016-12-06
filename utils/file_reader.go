package utils

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/beego/goyaml2"
	"bytes"
	"fmt"
)


// ReadYmlReader Read yaml file to arr.
// if json like, use json package, unless goyaml2 package.
func ReadFileToArray(path string) (cnf []interface{}, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil || len(buf) < 3 {
		return
	}

	if string(buf[0:1]) == "{" || string(buf[0:1]) == "["{
		fmt.Printf("Look like a Json, try json umarshal")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			fmt.Printf("It is Json Map")
			return
		}
	}

	data, err := goyaml2.Read(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Printf("Goyaml2 ERR>", string(buf), err)
		return
	}

	if data == nil {
		fmt.Printf("Goyaml2 output nil? Pls report bug\n" + string(buf))
		return
	}
	cnf, ok := data.([]interface{})
	if !ok {
		fmt.Printf("Not a Arr? >> ", string(buf), data)
		cnf = nil
	}
	return
}