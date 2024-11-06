# Minerva

Minerva is an open-source scanning tool designed to help developers identify security vulnerabilities, potential immprovements, and valuable insights within their code and applications. Named after the Roman goddess of wisdom, Minerva is designed to help developers quickly identify and address security vulnerabilities and other issues within their codebase.

## Key Features

- **Security Scanning**: Minerva scans your web applications for security vulnerabilities, including scanning scripts, stylesheets, and other files for potential security issues.
- **Minimalist Interface**: Since Minerva is designed to be a lightweight tool, it only provides a API for scanning and retrieving results. This allows developers to integrate Minerva into their existing workflows and tools.
- **Information Gathering**: Minerva can also be used to gather information about your codebase, including identifying network, dependencies, and other information.

## Prerequisites

- Go 1.23 or higher
- Postgres 16 or higher
- Redis 6 or higher
- S3 Bucket

All of these things are required to run Minerva. If you have Docker installed you can use our Docker Compose file to run Minerva.

## Installation

> [!IMPORTANT]
> Minerva is currently in development and only supports downloading the source code from GitHub. We are working on adding support for package managers and other installation methods.

To install Minerva, you can download the source code from GitHub and run the following command:

```bash
git clone https://github.com/codevault-llc/minerva.git

cd minerva

# Setup a .env file from the .env.example template.
cp .env.example .env

go mod download

go build -o minerva

./minerva
```

## Contributing to Minerva

Minerva encourages contributions from developers worldwide. Contributions can include:

- Adding new scanning modules or expanding existing ones.
- Reporting and resolving bugs.
- Writing or enhancing documentation.
- Proposing new features or optimizations.

## License

> [!NOTE]
> Minerva's License may change soon. Please check the [LICENSE](LICENSE) file for the most up-to-date information.

Minerva is licensed under the GPL-3.0 License. For more information, see the [LICENSE](LICENSE) file.
