package main

import (
	"fmt"
	"net/http"

	"github.com/pthomison/dbutils"
	"github.com/pthomison/dbutils/sqlite"
	"github.com/pthomison/errcheck"
	"github.com/pthomison/gormapi"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	rootCmd = &cobra.Command{
		Use:   "golang-gorm-api",
		Short: "golang-gorm-api",
		Run:   run,
	}
)

type APIObject struct {
	gorm.Model

	StringData  string  `json:"string_data,omitempty"`
	IntegerData int     `json:"integer_data,omitempty"`
	FloatData   float64 `json:"float_data,omitempty"`
	BooleanData bool    `json:"boolean_data,omitempty"`
}

func main() {
	err := rootCmd.Execute()
	errcheck.Check(err)
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("--- golang-gorm-api ---")

	client := sqlite.New(":memory:")

	DropAndCreateSamples(client)

	http.HandleFunc("/", gormapi.Index[APIObject](client))
	http.HandleFunc("/all", gormapi.All[APIObject](client))
	http.HandleFunc("/id/", gormapi.ID[APIObject](client))
	http.ListenAndServe(":5050", nil)

}

func DropAndCreateSamples(c dbutils.DBClient) {
	dbutils.Migrate(c, &APIObject{})
	dbutils.DeleteAll(c, &APIObject{})

	var objs []APIObject

	for i := 0; i < 10; i++ {
		o := APIObject{
			StringData:  fmt.Sprintf("%v", i),
			IntegerData: i,
			FloatData:   float64(i),
			BooleanData: i%2 == 0,
		}

		objs = append(objs, o)
	}

	dbutils.Create(c, objs)
}
