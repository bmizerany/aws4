package dydb_test

import (
	"fmt"
	"github.com/bmizerany/aws4/dydb"
	"log"
)

func init() {
	log.SetFlags(0)
}

func Example_listTables() {
	db := new(dydb.DB)

	var resp struct{ TableNames []string }
	if err := db.Do("ListTables", nil).Decode(&resp); err != nil {
		log.Fatal(err)
	}

	// Output:
	// ["Users", "Posts"]
	fmt.Printf("%q", resp.TableNames)
}
