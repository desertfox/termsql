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
	var tslString string
	if s.CaFile != "" && s.ClientCert != "" && s.ClientKey != "" {
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
		})
		tslString = "?tls=custom"
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s", s.User, s.Pass, s.Host, s.Port, s.Db, tslString))
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
