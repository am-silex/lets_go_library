## Overview

This repository contains an example REST API application written in Go. It's a backend for
a proverbial library app so often used for test projects. 

There are 2 entities:

###### author 
- ID - int 
- FirstName - string
- LastName - string
- Bio - string 
- DateOfBirth - string

###### book
- ID - int
- Title - string
- AuthorID - int
- Year - int
- ISBN - string

###### List of endpoints:

- POST/books — Add a new book
- GET /books — Get all books
- GET /books/{id} — Get book by ID
- PUT /books/{id} — Update book by ID
- DELETE /books/{id} — Delete book by ID
- POST/authors — Add new author
- GET /authors — Get all authors
- GET /authors/{id} — Get author by ID
- PUT /authors/{id} — Update author by ID
- DELETE /authors/{id} — Delete author by ID
- PUT /books/{book_id}/authors/{author_id} — update author and book in transaction

## Installing

This application is packed as 2 docker containers, so, 
building and installation is straightforward
1. clone repository

`git clone git@github.com:am-silex/lets_go_library.git`
2. build images and compile app

`docker-compose -f .\compose.yaml up --no-start`
3. start new containers

`docker-compose -f .\compose.yaml start`

It's all. Now, REST API available at 8080 port.