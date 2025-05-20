package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"tdl/types"

	"github.com/olivere/elastic/v7"
)

const taskIndexName = "tasks"

type EsRepo struct {
	client *elastic.Client
}

func NewEsRepo(client *elastic.Client) *EsRepo {
	return &EsRepo{client: client}
}

func (r *EsRepo) IndexTask(task *types.Task) error {
	_, err := r.client.Index().
		Index(taskIndexName).
		Id(fmt.Sprintf("%d", task.ID)).
		BodyJson(task).
		Do(context.Background())
	return err
}

func (r *EsRepo) DeleteTask(taskID uint) error {
	_, err := r.client.Delete().
		Index(taskIndexName).
		Id(fmt.Sprintf("%d", taskID)).
		Do(context.Background())
	return err
}

func (r *EsRepo) Search(query string, status types.TaskStatus, userID uint) ([]types.Task, error) {
	// 构建查询条件
	boolQuery := elastic.NewBoolQuery()

	// 添加用户ID过滤
	boolQuery.Must(elastic.NewTermQuery("user_id", userID))

	// 添加状态过滤（如果提供）
	if status != "" {
		boolQuery.Must(elastic.NewTermQuery("status", status))
	} else {
		// 排除已删除任务
		boolQuery.MustNot(elastic.NewTermQuery("status", types.TaskStatusDeleted))
	}

	// 添加搜索条件
	if query != "" {
		multiMatchQuery := elastic.NewMultiMatchQuery(
			query,
			"title", // title 字段权重更高
			"description",
			"tags",
		).Type("best_fields")
		//.Fuzziness("AUTO")

		boolQuery.Must(multiMatchQuery)
	}

	searchResult, err := r.client.Search().
		Index(taskIndexName).
		Query(boolQuery).
		Size(100).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	// 解析搜索结果
	var tasks []types.Task
	for _, hit := range searchResult.Hits.Hits {
		var task types.Task
		if err := json.Unmarshal(hit.Source, &task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
