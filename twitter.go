package main

import (
	"database/sql"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"time"
    "github.com/BurntSushi/toml"
)

type Config struct {
	Twitter TwitterConfig
	Db DbConfig
}
type DbConfig {
	User string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	Host string `toml:"host"`
}
type TwitterConfig struct {
	ConsumerKey string `toml:"consumer_key"`
	ConsumerSecret string `toml:"consumer_secret"`
	AccessToken string `toml:"access_token"`
	AccessTokenSecret string `toml:"access_token_secret"`
}

func main() {

	var config Config
    _, err := toml.DecodeFile("config/config.tml", &config)
    if err != nil {
        panic(err)
    }

	anaconda.SetConsumerKey(config.TwitterConfig.ConsumerKey)
	anaconda.SetConsumerSecret(config.TwitterConfig.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.TwitterConfig.AccessToken, config.TwitterConfig.AccessTokenSecret)

	dsl := config.Db.User + ":" + config.Db.Password + "@unix(/tmp/mysql.sock)/" + config.Db.Database
	db, err := sql.Open("mysql", dsl)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close() // 関数がリターンする直前に呼び出される

	v := url.Values{}
	v.Set("count", "100")
	searchResult, _ := api.GetSearch("#ハッシュタグ", v)
	i := 0
	for _, tweet := range searchResult.Statuses {
		fmt.Println("-------------------------------")
		fmt.Println(i)
		//fmt.Println(tweet.Text)
		i++
		log.Printf("%+v", tweet.User.Name)

		query := "INSERT INTO tweets (content, username, tweeted) values(?, ?, ?)"
		//result, err := db.Exec(query, tweet.Text, tweet.User, tweet.CreatedAt)

		createdAt, err := time.Parse(time.RubyDate, tweet.CreatedAt)
		if err != nil {
			log.Fatal("time error: ", err)
		}

		result, err := db.Exec(query, tweet.Text, tweet.User.Name, createdAt)
		if err != nil {
			log.Fatal("insert error: ", err)
		}
		lastID, lerr := result.LastInsertId()
		if lerr != nil { // unsupport...
			fmt.Println("insert last id: %d", lastID)
		}
	}
}
