package models

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

func GetEvents(query map[string]interface{}, page int, pageSize int) (map[string]interface{}, error) {
	data, err := Storage.Query(query, page, pageSize)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	res["count"] = data["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"]
	res["data"] = data["hits"].(map[string]interface{})["hits"]

	return res, err
}

type bucketStringItem struct {
	DocCount int64  `json:"doc_count"`
	Key      string `json:"key"`
}

type bucketIntItem struct {
	DocCount int64 `json:"doc_count"`
	Key      int64 `json:"key"`
}

type bucketFloatItem struct {
	DocCount int     `json:"doc_count"`
	Key      float64 `json:"key"`
}

type bucketIntervalData struct {
	Buckets []bucketFloatItem `json:"buckets"`
}
type bucketIntervalItem struct {
	Key      int                `json:"key"`
	DocCount int                `json:"doc_count"`
	Interval bucketIntervalData `json:"interval_data"`
}

const maxItemLength = 10

func GetEventInfo(query map[string]interface{}) (map[string]interface{}, error) {
	data, err := Storage.QueryInfo(query)
	if err != nil {
		return nil, err
	}

	resDbs := []string{}
	db := gjson.GetBytes(data, "aggregations.group_by_db.buckets")
	if db.Exists() {
		var items []bucketStringItem
		err := json.Unmarshal([]byte(db.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			resDbs = append(resDbs, item.Key)
		}
	}

	resTypes := make(map[int64]int64)
	types := gjson.GetBytes(data, "aggregations.group_by_type.buckets")
	if types.Exists() {
		var items []bucketIntItem
		err := json.Unmarshal([]byte(types.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			resTypes[item.Key] = item.DocCount
		}
	}

	resOps := make(map[int64]int64)
	ops := gjson.GetBytes(data, "aggregations.group_by_op.buckets")
	if ops.Exists() {
		var items []bucketIntItem
		err := json.Unmarshal([]byte(ops.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			resOps[item.Key] = item.DocCount
		}
	}

	resOpInterval := make(map[int]interface{})
	opInterval := gjson.GetBytes(data, "aggregations.group_op_interval.buckets")
	if opInterval.Exists() {
		var items []bucketIntervalItem
		err := json.Unmarshal([]byte(opInterval.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			var docCount []int
			for _, v := range item.Interval.Buckets {
				docCount = append(docCount, v.DocCount)
			}
			resOpInterval[item.Key] = docCount
		}
	}

	resIP := make(map[string]interface{})
	ips := gjson.GetBytes(data, "aggregations.group_by_src.buckets")
	if ips.Exists() {
		var ipList []string
		var docCount []int64
		var items []bucketStringItem
		err := json.Unmarshal([]byte(ips.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			ipList = append(ipList, item.Key)
			docCount = append(docCount, item.DocCount)
		}

		if len(ipList) == 0 {
			ipList = []string{}
			docCount = []int64{}
		}

		resIP["item"] = ipList
		resIP["num"] = docCount
	}

	resUser := make(map[string]interface{})
	users := gjson.GetBytes(data, "aggregations.group_by_user.buckets")
	if users.Exists() {
		var userList []string
		var docCount []int64
		var items []bucketStringItem
		err := json.Unmarshal([]byte(users.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			userList = append(userList, item.Key)
			docCount = append(docCount, item.DocCount)
		}

		if len(userList) == 0 {
			userList = []string{}
			docCount = []int64{}
		}

		resUser["item"] = userList
		resUser["num"] = docCount
	}

	resFingerprints := make(map[string]interface{})
	fingerprints := gjson.GetBytes(data, "aggregations.group_by_fingerprint.buckets")
	if fingerprints.Exists() {
		var list []string
		var docCount []int64
		var items []bucketStringItem
		err := json.Unmarshal([]byte(fingerprints.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			list = append(list, item.Key)
			docCount = append(docCount, item.DocCount)
		}
		
		if len(list) == 0 {
			list = []string{}
			docCount = []int64{}
		}

		resFingerprints["item"] = list
		resFingerprints["num"] = docCount
	}

	resTypeInterval := make(map[int]interface{})
	typeInterval := gjson.GetBytes(data, "aggregations.group_type_interval.buckets")
	if typeInterval.Exists() {
		var items []bucketIntervalItem
		err := json.Unmarshal([]byte(typeInterval.Raw), &items)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			var docCount []int
			for _, v := range item.Interval.Buckets {
				docCount = append(docCount, v.DocCount)
			}
			resTypeInterval[item.Key] = docCount
		}
	}

	res := make(map[string]interface{})
	res["db"] = resDbs
	res["types"] = resTypes
	res["ops"] = resOps
	res["opinfo"] = resOpInterval
	res["ip"] = resIP
	res["user"] = resUser
	res["fingerprint"] = resFingerprints
	res["typeinfo"] = resTypeInterval

	return res, nil
}

func SaveEvent(body interface{}) error {
	return Storage.Save(body)
}
