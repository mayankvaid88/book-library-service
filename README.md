# Book Store

Application exposes Simple CRUD APIs to manage the books.

---

## Prerequisites

Make sure these tools are installed:

- go compiler
- go mockgen binary (for mocking)
- go swag for swagger documentations
- docker
- docker-compose

## Important commands

Run unit test cases
make unit_test

Run integration test cases
make int_test

Build docker image 
make build_img

App vet for static code analysis
make vet

To run the app inside docker along with database
make run

To generate the swagger documentations
make swag

To generate mocks used in unit test cases
make mocks

Database credentials has to be provided in the .env file which will be sourced while running the container. Replace the secrets with actual value in .env_sample and rename it to .env.

