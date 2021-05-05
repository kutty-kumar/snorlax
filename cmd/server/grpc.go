package main

import (
	"log"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/infobloxopen/atlas-app-toolkit/requestid"
	"github.com/kutty-kumar/charminder/pkg"
	charminder "github.com/kutty-kumar/charminder/pkg"
	"github.com/kutty-kumar/ho_oh/user_service_v1"
	"github.com/kutty-kumar/user_service/pkg/domain/entity"
	"github.com/kutty-kumar/user_service/pkg/repo"
	"github.com/kutty-kumar/user_service/pkg/svc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

func NewGRPCServer(logger *logrus.Logger, dbConnectionString string) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    time.Duration(viper.GetInt("config.keepalive.time")) * time.Second,
				Timeout: time.Duration(viper.GetInt("config.keepalive.timeout")) * time.Second,
			},
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// logging middleware
				grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)),

				// Request-Id interceptor
				requestid.UnaryServerInterceptor(),

				// Metrics middleware
				grpc_prometheus.UnaryServerInterceptor,

				// validation middleware
				grpc_validator.UnaryServerInterceptor(),

				// collection operators middleware
				gateway.UnaryServerInterceptor(),
			),
		),
	)

	dbLogger := gLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gLogger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  gLogger.Info, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,        // Disable color
		},
	)

	// register database
	db, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return nil, err
	}
	createTables(db)

	// register service implementation with the grpcServer
	domainFactory := charminder.NewDomainFactory()
	domainFactory.RegisterMapping("user", func() charminder.Base {
		return &entity.User{}
	})
	dbOption := charminder.WithDb(db)
	externalIdSetter := func(externalId string, base pkg.Base) pkg.Base {
		base.SetExternalId(externalId)
		return base
	}
	setterOption := pkg.WithExternalIdSetter(externalIdSetter)
	userGormDao := charminder.NewBaseGORMDao(dbOption, charminder.WithCreator(domainFactory.GetMapping("user")), setterOption)
	userGormRepo := repo.NewUserGORMRepo(userGormDao)
	userSvc := svc.NewUserSvc(&userGormRepo)
	user_service_v1.RegisterUserServiceServer(grpcServer, &userSvc)

	return grpcServer, nil
}

func createTables(db *gorm.DB) {
	err := db.AutoMigrate(entity.User{})
	if err != nil {
		log.Fatalf("An error %v occurred while auto-migrating", err)
	}
}
