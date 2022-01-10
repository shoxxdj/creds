package main 

import (
	"fmt"
	"os"
	"errors"
	"flag"
	"github.com/spf13/viper"
	"github.com/ttacon/chalk"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main(){

	//Get config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.creds/")
	err := viper.ReadInConfig()
	if err != nil {
		if _, err := os.Stat(os.Getenv("HOME")+"/.creds");
		errors.Is(err,os.ErrNotExist){
			fmt.Println("need to create directory")
			os.Mkdir(os.Getenv("HOME")+"/.creds",0755)
			os.Create(os.Getenv("HOME")+"/.creds/config")


			viper.SetDefault("database.dblocation", os.Getenv("HOME")+"/.creds/"+"creds.db")
//			viper.SetDefault("database.dbname","default")

			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath("$HOME/.scope/")

			err :=viper.WriteConfig();
			if err!=nil{
				fmt.Println(err)
			}
		}else{
			fmt.Println("no idea")
		}
	}

	//Main

	//fmt.Println(viper.Get("database.dblocation"))
	//fmt.Println(chalk.Red, "Writing in colors", chalk.Cyan, "is so much fun", chalk.Reset)
	dbLocation:=	flag.String("dl",viper.GetString("database.dbLocation"),"Database location")
	saveConfig:=	flag.Bool("save",false,"Save configuration (need dbLocation to be defined to be efficient)")
	reset:=			flag.Bool("reset",false,"Reset configuration to default")

	fullDetails:=  flag.Bool("full",false,"Get full details (Location, creds and id)")
	loginToAdd:= 	flag.String("l","","Login to add in the database")
	passwordToAdd:=		flag.String("p","","Password to add in the database")
	credsToDel:=		flag.Int("d",-1,"The Creds ID to delete from the database")
	

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Creds : a binary to store creds for attackers. v:0.1\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "\t-%v: %v\n", f.Name,f.Usage) // f.Name, f.Value
		})
	}

	flag.Parse()


	if *saveConfig {
		viper.Set("database.dblocation",*dbLocation)
		err :=viper.WriteConfig();
		if err!=nil{
			fmt.Println(err)
		}
	}

	if *reset{
		viper.Set("database.dblocation", os.Getenv("HOME")+"/.creds/"+"creds.db")
		err :=viper.WriteConfig();
		if err!=nil{
			fmt.Println(err)
		}
		fmt.Println(chalk.Green,"[+]",chalk.Reset,"Configuration reset")
		os.Exit(1)
	}


	db, err := sql.Open("sqlite3", *dbLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS creds (id INTEGER PRIMARY KEY AUTOINCREMENT, login text,password text)
	`
	_,err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	if *credsToDel != -1{
		tx,err := db.Begin()
		if err!=nil{
			log.Fatal(err)
		}
		stmt,err := tx.Prepare("delete from creds where id=?")
		if err!=nil{
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(*credsToDel)
		tx.Commit()
		fmt.Println("Creds deleted")
		sqlStmt = "ALTER TABLE creds RENAME to old_creds";
		_,err = db.Exec(sqlStmt)
		if err!=nil{
			log.Fatal(err)
		}
		sqlStmt = "CREATE TABLE creds (id INTEGER PRIMARY KEY AUTOINCREMENT, login text,password text)"
		_,err = db.Exec(sqlStmt)
		if err!=nil{
			log.Fatal(err)
		}
		sqlStmt = "INSERT INTO creds(login,password) SELECT value from old_creds"
		_,err = db.Exec(sqlStmt)
		if err!=nil{
			log.Fatal(err)
		}
		sqlStmt = "DROP TABLE old_creds"
		_,err = db.Exec(sqlStmt)
		if err!=nil{
			log.Fatal(err)
		}
		fmt.Println(chalk.Red, "[-]", chalk.Reset, "Item deleted")
	}

	if *loginToAdd != "" && *passwordToAdd != ""{
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into creds(login,password) values(?,?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(*loginToAdd,*passwordToAdd)
		tx.Commit()
		fmt.Println(chalk.Green, "[+]", chalk.Reset, "Item added")
	}

	if *loginToAdd == "" && *credsToDel ==-1{
		sqlStmt = "Select id,login,password from creds"
		rows,err := db.Query(sqlStmt)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		if(*fullDetails){
			fmt.Println("|",*dbLocation,"|")
		}

		for rows.Next(){
			var id int
			var login string
			var password string
			err = rows.Scan(&id,&login,&password)
			if err != nil{
				log.Fatal(err)
			}
			if(*fullDetails){
				fmt.Println(id,"|",login,"|",password)
			}else{
				fmt.Println("|",login,"|",password)
			}
		}
	}
}
