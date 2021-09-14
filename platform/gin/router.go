package gin

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Routers contains the functions of http handler to clean payloads and pass it the service
type Routers interface {
	Serve(ms ...gin.HandlerFunc)
}

// Router data will be registered to http listener
type Router struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

type routing struct {
	host    string
	port    string
	routers []Router
}

// NewRouting is for initialize the handler
func NewRouting(host, port string, routers []Router) Routers {
	return &routing{
		host,
		port,
		routers,
	}
}

// Serve is to start serving the HTTP Listener for every domain
func (r *routing) Serve(ms ...gin.HandlerFunc) {
	server := gin.Default()
	server.Use(ms...)
	v1 := server.Group("/v1")

	if os.Getenv("ENV") == "development" {
		server.Use(corsMW())
	}
	for _, router := range r.routers {
		if len(router.Middlewares) != 0 {
			// Append the router to the middlware
			router.Middlewares = append(router.Middlewares, router.Handler)
			v1.Handle(router.Method, router.Path, router.Middlewares...)
		} else {
			v1.Handle(router.Method, router.Path, router.Handler)
		}
	}

	server.Static("/assets", "./assets")

	logrus.WithFields(logrus.Fields{
		"host": r.host,
		"port": r.port,
	}).Info("Starts Serving on HTTP")

	log.Fatal(http.ListenAndServe(r.host+":"+r.port, server))

}

func corsMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, OPTIONS, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
