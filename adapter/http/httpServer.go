package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gtkmk/finder_api/adapter/http/routes/comment"
	"github.com/gtkmk/finder_api/adapter/http/routes/document"
	"github.com/gtkmk/finder_api/adapter/http/routes/follow"
	"github.com/gtkmk/finder_api/adapter/http/routes/like"
	"github.com/gtkmk/finder_api/adapter/http/routes/post"
	"github.com/gtkmk/finder_api/adapter/http/routes/user"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/client"
	"github.com/gtkmk/finder_api/infra/notification"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/gtkmk/finder_api/core/domain/helper"
)

const (
	OriginLocalhostConst = "http://localhost"
)

const DefaultPortConst = ":8089"

type HttpServer struct {
	server            *http.Server
	app               *gin.Engine
	connection        port.ConnectionInterface
	uuidGenerator     port.UuidInterface
	passwordEncryptor port.EncryptionInterface
}

func NewHttpServer(
	connection port.ConnectionInterface,
	uuidGenerator port.UuidInterface,
	passwordEncryptor port.EncryptionInterface,
) *HttpServer {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	//TODO uncomment this line when we discover how to avoid this panic error http: connection has been hijacked
	//app.Use(middleware.NewTimeoutMiddleware().TimeoutMiddleware())

	serverConfiguration := &http.Server{
		Addr: DefaultPortConst,
	}

	return &HttpServer{
		serverConfiguration,
		app,
		connection,
		uuidGenerator,
		passwordEncryptor,
	}
}

func (httpServer *HttpServer) Start() error {
	httpServer.corsConfig()

	httpServer.registerRutes()

	fmt.Printf("server running in port %s\n", DefaultPortConst)

	httpServer.initialize()

	if err := httpServer.gracefulShutdown(); err != nil {
		return err
	}

	return nil
}

func (httpServer *HttpServer) registerRutes() {
	notificationService := notification.NewNotification(client.NewHttpClient("url"))

	user.NewUserRoutes(
		httpServer.app,
		httpServer.connection,
		httpServer.uuidGenerator,
		httpServer.passwordEncryptor,
		notificationService,
	).Register()

	post.NewPostRoutes(
		httpServer.app,
		httpServer.connection,
		httpServer.uuidGenerator,
		httpServer.passwordEncryptor,
		notificationService,
	).Register()

	like.NewLikeRoutes(
		httpServer.app,
		httpServer.connection,
		httpServer.uuidGenerator,
		httpServer.passwordEncryptor,
		notificationService,
	).Register()

	comment.NewCommentRoutes(
		httpServer.app,
		httpServer.connection,
		httpServer.uuidGenerator,
		httpServer.passwordEncryptor,
		notificationService,
	).Register()
	document.NewDocumentRoutes(
		httpServer.app,
		httpServer.connection,
		httpServer.uuidGenerator,
		httpServer.passwordEncryptor,
		notificationService,
	).Register()

	follow.NewFollowRoutes(
		httpServer.app,
		httpServer.connection,
		httpServer.uuidGenerator,
		httpServer.passwordEncryptor,
		notificationService,
	).Register()
}

func (httpServer *HttpServer) initialize() {
	go func() {
		err := httpServer.server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func (httpServer *HttpServer) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	fmt.Println("kill signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.server.Shutdown(ctx); err != nil {
		return helper.ErrorBuilder(helper.ErrorShuttingDownServerConst, err)
	}

	fmt.Println("Server shutdown gracefully")
	return nil
}

func (httpServer *HttpServer) corsConfig() {
	headers := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "Content-Length", "Accept-Language", "Accept-Encoding", "Connection", "Access-Control-Allow-Origin"})

	origins := handlers.AllowedOrigins([]string{
		OriginLocalhostConst,
	})

	methods := handlers.AllowedMethods([]string{helper.GET, helper.POST, helper.PUT, helper.PATCH, helper.DELETE, helper.OPTIONS})
	credentials := handlers.AllowCredentials()

	corsHandler := handlers.CORS(headers, origins, methods, credentials)(httpServer.app)
	httpServer.server.Handler = corsHandler
}
