package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	// network attributes
	router *echo.Echo
	addr   string
	// server dependencies
	deps *Depends
	// concurrency error handling
	errCh   chan error
	errPool *sync.Pool
}

// Factory for creating a new Server instance
func New(deps *Depends) *Server {
	svr := &Server{
		router: echo.New(),
		deps:   deps,
		addr:   deps.Config.Server.ADDR,
		errCh:  make(chan error),
		errPool: &sync.Pool{
			New: func() interface{} {
				return &echo.HTTPError{}
			},
		},
	}
	// initial server setup
	svr.router.HTTPErrorHandler = svr.handleErrors
	svr.setupMiddlewares()
	svr.setupRoutes()
	return svr
}

func (s *Server) setupRoutes() {
	s.router.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	authRoutes := s.router.Group("/auth")
	authRoutes.POST("/register", s.deps.AuthHdlr.RegisterHandler)
	authRoutes.POST("/login", s.deps.AuthHdlr.LoginHandler)
}

func (s *Server) setupMiddlewares() {
	s.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  s.deps.Config.Server.CORS.AllowOrigins,
		AllowHeaders:  s.deps.Config.Server.CORS.AllowedHeaders,
		ExposeHeaders: s.deps.Config.Server.CORS.AllowedHeaders,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowCredentials: true,
	}))
	s.router.Use(middleware.RequestID())
	logger := log.Logger
	s.router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// Log the request ID, URI, status, and latency
		LogRequestID: true, LogURI: true, LogStatus: true, LogLatency: true,
		// Customize the log format
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("request_id", v.RequestID).
				Str("uri", v.URI).
				Int("status", v.Status).
				Dur("latency", v.Latency).Send()
			return nil
		},
	}))
}

func (s *Server) handleErrors(err error, c echo.Context) {
	// predefined http error response pool
	httpErr := s.errPool.Get().(*echo.HTTPError)
	defer s.errPool.Put(httpErr)

	// set the default error response
	if he, ok := err.(*echo.HTTPError); ok {
		httpErr.Code = he.Code
		httpErr.Message = he.Message
	} else {
		httpErr.Code = http.StatusInternalServerError
		httpErr.Message = "Internal Server Error"
	}

	// log known internal errors
	if errors.As(err, &httpErr) {
		if httpErr.Code >= 500 {
			log.Error().Err(err).Send()
		}
		// TODO: implement custom error handling here
	} else {
		// log unknown internal errors
		log.Error().Err(err).Send()
	}

	// send the error response
	if !c.Response().Committed {
		c.Response().WriteHeader(httpErr.Code)
		if err := jsoniter.NewEncoder(c.Response().Writer).Encode(httpErr); err != nil {
			log.Error().Err(err).Send()
		}
	}
}

func (s *Server) handleShutdown() {
	defer close(s.errCh)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	<-sigch // watch for control-c force shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.router.Shutdown(ctx); err != nil {
		log.Error().Err(err).Send()
		s.errCh <- err
		return
	}

	if err := s.deps.Close(); err != nil {
		log.Error().Err(err).Send()
		s.errCh <- err
		return
	}
	log.Info().Msg("server shutdown gracefully...")
}

func (s *Server) Start() error {
	go s.handleShutdown()
	// start the server and watch for internal errors
	if err := s.router.Start(s.addr); err != http.ErrServerClosed {
		log.Error().Err(err).Send()
		return err
	}
	// handle external errors
	if err, ok := <-s.errCh; ok {
		return err
	}
	return nil
}
