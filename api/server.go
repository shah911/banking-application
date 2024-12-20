package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
)

//server serves HTTP requests for our banking services
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil , fmt.Errorf("cannot create token maker: %d", err) 
	} 
	server := &Server{
		config: config,
		store:  store,
		tokenMaker: tokenMaker,
		router: gin.Default(), 
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}


	server.router.POST("/users", server.createUser)
	server.router.POST("/users/login", server.loginUser)

	authRoutes := server.router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)


	return server, nil
}

//start runs the HTTP server on a specific address.
func(server *Server) Start (addess string) error {
	return server.router.Run(addess)
}

func errorResponse (err error) gin.H {
	return gin.H{"error": err.Error()}
}