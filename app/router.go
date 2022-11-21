package app

import (
	"erp-service/config"
	"erp-service/delivery/activity_log_delivery"
	"erp-service/delivery/bank_delivery"
	"erp-service/delivery/coin_delivery"
	"erp-service/delivery/dashboard_delivery"
	"erp-service/delivery/transaction_delivery"
	"erp-service/delivery/type_transaction_delivery"
	"erp-service/delivery/user_delivery"
	"erp-service/middleware"
	"erp-service/repository/activity_log_repository"
	"erp-service/repository/bank_repository"
	"erp-service/repository/coin_repository"
	"erp-service/repository/dashboard_repository"
	"erp-service/repository/transaction_repository"
	"erp-service/repository/type_repository"
	"erp-service/repository/user_repository"
	"erp-service/usecase/activity_log_usecase"
	"erp-service/usecase/bank_usecase"
	"erp-service/usecase/coin_usecase"
	"erp-service/usecase/dashboard_usecase"
	"erp-service/usecase/jwt_usecase"
	"erp-service/usecase/transaction_usecase"
	"erp-service/usecase/type_transaction_usecase"
	"erp-service/usecase/user_usecase"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func InitRouter(mysqlConn *gorm.DB) *gin.Engine {
	_, err := strconv.Atoi(config.CONFIG["TIMEOUT_MIDDLEWARE_IN_MS"])
	if err != nil {
		log.Println("failed to get timeout config")
	}

	validate := validator.New()

	logRepository := activity_log_repository.NewActivityLogRepository(mysqlConn)
	logUsecase := activity_log_usecase.NewActivityLogUsecase(logRepository)
	logDelivery := activity_log_delivery.NewActivityLogDelivery(logUsecase)

	bankRepository := bank_repository.NewBankRepository(mysqlConn)
	bankUsecase := bank_usecase.NewBankUsecase(bankRepository, validate)
	bankDelivery := bank_delivery.NewBankDelivery(bankUsecase, logUsecase)

	coinRepository := coin_repository.NewCoinRepository(mysqlConn)
	coinUsecase := coin_usecase.NewCoinUsecase(coinRepository, validate)
	coinDelivery := coin_delivery.NewCoinDelivery(coinUsecase, logUsecase)

	userRepository := user_repository.NewUserRepository(mysqlConn)
	jwtUsecase := jwt_usecase.NewJwtUsecase(userRepository)
	userUsecase := user_usecase.NewUserUsecase(userRepository, jwtUsecase, validate)
	userDelivery := user_delivery.NewUserDelivery(userUsecase, logUsecase)

	trxRepository := transaction_repository.NewTransactionRepository(mysqlConn)
	typeRepository := type_repository.NewTypeRepository(mysqlConn)
	typeUsecase := type_transaction_usecase.NewTypeTransactionUsecase(typeRepository)
	typeDelivery := type_transaction_delivery.NewTypeTransactionDelivery(typeUsecase)
	trxUsecase := transaction_usecase.NewTransactionUsecase(trxRepository, coinRepository, bankRepository, typeRepository, validate)
	trxDelivery := transaction_delivery.NewTransactionDelivery(trxUsecase, logUsecase)

	dashboardRepository := dashboard_repository.NewDashboardRepository(mysqlConn)
	dashboardUsecase := dashboard_usecase.NewDashboardUsecase(dashboardRepository, coinRepository)
	dashboardDelivery := dashboard_delivery.NewDashboardDelivery(dashboardUsecase)

	router := gin.Default()
	router.Use(middleware.CorsMiddleware())

	router.POST("/login", userDelivery.UserLogin)
	router.GET("/dashboard", dashboardDelivery.GetDashboard)

	userRoute := router.Group("/")
	userRoute.Use(middleware.UserAuth(jwtUsecase))
	{
		userRoute.GET("/transaction", trxDelivery.GetAllTransaction)
		userRoute.POST("/transaction", trxDelivery.AddTransaction)
		//coin
		userRoute.GET("/coin-balance", coinDelivery.GetDetailCoin)
		userRoute.GET("/type-transaction", typeDelivery.GetAllType)
	}

	adminRoute := router.Group("/")
	adminRoute.Use(middleware.AdminAuth(jwtUsecase))
	{
		//user
		adminRoute.GET("/user", userDelivery.GetAllUser)
		adminRoute.GET("/role", userDelivery.GetAllRole)
		adminRoute.POST("/user", userDelivery.AddUser)

		//bank
		adminRoute.GET("/bank", bankDelivery.GetAllBank)
		adminRoute.POST("/bank", bankDelivery.AddBank)
		adminRoute.PUT("/bank-balance", bankDelivery.UpdateBankBalance)
		adminRoute.PUT("/bank/:id", bankDelivery.UpdateBank)

		//log
		adminRoute.GET("/log", logDelivery.GetActivity)

		//coin
		adminRoute.GET("/coin", coinDelivery.GetCoin)
		adminRoute.PUT("/coin-balance", coinDelivery.UpdateCoinBalance)
	}

	return router
}
