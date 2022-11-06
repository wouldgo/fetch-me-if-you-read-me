package main

import (
	"errors"
	"fetch-me-if-you-read-me/imaginer"
	"fetch-me-if-you-read-me/model"
	"fetch-me-if-you-read-me/server"

	logging "fetch-me-if-you-read-me/logger"

	"flag"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Options struct {
	PostgresqlConfigurations *model.PostgresqlConfigurations
	Logger                   *logging.Logger
	Imaginer                 *imaginer.ImaginerConfs
	Server                   *server.ServerConfs
}

func parseOptions() (*Options, error) {
	var parentConfig zap.Config
	host := flag.String("host", "0.0.0.0", "Host where server will listen")
	port := flag.String("port", "3000", "Port where server will listen")
	imageColor := flag.String("image-color", "#00FFFF", "Image color")

	logEnvironment := flag.String("log-environment", "", "Log environment")

	var logLevel logging.LoggingLevel
	flag.Var(&logLevel, "log-level", "log level")

	postgresqlAdministrator := flag.String("postgresql-administrator", "", "postgresql database administrator username")
	postgresqlAdministratorPassword := flag.String("postgresql-administrator-password", "", "postgresql database administrator password")
	postgresqlHost := flag.String("postgresql-host", "", "hostname of postgresql server")
	postgresqlDatabase := flag.String("postgresql-database", "", "postgresql database name")
	postgresqlUsername := flag.String("postgresql-username", "", "postgresql user")
	postgresqlPassword := flag.String("postgresql-password", "", "postgresql password")
	postgresqlThreads := flag.Int("postgresql-threads", 1, "number of thread for postgresql client")

	flag.Parse()

	hostEnv, hostEnvSet := os.LookupEnv("HOST")
	portEnv, portEnvSet := os.LookupEnv("PORT")
	imageColorEnv, imageColorSet := os.LookupEnv("IMAGE_COLOR")

	logLevelEnv, logLevelEnvSet := os.LookupEnv("LOG_LEVEL")
	logEnvironmentEnv, logEnvironmentEnvSet := os.LookupEnv("LOG_ENVIRONMENT")

	postgresqlAdministratorEnv, postgresqlAdministratorEnvSet := os.LookupEnv("POSTGRESQL_ADMINISTRATOR")
	postgresqlAdministratorPasswordEnv, postgresqlAdministratorPasswordEnvSet := os.LookupEnv("POSTGRESQL_ADMINISTRATOR_PASSWORD")
	postgresqlHostEnv, postgresqlHostEnvSet := os.LookupEnv("POSTGRESQL_HOST")
	postgresqlDatabaseEnv, postgresqlDatabaseEnvSet := os.LookupEnv("POSTGRESQL_DATABASE")
	postgresqlUsernameEnv, postgresqlUsernameEnvSet := os.LookupEnv("POSTGRESQL_USERNAME")
	postgresqlPasswordEnv, postgresqlPasswordEnvSet := os.LookupEnv("POSTGRESQL_PASSWORD")
	postgresqlThreadsEnv, postgresqlThreadsEnvSet := os.LookupEnv("POSTGRESQL_THREADS")

	if hostEnvSet {
		host = &hostEnv
	}

	if portEnvSet {
		port = &portEnv
	}

	if imageColorSet {
		imageColor = &imageColorEnv
	}

	rgbaColor, err := parseHexColor(*imageColor)

	if err != nil {
		return nil, errors.New("image color is not a valid hex value")
	}

	imaginerConf := imaginer.ImaginerConfs{
		Color: &rgbaColor,
	}

	serverConf := server.ServerConfs{
		Host: *host,
		Port: *port,
	}

	if logLevelEnvSet {
		logLevel = logging.LoggingLevelFrom(logLevelEnv)
	}

	if logEnvironmentEnvSet {
		logEnvironment = &logEnvironmentEnv
	}

	if strings.EqualFold(*logEnvironment, "production") {
		parentConfig = zap.NewProductionConfig()
	} else {
		parentConfig = zap.NewDevelopmentConfig()
	}

	config := zap.Config{
		Level:            logLevel.ToZap(),
		Encoding:         "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    parentConfig.EncoderConfig,
	}

	logger, err := config.Build()

	if err != nil {
		return nil, err
	}

	defer logger.Sync()
	sugar := logger.Sugar()

	if postgresqlAdministratorEnvSet {
		postgresqlAdministrator = &postgresqlAdministratorEnv
	}

	if postgresqlAdministratorPasswordEnvSet {
		postgresqlAdministratorPassword = &postgresqlAdministratorPasswordEnv
	}

	if postgresqlHostEnvSet {
		postgresqlHost = &postgresqlHostEnv
	}

	if postgresqlDatabaseEnvSet {
		postgresqlDatabase = &postgresqlDatabaseEnv
	}

	if postgresqlUsernameEnvSet {
		postgresqlUsername = &postgresqlUsernameEnv
	}

	if postgresqlPasswordEnvSet {
		postgresqlPassword = &postgresqlPasswordEnv
	}

	if postgresqlThreadsEnvSet {
		postgresqlThreadsFromEnv, err := strconv.ParseInt(postgresqlThreadsEnv, 10, 32)
		if err != nil {
			return nil, err
		}

		*postgresqlThreads = int(postgresqlThreadsFromEnv)
	}

	if postgresqlAdministrator == nil ||
		postgresqlAdministratorPassword == nil ||
		postgresqlHost == nil ||
		postgresqlDatabase == nil ||
		postgresqlUsername == nil ||
		postgresqlPassword == nil ||
		strings.EqualFold(*postgresqlAdministrator, "") ||
		strings.EqualFold(*postgresqlAdministratorPassword, "") ||
		strings.EqualFold(*postgresqlHost, "") ||
		strings.EqualFold(*postgresqlDatabase, "") ||
		strings.EqualFold(*postgresqlUsername, "") ||
		strings.EqualFold(*postgresqlPassword, "") {

		return nil, errors.New("Postgresql configuration is not set")
	}

	applicationName, err := os.Executable()
	if err != nil {

		return nil, err
	}

	return &Options{
		PostgresqlConfigurations: &model.PostgresqlConfigurations{
			Administrator:         postgresqlAdministrator,
			AdministratorPassword: postgresqlAdministratorPassword,
			Host:                  postgresqlHost,
			Database:              postgresqlDatabase,
			Username:              postgresqlUsername,
			Password:              postgresqlPassword,
			Threads:               postgresqlThreads,
			ApplicationName:       applicationName,
		},
		Logger: &logging.Logger{
			Log: sugar,
		},
		Imaginer: &imaginerConf,
		Server:   &serverConf,
	}, nil
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}
