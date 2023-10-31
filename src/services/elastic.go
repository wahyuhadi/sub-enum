package services

import (
	"errors"
	"log"
	"sub/model"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

var (
	// this only dummy credential folks
	Uelas = "elastic"
	Pelas = "changeme"
	Helas = "http://192.168.1.17 :9200"
)

func Els(data model.SubDomainMetaData) error {
	cfg := elasticsearch.Config{
		Addresses: []string{
			Helas,
		},
		Username: Uelas,
		Password: Pelas,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return errors.New("error connection")
	}

	res, err := es.Index("enumeration", esutil.NewJSONReader(&data))
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	defer res.Body.Close()
	log.Println(res)

	return nil
}
