package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/olivere/elastic.v1"
)

const (
	DEFAULT_ELASTICSEARCH_HOST = "http://localhost:9200"
	DEFAULT_INDEX_PREFIX       = "logstash-"
	DEFAULT_DATE_FORMAT        = "2006.01.02"
	DEFAULT_NUM_RESULTS        = 500
)

type fieldName []string

var (
	ec             *elastic.Client
	numResults     *int
	indexPrefix    *string
	stringQuery    = "*"
	jsonFlag       *bool
	fieldSeparator           = " "
	fieldNameFlag  fieldName = []string{"@timestamp", "message"}
)

func init() {
	var err error
	var host string
	host = os.Getenv("ELASTICSEARCH_HOST")
	if host == "" {
		host = DEFAULT_ELASTICSEARCH_HOST
	}
	ec, err = elastic.NewClient(http.DefaultClient, host)
	if err != nil {
		log.Fatal("can't connect to elasticsearch", host, err)
	}
}

func (s *fieldName) String() string {
	return fmt.Sprint(*s)
}

func (s *fieldName) Set(value string) error {
	if len(*s) > 0 {
		*s = []string{}
	}
	for _, f := range strings.Split(value, ",") {
		*s = append(*s, f)
	}
	return nil
}

func config() {
	numResults = flag.Int("m", DEFAULT_NUM_RESULTS, "max number of results")
	indexPrefix = flag.String("prefix", DEFAULT_INDEX_PREFIX, "index name prefix (before date)")
	jsonFlag = flag.Bool("j", false, "results in json")
	flag.Var(&fieldNameFlag, "f", "fields to return (comma separated)")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		stringQuery = strings.Join(args, " ")
	}
}

func main() {
	config()
	duration, _ := time.ParseDuration("-0h")
	indexDate := time.Now().Add(duration).UTC().Format(DEFAULT_DATE_FORMAT)
	index := *indexPrefix + indexDate
	// fmt.Println(index)
	// fmt.Println("Query - ",stringQuery)
	// fmt.Println(*numResults)
	// fmt.Println(*numResults)
	// fmt.Println(fieldNameFlag)

	q := elastic.NewQueryStringQuery(stringQuery)
	res, err := ec.Search().
		Index(index).
		Query(q).
		From(0).Size(*numResults).
		Sort("@timestamp", true).
		Pretty(true).
		Do()

	if err != nil {
		log.Fatal("issue with query", err)
	}

	if res.Hits != nil {
		fmt.Printf("Found a total of %d messages\n", res.Hits.TotalHits)
		for _, hit := range res.Hits.Hits {
			// hit.Index contains the name of the index

			if *jsonFlag {
				fmt.Printf("%s\n", *hit.Source)
			} else {
				// Deserialize hit.Source into a message (could also be just a map[string]interface{}).
				var m map[string]interface{}
				err := json.Unmarshal(*hit.Source, &m)
				if err != nil {
					log.Fatal("Bad json")
				}

				var fields []string
				for _, name := range fieldNameFlag {
					fields = append(fields, fmt.Sprintf("%v", m[name]))
				}
				fmt.Printf("%s\n", strings.Join(fields, fieldSeparator))
			}
		}
	}
}
