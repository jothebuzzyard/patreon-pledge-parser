# Patreon CSV parser

## Overview
This is a Go project that can take an exported patreon CSV file and parse it into sepperate tier files. It automatically filters out non-paying patreons or those that have expires.

## For users
1. Download the latest version from the release space to the right.
2. Put the file in the same folder as your CSV file.
3. Either name your csv file pledges.csv, or input the path when the applicaiton asks
4. Run the app. Windows will complain, because I am not paying them to sign it. Just click run anyway. If you want to feel more secure you can check the code and compile this repo yourself.
5. Once done a folder "output" will be created. In there you will find multiple txt files. One per tier.

### Configuration: `settings.conf`

This exporter can be customized using a `settings.conf` file placed in the same directory as the executable. If no `settings.conf` is found, default values are used.

#### How to Create `settings.conf` on Windows

1. Open Notepad.
2. Paste the example configuration below.
3. Save the file as `settings.conf` in the same folder as the exporter executable.

#### Settings Reference

| Setting                | Values / Format                  | Description                                                                                 |
|------------------------|----------------------------------|---------------------------------------------------------------------------------------------|
| EXPORT_SVG             | `true` or `false`                | Enable or disable SVG export.                                                               |
| EXPORT_TXT             | `true` or `false`                | Enable or disable TXT export.                                                               |
| OUTPUT_DIR             | Directory name                   | Output folder for generated files.                                                          |
| DEFAULT_CSV_FILE       | Filename                         | Default CSV file to process (if not found user will be prompted for the filename).          |
| SVG_WIDTH              | Whole number                     | Width of the SVG output in pixels.                                                          |
| SVG_MARGIN_TO_EDGE     | Whole number                     | Margin from edge of SVG in pixels.                                                          |
| SVG_COLUMN_GAP         | Whole number                     | Gap between columns in SVG in pixels.                                                       |
| SVG_FONTSIZE           | Whole number                     | Font size in SVG output (pixels).                                                           |
| SVG_LINEHEIGHT         | Whole number                     | Line height in SVG output (pixels).                                                         |
| SVG_COLUMNS            | Whole number                     | Number of columns in SVG output.                                                            |
| SVG_FONTFAMILY         | Font family string               | Font family for SVG text (e.g., `Arial, sans-serif`).                                       |
| SVG_COLUMN_COLORS      | Comma-separated color hex values | Colors for each column (e.g., `#3aff22,#c622ff,#a8ff21`).                                   |
| SVG_RANDOMIZE_COLORS   | `true` or `false`                | Randomize name colors for SVG output.                                                       |

#### Example `settings.conf` (Default Values)

```
# settings.conf - Example configuration

EXPORT_SVG=true
EXPORT_TXT=true
OUTPUT_DIR=output
DEFAULT_CSV_FILE=pledges.csv

SVG_WIDTH=1161
SVG_MARGIN_TO_EDGE=26
SVG_COLUMN_GAP=54
SVG_FONTSIZE=16
SVG_LINEHEIGHT=20
SVG_COLUMNS=3
SVG_FONTFAMILY=Trebuchet MS, Arial, sans-serif
SVG_COLUMN_COLORS=#3aff22,#c622ff,#a8ff21
SVG_RANDOMIZE_COLORS=true
```

Copy and edit this file as needed to customize the exporter's behavior.


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
