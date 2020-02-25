package graphiql

import (
	"github.com/graphql-go/graphql"
)

type GraphAuth struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (hello *GraphAuth) Query(id int, name string) (hellos []GraphAuth, err error) {
	allHello := []GraphAuth{
		{0, "vinli"},
		{1, "daisy"},
	}
	for _, v := range allHello {
		if id >= 0 && len(name) <= 0 {
			if v.Id == id {
				hellos = append(hellos, v)
			}
		}

		if id < 0 && len(name) > 0 {
			if v.Name == name {
				hellos = append(hellos, v)
			}
		}

		if id >= 0 && len(name) > 0 {
			if v.Name == name && v.Id == id {
				hellos = append(hellos, v)
			}
		}

		if id < 0 && len(name) <= 0 {
			hellos = append(hellos, v)
		}
	}
	return
}
func GraphQueAuth() *graphql.Field {

	var authType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Auth",
		Description: "GraphAuth Model",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	var query = graphql.Field{
		Name:        "Auth",
		Description: "Query GraphAuth",
		Type:        graphql.NewList(authType),
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (result interface{}, err error) {
			id, _ := p.Args["id"].(int)
			name, _ := p.Args["name"].(string)

			return (&GraphAuth{}).Query(id, name)
		},
	}
	return &query
}
func GraphMutAuth() *graphql.Field {
	todoType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Todo",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	return &graphql.Field{
		Type: todoType, // the return type for this field
		Args: graphql.FieldConfigArgument{
			"username": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (result interface{}, err error) {
			return
		},
	}
}
