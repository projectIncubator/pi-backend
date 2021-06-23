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
To run make on Windows, consider using WSL
