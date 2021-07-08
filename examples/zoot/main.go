package main

import (
	_ "embed"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sirupsen/logrus"
)

//go:embed schema.graphql
var schema string

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := newResolver()
	parsedSchema := graphql.MustParseSchema(schema, resolver, graphql.UseFieldResolvers())

	r := mux.NewRouter()
	r.Handle("/query", &relay.Handler{Schema: parsedSchema})

	logrus.WithField("port", port).Info("starting zoot")
	_ = http.ListenAndServe(":"+port, r)
}

type service struct {
	Name    string
	Version string
	Schema  string
}

type foo struct {
	ID   graphql.ID
	Zoot bool
	Bar  bool
}

type resolver struct {
	Service service
}

func newResolver() *resolver {
	return &resolver{
		Service: service{
			Name:    "graph-gophers-service",
			Version: "0.1",
			Schema:  schema,
		},
	}
}

func (r *resolver) Foo(args struct {
	ID graphql.ID
}) (*foo, error) {
	logrus.WithField("id", args.ID).Info("get a Foo")
	return &foo{
		ID:   args.ID,
		Zoot: true,
	}, nil
}

func (r *resolver) Foos(args struct {
	IDs []graphql.ID
}) ([]*foo, error) {
	logrus.WithField("ids", args.IDs).Info("get some Foos")
	foos := make([]*foo, len(args.IDs))
	for i := range foos {
		foos[i] = &foo{
			ID:   args.IDs[i],
			Zoot: true,
			Bar:  false,
		}
	}
	return foos, nil
}
