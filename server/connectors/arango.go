package connectors

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"search/config"
	model "search/server/models"
	"strconv"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var (
	db    driver.Database
	col   driver.Collection
	views driver.DatabaseViews
)

// Initialize creates/opens database
func Initialize(config *config.Config) bool {
	var (
		client driver.Client
		ctx    = context.Background()
	)
	if conn, err := http.NewConnection(http.ConnectionConfig{Endpoints: []string{os.Getenv("DATABASE_URL")}}); err != nil {
		log.Println("FAILED::could Not connect to database because of", err.Error())
		return false
	} else if client, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config.Database.DBUser, config.Database.DBPassword),
	}); err != nil {
		log.Println("FAILED::could not create client because of", err.Error())
		return false
	}
	if exist, _ := client.DatabaseExists(ctx, config.Database.DBName); exist {
		db, _ = client.Database(ctx, config.Database.DBName)
	} else {
		db, _ = client.CreateDatabase(ctx, config.Database.DBName, nil)
	}
	loadData()
	createView()
	return true
}

// QueryDocument runs query and returns a cursor to iterate over the returned document
func QueryDocument(query string, bindVars map[string]interface{}) (driver.Cursor, error) {
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := db.Query(ctx, query, bindVars)
	defer cursor.Close()
	return cursor, err
}

func importDocuments(collection string, documents interface{}) error {
	err := openCollection(collection)
	if err != nil {
		return err
	}
	var options driver.ImportDocumentOptions
	options.Complete = true
	options.OnDuplicate = driver.ImportOnDuplicateUpdate
	_, err = col.ImportDocuments(context.Background(), documents, &options)
	if err != nil {
		log.Println("FAILED::could not import documents because of", err.Error())
		return err
	}
	return nil
}

func openCollection(name string) error {
	ctx := context.Background()
	exist, err := db.CollectionExists(ctx, name)
	if err != nil {
		log.Println("FAILED::could not check if collection exist because of ", err.Error())
		return err
	}
	if exist {
		col, err = db.Collection(ctx, name)
		if err != nil {
			log.Println("FAILED::could open collection because of ", err.Error())
			return err
		}
	} else {
		col, err = db.CreateCollection(ctx, name, nil)
		if err != nil {
			log.Println("FAILED::could not create collection because of ", err.Error())
			return err
		}
	}
	return nil
}

func createView() {
	ctx := context.Background()
	exist, _ := db.ViewExists(ctx, "restaurantSearch")
	if !exist {
		var options driver.ArangoSearchViewProperties
		options.Links = map[string]driver.ArangoSearchElementProperties{
			"restaurants": driver.ArangoSearchElementProperties{
				Fields: map[string]driver.ArangoSearchElementProperties{
					"name": driver.ArangoSearchElementProperties{
						Analyzers: []string{"text_en"},
					},
					"cuisines": driver.ArangoSearchElementProperties{
						Analyzers: []string{"text_en"},
					},
					"address": driver.ArangoSearchElementProperties{
						Analyzers: []string{"text_en"},
					},
					"city": driver.ArangoSearchElementProperties{
						Analyzers: []string{"text_en"},
					},
				},
			},
		}
		_, err := db.CreateArangoSearchView(ctx, "restaurantSearch", &options)
		if err != nil {
			log.Println("FAILED::unable to create view because of ", err.Error())
		}
	}
}

func loadData() {
	var (
		restaurants []model.Restaurant
		parsedData  []map[string]interface{}
	)
	restaurant := new(model.Restaurant)
	data, err := os.Open("./server/connectors/data.json")
	if err != nil {
		log.Println("FAILED::could not open data because of ", err.Error())
	}
	defer data.Close()
	dataByte, _ := ioutil.ReadAll(data)
	dataByte = bytes.TrimPrefix(dataByte, []byte("\xef\xbb\xbf"))
	err = json.Unmarshal(dataByte, &parsedData)
	if err != nil {
		log.Println("FAILED::could not unmarshal data because of ", err.Error())
	}
	for _, value := range parsedData {
		rating, _ := strconv.ParseFloat(value["rating"].(string), 32)
		restaurantID, _ := strconv.ParseInt(value["restaurantID"].(string), 10, 32)
		restaurant.Address = value["address"].(string)
		restaurant.City = value["city"].(string)
		restaurant.Cuisines = value["cuisines"].(string)
		restaurant.Image = value["image"].(string)
		restaurant.IsVeg, _ = strconv.ParseBool(value["isVeg"].(string))
		restaurant.Name = value["name"].(string)
		restaurant.Rating = float32(rating)
		restaurant.RestaurantID = int32(restaurantID)
		restaurant.URL = value["url"].(string)
		restaurants = append(restaurants, *restaurant)
	}
	importDocuments("restaurants", restaurants)
}
