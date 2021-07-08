package main

import (
	_ "embed"
	"fmt"
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

	logrus.WithField("port", port).Info("startin to many foos")
	_ = http.ListenAndServe(":"+port, r)
}

type service struct {
	Name    string
	Version string
	Schema  string
}

type foo struct {
	ID           graphql.ID
	GraphGophers bool
}

type resolver struct {
	Service service
}

func newResolver() *resolver {
	return &resolver{
		Service: service{
			Name:    "foo1",
			Version: "0.1",
			Schema:  schema,
		},
	}
}

func (r *resolver) Foo(args struct {
	ID graphql.ID
}) (*foo, error) {
	return &foo{
		ID: args.ID,
	}, nil
}

func (r *resolver) GetFoo(args struct {
	ID graphql.ID
}) (*foo, error) {
	return r.Foo(args)
}

func (r *resolver) ToManyFoos() ([]*foo, error) {
	logrus.Info("TooManyFoos")
	foos := make([]*foo, 20)
	for i := 0; i < 20; i++ {
		foos[i] = &foo{
			ID: graphql.ID(fmt.Sprintf("%d", i%10)),
		}
	}
	return foos, nil
}
