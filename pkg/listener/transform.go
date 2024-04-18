package listener

import (
	"encoding/base64"
	"fmt"
	"strings"
)

var processors map[string](func(any) any) = make(map[string](func(any) any))

func init() {
	processors["stringArray"] = stringArray

}
func stringArray(str any) any {
	raw, ok := str.(string)
	if !ok {
		return nil
	}
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		fmt.Print("ERROR stringArray", err.Error())
	}
	splited := strings.Split(string(decoded[1:len(decoded)-1]), ",")

	return splited
}
