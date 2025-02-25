# Contributing to av-scanner

Thank you for considering contributing to av-scanner! We welcome contributions from the community and are excited to see what you can bring to the project.

## How to Contribute

### Reporting Issues

If you encounter any bugs or have feature requests, please open an issue on the [GitHub Issues](https://github.com/gooygeek/av-scanner/issues) page. Provide as much detail as possible to help us understand and address the issue.

### Forking the Repository

1. Fork the repository by clicking the "Fork" button on the top right of the repository page.
2. Clone your forked repository to your local machine:
    ```shell
    git clone https://github.com/yourusername/av-scanner.git
    cd av-scanner
    ```

### Creating a Branch

Create a new branch for your work to keep your changes organized and separate from the main branch:

```shell
git checkout -b feature/your-feature-name
```

### Making Changes

Make your changes to the codebase.

Ensure your code follows the project's coding standards and passes all tests.
If you've added new functionality, update the documentation accordingly.

### Formatting and Documentation

Before submitting your changes, it is appriciated to format the code and to generate documentation.

Formatting uses the standard `go fmt` tool, and documentation is generated using the `av-scanner` binary once it's been built.

This is as easy as running the following commands:

```shell
# Format the code
go fmt -C ./src

# Generate documentation
# Done after building the binary
./bin/av-scanner docs --format md
```

### Running Tests

Before submitting your changes, run the tests to ensure everything is working correctly:

```shell
go test -C ./src
```

### Committing Changes

Commit your changes with a clear and concise commit message:

```shell
git add .
git commit -m "Add feature: description of your feature"
```

### Pushing Changes

Push your changes to your forked repository:

```shell
git push origin feature/your-feature-name
```

### Creating a Pull Request

Go to the original repository on GitHub and click the "New Pull Request" button.
Select your branch from the "compare" dropdown.
Provide a clear and detailed description of your changes in the pull request.
Submit the pull request.

## Code of Conduct

This project does not have a formal code of conduct. However, we expect contributors to adhere to the general principles of respect, kindness, and professionalism when interacting with others in the community.

## Thank You!

Thank you for your interest in contributing to av-scanner! Your contributions are greatly appreciated.