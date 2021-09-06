package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type server struct {
	db       database
	router   *mux.Router
	conf     *Config
	n        *negroni.Negroni
	jmw      *jwtmiddleware.JWTMiddleware
	userpass map[string]string
}

func (s *server) serve(endpoint string) {
	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Accept-Language", "Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})
	handler := cor.Handler(s.router)
	s.n.Use(negronilogrus.NewMiddleware())
	s.n.UseHandler(handler)
	log.Fatal(http.ListenAndServe(endpoint, s.n))
}

func (s *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Index page / requested")
		fmt.Fprintf(w, "Quasar is running!")
	}
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

func (s *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		if pass, ok := s.userpass[user.Username]; ok {
			if pass != user.Password {
				http.Error(w, "Password Incorrect", 403)
				log.Error(fmt.Sprintf("Incorrect password for %s", user.Username))
				return
			}
		} else {
			http.Error(w, fmt.Sprintf("No user %s exists.", user.Username),
				403)
			log.Error(fmt.Sprintf("No user %s exists.", user.Username))
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
		})
		tokenString, err := token.SignedString([]byte(s.conf.Authsecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error(err)
			return
		}
		err = json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error(err)
			return
		}
		log.Println("Returning token: ", tokenString)
	}
}

func (s *server) routes() {
	s.router.HandleFunc("/", s.handleHome()).Methods("GET")
	// networks endpoints
	s.router.Handle("/api/login", negroni.New(
		negroni.Wrap(http.HandlerFunc(s.handleLogin())),
	)).Methods("POST")
	s.router.Handle("/api/networks/all", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleGetAllNetworks())),
	)).Methods("GET")
	s.router.Handle("/api/networks/new", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleNewNetwork())),
	)).Methods("POST")
	s.router.Handle("/api/networks/{NETWORK}/cert", negroni.New(
		negroni.Wrap(http.HandlerFunc(s.handleGetNetworkCert())),
	)).Methods("GET")
	s.router.Handle("/api/networks/{NETWORK}/info", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleNetworkInfo())),
	)).Methods("GET")
	s.router.Handle("/api/networks/{NETWORK}/delete", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleDeleteNetwork())),
	)).Methods("DELETE")
	s.router.Handle("/api/networks/{NETWORK}/update", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleUpdateNetwork())),
	)).Methods("POST")
	// node management endpoints
	s.router.Handle("/api/networks/{NETWORK}/nodes/all", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleGetAllNodes())),
	)).Methods("GET")
	s.router.Handle("/api/networks/{NETWORK}/nodes/{NODENAME}/approve", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleApproveNode())),
	)).Methods("POST")
	s.router.Handle("/api/networks/{NETWORK}/nodes/{NODENAME}/info", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleNodeInfo())),
	)).Methods("GET")
	s.router.Handle("/api/networks/{NETWORK}/nodes/{NODENAME}/disable", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleDisableNode())),
	)).Methods("POST")
	s.router.Handle("/api/networks/{NETWORK}/nodes/{NODENAME}/update", negroni.New(
		negroni.HandlerFunc(s.jmw.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(s.handleUpdateNode())),
	)).Methods("POST")
	// neutron endpoints
	s.router.HandleFunc("/api/neutron/join", s.handleJoinNetwork()).Methods("POST")
	s.router.HandleFunc("/api/neutron/config", s.handleGetConfig()).Methods("GET")
	s.router.HandleFunc("/api/neutron/leave", s.handleLeaveNetwork()).Methods("POST")
}

func runServe(configPath string,
	listenAddress string,
	listenPort int,
) {
	s := new(server)
	var err error
	s.conf, err = NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if listenAddress == "" {
		listenAddress = s.conf.Quasar.Listen.Host
	}
	if listenPort == 0 {
		listenPort = s.conf.Quasar.Listen.Port
	}

	endpoint := listenAddress + ":" + fmt.Sprint(listenPort)

	log.WithFields(log.Fields{
		"config": configPath,
	}).Info("Loaded config")

	if s.conf.Database.Type == "bolt" {
		s.db = new(boltdbi)
		err = s.db.connect(s.conf.Database.Source)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.WithFields(log.Fields{
			"requested_db": s.conf.Database.Type,
		}).Fatal("Currently only bolt is supported as a database type.")
	}
	s.conf.Authsecret = os.Getenv("QUASAR_AUTHSECRET")
	if s.conf.Authsecret == "" {
		log.Fatal("Environment variable QUASAR_AUTHSECRET cannot be empty.")
	}
	adminpassword := os.Getenv("QUASAR_ADMINPASS")
	if adminpassword == "" {
		log.Fatal("Environment variable QUASAR_ADMINPASS cannot be empty.")
	}

	log.WithFields(log.Fields{
		"endpoint": endpoint,
	}).Info("Starting Quasar server")
	s.router = mux.NewRouter().StrictSlash(true)
	s.n = negroni.New()
	s.userpass = make(map[string]string)
	s.userpass["admin"] = adminpassword
	s.jmw = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(s.conf.Authsecret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	s.routes()
	s.serve(endpoint)
}
