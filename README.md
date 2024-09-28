# Random Generation API

A simple, public API to generate random words or numbers. This service is free
to use and available as open-source software. Do you need a random word or a
UUID in a script? Can't pull in dependencies easily? Call this API to grab one
quickly.

## Features

The API provides endpoints for generating:

- Random integers
- Random floating-point numbers
- Dice rolls with various formats
- ULIDs (Universally Unique Lexicographically Sortable Identifier)
- UUIDs (versions 4 and 7)
- Nano IDs
- Random words from various categories like animals or vegetables

## Usage

The API is accessible at `https://rnd.bgenc.dev`. You can make requests to the various endpoints, each with its own set of optional parameters.

Example:

```
https://rnd.bgenc.dev/v1/int?min=1&max=100
```

For detailed information on each endpoint and its parameters, please refer to
the API documentation.

TODO: API documentation site in progress, please hold!

### Rate Limits

There are rate limits in place to prevent abuse. If you need higher rate limits or uptime guarantees, please contact rnd@bgenc.dev to discuss options.

## Local Development

To run the project locally:

1. Clone the repository
2. Install dependencies (Go 1.22 or later and npm are required)
3. Run the following commands:

```bash
npm install
npm run dev
```

This will start both the UI and the backend. The UI will be available at `localhost:5173`. The backend will restart if there are any changes, and the frontend will hot reload.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

### Contributing

Issue reports, and pull requests are always welcome! If you encounter any bugs,
or if there are other generators you'd like to see here, let us know. If you'd
like to contribute, reach out and we can help you get started.

All contributions require a Contributor License Agreement.

## Contact

For inquiries, please contact rnd@bgenc.dev. Increased rate limits, uptime
guarantees, and licensing options are available.

## Disclaimer

While we do our best to keep the service running, there are no guarantees of uptime or support. This service is provided as-is for small projects and scripts.
