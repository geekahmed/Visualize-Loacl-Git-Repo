package main

import (
	"flag"
)

func main(){
	var folder string
	var email string

	flag.StringVar(&folder, "add", "", "Add a new folder to scan for Git repos")
	flag.StringVar(&email, "email", "geekahmed1@gmail.com", "The email to scan its git repos")

	flag.Parse()
	if folder != "" {
		scan(folder)
		return
	}
	stats(email)

}