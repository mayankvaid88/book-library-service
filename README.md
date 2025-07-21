# Book Store 📚

This application provides simple **CRUD (Create, Read, Update, Delete)** APIs for managing books.

---

## Prerequisites

Before you get started, ensure you have the following tools installed on your system:

* **Go Compiler**: Essential for building and running the Go application.
* **Go Mockgen**: Used for generating mock interfaces for testing purposes.
* **Go Swag**: Utilized to generate **Swagger API documentation**.
* **Docker**: For containerizing the application and its dependencies.
* **Docker Compose**: To orchestrate and run multi-container Docker applications, like this one with its database.

---

## Important Commands

Here are some crucial `make` commands to help you manage and interact with the application:

* **`make unit_test`**: Runs all the **unit test cases** for the application.
* **`make int_test`**: Executes the **integration test cases**.
* **`make build_img`**: Builds the **Docker image** for the application.
* **`make vet`**: Performs **static code analysis** (vetting) on the application's source code.
* **`make run`**: Starts the application inside **Docker along with its database**.
* **`make swag`**: Generates the **Swagger API documentation**.
* **`make mocks`**: Generates the **mock files** used in unit test cases.
* **`make coverage`**: Evaluates the **coverage** of unit test cases.

---

## Database Configuration

**Database credentials** must be provided in a **`.env`** file under infra folder. This file will be sourced when the application container is run.

**Important**: Before containerizing the application, **rename `.env_sample` to `.env`** and replace the placeholder secrets with your actual database credentials.