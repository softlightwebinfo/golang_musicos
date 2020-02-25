package settings

import (
	"context"
	"github.com/machinebox/graphql"
	"log"
)

type GraphClient struct {
	client  *graphql.Client
	request *graphql.Request
}

func (a *GraphClient) GraphClient() {
	config := GetGraphqlConfig()
	a.client = graphql.NewClient(config.Client)
}

func (a *GraphClient) Request(request string) {
	a.request = graphql.NewRequest(request)
	a.request.Header.Set("x-hasura-admin-secret", "pass")

}

func (a *GraphClient) SetVar(key, value string) {
	a.request.Var(key, value)
}
func (a *GraphClient) Run(ctx context.Context, respData interface{}) {
	if err := a.client.Run(ctx, a.request, &respData); err != nil {
		log.Fatal("error", err)
	}
}
