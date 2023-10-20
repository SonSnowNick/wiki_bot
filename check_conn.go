package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

func main() {
	coll := get_coll_names()
	fmt.Println(coll)
}

func get_coll_names() []string {
	connStr := "user=postgres password=... host=ip port=5432 dbname=postgres sslmode=disable"

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = '...';
	`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var columnNames []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			log.Fatal(err)
		}
		columnNames = append(columnNames, columnName)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return columnNames
}
