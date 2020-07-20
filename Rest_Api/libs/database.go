package libs

import (
	"database/sql"
	"fmt"

	//database library
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "amiramir12008"
	dbname   = "quotes_db"
)

// Quote model
type Quote struct {
	ID        int     `json:"id"`
	QuoteText string  `json:"quote"`
	Likes     int     `json:"likes"`
	Subject   string  `json:"subject"`
	Writer    *Writer `json:"writer"`
}

//Writer model
type Writer struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

//LikeRes model
type LikeRes struct {
	RowEffected int64 `json:"roweffected"`
}

//QuoteDB : quotes from database
type QuoteDB struct {
	ID        int
	QuoteText string
	Likes     int
	SubjectID int
	WriterID  int
}

//DbConnection : connection to database
func DbConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

//GetAllQuotes : get all quotes
func GetAllQuotes(db *sql.DB) []Quote {
	var result []Quote
	quoteDB := getAllDBQuotes(db)
	for _, quote := range quoteDB {
		fQuote := Quote{ID: quote.ID, QuoteText: quote.QuoteText, Likes: quote.Likes, Subject: getSubject(db, quote.SubjectID), Writer: getWriter(db, quote.WriterID)}
		result = append(result, fQuote)
	}
	db.Close()
	return result
}

//UpdateLikes : update likes count in database
func UpdateLikes(db *sql.DB, toAdd int64, id string) LikeRes {
	query := "UPDATE quotes SET likes = likes + $1 WHERE id=$2"
	res, err := db.Exec(query, toAdd, id)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	db.Close()
	result := LikeRes{RowEffected: count}
	return result
}

func getAllDBQuotes(db *sql.DB) []QuoteDB {
	var result []QuoteDB
	rows, err := db.Query("SELECT * FROM quotes ORDER BY id")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var quoteDB QuoteDB
		err := rows.Scan(&quoteDB.ID, &quoteDB.QuoteText, &quoteDB.Likes, &quoteDB.SubjectID, &quoteDB.WriterID)
		if err == nil {
			result = append(result, quoteDB)
		}
	}

	return result
}

func getSubject(db *sql.DB, subjectID int) string {
	var result string
	row := db.QueryRow("SELECT subject FROM subjects WHERE id=$1;", subjectID)
	switch err := row.Scan(&result); err {
	case sql.ErrNoRows:
		return "null"
	case nil:
		return result
	default:
		panic(err)
	}
}

func getWriter(db *sql.DB, writerID int) *Writer {
	var result Writer
	row := db.QueryRow("SELECT name, lastname FROM writers WHERE id=$1;", writerID)
	switch err := row.Scan(&result.FirstName, &result.LastName); err {
	case sql.ErrNoRows:
		return nil
	case nil:
		return &result
	default:
		panic(err)
	}
}
