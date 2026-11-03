package db

var q* Queries

func Init(db DBTX) {
	q = New(db)
}

func GetQueries() *Queries {
	return q
}
