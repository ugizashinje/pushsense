package sender

import (
	"context"
	"fmt"
	"strings"
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
		fmt.Println("No collection with name ", collection)
	}
	return client.Collections().Create(context.Background(), &schema)
}
func Send(collection string, dbdata []map[string]any) error {
	params := &api.ImportDocumentsParams{
		Action:    pointer.String("create"),
		BatchSize: pointer.Int(1000),
	}
	data := []any{}
	for _, i := range dbdata {
		data = append(data, i)
	}
	res, err := client.Collection(collection).Documents().Import(context.Background(), data, params)
	if err == nil {
		fmt.Printf("Uploaded %d documents \n", len(res))
	}
	return err
}
func Delete(collection string, dbdata []string) error {
	filter := "id:[" + strings.Join(dbdata, ",") + "]"
	params := &api.DeleteDocumentsParams{
		FilterBy: &filter,
	}

	_, err := client.Collection(collection).Documents().Delete(context.Background(), params)
	return err
}
