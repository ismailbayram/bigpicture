# BigPicture: Validate Architecture
**Do not tell the rules, define them.**

BigPicture is a tool to validate the architecture of project. 
It can be used in Continuous Integration (CI) pipelines to validate the architecture of the project
like in the `.github/workflows/codequality.yml`.


# Supported Languages
- Go
- Python
- Java (Under Development)
- C# (Under Development)
- JS (Under Development)

# Installation
## Install with Go
```bash
go install github.com/ismailbayram/bigpicture@1.0.0
```

# Usage
## Server
Runs a server on port 44525. Architecture of the project can be seen on the browser.
```bash
bigpicture server
```

## Validate
Validates the architecture of the project according to the rules defined in the `.bigpicture.json` file.
```bash
bigpicture validate
```

## .bigpicture.json File
`.bigpicture.json` file is used to define the rules. It should be in the root directory of the project.
```json
{
    "port": 44525,
    "ignore": [
        "web"
    ],
    "validators": [
        ...
    ]
}
```
**port**: Port number of the server. Default value is 44525.

**ignore**: Directories to ignore. Default value is empty. For instance in this project `web` directory includes
non-go files, thus it should be ignored.

**validators**: List of validators. Default value is empty. See [Validators](#validators) section for more information.

## Validators
### NoImportValidator
Checks if the package imports the given package. **It can be used in layered architectures.**

**Example 1**:
For instance, in this project, `/internal/config` package can not import any other package. 
```json
{
    "type": "no_import",
    "args": {
        "from": "/internal/config",
        "to": "*"
    }
}
```
**Example 2**:
For instance, in this project, `/internal/validator` package can not import any package in the `/internal/browser` package. 
```json
{
    "type": "no_import",
    "args": {
        "from": "/internal/validator",
        "to": "/internal/browser"
    }
}
```

### LineCountValidator
Checks if the package has files which have more than the given number of lines.

**Example**:
For instance, in this project, `/internal/browser` package can not have files which have more than 100 lines. 
```json
{
    "type": "line_count",
    "args": {
        "module": "/internal/browser",
        "max": 100,
        "ignore": ["*_test.go", "test/*.go"]
    }
}
```

### FunctionValidator
Checks if the package has functions which have more than the given number of lines.

**Example**:
For instance, in this project, `/internal/browser` package can not have functions which have more than 10 lines. 
```json
{
    "type": "function",
    "args": {
        "module": "/internal",
        "max_line_count": 50,
        "ignore": ["*_test.go", "test/*.go"]
    }
}
```

### InstabilityValidator
Checks if the instability metric of a package according to its directory is more than the given number.

**Instability Calculation**:

Package A is imported by 3 packages and it imports 2 packages. Instability metric of the package A is
`2 / (2 + 3) = 0.4`.

**Example**:
For instance, in this project, `/internal/graph` package can not have instability metric more than 0.5. 
```json
{
    "type": "instability",
    "args": {
        "module": "/internal/graph",
        "max": 0.5
    }
}
```



# Contribution
There are many ways in which you can participate in this project, for example:

- Implement a new validator in `/internal/validator` directory.
- Implement a new language support in `/internal/browser` directory.
