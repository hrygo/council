# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2024-12-14

### Added
- Multi-persona debate workflow (Affirmative, Negative, Adjudicator)
- Multi-provider support: DeepSeek, Google Gemini, Alibaba DashScope
- Provider selection via CLI flags (`--pro-provider`, `--con-provider`, `--judge-provider`)
- Model selection per role (`--pro-model`, `--con-model`, `--judge-model`)
- Google Gemini official SDK integration
- Streaming output support
- CLI with file, stdin, and interactive input modes
- Colored terminal output with progress indicators
- Context cancellation support (Ctrl+C graceful shutdown)
