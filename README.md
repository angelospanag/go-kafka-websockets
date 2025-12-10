# go-kafka-websockets

Experimenting with a full flow of sending messages from an Apache Kafka consumer, to a Go websockets endpoint, to a UI
using a JavaScript websockets connection.

<!-- TOC -->
* [go-kafka-websockets](#go-kafka-websockets)
  * [Prerequisites](#prerequisites)
    * [macOS using `brew`](#macos-using-brew)
    * [Create a `.env` file at the root of the project](#create-a-env-file-at-the-root-of-the-project)
  * [Running](#running)
    * [Start Apache Kafka (using MacOS and `brew`)](#start-apache-kafka-using-macos-and-brew)
    * [Start Kafka producer](#start-kafka-producer)
    * [Start Kafka consumer](#start-kafka-consumer)
    * [Start server](#start-server)
      * [Development server](#development-server)
      * [Production server](#production-server)
  * [Available Tasks](#available-tasks)
    * [🚀 Run the application in dev mode](#-run-the-application-in-dev-mode)
    * [🎨 Format the code](#-format-the-code)
    * [🔍 Lint the codebase](#-lint-the-codebase)
    * [⬆️ Update and tidy dependencies](#-update-and-tidy-dependencies)
    * [🔧 Build the Go binary](#-build-the-go-binary)
    * [🚀 Build and run the application](#-build-and-run-the-application)
    * [🧹 Clean build artifacts](#-clean-build-artifacts)
<!-- TOC -->

## Prerequisites

- [Go 1.25](https://go.dev)
- [golangci-lint](https://golangci-lint.run/) for linting
- [Task](https://taskfile.dev/) for running tasks

### macOS using `brew`

```bash
brew install go@1.25 golangci-lint go-task
```

### Create a `.env` file at the root of the project

```dotenv
TOPICS=quickstart-events
BOOTSTRAP_SERVERS=localhost:9092
GROUP_ID=my-group
AUTO_OFFSET_RESET=latest
```

## Running

### Start Apache Kafka (using MacOS and `brew`)

```bash
brew services start kafka
```

### Start Kafka producer

```bash
kafka-console-producer --topic quickstart-events --bootstrap-server localhost:9092
```

### Start Kafka consumer

```bash
kafka-console-consumer --topic quickstart-events --from-beginning --bootstrap-server localhost:9092
```

### Start server

#### Development server

```bash
task dev
```

#### Production server

```bash
task run
```

Visit http://localhost:8000 and send some text from the Kafka console producer. The text will appear on your screen,
after being picked up by the Kafka consumer of the backend, and sent through a websockets connection to the UI.

## Available Tasks

Below are the available tasks defined in [`Taskfile.yml`](./Taskfile.yml):

### 🚀 Run the application in dev mode

```sh
task dev
```

### 🎨 Format the code

```sh
task fmt
```

### 🔍 Lint the codebase

```sh
task lint
```

### ⬆️ Update and tidy dependencies

```sh
task update-deps
```

### 🔧 Build the Go binary

```sh
task build
```

### 🚀 Build and run the application

```sh
task run
```

### 🧹 Clean build artifacts

```sh
task clean
```
