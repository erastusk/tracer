# Tracer 🔍

A comprehensive connectivity and messaging diagnosis tool for distributed systems. Tracer provides CLI-based testing capabilities for Kafka clusters (MSK, Confluent), Redis instances, and other messaging platforms with support for various authentication methods.

## Features

### 🚀 Messaging Platforms
- **Amazon MSK**: SASL/SCRAM authentication with TLS encryption
- **Confluent Cloud**: OAuth 2.0 authentication with token management
- **Redis**: Standard and TLS connections with optional SSL verification

### 🔐 Authentication Methods
- SASL/SCRAM-SHA-512 for MSK clusters
- OAuth 2.0 client credentials flow for Confluent Cloud
- AWS Secrets Manager integration for credential management
- Interactive CLI prompts with secure password input

### 🛠️ Capabilities
- **Connectivity Testing**: TCP dial testing with configurable timeouts
- **Message Production**: Send test messages with automatic partitioning
- **Message Consumption**: Real-time message consumption with consumer groups
- **Topic Management**: List and validate Kafka topics
- **Configuration Caching**: Persistent storage of connection parameters

## Installation

### From Source
```bash
git clone https://github.com/erastusk/tracer.git
cd tracer
go mod download
go build -o tracer
```

### Using Make
```bash
make build
```

### Docker
```bash
docker build -t tracer .
docker run -it tracer
```

## Usage

### Basic Commands

```bash
# Start interactive CLI
./tracer

# Test MSK connectivity
./tracer msk

# Test Confluent Cloud connectivity
./tracer confluent

# Test Redis connectivity  
./tracer redis
```

### Command Line Flags

```bash
# Specify topic for Kafka operations
./tracer msk --topic=my-topic

# Use cached credentials
./tracer confluent --use-cache

# Override default configuration
./tracer redis --endpoint=redis.example.com:6379
```

### Interactive Mode

The CLI provides interactive prompts for:
- Connection endpoints
- Authentication credentials
- Topic selection
- Operation type (produce/consume)

## Configuration

### Environment Variables

```bash
export AWS_REGION=us-east-1
export AWS_PROFILE=default
```

### AWS Secrets Manager

Store credentials in AWS Secrets Manager with the following structure:

```json
{
  "username": "your-username",
  "password": "your-password",
  "client_id": "oauth-client-id",
  "client_secret": "oauth-client-secret"
}
```

### Cache Directory

Configuration is cached in `~/.tracer/` for reuse between sessions.

## Architecture

```
tracer/
├── cmd/                    # Cobra CLI commands
├── internal/
│   ├── cache/             # Configuration persistence
│   ├── confluent/         # Confluent Cloud client
│   ├── msk/              # Amazon MSK client
│   ├── oauth/            # OAuth 2.0 token provider
│   ├── prompts/          # Interactive CLI prompts
│   ├── redis/            # Redis connectivity
│   ├── secrets/          # AWS Secrets Manager
│   ├── types/            # Common interfaces and types
│   └── utils/            # Utilities and helpers
├── mocks/                # Generated mocks for testing
└── certs/               # TLS certificates
```

## Development

### Prerequisites
- Go 1.21+
- AWS CLI configured
- Docker (optional)

### Building
```bash
# Build binary
make build

# Run tests
make test

# Generate mocks
make gen

# Clean build artifacts
make clean
```

### Testing
```bash
# Run all tests with coverage
go test -v -cover ./...

# Test specific package
go test -v ./internal/msk
```

### Code Generation
```bash
# Generate mocks for interfaces
go generate ./...
```

## Docker Support

The Docker image includes comprehensive networking tools for debugging:

- `tcpdump`, `tshark` for packet analysis
- `nmap`, `nping` for network scanning
- `curl`, `httpie` for HTTP testing
- `jq` for JSON processing
- `iperf3` for bandwidth testing

```bash
# Run in container
docker run -it tracer bash

# Execute tracer in container
docker run -it tracer ./tracer
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style
- Follow Go conventions and best practices
- Add unit tests for new functionality
- Update documentation for API changes
- Use meaningful commit messages

## Examples

### MSK Connection Test
```bash
$ ./tracer msk
? Please enter endpoint: msk-cluster.kafka.us-east-1.amazonaws.com:9098
? Please enter username: myuser
? Please enter password: ********
✅ Successfully connected to MSK cluster
✅ Topic validation passed
📨 Test message sent successfully
```

### Confluent Cloud OAuth
```bash
$ ./tracer confluent
? Please enter endpoint: pkc-xxxxx.us-east-1.aws.confluent.cloud:9092
🔐 Authenticating with OAuth 2.0...
✅ OAuth token obtained successfully
✅ Connected to Confluent Cloud
📋 Available topics: [events, logs, metrics]
```

## Troubleshooting

### Common Issues

**Connection Timeouts**
- Verify network connectivity and firewall rules
- Check endpoint URL format and port
- Ensure security groups allow traffic

**Authentication Failures**
- Verify credentials in AWS Secrets Manager
- Check OAuth client configuration
- Validate username/password for SASL

**TLS Certificate Issues**
- Ensure proper CA certificates are embedded
- Verify TLS version compatibility (minimum TLS 1.2)

### Debug Mode
```bash
# Enable verbose logging
./tracer --debug msk

# Check connectivity only
./tracer --test-connection-only redis
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Sarama](https://github.com/IBM/sarama) for Kafka client
- [Survey](https://github.com/AlecAivazis/survey) for interactive prompts
- [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2) for AWS integration

---

**Built with ❤️ for distributed systems debugging**