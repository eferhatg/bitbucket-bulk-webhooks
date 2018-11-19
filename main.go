package main

import (
	"fmt"
	"log"
	"time"

	"github.com/eferhatg/git-batch-webhook/provider"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	fmt.Println("Environment arguments loaded")
	btb := provider.NewBitbucket()
	err = btb.Authenticate()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if btb.Token.AccessToken == "" {
		log.Fatal("Couldn't be authenticated to bitbucket api")
		return
	}
	fmt.Println("Authenticated to bitbucket api")
	err = btb.CheckPermissions()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println("Permissions checked and no problem found")
	err = btb.CheckScopes()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println("Scopes checked and no problem found")
	btb.GetRepositories()

	for _, repo := range btb.Repositories {
		time.Sleep(1 * time.Second)
		btb.AddWebHook(repo)
		fmt.Println("Added:" + repo.Fullname)
	}

}
