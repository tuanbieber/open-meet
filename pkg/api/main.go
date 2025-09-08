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
	gin.SetMode(gin.ReleaseMode)

	r.Use(gin.Logger())
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())

	room := r.Group("/rooms").Use(middleware.Authentication(), middleware.RateLimit())
	{
		room.POST("", svc.CreateRoomHandler)
		room.GET("/:room_name", svc.GetRoomHandler)

	}

	oauth := r.Group("/")
	{
		oauth.POST("/callback", svc.CallbackHandler)
	}

	participant := r.Group("/").Use(middleware.Authentication())
	{
		participant.POST("/livekit-tokens", svc.LiveKitTokenHandler)
	}

	return r, nil
}
