# News app

News app is a simple application where users can browse and read articles from various RSS providers. 
This project is the backend part of it. The only purpose of the project is to demonstrate the basic Go code architecture and design.

The implementation consists of a single service. It provides a REST API to the client (mobile app) to interact with the domain. This service is also responsible for retrieving new feeds from the providers and refreshing them in its storage, this should be done periodically by some external CRON job to ensure up-to-date content.

## REST API

### Get Articles

  Return list of all articles, optionally it can be filtered by category or provider. Articles are ordered by publication date.

* **URL**

  /v1/feed/articles

* **Method:**
  
  `GET`
  
* **Query Params**

   **Optional:**
 
   `categories=[list of strings]` - filter by categories
   
   `providers=[list of strings]` - filter by providers

### Get Article

  Return a particular article by its id.

* **URL**

  /v1/feed/articles/{id}

* **Method:**
  
  `GET`
  
* **URL Params**

   **Required:**
 
   `id=[string]` - id of the article

* **Error Response:**

  * Article with such id not found <br />
    **Code:** 404 NOT FOUND

### Get Categories

  Return list of all categories.

* **URL**

  /v1/feed/categories

* **Method:**
  
  `GET`

### Get Providers

  Return list of all providers.

* **URL**

  /v1/feed/providers

* **Method:**
  
  `GET`

### Refresh Feed

  Download all available feeds from the providers and replace existing ones in the storage.

* **URL**

  /feedrefresh

* **Method:**
  
  `POST`

## How to run

1. Download the repo
2. Install all dependencies: <br> `go get ./...`
3. Run binary that was compiled by the previous command. It should be located in your *$GOPATH/bin*: <br> `$GOPATH/bin/feed-server`
4. To fill the server with data manually call the */feedrefresh* endpoint: <br> `curl --location --request POST http://localhost:8080/feedrefresh`
5. The API is ready to use

## Implementation details

* Code follows clean (hexagonal) architecture.
* Code is written to be highly testable, but tests itself are missing due to the simplicity of this project.
* Project does not use any framework.
* Service stores all data obtained from the provides in memory, because of that the API and refresh logic run in the same process.
* Dependencies management is done by Go modules.
* Project has golangci-lint enabled to ensure high quality of the code.
