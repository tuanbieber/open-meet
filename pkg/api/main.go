package api

import (
	"open-meet/pkg/logger"
	"open-meet/pkg/middleware"
	"open-meet/pkg/store"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
)

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	AllowedOrigins     string
	LiveKitServer      string
	LiveKitAPIKey      string
	LiveKitAPISecret   string
}

type Service struct {
	Config *Config
	Log    logr.Logger
	Store  store.Store
	Cache  any
}

func NewEngine(config *Config) (*gin.Engine, error) {
	log, err := logger.NewDevelopmentLogger()
	if err != nil {
		return nil, err
	}

	st, err := store.NewStore()
	if err != nil {
		log.Error(err, "failed to create store")
		return nil, err
	}

	svc := &Service{
		Config: config,
		Log:    log,
		Store:  st,
		Cache:  nil,
	}

	r := gin.Default()

	r.Use(gin.Logger())
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())

	v1 := r.Group("/")
	{
		// Room endpoints with middleware chain
		rooms := v1.Group("/rooms")
		rooms.Use(
			middleware.Authentication(),
			middleware.RateLimitRoom(),
		)
		{
			rooms.POST("", svc.CreateRoomHandler)
			rooms.GET("/:room_name", svc.GetRoomHandler)

		}

		oauth := v1.Group("/")
		{
			oauth.POST("/callback", svc.CallbackHandler)
		}
	}

	return r, nil
}
