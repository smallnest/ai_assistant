package aiassistant

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

type NrIdcPair struct {
	SentFrom string
	SpanID   string
	SrcIDC   string
	DstIDC   string
	Failed   uint32
	Total    uint32
	CreateAt time.Time
	Latency  int64
}

func main() {
	// Connect to the database
	db, err := sql.Open("clickhouse", "tcp://localhost:9000?database=default")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the table if it doesn't exist
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS nr_idc_pair_all (
			sent_from String,
			span_id String,
			src_idc String,
			dst_idc String,
			failed UInt32,
			total UInt32,
			create_at DateTime,
			latency Int64 DEFAULT 0
		) ENGINE = MergeTree(create_at, (sent_from, span_id), 8192)
	`); err != nil {
		panic(err)
	}

	// Insert a new row
	nrIdcPair := NrIdcPair{
		SentFrom: "example",
		SpanID:   "123",
		SrcIDC:   "beijing",
		DstIDC:   "shanghai",
		Failed:   0,
		Total:    1,
		CreateAt: time.Now(),
		Latency:  100,
	}
	if _, err := db.Exec(`
		INSERT INTO nr_idc_pair_all (
			sent_from, span_id, src_idc, dst_idc, failed, total, create_at, latency
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, nrIdcPair.SentFrom, nrIdcPair.SpanID, nrIdcPair.SrcIDC, nrIdcPair.DstIDC, nrIdcPair.Failed, nrIdcPair.Total, nrIdcPair.CreateAt, nrIdcPair.Latency); err != nil {
		panic(err)
	}

	// Update an existing row
	if _, err := db.Exec(`
		UPDATE nr_idc_pair_all SET failed = ? WHERE sent_from = ? AND span_id = ?
	`, 1, "example", "123"); err != nil {
		panic(err)
	}

	// Delete a row
	if _, err := db.Exec(`
		DELETE FROM nr_idc_pair_all WHERE sent_from = ? AND span_id = ?
	`, "example", "123"); err != nil {
		panic(err)
	}

	// Query the table
	rows, err := db.Query(`
		SELECT sent_from, span_id, src_idc, dst_idc, failed, total, create_at, latency
		FROM nr_idc_pair_all
	`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var nrIdcPair NrIdcPair
		if err := rows.Scan(&nrIdcPair.SentFrom, &nrIdcPair.SpanID, &nrIdcPair.SrcIDC, &nrIdcPair.DstIDC, &nrIdcPair.Failed, &nrIdcPair.Total, &nrIdcPair.CreateAt, &nrIdcPair.Latency); err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", nrIdcPair)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
}
