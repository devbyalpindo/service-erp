package main

import (
	"context"
	"erp-service/app"
	"erp-service/config"
	"erp-service/helper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	config.InitConfig()
	loc, err := time.LoadLocation("Asia/Jakarta")
	helper.PanicIfError(err)
	time.Local = loc

	mysqlConn, errMysql := config.ConnectMySQL()
	if errMysql != nil {
		log.Print("error mysql connection:", errMysql)
	}

	if errMysql == nil {
		router := app.InitRouter(mysqlConn)
		log.Println("routes Initialized")

		port := config.CONFIG["PORT"]
		srv := &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}

		log.Println("Server Initialized")
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Println("Shutdown Server...")

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exiting")
	}
}
