package main

import (
	"encoding/json"
	pb "github.com/alessioVnt/frontend/pb"
	"io"
	"io/ioutil"
	"net/http"
)

type platformDetails struct {
	css      string
	provider string
}

//User service handlers

//Restaurant service handlers

func (fe *frontendServer) restaurantListHandler(w http.ResponseWriter, r *http.Request) {
	restaurants, err := pb.NewRestaurantServiceClient(fe.restaurantSvcConn).
		GetAllRestaurants(r.Context(), &pb.RestaurantsRequest{
			Message: "prova",
		})
	if err != nil {
		print("Can't get restaurants list")
		return
	}
	for {
		restaurant, err := restaurants.Recv()
		if err == io.EOF {
			return
		}
		print(restaurant.String() + "\n")
	}
}

func (fe *frontendServer) addRestaurantHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		print("error")
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	name := result["name"]
	city := result["city"]
	address := result["address"]
	redTags := result["TAG"].(map[string]interface{})
	tags := &pb.TAG{Tag1: redTags["tag1"].(string), Tag2: redTags["tag2"].(string), Tag3: redTags["tag3"].(string)}

	response, err := pb.NewRestaurantServiceClient(fe.restaurantSvcConn).
		AddRestaurant(r.Context(), &pb.AddRestaurantRequest{
			Name:    name.(string),
			City:    city.(string),
			Address: address.(string),
			TAG:     tags,
		})
	if err != nil {
		print("Error in adding restaurant")
		return
	}
	print(response)

	//for key, value := range redTags {
	// Each value is an interface{} type, that is type asserted as a string
	//	fmt.Println(key, value.(string))
	//}

}

//Mail service handlers

//Recommendation service handlers

//Cart service handlers

//Checkout service handlers
