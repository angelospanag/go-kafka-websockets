# go-kafka-websockets

Experimenting with a full flow of sending messages from an Apache Kafka consumer, to a Go websockets endpoint, to a UI using a JavaScript websockets connection.

## Getting Started

[mise](https://mise.jdx.dev/) manages the pinned toolchain (Go 1.26, golangci-lint).

```bash
# macOS / Linux
curl https://mise.run | sh

# Windows
winget install jdx.mise
```

Activate mise in your shell (`~/.zshrc`):

```zsh
eval "$(mise activate zsh)"
```

Then, in the repo:

```bash
mise trust    # one-time
mise install  # downloads Go and golangci-lint
```

Create a `.env` file at the root of the project:

```dotenv
TOPICS=quickstart-events
BOOTSTRAP_SERVERS=localhost:9092
GROUP_ID=my-group
AUTO_OFFSET_RESET=latest
```

## Running

Install and start Kafka (macOS):

```bash
brew install kafka
brew services start kafka
```

Start a Kafka producer to send messages:

```bash
kafka-console-producer --topic quickstart-events --bootstrap-server localhost:9092
```

Start the server:

```bash
mise run dev
```

Visit http://localhost:8000 and send some text from the Kafka console producer. The text will appear on your screen after being picked up by the Kafka consumer in the backend and forwarded through the WebSocket connection.

## Development

| Command           | Description                                 |
|-------------------|---------------------------------------------|
| `mise run dev`    | Run without building a binary               |
| `mise run build`  | Build the Go binary                         |
| `mise run test`   | Run tests                                   |
| `mise run fmt`    | Format code via `golangci-lint fmt`         |
| `mise run lint`   | Lint via `golangci-lint run`                |
| `mise run vuln`   | Scan dependencies for known vulnerabilities |
| `mise run deps`   | Update and tidy dependencies                |
| `mise run clean`  | Remove build artifacts                      |
