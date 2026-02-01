package container

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"

	dbConfig "github.com/DucTran999/dbkit/config"

	"auth-service/config"
	"auth-service/pkg/projectpath"

	"github.com/DucTran999/cachekit"
	"github.com/DucTran999/dbkit"
	"github.com/DucTran999/jwtkit"
	"github.com/DucTran999/shared-pkg/logger"
)

func newLogger(cfg *config.EnvConfiguration) (logger.ILogger, error) {
	return logger.NewLogger(logger.Config{
		Environment: cfg.ServiceEnv,
		LogToFile:   cfg.LogToFile,
		FilePath:    cfg.LogFilePath,
	})
}

func newAuthDBConnection(config *config.EnvConfiguration) (dbkit.Connection, error) {
	pgConf := dbConfig.PostgreSQLConfig{
		Config: dbConfig.Config{
			Host:     config.DBHost,
			Port:     config.DBPort,
			Username: config.DBUsername,
			Password: config.DBPasswd,
			Database: config.DBDatName,
			TimeZone: config.DBTimezone,
		},
		PoolConfig: dbConfig.PoolConfig{
			MaxOpenConnection: config.DBMaxOpenConnections,
			MaxIdleConnection: config.DBMaxIdleConnections,
			ConnMaxIdleTime:   time.Duration(config.DBMaxConnectionIdleTime) * time.Second,
		},
		SSLMode: dbConfig.PgSSLDisable,
	}

	conn, err := dbkit.NewPostgreSQLConnection(pgConf)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newSigner(config *config.EnvConfiguration) (jwtkit.JWT, error) {
	root := projectpath.MustRoot()

	// Read the private key PEM file
	sanitizedPrivate := filepath.Base(config.PrivPem)
	privPem, err := os.ReadFile(path.Join(root, "./keys/"+sanitizedPrivate)) // #nosec G304
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Parse the private key
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privPem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Read the public key PEM file
	sanitizedPublic := filepath.Base(config.PubPem)
	pubPem, err := os.ReadFile(path.Join(root, "./keys/"+sanitizedPublic)) // #nosec G304
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}

	// Parse the public key
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(pubPem)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	cfg := jwtkit.Config{
		Alg:        jwtkit.SigningAlgorithm(config.SignMethod),
		RSAPrivate: signKey,
		RSAPublic:  verifyKey,
	}

	return jwtkit.NewJWT(cfg)
}

func newRedisCache(config *config.EnvConfiguration) (cachekit.RemoteCache, error) {
	cacheConf := cachekit.RedisConfig{
		Host:     config.RedisHost,
		Port:     config.RedisPort,
		Password: config.RedisPasswd,
	}

	return cachekit.NewRedisCache(cacheConf)
}
