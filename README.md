# Introduction

Go-Blog was created as a task to learn about go.

# Environment setup

* Production 
    * .env
* Development
    * .env.development.local
* Test
    * .env.test

.env.example
```shell
# Application port
APP_BASE_URL = url.dev
APP_PORT = 8080

# Database settings
DB_USER = user
DB_PASSWORD = password
DB_TABLE = table
DB_URL = localhost
DB_PORT = 3306
```

# Running the server

```sh
go run !(*_test).go
```

Additional flags
```sh
go run !(*_test).go envSet -env=production
```
**envSet** is a flagSet which holds additional flags that can be set
#### envSet flags
* env
    * production|development|test - by default, the env is set to production

# Running tests

Running tests from the root directory
This runs all tests in the project
-p=1 is used to run tests one by one so that the test database would not be dropped/edited by another test.

```sh
go test -p=1 ./...
```

# Endpoints

```shell
# GET
# lists support url params: page, limit, order
# and a combination of all of them
/v1/status #returns app status
/v1/articles #returns a list of articles
/v1/article/:id #returns an article by id
/v1/article/:id/comments #returns a list of comments
# POST
/v1/comment/save #saves a comment to an article
```
A valid request needs to provide **api_key** in the header for authorization (except **/v1/status**).

# Modules used
* [GoDotEnv](https://github.com/joho/godotenv) v1.4.0
* [HttpRouter](https://github.com/julienschmidt/httprouter) v1.3.0
* [Alice](https://github.com/justinas/alice) v1.2.0
* [Testify](https://github.com/stretchr/testify) v1.7.0
* [Go Networking](https://pkg.go.dev/golang.org/x/net) v0.0.0-20211216030914-fe4d6282115f
* [GORM/Mysql](https://gorm.io/) v1.2.1
* [GORM](https://gorm.io/) v1.22.4

# TODO
### Endpoints
* comments
    - [ ] \(optional) edit comment
    - [ ] \(optional) remove comment
* articles
    - [ ] \(optional) save article
    - [ ] \(optional) edit article
    - [ ] \(optional)remove article
- [ ] \(optional) User registration
- [ ] \(optional) User login

### Backend functionality
- [ ] \(optional) possibility to enable/disable migration|seed
- [ ] \(optional) switch httprouter to gin

### React frontend
- [ ] \(optional) create a react frontend app
