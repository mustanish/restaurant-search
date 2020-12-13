package handlers

import (
	"log"
	"net/http"
	"search/server/connectors"
	"search/server/constants"
	"search/server/responses"
	"strconv"
	"strings"

	"github.com/arangodb/go-driver"
	"github.com/go-chi/render"
)

// SearchRestaurants is used to search a restaurant
func SearchRestaurants(res http.ResponseWriter, req *http.Request) {
	var (
		query       strings.Builder
		restaurants []responses.Restaurant
	)
	isVeg, _ := strconv.ParseBool(req.URL.Query().Get("veg"))
	restaurant, bindVars, response := new(responses.Restaurant), make(map[string]interface{}), make(map[string]interface{})
	query.WriteString("FOR doc IN restaurantSearch SEARCH ")
	query.WriteString("ANALYZER(STARTS_WITH(doc.name, @searchString), @analyzer) OR ")
	query.WriteString("PHRASE(doc.cuisines, @searchString, @analyzer) OR ")
	query.WriteString("PHRASE(doc.address, @searchString, @analyzer) OR ")
	query.WriteString("PHRASE(doc.city, @searchString, @analyzer) ")
	if isVeg {
		query.WriteString("FILTER doc.veg == true ")
	}
	query.WriteString("SORT TFIDF(doc) DESC LIMIT 30 RETURN doc")
	bindVars["searchString"], bindVars["analyzer"] = req.URL.Query().Get("query"), "text_en"
	cursor, err := connectors.QueryDocument(query.String(), bindVars)
	if err != nil {
		render.Render(res, req, responses.NewHTTPError(http.StatusServiceUnavailable, constants.Unavailable))
		return
	}
	for {
		_, err := cursor.ReadDocument(nil, restaurant)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Println("FAILED::could not read restaurant document because of", err.Error())
		} else {
			restaurants = append(restaurants, *restaurant)
		}
	}
	response["restaurants"] = restaurants
	response["count"] = len(restaurants)
	render.Render(res, req, responses.NewHTTPSucess(http.StatusOK, response))
}
