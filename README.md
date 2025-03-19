# AEM GoSync

AEM GoSync is a Go application that monitors specified directories for changes, creates an AEM content package, and deploys it to an AEM instance. This project is designed to streamline the process of managing content in Adobe Experience Manager (AEM).

> **Note:** This project is currently under development. Some features may not work as expected, and others may not be implemented yet.

## Features

- Monitors directory paths for changes.
- Creates AEM content packages.
- Deploys packages to AEM instances.

## Prerequisites

- Go 1.16 or higher
- Access to an AEM instance

## Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd aem-gosync
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Usage

To run the application, use the following command:

```
go run cmd/main.go <directory-path-1> <directory-path-2> ...
```

Replace `<directory-path-1>`, `<directory-path-2>`, etc., with the paths of the directories you want to monitor.

## Example

```
go run cmd/main.go /path/to/directory1 /path/to/directory2
```

This command will start monitoring the specified directories for any changes.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.