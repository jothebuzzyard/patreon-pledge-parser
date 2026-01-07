# Patreon CSV parser

## Overview
This is a Go project that can take an exported patreon CSV file and parse it into sepperate tier files. It automatically filters out non-paying patreons or those that have expires.

## For users
1. Download the latest version from the release space to the right.
2. Put the file in the same folder as your CSV file.
3. Either name your csv file pledges.csv, or input the path when the applicaiton asks
4. Run the app. Windows will complain, because I am not paying them to sign it. Just click run anyway. If you want to feel more secure you can check the code and compile this repo yourself.
5. Once done a folder "output" will be created. In there you will find multiple txt files. One per tier.

## Development Environment
This project uses a development container to provide a consistent development environment. The container is configured using the files located in the `.devcontainer` directory.

## Setup Instructions
1. **Clone the repository**:
   ```
   git clone <repository-url>
   cd go-project
   ```

2. **Open in a development environment**:
   Open the project in your preferred development environment that supports development containers.

3. **Build the container**:
   The development container will automatically build using the `Dockerfile` specified in `.devcontainer/devcontainer.json`.

## Usage
- The main application logic can be found in `src/main.go`.
- To run the application, use the following command inside the container:
  ```
  go run src/main.go
  ```

## Dependencies
Dependencies are managed using Go modules. The module and its dependencies are defined in `go.mod` and `go.sum`.

## Building a release
```
GOOS=windows GOARCH=amd64 go build -o patreon-pledge-parser.exe ./src/main.go
```

## Contributing
Feel free to contribute to this project by submitting issues or pull requests. Please ensure that your code adheres to the project's coding standards.

## License
This project is licensed under the MIT License. See the LICENSE file for more details.
