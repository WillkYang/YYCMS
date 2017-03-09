package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/beego/goyaml2"
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

	if string(buf[0:1]) == "{" || string(buf[0:1]) == "[" {
		fmt.Println("Look like a Json, try json umarshal")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			fmt.Println("It is Json Map")
			return
		}
	}

	data, err := goyaml2.Read(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Print("Goyaml2 ERR>", string(buf), err)
	}

	if data == nil {
		fmt.Print("Goyaml2 output nil? Pls report bug\n" + string(buf))
		return
	}
	cnf, ok := data.([]interface{})
	if !ok {
		fmt.Print("Not a Arr? >> ", string(buf), data)
		cnf = nil
	}
	return
}


// ReadYmlReader Read yaml file to arr.
// if json like, use json package, unless goyaml2 package.
func ReadFileToMap(path string) (cnf map[string]interface{}, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil || len(buf) < 3 {
		return
	}

	if string(buf[0:1]) == "{" || string(buf[0:1]) == "[" {
		fmt.Printf("Look like a Json, try json umarshal")
		err = json.Unmarshal(buf, &cnf)
		if err == nil {
			fmt.Printf("It is Json Map")
			return
		}
	}

	data, err := goyaml2.Read(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Print("Goyaml2 ERR>", string(buf), err)
	}

	if data == nil {
		fmt.Print("Goyaml2 output nil? Pls report bug\n" + string(buf))
		return
	}
	cnf, ok := data.(map[string]interface{})
	if !ok {
		fmt.Print("Not a Map? >> ", string(buf), data)
		cnf = nil
	}
	return
}
