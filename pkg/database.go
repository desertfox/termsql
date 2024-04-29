package termsql

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func Connect(s Server) (*sql.DB, error) {
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
