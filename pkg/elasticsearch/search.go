package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"log"
)

// TaskDocument 定义任务文档结构
type TaskDocument struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// SearchService 搜索服务
type SearchService struct {
	Client *elastic.Client
}

// NewSearchService 创建搜索服务
func NewSearchService(url string) (*SearchService, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	return &SearchService{Client: client}, nil
}

// SearchTasks 搜索任务
func (s *SearchService) SearchTasks(ctx context.Context, query, status string) ([]TaskDocument, error) {
	// 创建布尔查询
	boolQuery := elastic.NewBoolQuery()

	// 添加关键词搜索条件
	if query != "" {
		boolQuery.Must(elastic.NewMultiMatchQuery(query, "title", "description"))
	}

	// 添加状态过滤条件
	if status != "" {
		boolQuery.Filter(elastic.NewTermQuery("status", status))
	}

	// 执行搜索
	searchResult, err := s.Client.Search().
		Index("task_index").
		Query(boolQuery).
		Highlight(elastic.NewHighlight().Field("title").Field("description")).
		Size(100).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	// 处理结果
	var tasks []TaskDocument
	for _, hit := range searchResult.Hits.Hits {
		var task TaskDocument
		err := json.Unmarshal(hit.Source, &task)
		if err != nil {
			log.Printf("Error unmarshalling task: %v", err)
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
