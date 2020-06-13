package main

import (
	"fmt"
	"regexp"

	phonedb "github.com/saurabh-sikchi/go_phone-number-normalizer/db"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "dev"
	password = "saurabh"
	dbname   = "gophercises_phone"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phonedb.Migrate("postgres", psqlInfo))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(db.Seed())

	phones, err := db.AllPhones()
	must(err)

	for _, p := range phones {
		number := normalize(p.Number)
		if p.Number != number {
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				// delete this number
				must(db.DeletePhone(p.ID))
			} else {
				p.Number = number
				must(db.UpdatePhone(&p))
			}
		}
	}

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// Unused
// func getPhone(db *sql.DB, id int) (string, error) {
// 	var phone string

// 	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id = $1", id).Scan(&phone)

// 	if err != nil {
// 		return "", err
// 	}

// 	return phone, nil
// }

// func normalize(phone string) string {
// 	var buf bytes.Buffer
// 	for _, c := range phone {
// 		if c >= '0' && c <= '9' {
// 			buf.WriteRune(c)
// 		}
// 	}
// 	return buf.String()
// }

func normalize(phone string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(phone, "")
}
