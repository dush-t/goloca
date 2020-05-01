package db

func CreateStore() Store {
	postgres := PostgresConn{}
	postgres.Connect()
	return &postgres
}
