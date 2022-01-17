package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jixindatech/sqlaudit/pkg/config"
	"github.com/jixindatech/sqlaudit/pkg/core/golog"
	"go.uber.org/zap"
	"math"
)

const maxItemGroup = 10

type EsStorage struct {
	config *config.EsConfig
	client *elasticsearch.Client
}

func (e *EsStorage) InitStorage(cfg *config.EsConfig) error {
	var err error

	if e.client != nil {
		return nil
	}

	e.config = cfg
	e.client, err = elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{e.config.Host},
			Username:  e.config.User,
			Password:  e.config.Password,
		},
	)
	if err != nil {
		return err
	}

	res, err := e.client.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (e *EsStorage) Save(body interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return err
	}
	res, err := e.client.Index(e.config.Index, &buf)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (e *EsStorage) Query(_query map[string]interface{}, page, pageSize int) (map[string]interface{}, error) {
	var queries []esquery.Mappable
	name := _query["name"].(string)
	if len(name) > 0 {
		queries = append(queries, esquery.Term("name", name))
	}

	ruleType := _query["type"].(int)
	if ruleType > 0 {
		queries = append(queries, esquery.Term("type", ruleType))
	}

	fingerprint := _query["fingerprint"].(string)
	if len(fingerprint) > 0 {
		queries = append(queries, esquery.Term("fingerprint", fingerprint))
	}

	sql := _query["sql"].(string)
	if len(sql) > 0 {
		queries = append(queries, esquery.Match("sql", sql))
	}

	db := _query["db"].(string)
	if len(db) > 0 {
		queries = append(queries, esquery.Term("db", db))
	}

	ip := _query["ip"].(string)
	if len(ip) > 0 {
		queries = append(queries, esquery.Term("src", ip))
	}

	user := _query["user"].(string)
	if len(user) > 0 {
		queries = append(queries, esquery.Term("user", user))
	}

	op := _query["op"].(int)
	if op != 0 {
		queries = append(queries, esquery.Term("op", op))
	}

	start := _query["start"].(int64) / 1000
	end := _query["end"].(int64) / 1000
	if start > 0 && end > 0 && end > start {
		queries = append(queries, esquery.Range("time").Gte(start).Lte(end))
	} else {
		return nil, errors.New("query time has error")
	}

	query := esquery.Search().Sort("time", "desc").Query(
		esquery.Bool().Must(queries...),
	)
	res, err := query.Run(
		e.client,
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex(e.config.Index),
		e.client.Search.WithFrom((page-1)*pageSize),
		e.client.Search.WithSize(pageSize),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		golog.Error("es", zap.String("err", err.Error()))
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("elasticsearch body has error")
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}

func (e *EsStorage) QueryInfo(_query map[string]interface{}) ([]byte, error) {
	var queries []esquery.Mappable

	db := _query["db"].(string)
	if len(db) > 0 {
		queries = append(queries, esquery.Term("db", db))
	}

	start := _query["start"].(int64)
	end := _query["end"].(int64)
	if start > 0 && end > 0 && end > start {
		queries = append(queries, esquery.Range("time").Gte(start).Lte(end))
	} else {
		return nil, errors.New("query time has error")
	}

	interval := (end - start) / maxItemGroup

	res := esquery.Search().Query(esquery.Bool().Must(queries...))
	if len(db) > 0 {
		res = res.Aggs(
			esquery.TermsAgg("group_by_type", "type"),
			esquery.TermsAgg("group_by_op", "op"),
			esquery.TermsAgg("group_by_src", "src"),
			esquery.TermsAgg("group_by_user", "user"),
			esquery.TermsAgg("group_by_fingerprint", "fingerprint").Size(math.MaxInt32),
			esquery.TermsAgg("group_op_interval", "op").Aggs(
				esquery.CustomAgg("interval_data", map[string]interface{}{
					"histogram": map[string]interface{}{
						"field":    "time",
						"interval": interval,
						"extended_bounds": map[string]interface{}{
							"min": start,
							"max": end,
						},
					},
				})),
			esquery.TermsAgg("group_type_interval", "type").Aggs(
				esquery.CustomAgg("interval_data", map[string]interface{}{
					"histogram": map[string]interface{}{
						"field":    "time",
						"interval": interval,
						"extended_bounds": map[string]interface{}{
							"min": start,
							"max": end,
						},
					},
				})),
		)
	} else {
		res = res.Aggs(
			esquery.TermsAgg("group_by_db", "db"),
			esquery.TermsAgg("group_by_type", "type"),
			esquery.TermsAgg("group_by_src", "src"),
			esquery.TermsAgg("group_by_user", "user"),
			esquery.TermsAgg("group_by_fingerprint", "fingerprint").Size(math.MaxInt32),
			esquery.TermsAgg("group_by_op", "op"),
			esquery.TermsAgg("group_op_interval", "op").Aggs(
				esquery.CustomAgg("interval_data", map[string]interface{}{
					"histogram": map[string]interface{}{
						"field":    "time",
						"interval": interval,
						"extended_bounds": map[string]interface{}{
							"min": start,
							"max": end,
						},
					},
				})),
			esquery.TermsAgg("group_type_interval", "type").Aggs(
				esquery.CustomAgg("interval_data", map[string]interface{}{
					"histogram": map[string]interface{}{
						"field":    "time",
						"interval": interval,
						"extended_bounds": map[string]interface{}{
							"min": start,
							"max": end,
						},
					},
				})),
		)
	}
	response, err := res.Size(0).Run(
		e.client,
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex(e.config.Index),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, errors.New("elasticsearch body has error")
	}

	var b bytes.Buffer
	_, err = b.ReadFrom(response.Body)

	return b.Bytes(), err
}

func (e *EsStorage) QueryFingerPrintInfo(_query map[string]interface{}, page, pageSize int) (map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(_query); err != nil {
		return nil, fmt.Errorf("Error encoding query: %s", err)
	}

	response, err := e.client.Search(
		e.client.Search.WithContext(context.Background()),
		e.client.Search.WithIndex(e.config.Index),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithFrom((page-1)*pageSize),
		e.client.Search.WithSize(pageSize),
		e.client.Search.WithTrackTotalHits(true),
		e.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, errors.New("elasticsearch body has error")
	}

	var r map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}

	return r, nil
}
