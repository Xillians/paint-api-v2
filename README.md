# Paint API v2
![Coverage](https://img.shields.io/badge/Coverage-78.1%25-brightgreen)

## Table of Contents

- [Paint API v2](#paint-api-v2)
  - [Table of Contents](#table-of-contents)
  - [Quick start](#quick-start)
  - [Installation](#installation)
  - [Local development](#local-development)
    - [Testing](#testing)
    - [Linting](#linting)
  - [Database](#database)
  - [API documentation](#api-documentation)
  - [Deployment](#deployment)
  - [Production](#production)

## Quick start

This application uses a taskfile to run the setup. If you simply want to get it
to run, you need to set your environment variables (see .env.example for 
the list of variables). then you can do the following tasks:

```sh
task install
task run
```

## Installation

If you have never used task, you can look at [this guide](https://taskfile.dev/)
to find out more about that. When you have task available you can use the tasks
inside the taskfile, such as `install` to get the app to install.

## Local development

when you have installed and got the app running via the quickstart explanation,
you can start running with the database. You can follow the testing and linting
segment to verify the code you make is following guidelines.

### Testing

There are 3 tasks you can use in this project:
```sh
task test
task test:coverage
task test:coverage:html
```

regular test will check if all tests are fine. the two latter are to check what
the coverage overall is and get an html output to see where the missing coverage
is.

### Linting

This app uses `task lint` to check if the app is following linting guidelines.
Doing this often is recommended.

## Database

The database can can be found on turso.tech and this is the layout:

![Database layout](./docs/db%20layout.svg)

diagram drawn via this [url](https://dbdiagram.io/d/65e09303cd45b569fb380c37)

## API documentation

This application uses openAPI which you can either find on the 
[public url](https://paint-api-v2.fly.dev/docs) or localhost:8080/docs

## Deployment

This project uses automated deploy pipelines to push the dockerized app 
on fly.io. if you want to do manual deployments, you can do so with their 
[CLI](https://fly.io/docs/flyctl/).

## Production

This app can be found on [here](https://paint-api-v2.fly.dev). 