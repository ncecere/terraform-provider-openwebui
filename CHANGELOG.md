# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-12-20

### Added
- Initial release of the OpenWebUI Terraform Provider
- User management capabilities:
  - User data source for querying user information
  - Support for user lookup by ID, email, or name
  - Access to user metadata and settings

- Group management:
  - Group resource for creating and managing user groups
  - Group data source for querying existing groups
  - Comprehensive permission system for workspace and chat features
  - User-group membership management

- Knowledge base management:
  - Knowledge base resource for creating and configuring knowledge bases
  - Knowledge base data source for querying existing knowledge bases
  - Support for public and private access controls
  - Custom metadata and tagging capabilities

- Model management:
  - Model resource for deploying and configuring AI models
  - Model data source for querying existing models
  - Extensive model parameter configuration options
  - Access control and capability management
  - Support for model metadata and tags

### Documentation
- Comprehensive provider documentation
- Detailed resource and data source documentation
- Example configurations for all resources
- Development and contribution guidelines
- Quick start guide and usage examples

### Infrastructure
- Terraform Plugin Framework implementation
- Modular client architecture
- Robust error handling
- Comprehensive test suite
- CI/CD pipeline for automated testing and releases

[1.0.0]: https://github.com/ncecere/terraform-provider-openwebui/releases/tag/v1.0.0
