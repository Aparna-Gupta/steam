package master

import (
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/gorilla/context"
	"github.com/h2oai/steamY/lib/fs"
	"github.com/h2oai/steamY/lib/proxy"
	"github.com/h2oai/steamY/lib/rpc"
	"github.com/h2oai/steamY/master/auth"
	"github.com/h2oai/steamY/master/data"
	"github.com/h2oai/steamY/master/web"
)

const (
	defaultWebAddress          = ":9000"
	defaultClusterProxyAddress = ":9001"
	defaultCompilationAddress  = ":8080"
	defaultScoringServiceHost  = ""
)

type Opts struct {
	WebAddress                string
	WorkingDirectory          string
	ClusterProxyAddress       string
	CompilationServiceAddress string
	ScoringServiceHost        string
	EnableProfiler            bool
	YarnKerberosEnabled       bool
	YarnUserName              string
	YarnKeytab                string
	DBName                    string
	DBUserName                string
	DBSSLMode                 string
	SuperuserName             string
	SuperuserPassword         string
}

var DefaultOpts = &Opts{
	defaultWebAddress,
	path.Join(".", fs.VarDir, "master"),
	defaultClusterProxyAddress,
	defaultCompilationAddress,
	defaultScoringServiceHost,
	false,
	false,
	"",
	"",
	"steam",
	"steam",
	"disable",
	"",
	"",
}

type AuthProvider interface {
	Secure(handler http.Handler) http.Handler
	Logout() http.Handler
}

func Run(version, buildDate string, opts *Opts) {
	log.Printf("steam v%s build %s\n", version, buildDate)

	// --- external ip for base and proxy ---
	webAddress := opts.WebAddress
	proxyAddress := opts.ClusterProxyAddress

	// --- set up wd ---
	wd, err := fs.MkWorkingDirectory(opts.WorkingDirectory)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Working directory:", wd)

	// --- www root ---
	wwwroot := fs.GetWwwRoot(wd)
	if _, err := os.Stat(path.Join(wwwroot, "index.html")); err != nil {
		log.Fatalf("Web root not found at %s: %v\n", wwwroot, err)
	}
	log.Println("WWW root:", wwwroot)

	// --- init storage ---

	db, err := data.Connect(opts.DBUserName, opts.DBName, opts.DBSSLMode)
	if err != nil {
		log.Fatalf("Failed connecting to database %s as user %s (SSL=%s): %s\n", opts.DBName, opts.DBUserName, opts.DBSSLMode, err)
	}

	isPrimed, err := data.IsPrimed(db)
	if err != nil {
		log.Fatalln("Failed database version check:", err)
	}

	if !isPrimed {
		if opts.SuperuserName == "" || opts.SuperuserPassword == "" {
			log.Fatalln("Starting Steam for the first time requires both --superuser-name and --superuser-password arguments to \"steam serve master\".")
		}

		if err := auth.ValidateUsername(opts.SuperuserName); err != nil {
			log.Fatalln("Invalid superuser username:", err)
		}

		if err := auth.ValidatePassword(opts.SuperuserPassword); err != nil {
			log.Fatalln("Invalid superuser password:", err)
		}

		if err := data.Prime(db); err != nil {
			log.Fatalln("Failed priming database:", err)
		}
	}

	ds, err := data.NewDatastore(db)
	if err != nil {
		log.Fatalln("Failed initializing from database:", err)
	}

	if !isPrimed {
		passwordHash, err := auth.HashPassword(opts.SuperuserPassword)
		if err != nil {
			log.Fatalln("Failed hashing superuser password:", err)
		}

		if _, _, err := ds.CreateSuperuser(opts.SuperuserName, passwordHash); err != nil {
			log.Fatalln("Failed superuser identity setup:", err)
		}

		_, err = ds.NewPrincipal(opts.SuperuserName)
		if err != nil {
			log.Fatalln("Failed reading superuser principal:", err)
		}
	}

	// --- create basic auth service ---
	defaultAz := newDefaultAz(ds)

	// --- create proxy services ---
	rp := proxy.NewRProxy()
	rpFailch := make(chan error)

	go func() {
		if err := http.ListenAndServe(proxyAddress, rp); err != nil {
			rpFailch <- err
		}
	}()

	// --- create front end api services ---

	webServeMux := http.NewServeMux()
	webServiceImpl := web.NewService(
		defaultAz,
		wd,
		ds,
		opts.CompilationServiceAddress,
		opts.ScoringServiceHost,
		opts.YarnKerberosEnabled,
		opts.YarnUserName,
		opts.YarnKeytab,
	)

	// TODO add CLI args for other other providers
	// TODO conditionally instantiate auth provider based on CLI arg
	authProvider := newBasicAuthProvider(defaultAz, webAddress)

	webServeMux.Handle("/logout", authProvider.Logout())
	webServeMux.Handle("/web", authProvider.Secure(rpc.NewServer(rpc.NewService("web", webServiceImpl))))
	webServeMux.Handle("/upload", authProvider.Secure(newUploadHandler(defaultAz, wd, webServiceImpl.Service)))
	webServeMux.Handle("/", authProvider.Secure(http.FileServer(http.Dir(path.Join(wd, "/www")))))

	if opts.EnableProfiler {
		// --- pprof registrations (no auth) ---
		webServeMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		webServeMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		webServeMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		webServeMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	// --- start http servers ---

	failch := make(chan error)
	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Web server listening at", webAddress)
		prefix := ""
		if len(webAddress) > 1 && webAddress[:1] == ":" {
			prefix = "localhost"
		}
		log.Printf("Point your web browser to http://%s%s/\n", prefix, webAddress)
		if err := http.ListenAndServe(webAddress, context.ClearHandler(webServeMux)); err != nil {
			failch <- err
		}
	}()

	for {
		select {
		case err := <-rpFailch:
			log.Fatalln("HTTP proxy server startup failed:", err)
			return
		case err := <-failch:
			log.Fatalln("HTTP server startup failed:", err)
			return
		case sig := <-sigch:
			log.Println(sig)
			log.Println("Shut down gracefully.")
			return
		}
	}
}
