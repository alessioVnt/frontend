package main

import (
	"context"
	"io"

	pb "github.com/alessioVnt/frontend/pb"
)

//User service RPCs

func (fe *frontendServer) getUserByID(ctx context.Context, id int32) {
	user, err := pb.NewSdccUserServiceClient(fe.userSvcConn).
		FindByID(ctx, &pb.IDMessage{
			Id: id,
		})
	if err != nil {
		print("Error in getting user")
		return
	}
	print("User found!" + "\n")
	userName := user.Username
	userAddress := user.Address
	userMail := user.Mail
	print(userName + "\n")
	print(userAddress + "\n")
	print(userMail + "\n")
}

//Restaurant service RPCs

func (fe *frontendServer) getAllRestaurants(ctx context.Context) {
	restaurants, err := pb.NewRestaurantServiceClient(fe.restaurantSvcConn).
		GetAllRestaurants(ctx, &pb.RestaurantsRequest{
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

func (fe *frontendServer) addRestaurant(ctx context.Context, name string, city string, address string, tags *pb.TAG) {
	response, err := pb.NewRestaurantServiceClient(fe.restaurantSvcConn).
		AddRestaurant(ctx, &pb.AddRestaurantRequest{
			Name:    name,
			City:    city,
			Address: address,
			TAG:     tags,
		})
	if err != nil {
		print("Error in adding restaurant")
		return
	}
	print(response)
}

//Mail service handlers

func (fe *frontendServer) sendMail(ctx context.Context, tag string, id string) {
	isSuccessful, err := pb.NewSdccMailServiceClient(fe.mailSvcConn).
		SendMail(ctx, &pb.sen)
}

//Recommendation service handlers

func (fe *frontendServer) getRecommendations(ctx context.Context, id string) []string {
	response, err := pb.NewRecommendationServiceClient(fe.recommendationSvcConn).
		GetRecommendations(ctx, &pb.GetRecommendationsRequest{
			UserID: id,
		})

	if err != nil {
		print("Error getting recommendations for user")
		return nil
	}

	for rec := range response.RecommendationList {
		print(rec)
	}

	return response.RecommendationList
}

//Order service handlers

func (fe *frontendServer) getOrder(ctx context.Context, id string) pb.Cart {
	order, err := pb.NewOrderServiceClient(fe.orderSvcConn).
		GetCart(ctx, &pb.GetCartRequest{
			UserId: id,
		})

	if err != nil {
		print("error in getting user order")
		return pb.Cart{}
	}

	return *order
}

func (fe *frontendServer) emptyCart(ctx context.Context, id string) {
	_, err := pb.NewOrderServiceClient(fe.orderSvcConn).
		EmptyCart(ctx, &pb.EmptyCartRequest{
			UserId: id,
		})

	if err != nil {
		print("error in emptying user order")
		return
	}
	print("order succesfully cleared \n")
	return

}

func (fe *frontendServer) addToCart(ctx context.Context, userID string, restaurantID string, itemToAdd *pb.CartItem) {
	_, err := pb.NewOrderServiceClient(fe.orderSvcConn).
		AddItem(ctx, &pb.AddItemRequest{
			UserId:       userID,
			RestaurantId: restaurantID,
			Item:         itemToAdd,
		})

	if err != nil {
		print("Error in adding item to order \n")
		return
	}
}

//Checkout service handlers

func (fe *frontendServer) executeCheckout(ctx context.Context, userID string, restaurantID string, menuItems []*pb.MenuItem, cardNumber string, cvc string, expiration string) bool {
	response, err := pb.NewCheckoutServiceClient(fe.checkoutSvcConn).
		ExecuteTransaction(ctx, &pb.TransactionInfo{
			UserID:         userID,
			RestaurantID:   restaurantID,
			Items:          menuItems,
			CardNumber:     cardNumber,
			Cvc:            cvc,
			CardExpiration: expiration,
		})

	if err != nil {
		print("Error in executing transaction \n")
		return false
	}
	return response.IsSuccessful
}
