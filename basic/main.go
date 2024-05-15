package main

import (
	"fmt"

	"github.com/9ssi7/age"
)

var dsn string = "host=127.0.0.1 port=5455 dbname=postgres user=postgres password=postgres sslmode=disable"
var graphName string = "myGraph"

type Person struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
}

type WorksWith struct {
	ID     int64 `json:"id"`
	Weight int   `json:"weight"`
}

func main() {
	ag := age.New(age.Config{
		GraphName: graphName,
		Dsn:       dsn,
	})
	_, err := ag.Prepare()

	if err != nil {
		panic(err)
	}

	tx, err := ag.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(0, "CREATE (n:Person {name: '%s'})", "Joe")
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(0, "CREATE (n:Person {name: '%s', age: %d})", "Smith", 10)
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(0, "CREATE (n:Person {name: '%s', weight:%f})", "Jack", 70.3)
	if err != nil {
		panic(err)
	}

	tx.Commit()

	tx, err = ag.Begin()
	if err != nil {
		panic(err)
	}

	cursor, err := tx.Exec(1, "MATCH (n:Person) RETURN n")
	if err != nil {
		panic(err)
	}

	count := 0
	for cursor.Next() {
		entities, err := cursor.GetRow()
		if err != nil {
			panic(err)
		}
		count++
		vertex := entities[0].(*age.Vertex)
		fmt.Println(count, "]", vertex.Id(), vertex.Label(), vertex.Props())
	}

	fmt.Println("Vertex Count:", count)

	_, err = tx.Exec(0, "MATCH (a:Person), (b:Person) WHERE a.name='%s' AND b.name='%s' CREATE (a)-[r:workWith {weight: %d}]->(b)",
		"Jack", "Joe", 3)
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(0, "MATCH (a:Person {name: '%s'}), (b:Person {name: '%s'}) CREATE (a)-[r:workWith {weight: %d}]->(b)",
		"Joe", "Smith", 7)
	if err != nil {
		panic(err)
	}

	tx.Commit()

	tx, err = ag.Begin()
	if err != nil {
		panic(err)
	}

	cursor, err = tx.Exec(1, "MATCH p=()-[:workWith]-() RETURN p")
	if err != nil {
		panic(err)
	}

	count = 0
	for cursor.Next() {
		entities, err := cursor.GetRow()
		if err != nil {
			panic(err)
		}
		count++

		path := entities[0].(*age.Path)

		vertexStart := path.GetAsVertex(0)
		edge := path.GetAsEdge(1)
		vertexEnd := path.GetAsVertex(2)

		fmt.Println(count, "]", vertexStart, edge.Props(), vertexEnd)
	}

	// Query with return many columns
	cursor, err = tx.Exec(3, "MATCH (a:Person)-[l:workWith]-(b:Person) RETURN a, l, b")
	if err != nil {
		panic(err)
	}

	count = 0
	for cursor.Next() {
		row, err := cursor.GetRow()
		if err != nil {
			panic(err)
		}

		count++

		person := &Person{}
		err = age.ParseStruct(row[0], person)
		if err != nil {
			fmt.Println("Error on ParseStruct:", err)
		}

		workWith := &WorksWith{}
		err = age.ParseStruct(row[1], workWith)
		if err != nil {
			fmt.Println("Error on ParseStruct:", err)
		}

		fmt.Println("ROW.person.name:", person.Name)
		fmt.Println("ROW.workWith.weight:", workWith.Weight)

		v1 := row[0].(*age.Vertex)
		edge := row[1].(*age.Edge)
		v2 := row[2].(*age.Vertex)

		fmt.Println("ROW ", count, ">>", "\n\t", v1, "\n\t", edge, "\n\t", v2)
	}

	_, err = tx.Exec(0, "MATCH (n:Person) DETACH DELETE n RETURN *")
	if err != nil {
		panic(err)
	}
	tx.Commit()
}
