package qiita

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ktsujichan/qiita-sdk-go/qiita"
)

type Qiita struct {
	*qiita.Client
}

type SearchCondition struct {
	User    string `url:"user,omitempty"`
	Title   string `url:"title,omitempty"`
	Body    string `url:"body,omitempty"`
	Code    string `url:"code,omitempty"`
	Tag     string `url:"tag,omitempty"`
	NoTag   string `url:"-tag,omitempty"`
	Stocks  string `url:"stocks,omitempty"`
	Created string `url:"created,omitempty"`
	Updated string `url:"updated,omitempty"`
	Page    uint   `url:"-"`
	PerPage uint   `url:"-"`
}

type Result struct {
	qiita.Items
}

func New(accessToken string) *Qiita {
	c, err := qiita.NewClient(accessToken, *qiita.NewConfig())

	if err != nil {
		panic(fmt.Sprintf("cannot create qiita client: %s", err))
	}

	return &Qiita{c}
}

func (s *SearchCondition) ParseQuery() string {
	elem := reflect.ValueOf(s).Elem()
	size := elem.NumField()

	var query []string

	for i := 0; i < size; i++ {
		tag := elem.Type().Field(i).Tag.Get("url")

		if tag == "-" {
			continue
		}

		value := elem.Field(i).Interface().(string)
		// if value is no set and omitempty, no encode
		if value == "" && strings.Contains(tag, "omitempty") {
			continue
		}

		key := strings.Split(tag, ",")[0]
		query = append(query, fmt.Sprintf("%s:%s", key, value))
	}

	return strings.Join(query, " ")
}

func (q *Qiita) SearchItems(cond *SearchCondition) (Result, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := cond.ParseQuery()

	items, err := q.ListItems(ctx, cond.Page, cond.PerPage, query)

	if err != nil {
		return Result{}, err
	}

	return Result{*items}, nil
}
