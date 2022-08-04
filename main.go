package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	utils "github.com/pthomison/golang-utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	rootCmd = &cobra.Command{
		Use:   "golang-gorm-api",
		Short: "golang-gorm-api",
		Run:   run,
	}

	dbClient = &utils.DBClient{}
)

type APIObject struct {
	gorm.Model

	ID uint `json:"id,omitempty" gorm:"primaryKey"`

	StringData  string  `json:"string_data,omitempty"`
	IntegerData int     `json:"integer_data,omitempty"`
	FloatData   float64 `json:"float_data,omitempty"`
	BooleanData bool    `json:"boolean_data,omitempty"`
}

func main() {
	dbClient.RegisterFlags(rootCmd)

	err := rootCmd.Execute()
	utils.Check(err)
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("--- golang-gorm-api ---")

	dbClient.InitializeClient(logger.Silent)

	DropAndCreateSamples()

	http.HandleFunc("/", Index)
	http.HandleFunc("/all", All)
	http.HandleFunc("/id/", ID)
	http.ListenAndServe(":5050", nil)

}

func DropAndCreateSamples() {
	dbClient.DB.Migrator().DropTable(&APIObject{})
	dbClient.DB.AutoMigrate(&APIObject{})

	objs := []APIObject{}

	for i := 0; i < 10; i++ {
		o := APIObject{
			StringData:  fmt.Sprintf("%v", i),
			IntegerData: i,
			FloatData:   float64(i),
			BooleanData: i%2 == 0,
		}

		objs = append(objs, o)
	}

	utils.Create(dbClient, objs)
}

func Index(w http.ResponseWriter, r *http.Request) {
	objs := []APIObject{}

	objs = utils.SelectAll[APIObject](dbClient, []string{"id"})

	json, err := json.Marshal(objs)
	utils.Check(err)

	_, err = w.Write(json)
	utils.Check(err)
}

func All(w http.ResponseWriter, r *http.Request) {
	objs := []APIObject{}

	objs = utils.SelectAll[APIObject](dbClient, []string{})

	json, err := json.Marshal(objs)
	utils.Check(err)

	_, err = w.Write(json)
	utils.Check(err)
}

func ID(w http.ResponseWriter, r *http.Request) {
	id_str := strings.TrimPrefix(r.URL.Path, "/id/")
	id, err := strconv.Atoi(id_str)
	utils.Check(err)

	fmt.Println(id)

	objs := []APIObject{}
	objs = utils.SelectWhere[APIObject](dbClient, []string{}, "id = ?", id)

	json, err := json.Marshal(objs)
	utils.Check(err)

	_, err = w.Write(json)
	utils.Check(err)
}
