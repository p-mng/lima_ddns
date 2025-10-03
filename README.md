# lima_ddns

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/p-mng/lima_ddns/go.yml) ![GitHub Tag](https://img.shields.io/github/v/tag/p-mng/lima_ddns)

A daemon to dynamically update DNS nameserver records at [lima-city.de](https://www.lima-city.de/) using the [official REST API](https://www.lima-city.de/docs/). Requires a valid API key and a paid domain. Can be run using Docker or as a standalone binary.

## Usage

1. Create an API key at <https://www.lima-city.de/usercp/api_keys>. The key should have all `dns.*` permissions.
2. Create a config file using the example config file (`config.sample.yml`).
3. Run the daemon using the provided `docker-compose.yml`.
