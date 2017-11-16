package orderfile

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

func ReadFile(orderfilename string) ([]interface{}, error) {
	filecontent, err := ioutil.ReadFile(orderfilename)
	if err != nil {
		return nil, err
	}
	var orders []interface{}
	err = yaml.Unmarshal(filecontent, &orders)
	return orders, nil
}
