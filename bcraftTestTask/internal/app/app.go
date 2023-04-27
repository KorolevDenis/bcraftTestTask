package app

import (
	"bcraftTestTask/internal/DBController"
	"bcraftTestTask/internal/auth"
	"bcraftTestTask/internal/controllers"
	"bcraftTestTask/internal/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type App struct {
	router     *mux.Router
	recipeCtr  *controllers.RecipeController
	serviceCtr *controllers.ServiceController
	authCtr    *controllers.AuthController
}

func NewApp() (*App, error) {

	recipeCtr := controllers.NewRecipeController(repository.NewRecipeRepository(DBController.NewDBController().GetDB()))
	authCtr := controllers.NewAuthController(repository.NewAuthRepository(DBController.NewDBController().GetDB()))
	serviceCtr := controllers.NewServiceController()

	app := &App{
		recipeCtr:  recipeCtr,
		serviceCtr: serviceCtr,
		authCtr:    authCtr,
	}

	router, err := app.newRouter()
	if err != nil {
		return nil, err
	}
	app.router = router

	return app, nil
}

func (a App) Start() error {
	log.Println("Start REST API")
	err := http.ListenAndServe(":8000", a.router)

	if err != nil {
		return err
	}

	return nil
}

func (a App) newRouter() (*mux.Router, error) {
	//isReady value needed for the readiness probe
	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		time.Sleep(10 * time.Second)
		isReady.Store(true)
	}()

	router := mux.NewRouter()
	router.Use(auth.JwtAuthentication)

	router.HandleFunc("/probes/liveness",
		a.serviceCtr.Liveness)

	router.HandleFunc("/probes/readiness",
		a.serviceCtr.Readiness(isReady))

	router.HandleFunc("/recipe",
		a.recipeCtr.GetRecipes).Methods("GET")

	router.HandleFunc("/recipe/{id:[0-9]+}",
		a.recipeCtr.GetRecipeById).Methods("GET")

	router.HandleFunc("/recipe",
		a.recipeCtr.PutRecipe).Methods("PUT")

	router.HandleFunc("/recipe/{id:[0-9]+}",
		a.recipeCtr.DeleteRecipeById).Methods("DELETE")

	router.HandleFunc("/user/new",
		a.authCtr.CreateAccount).Methods("POST")

	router.HandleFunc("/user/login",
		a.authCtr.Login).Methods("POST")

	return router, nil
}
