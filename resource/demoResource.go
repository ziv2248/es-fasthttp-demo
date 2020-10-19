package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	. "webapiserver/internal"

	"github.com/olivere/elastic/v7"
	"github.com/valyala/fasthttp"
	"gitlab.bcowtech.de/bcow-go/httparg"
)

const (
	PersonIndex string = "user"
)

type (
	// DemoResource type
	DemoResource struct {
		ServiceProvider *ServiceProvider
	}
	// Person type
	Person struct {
		FirstName string
		LastName  string
		Id        string
	}
	// argument types
	GetDemoArgs struct {
		LastName string `query:"lastName"`
	}

	PostDemoArgs struct {
		FirstName string `query:"firstName"`
		LastName  string `query:"lastName"`
	}

	PutDemoArgs struct {
		ID       string `query:"id"`
		LastName string `query:"lastName"`
	}

	DeleteDemoArgs struct {
		ID string `query:"id"`
	}
)

func (r *DemoResource) Init() {}

func (r *DemoResource) GET(ctx *fasthttp.RequestCtx) {
	defer func() {
		recover()
	}()

	args := GetDemoArgs{}
	httparg.NewProcessor(&args, httparg.ProcessorOption{
		ErrorHandleFunc: func(err error) {
			fmt.Printf("%% error: %+v\n", err)
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			panic(err)
		},
	}).ProcessQueryString(ctx.QueryArgs().String())

	queryData := elastic.NewBoolQuery()
	/*
	 *  term:  完全匹配搜尋
	 *  match: 分詞搜尋
	 *  ＊要轉小寫
	 */
	queryData.Must(elastic.NewTermQuery("LastName", strings.ToLower(args.LastName)))
	/*
	 *  Get():    獲取文檔(Id)
	 *  Search(): 搜尋文檔(Query)
	 */
	temp := r.ServiceProvider.ESClient.Search().Index(PersonIndex).Query(queryData)
	esResponse, _ := temp.Do(context.Background())
	var response []Person
	for _, value := range esResponse.Hits.Hits {
		resPersonBody := Person{}
		json.Unmarshal([]byte(value.Source), &resPersonBody)
		resPersonBody.Id = value.Id
		response = append(response, resPersonBody)
	}
	responseStr, _ := json.Marshal(response)
	if len(response) == 0 {
		responseStr, _ = json.Marshal(Person{FirstName: "", LastName: "", Id: ""})
	}
	ctx.Success("text/plain", []byte(responseStr))
}

func (r *DemoResource) Post(ctx *fasthttp.RequestCtx) {
	defer func() {
		recover()
	}()

	args := PostDemoArgs{}
	httparg.NewProcessor(&args, httparg.ProcessorOption{
		ErrorHandleFunc: func(err error) {
			fmt.Printf("%% error: %+v\n", err)
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			panic(err)
		},
	}).ProcessQueryString(ctx.QueryArgs().String())

	PersonDetail := Person{FirstName: "Demo," + args.FirstName, LastName: args.LastName}
	esResponse, err := r.ServiceProvider.ESClient.Index().
		Index(PersonIndex).
		BodyJson(PersonDetail).
		Do(context.Background())
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		panic(err)
	}
	fmt.Printf("Indexed user %s to index %s, type %s\n", esResponse.Id, esResponse.Index, esResponse.Type)

	ctx.Success("text/plain", []byte("Success!!"))
}

func (r *DemoResource) Put(ctx *fasthttp.RequestCtx) {
	defer func() {
		recover()
	}()

	args := PutDemoArgs{}
	httparg.NewProcessor(&args, httparg.ProcessorOption{
		ErrorHandleFunc: func(err error) {
			fmt.Printf("%% error: %+v\n", err)
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			panic(err)
		},
	}).ProcessQueryString(ctx.QueryArgs().String())

	esResponse, err := r.ServiceProvider.ESClient.Update().
		Index(PersonIndex).
		Id(args.ID).
		Doc(map[string]interface{}{"LastName": args.LastName}).
		Do(context.Background())
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		panic(err)
	}
	fmt.Println(esResponse.Result)
	ctx.Success("text/plain", []byte("Success!!"))
}

func (r *DemoResource) Delete(ctx *fasthttp.RequestCtx) {
	defer func() {
		recover()
	}()

	args := DeleteDemoArgs{}
	httparg.NewProcessor(&args, httparg.ProcessorOption{
		ErrorHandleFunc: func(err error) {
			fmt.Printf("%% error: %+v\n", err)
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
			panic(err)
		},
	}).ProcessQueryString(ctx.QueryArgs().String())

	esResponse, err := r.ServiceProvider.ESClient.Delete().
		Index(PersonIndex).
		Id(args.ID).
		Do(context.Background())
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		panic(err)
	}
	fmt.Println(esResponse.Result)
	ctx.Success("text/plain", []byte("Success!!"))
}
