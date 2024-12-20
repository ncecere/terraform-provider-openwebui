# Contributing to OpenWebUI Terraform Provider

Thank you for your interest in contributing to the OpenWebUI Terraform Provider! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How to Contribute

### Reporting Issues

Before submitting an issue:
1. Check the [existing issues](https://github.com/ncecere/terraform-provider-openwebui/issues) to avoid duplicates
2. Use the provided issue template
3. Include as much relevant information as possible:
   - Provider version
   - Terraform version
   - OpenWebUI version
   - Complete error messages
   - Minimal reproduction steps

### Submitting Changes

1. **Fork the Repository**
   - Create your own fork of the repository
   - Clone your fork locally

2. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```
   
   Branch naming conventions:
   - `feature/` for new features
   - `fix/` for bug fixes
   - `docs/` for documentation changes
   - `refactor/` for code refactoring

3. **Make Your Changes**
   - Follow the [Development Guide](DEVELOPMENT.md)
   - Adhere to the existing code style
   - Add tests for new functionality
   - Update documentation as needed

4. **Commit Your Changes**
   - Use clear, descriptive commit messages
   - Follow conventional commits format:
     ```
     type(scope): description
     
     [optional body]
     [optional footer]
     ```
   - Types: feat, fix, docs, style, refactor, test, chore
   - Example:
     ```
     feat(knowledge): add support for private knowledge bases
     
     - Implement private access control
     - Add documentation for private knowledge bases
     - Add tests for access control
     
     Closes #123
     ```

5. **Submit a Pull Request**
   - Push your changes to your fork
   - Create a pull request against the main repository
   - Use the provided pull request template
   - Link relevant issues
   - Provide a clear description of your changes

### Pull Request Process

1. **Initial Submission**
   - Ensure all tests pass
   - Update relevant documentation
   - Add your changes to CHANGELOG.md
   - Fill out the pull request template completely

2. **Review Process**
   - Maintainers will review your code
   - Address any feedback or requested changes
   - Keep the pull request updated

3. **Acceptance**
   - All checks must pass
   - Required reviews must be approved
   - Changes must be up to date with the base branch

## Development Guidelines

### Code Style

- Follow Go best practices and conventions
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and concise

### Testing

- Add unit tests for new functionality
- Add acceptance tests for new resources
- Ensure all tests pass locally before submitting
- Include test cases for edge cases and error conditions

### Documentation

- Update relevant documentation files
- Add examples for new features
- Follow the existing documentation style
- Keep examples clear and concise

## Release Process

1. **Version Updates**
   - Update version numbers according to [Semantic Versioning](https://semver.org/)
   - Update CHANGELOG.md with all changes

2. **Documentation**
   - Ensure all documentation is up to date
   - Update examples if needed
   - Review API references

3. **Testing**
   - Run all test suites
   - Perform manual testing of new features
   - Verify examples work as expected

## Getting Help

- Check the [Development Guide](DEVELOPMENT.md)
- Review existing issues and pull requests
- Join community discussions
- Contact maintainers for clarification

## Recognition

Contributors will be recognized in:
- CHANGELOG.md for their contributions
- GitHub repository insights
- Release notes when applicable

## Additional Notes

### Security Issues

For security-related issues:
1. Do NOT open a public issue
2. Contact the maintainers directly
3. Follow responsible disclosure practices

### Legal

By contributing to this project, you agree to license your contributions under the same license as the project.

## Questions?

If you have questions about contributing:
1. Review the documentation
2. Check existing issues
3. Open a new issue with the question label

Thank you for contributing to the OpenWebUI Terraform Provider!
