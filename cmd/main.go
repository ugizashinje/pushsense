package main

import "github.com/ugizashinje/pushsense/pkg/listener"

func init() {
}
func main() {

	listener.Start()

	select {}
}
