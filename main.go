package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         "127.0.0.1:8080",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	// Register handlers.
	handleFunc("/rolldice", rolldice)

	// Add HTTP instrumentation for the whole server.
	handler := otelhttp.NewHandler(mux, "/")
	return handler
}

// package main

// import (
// 	"context"
// 	"errors"
// 	"log"
// 	"net"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"time"

// 	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
// )

// func main() {
// 	if err := run(); err != nil {
// 		log.Fatalln(err)
// 	}
// }

// func run() (err error) {
// 	// Handle SIGINT (CTRL+C) gracefully.
// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer stop()

// 	// Set up OpenTelemetry.
// 	otelShutdown, err := setupOTelSDK(ctx)
// 	if err != nil {
// 		return
// 	}
// 	// Handle shutdown properly so nothing leaks.
// 	defer func() {
// 		err = errors.Join(err, otelShutdown(context.Background()))
// 	}()

// 	// Start HTTP server.
// 	srv := &http.Server{
// 		Addr:         ":8080",
// 		BaseContext:  func(_ net.Listener) context.Context { return ctx },
// 		ReadTimeout:  time.Second,
// 		WriteTimeout: 10 * time.Second,
// 		Handler:      newHTTPHandler(),
// 	}
// 	srvErr := make(chan error, 1)
// 	go func() {
// 		srvErr <- srv.ListenAndServe()
// 	}()

// 	// Wait for interruption.
// 	select {
// 	case err = <-srvErr:
// 		// Error when starting HTTP server.
// 		return
// 	case <-ctx.Done():
// 		// Wait for first CTRL+C.
// 		// Stop receiving signal notifications as soon as possible.
// 		stop()
// 	}

// 	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
// 	err = srv.Shutdown(context.Background())
// 	return
// }

// func newHTTPHandler() http.Handler {
// 	mux := http.NewServeMux()

// 	// handleFunc is a replacement for mux.HandleFunc
// 	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
// 	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
// 		// Configure the "http.route" for the HTTP instrumentation.
// 		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
// 		mux.Handle(pattern, handler)
// 	}

// 	// Register handlers.
// 	handleFunc("/rolldice/", rolldice)
// 	handleFunc("/rolldice/{player}", rolldice)

// 	// Add HTTP instrumentation for the whole server.
// 	handler := otelhttp.NewHandler(mux, "/")
// 	return handler
// }

// // package main

// // import (
// // 	"net/http"

// // 	"github.com/gin-gonic/gin"
// // )

// // // album represents data about a record album.
// // type album struct {
// // 	ID     string  `json:"id"`
// // 	Title  string  `json:"title"`
// // 	Artist string  `json:"artist"`
// // 	Price  float64 `json:"price"`
// // }

// // // albums slice to seed record album data.
// // var albums = []album{
// // 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// // 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// // 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// // }

// // // getAlbums responds with the list of all albums as JSON.
// // func getAlbums(c *gin.Context) {
// // 	c.IndentedJSON(http.StatusOK, albums)
// // }

// // // postAlbums adds an album from JSON received in the request body.
// // func postAlbums(c *gin.Context) {
// // 	var newAlbum album

// // 	// Call BindJSON to bind the received JSON to
// // 	// newAlbum.
// // 	if err := c.BindJSON(&newAlbum); err != nil {
// // 		return
// // 	}

// // 	// Add the new album to the slice.
// // 	albums = append(albums, newAlbum)
// // 	c.IndentedJSON(http.StatusCreated, newAlbum)
// // }

// // // getAlbumByID locates the album whose ID value matches the id
// // // parameter sent by the client, then returns that album as a response.
// // func getAlbumByID(c *gin.Context) {
// // 	id := c.Param("id")

// // 	// Loop over the list of albums, looking for
// // 	// an album whose ID value matches the parameter.
// // 	for _, a := range albums {
// // 		if a.ID == id {
// // 			c.IndentedJSON(http.StatusOK, a)
// // 			return
// // 		}
// // 	}
// // 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// // }

// // func main() {
// // 	router := gin.Default()
// // 	router.GET("/albums", getAlbums)
// // 	router.GET("/albums/:id", getAlbumByID)
// // 	router.POST("/albums", postAlbums)

// // 	router.Run("localhost:8080")
// // }
