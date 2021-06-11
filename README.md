# Project Incubator

This repo contains the source code for the Project Incubator backend, which is built in Golang.

## Setting up a development environment

### Requirements
- You will need GO to be installed. For convenience, we will be using the Goland IDE.
- It will help if you have Postman installed for checking of the REST endpoints.

**Important**: If this is a first-time setup, you must run the following command from the project's root directory.
```sh
# Installs all go dependencies within the go.mod file
$ go get
```  
If you at any point in time you need to do "go get NAME_OF_MODULE", you should add this to the go.mod as well.  

### Starting the app

```sh
# start the backend using the provided Makefile
$ make run
```

This let's you run the app on localhost. It will do so on port 8000.

```sh
# On Windows to configure env vars in powershell, use the following
$ $Env:DATABASE_URL='postgres://uzwwzorwjgovha:7b761af65f1da20457023b2938c433e099ead5c52397438cda64ee69d1096a91@ec2-52-72-65-76.compute-1.amazonaws.com:5432/d5oipjlaqis6ph'
$ $Env:PORT=8000

# To run without a Makefile use
$ go run .
```
