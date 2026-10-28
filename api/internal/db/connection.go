package db

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GDB *gorm.DB

func Connection() error {
	dsn, err := supabaseDSN(os.Getenv("SUPABASE_DB_URL"))
	if err != nil {
		return err
	}

	connection, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Warn),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return errors.New("no se pudo conectar con Supabase PostgreSQL; verifica SUPABASE_DB_URL y la conectividad de red")
	}

	GDB = connection
	return nil
}

func supabaseDSN(connectionString string) (string, error) {
	if connectionString == "" {
		return "", errors.New("SUPABASE_DB_URL es obligatorio; configura la cadena PostgreSQL de Supabase")
	}
	parsed, err := url.Parse(connectionString)
	if err != nil {
		return "", fmt.Errorf("SUPABASE_DB_URL no es una URL válida: %w", err)
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return "", errors.New("SUPABASE_DB_URL debe ser una URL postgres:// o postgresql://")
	}
	if parsed.Host == "" {
		return "", errors.New("SUPABASE_DB_URL debe incluir el host de PostgreSQL")
	}

	query, err := url.ParseQuery(parsed.RawQuery)
	if err != nil {
		return "", errors.New("SUPABASE_DB_URL contiene parámetros de conexión no válidos")
	}
	sslMode := query.Get("sslmode")
	if sslMode == "" {
		query.Set("sslmode", "require")
	} else if sslMode == "disable" || sslMode == "allow" || sslMode == "prefer" {
		return "", errors.New("SUPABASE_DB_URL debe exigir TLS; usa sslmode=require o una verificación de certificado")
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
