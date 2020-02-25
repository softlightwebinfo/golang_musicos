package graphiql

import "github.com/graphql-go/graphql"

type Hello struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (hello *Hello) Query(id int, name string) (hellos []Hello, err error) {
	allHello := []Hello{
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
func GraphHello() *graphql.Field {

	var helloType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "GraphAuth",
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

	var queryHello = graphql.Field{
		Name:        "QueryHello",
		Description: "Query GraphAuth",
		Type:        graphql.NewList(helloType),
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

			return (&Hello{}).Query(id, name)
		},
	}
	return &queryHello
}
