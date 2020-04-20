package sip

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" //sqlite3 driver import
	"github.com/superirale/sipserver/utils"
)

var locations []AOR

// AOR struct
type AOR struct {
	Username        string
	PhysicalAddress string
}

// NewAOR Create an AOR struct object
func NewAOR(data map[string]string) *AOR {
	aor := new(AOR)
	aor.Username = data["username"]
	aor.PhysicalAddress = data["physicalAddress"]

	return aor
}

//SaveAOR saves aor to DB
func SaveAOR(aor *AOR) (int64, error) {

	db, err := sql.Open("sqlite3", "./sipdb.db")
	stmt, err := db.Prepare("INSERT INTO AOR(username, physical_address) values(?,?)")
	utils.CheckError(err)

	res, err := stmt.Exec(aor.Username, aor.PhysicalAddress)
	utils.CheckError(err)

	return res.LastInsertId()
}

// GetAOR function returns a client locations
func GetAOR(username string) []string {

	var results []string

	db, err := sql.Open("sqlite3", "./sipdb.db")
	rows, err := db.Query("SELECT physical_address FROM AOR WHERE username=?", username)
	utils.CheckError(err)

	defer rows.Close()

	for rows.Next() {
		aor := AOR{}
		err = rows.Scan(&aor.PhysicalAddress)
		utils.CheckError(err)

		results = append(results, aor.PhysicalAddress)
	}

	return results
}

func isAORExists(aor *AOR) bool {

	var username string
	var physicalAddress string
	var check bool

	db, err := sql.Open("sqlite3", "./sipdb.db")
	queryString := "SELECT username, physical_address FROM AOR WHERE username=? AND physical_address=?"
	row := db.QueryRow(queryString, aor.Username, aor.PhysicalAddress)
	err2 := row.Scan(&username, &physicalAddress)

	if err2 != nil {
		if err == sql.ErrNoRows {
			check = false
		}
	}
	if username != "" {
		check = true
	} else {
		check = false
	}
	return check
}

// RemoveAOR function to remove AOR from db
func RemoveAOR(aor *AOR) (int64, error) {

	db, err := sql.Open("sqlite3", "./sipdb.db")
	stmt, err := db.Prepare("DELETE FROM AOR WHERE username=? AND physical_address=?")
	utils.CheckError(err)

	res, err := stmt.Exec(aor.Username, aor.PhysicalAddress)
	utils.CheckError(err)

	return res.RowsAffected()
}
