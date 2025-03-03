package blevesearchpackage

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"log"
)

type BleveRepo struct {
	index bleve.Index
}

func NewBleveRepo(indexPath string) *BleveRepo {
	index, err := bleve.Open(indexPath)
	if err != nil {
		// Nếu chưa có, tạo mới
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &BleveRepo{index: index}
}

// IndexData Index toàn bộ dữ liệu từ MongoDB
func (r *BleveRepo) IndexData(data []map[string]interface{}) error {
	for _, doc := range data {
		docID := fmt.Sprintf("%s_%v", doc["type"], doc["_id"])
		if err := r.index.Index(docID, doc); err != nil {
			return err
		}
	}
	return nil
}

// Search Tìm kiếm dữ liệu
func (r *BleveRepo) Search(query string, limit int) ([]string, error) {
	searchQuery := bleve.NewMatchQuery(query)
	req := bleve.NewSearchRequest(searchQuery)
	req.Size = limit

	results, err := r.index.Search(req)
	if err != nil {
		return nil, err
	}

	var hits []string
	for _, hit := range results.Hits {
		hits = append(hits, hit.ID)
	}
	return hits, nil
}
