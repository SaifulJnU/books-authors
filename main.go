package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/saifujnu/books-authors/auth"
	"github.com/saifujnu/books-authors/config"
	"github.com/saifujnu/books-authors/controllers"
	"github.com/saifujnu/books-authors/db/mongo"
)

var (
	Logger *zap.Logger

	// Define Prometheus metrics
	apiRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_api_requests_total",
			Help: "Total number of API requests.",
		},
		[]string{"method"},
	)

	//--------------------------boook and author--------------------------------
	successfulBookAuthorsFetch = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "myapp_fetch_books_and_authors",
			Help: "Total number of successful fetching",
		},
	)
	// histogram for request duration
	requestDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "myapp_request_duration_seconds",
			Help: "Histogram of request duration",
		},
		[]string{"route"},
	)

	// gauge metric for system status
	systemStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "myapp_system_status",
			Help: "Current system status",
		},
	)

	//This is working
	// Define a custom Prometheus metric for successful logins
	successfulLogins = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "myapp_successful_logins_total",
			Help: "Total number of successful logins.",
		},
	)

	// Define Prometheus metric for Books API requests
	booksAPIRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_books_api_requests_total",
			Help: "Total number of Books API requests.",
		},
		[]string{"method"},
	)
)

func InitializeLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment() // NewDevelopment also works for info debug and log
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config.SetEnvionment()

	var errLogger error
	Logger, errLogger = InitializeLogger() //error handleing of logger function
	if errLogger != nil {
		panic("Failed to initialize Zap logger: " + errLogger.Error())
	}
}

func main() {
	m, err := mongo.Connect()
	if err != nil {
		Logger.Error("Failed to connect to MongoDB", zap.Error(err))
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(ginzap.Ginzap(Logger, time.RFC3339, true)) //wrapping zan with gin now it will give us logger as json
	router.Use(ginzap.RecoveryWithZap(Logger, true))

	authorController := controllers.NewAuthorController(m, Logger)
	bookController := controllers.NewBookController(m, Logger)
	authController := controllers.NewAuthController(m, Logger)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", authController.Signup)
		authRoutes.POST("/login", func(c *gin.Context) {

			successfulLogins.Inc()
			authController.Login(c)
		})
	}

	// Middleware to increment API request count
	router.Use(func(c *gin.Context) {
		apiRequests.WithLabelValues(c.Request.Method).Inc()
		c.Next()
	})

	start := time.Now()

	bookRoutes := router.Group("/books")

	bookRoutes.Use(auth.JWTMiddleware())
	{
		bookRoutes.GET("/", bookController.GetBooks)
		bookRoutes.GET("/:id", bookController.GetBookByID)
		bookRoutes.POST("/", bookController.CreateBook)
		bookRoutes.PUT("/:id", bookController.UpdateBook)
		bookRoutes.DELETE("/:id", bookController.DeleteBook)
		bookRoutes.GET("/books-and-authors", func(c *gin.Context) {
			successfulBookAuthorsFetch.Inc()
			bookController.GetAllBooksAndAuthors(c)

		})
		duration := time.Since(start)
		requestDurationHistogram.WithLabelValues("books-and-authors").Observe(duration.Seconds())
		systemStatus.Set(1)

		bookRoutes.GET("/books-by-author/:authorName", bookController.GetBooksByAuthorName)
	}

	// Middleware to increment Books API request count
	bookRoutes.Use(func(c *gin.Context) {
		booksAPIRequests.WithLabelValues(c.Request.Method).Inc()
		c.Next()
	})

	authorRoutes := router.Group("/authors")
	authorRoutes.Use(auth.JWTMiddleware())
	{
		authorRoutes.GET("/", authorController.GetAuthors)
		authorRoutes.GET("/:id", authorController.GetAuthorByID)
		authorRoutes.POST("/", authorController.CreateAuthor)
		authorRoutes.PUT("/:id", authorController.UpdateAuthor)
		authorRoutes.DELETE("/:id", authorController.DeleteAuthor)
	}

	// Register the custom metrics to be exposed
	prometheus.MustRegister(successfulLogins)
	prometheus.MustRegister(booksAPIRequests)
	//prometheus.MustRegister(successfulBookAuthorsFetch)
	prometheus.MustRegister(successfulBookAuthorsFetch, requestDurationHistogram, systemStatus)

	// Expose /metrics endpoint for Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	Logger.Info("Server started on :8080")

	err1 := router.Run(":8080")
	if err1 != nil {
		Logger.Error("Failed to start server", zap.Error(err1))
		os.Exit(1)
	}
}
