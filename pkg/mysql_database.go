package termsql

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Pinger interface {
	Ping() error
}

func MySQLConnect(s Server) (*sql.DB, error) {
	if s.Port == 0 {
		s.Port = 3306
	}

	config := mysql.Config{
		User:                 s.User,
		Passwd:               s.Pass,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%d", s.Host, s.Port),
		DBName:               s.Db,
		AllowNativePasswords: true,
		TLSConfig:            "skip-verify",
	}

	if DEBUG {
		fmt.Println(config.FormatDSN())
	}

	if s.CaFile != "" && s.ClientCert != "" && s.ClientKey != "" {
		config.TLSConfig = "custom"
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(s.CaFile)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, fmt.Errorf("failed to append PEM")
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(s.ClientCert, s.ClientKey)
		if err != nil {
			return nil, err
		}
		clientCert = append(clientCert, certs)
		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:      rootCertPool,
			Certificates: clientCert,
			MinVersion:   tls.VersionTLS12,
			MaxVersion:   tls.VersionTLS12,
		})
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func PingDB[T Pinger](db T) error {
	return db.Ping()
}

func RunQueryDynamic(db *sql.DB, q *Query, params ...string) (map[string]string, error) {
	rows, err := db.Query(q.Query, params)
	if err != nil {
		return nil, fmt.Errorf("error running query: %s", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %s", err)
	}

	var (
		results    = make(map[string]string, len(columns))
		result     = make([]string, len(columns))
		resultPtrs = make([]interface{}, len(columns))
	)

	for i := 0; i < len(columns); i++ {
		resultPtrs[i] = &result[i]
	}

	for rows.Next() {
		if err := rows.Scan(resultPtrs...); err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err)
		}

		for i, col := range result {
			results[columns[i]] = col
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %s", err)
	}

	return results, nil
}
