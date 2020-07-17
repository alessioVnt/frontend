package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"time"
)

const (
	port = "8080"
)

type frontendServer struct {
	restaurantSvcAddr string
	restaurantSvcConn *grpc.ClientConn

	userSvcAddr string
	userSvcConn *grpc.ClientConn

	mailSvcAddr string
	mailSvcConn *grpc.ClientConn

	recommendationSvcAddr string
	recommendationSvcConn *grpc.ClientConn

	orderSvcAddr string
	orderSvcConn *grpc.ClientConn
}

func main() {
	ctx := context.Background()
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	srvPort := port
	if os.Getenv("PORT") != "" {
		srvPort = os.Getenv("PORT")
	}
	addr := os.Getenv("LISTEN_ADDR")
	svc := new(frontendServer)
	mustMapEnv(&svc.restaurantSvcAddr, "RESTAURANT_SERVICE_ADDR")
	mustMapEnv(&svc.userSvcAddr, "USER_SERVICE_ADDR")
	mustMapEnv(&svc.mailSvcAddr, "MAIL_SERVICE_ADDR")
	mustMapEnv(&svc.recommendationSvcAddr, "RECOMMENDATION_SERVICE_ADDR")
	mustMapEnv(&svc.orderSvcAddr, "CART_SERVICE_ADDR")

	mustConnGRPC(ctx, &svc.restaurantSvcConn, svc.restaurantSvcAddr)
	mustConnGRPC(ctx, &svc.userSvcConn, svc.userSvcAddr)
	mustConnGRPC(ctx, &svc.mailSvcConn, svc.mailSvcAddr)
	mustConnGRPC(ctx, &svc.recommendationSvcConn, svc.recommendationSvcAddr)
	mustConnGRPC(ctx, &svc.orderSvcConn, svc.orderSvcAddr)

	r := mux.NewRouter()
	r.HandleFunc("/restaurants", svc.restaurantListHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/addRestaurant", svc.addRestaurantHandler).Methods(http.MethodPost)
	r.HandleFunc("/getUser", svc.getUserByIDHandler).Methods(http.MethodGet)

	var handler http.Handler = r

	log.Infof("starting server on " + addr + ":" + srvPort)
	log.Fatal(http.ListenAndServe(addr+":"+srvPort, handler))
}

func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		// panic(fmt.Sprintf("environment variable %q not set", envKey))
		//TEMPORANEO
		if envKey == "RESTAURANT_SERVICE_ADDR" {
			v = "localhost:50051"
		}
		if envKey == "USER_SERVICE_ADDR" {
			//DA CONTROLLARE
			v = "localhost:50051"
		}
	}
	*target = v
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) {
	var err error
	*conn, err = grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second*3))
	if err != nil {
		panic(errors.Wrapf(err, "grpc: failed to connect %s", addr))
	}
}
