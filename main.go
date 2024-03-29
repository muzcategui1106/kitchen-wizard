package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/muzcategui1106/kitchen-wizard/pkg/api"
	"github.com/muzcategui1106/kitchen-wizard/pkg/db/model"
	"github.com/muzcategui1106/kitchen-wizard/pkg/logger"
	rest_middleware "github.com/muzcategui1106/kitchen-wizard/pkg/protocol/rest/middleware"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/db/postgres"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/oidc"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/storage/object"
	"github.com/muzcategui1106/kitchen-wizard/pkg/util/tracing"
)

func main() {
	var logLevel int
	var logTimeFormat string
	var tracingCollectorAddress string
	var dexProviderURL string
	var oidcClientID string
	var OidcClientSecret string
	var oidcRedirectURI string
	var postgresDBHostname string
	var postgresDBPort string
	var postgresDbUsername string
	var postgresDBPassword string
	var s3StoreEndpoint string
	var s3StoreAccessKey string
	var s3StoreSecretKey string

	flag.IntVar(&logLevel, "log-level", 0, "Global log level")
	flag.StringVar(&logTimeFormat, "log-time-format", "2006-01-02T15:04:05Z07:00",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.StringVar(&tracingCollectorAddress, "otp-collector-address", "http://collector-collector.observability.svc:14268/api/traces", "open tracing collector address")
	flag.StringVar(&dexProviderURL, "dex-provider-url", os.Getenv("DEX_PROVIDER_URL"), "provider URL for dex")
	flag.StringVar(&oidcClientID, "oidc-client-id", os.Getenv("OIDC_CLIENT_ID"), "oidc client id to be used for openid flows")
	flag.StringVar(&OidcClientSecret, "oidc-client-secret", os.Getenv("OIDC_CLIENT_SECRET"), "oidc client secret to be used for oauth flows")
	flag.StringVar(&oidcRedirectURI, "oidc-redirect-url", os.Getenv("OIDC_REDIRECT_URL"), "the url with schema that we will use for call backs. Do not include URL path that is hardcoded by default to /v1/auth/oidc/callback")
	flag.StringVar(&postgresDBHostname, "postgres-db-hostname", os.Getenv("POSTGRES_DB_HOSTNAME"), "the hostname of the postgres db")
	flag.StringVar(&postgresDBPort, "postgres-db-port", os.Getenv("POSTGRES_DB_PORT"), "the port for the DB")
	flag.StringVar(&postgresDbUsername, "postgres-db-username", "", "the username for postgres database")
	flag.StringVar(&postgresDBPassword, "postgres-db-password", os.Getenv("POSTGRES_DB_PASSWORD"), "the psssword for the db. This should be passed as an environmental variable for security purposes")
	flag.StringVar(&s3StoreEndpoint, "store-s3-endpoint", "https://minio.kitchen-wizard.svc.cluster.local:443", "the s3 endpoint of the object store")
	flag.StringVar(&s3StoreAccessKey, "store-s3-access-key", os.Getenv("S3_STORE_ACCESS_KEY"), "the s3 endpoint access key ")
	flag.StringVar(&s3StoreSecretKey, "store-s3-secret-key", os.Getenv("S3_STORE_SECRET_KEY"), "the s3 endpoint access key ")

	flag.Parse()

	mainContext := context.Background()

	if err := logger.Init(logLevel, logTimeFormat); err != nil {
		log.Fatalf("failed to initialize logging: %v", err)
	}

	if err := tracing.InitJaegerTracer(mainContext, tracingCollectorAddress); err != nil {
		logger.Log.Sugar().Warnf("could not setup tracing, erro was %v", err)
	}

	dbConn, err := postgres.NewClient(postgresDBHostname, postgresDBPort, postgresDbUsername, postgresDBPassword)
	if err != nil {
		logger.Log.Sugar().Fatal("exiting as it could not connecto to DB")
	}
	if err = model.AutoMigrateSchemas(dbConn); err != nil {
		logger.Log.Sugar().Fatal(err)
	}

	// create object store client
	logger.Log.Debug("creating s3 client")
	s3Client, err := object.CreateClient(s3StoreEndpoint, s3StoreAccessKey, s3StoreSecretKey)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	logger.Log.Debug("created s3 client")
	if err = s3Client.CreateNeccesaryBuckets(); err != nil {
		logger.Log.Fatal(err.Error())
	}

	// create oidc provider config to enable oidc auth
	// creating oidc client and verifier
	oidcPorviderConfig := oidc.ProviderConfig{
		ProviderURL:      dexProviderURL,
		OidcClientID:     oidcClientID,
		OidcClientSecret: OidcClientSecret,
		OidcRedirectURL:  oidcRedirectURI,
	}

	oauth2Config, verifier, err := oidc.CreateOIDCClient(mainContext, oidcPorviderConfig)
	if err != nil {
		logger.Log.Sugar().Fatalf("could not start oidc providerConfig due to %s", err.Error())
	}

	// start http server
	apiConfig := api.ApiServerConfig{
		DBConn:            dbConn,
		ObjectStoreClient: s3Client,
	}
	ApiServer, err := api.NewApiServer(
		mainContext,
		apiConfig,
		api.WithCors(),
		api.WithMiddleware(rest_middleware.StructuredLogger(logger.Log)),
		api.WithTracing(),
		api.WithSessionManagement(),
		api.WithOIDCAuth(oauth2Config, *verifier),
	)
	if err != nil {
		log.Fatalf("could not initialize http server: %v", err)
	}
	go func() {
		ApiServer.Run(mainContext)
	}()

	// run forerver
	stop := make(chan struct{}, 1)
	<-stop
}
