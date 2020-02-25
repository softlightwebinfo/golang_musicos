package libs

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

type GraphqlModel struct {
}

func (a *GraphqlModel) GraphqlHandler(schema *graphql.Schema) gin.HandlerFunc {
	h := handler.New(&handler.Config{
		Schema:   schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
func (a *GraphqlModel) Init(route *gin.Engine, schema *graphql.Schema) {
	route.POST("/graphql", a.GraphqlHandler(schema))
	route.GET("/graphql", a.GraphqlHandler(schema))
}
