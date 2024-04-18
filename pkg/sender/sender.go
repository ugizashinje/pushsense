package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/ugizashinje/pushsense/conf"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

var client *typesense.Client

func init() {
	client = typesense.NewClient(
		typesense.WithServer(conf.Typesense.Url),
		typesense.WithAPIKey(conf.Typesense.ApiKey),
		typesense.WithConnectionTimeout(5*time.Second),
		typesense.WithCircuitBreakerMaxRequests(50),
		typesense.WithCircuitBreakerInterval(2*time.Minute),
		typesense.WithCircuitBreakerTimeout(1*time.Minute),
	)

}
func CreateCollection(collection string, schema api.CreateCollectionJSONRequestBody) (*api.CollectionResponse, error) {
	_, err := client.Collection(collection).Delete(context.Background())
	if err != nil {
		fmt.Println("Failed to delete collection ", collection)
	}
	return client.Collections().Create(context.Background(), &schema)
}
func Send(collection string, dbdata []map[string]any) error {
	params := &api.ImportDocumentsParams{
		Action:    pointer.String("create"),
		BatchSize: pointer.Int(100),
	}
	data := []any{}
	for _, i := range dbdata {
		data = append(data, i)
	}
	_, err := client.Collection(collection).Documents().Import(context.Background(), data, params)

	return err
}
