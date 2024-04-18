package listener

import (
	"github.com/lib/pq"
)

var processors map[string](func(any) any) = make(map[string](func(any) any))

func init() {
	processors["stringArray"] = stringArray

}
func stringArray(str any) any {
	array := pq.StringArray{}
	array.Scan(str)
	return array
}
