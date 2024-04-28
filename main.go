package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

var (
	tsqlDirectory string
	tsqlSevers    string
	selectGroup   string
	selectServer  int
	queryGroup    string
	queryName     string
	query         string
)

func init() {
	flag.StringVar(&tsqlDirectory, "d", ".tsql-cli", "tsql directory")
	flag.StringVar(&tsqlSevers, "s", "servers.yaml", "tsql servers")
	//Direct to server
	flag.StringVar(&selectGroup, "g", "", "server group name")
	flag.IntVar(&selectServer, "p", 0, "server position")
	flag.StringVar(&query, "q", "", "query")
	//Hot query
	flag.StringVar(&queryGroup, "qg", "", "query group name")
	flag.StringVar(&queryName, "qn", "", "query group query name")

	flag.Parse()
}

func main() {
	if os.Getenv("TSQL_DIRECTORY") != "" {
		tsqlDirectory = os.Getenv("TSQL_DIRECTORY")
	}

	if selectGroup == "" && os.Getenv("TSQL_DEFAULT_GROUP") != "" {
		selectGroup = os.Getenv("TSQL_DEFAULT_GROUP")
	}

	if selectServer == 0 && os.Getenv("TSQL_DEFAULT_SERVER_POSITION") != "" {
		selectServer, _ = strconv.Atoi(os.Getenv("TSQL_DEFAULT_SERVER_POSITION"))
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	tsqlDirectory = filepath.Join(home, tsqlDirectory)

	_, err = os.Stat(tsqlDirectory)
	if err != nil && os.IsNotExist(err) {
		panic(fmt.Errorf("no directory found: %s %s", tsqlDirectory, err))
	}

	serverList, err := LoadServerList(filepath.Join(tsqlDirectory, tsqlSevers))
	if err != nil {
		panic(err)
	}

	queryMap, err := LoadQueryMapDirectory(tsqlDirectory, tsqlSevers)
	if err != nil {
		panic(err)
	}

	if len(queryMap) == 0 {
		panic(fmt.Errorf("no queries found in %s", tsqlDirectory))
	}

	if queryGroup != "" && queryName != "" {
		q, err := queryMap.FindQuery(queryGroup, queryName)
		if err != nil {
			panic(err)
		}
		selectGroup = q.DatabaseGroup
		selectServer = q.DatabasePos
		query = q.Query
	}

	if selectGroup == "" {
		fmt.Println("Hotkeys:")
		for group, queries := range queryMap {
			fmt.Printf("group: %s\n", group)
			for _, query := range queries {
				fmt.Printf("  name: %s\n", query.Name)
				fmt.Printf("  query: %s\n", query.Query)
			}
		}
		os.Exit(0)
	}

	if selectGroup == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	server, err := serverList.FindServer(selectGroup, selectServer)
	if err != nil {
		panic(err)
	}

	db, err := connectDatabase(server)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			panic(err)
		}

		result := make([]any, len(columns))
		resultPtrs := make([]any, len(columns))

		for i := 0; i < len(columns); i++ {
			resultPtrs[i] = &result[i]
		}

		err = rows.Scan(resultPtrs...)
		if err != nil {
			panic(err)
		}

		for i, col := range result {
			val := col
			if b, ok := col.([]byte); ok {
				val = string(b)
			}
			fmt.Printf("Column %s: %v\n", columns[i], val)
		}
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}

}

func connectDatabase(s Server) (*sql.DB, error) {
	var tslString string
	if s.CaFile != "" && s.ClientCert != "" && s.ClientKey != "" {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(s.CaFile)
		if err != nil {
			return nil, err
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatal("Failed to append PEM.")
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
