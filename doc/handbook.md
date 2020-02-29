# Introduction
This document is the handbook for the privatebin command line
interface.

# CLI
## Configuration
Top level object:
| key                  | type       | default   | description                                         |
|----------------------|------------|-----------|-----------------------------------------------------|
| `open_discussion`    | bool       | false     | the default value of open discussion for a paste    |
| `burn_after_reading` | bool       | false     | the default value of burn after reading for a paste |
| `formatter`          | string     | plaintext | the default formatter for a paste                   |
| `expire`             | string     | 1day      | the default time to live for a paste                |
| `bin`                | array<bin> | n/a       | the list of bin instances                           |

Bin object:
| key                  | type   | default | description                                         |
|----------------------|--------|---------|-----------------------------------------------------|
| `name`               | string | n/a     | the name of the bin instance                        |
| `host`               | string | n/a     | the url of the bin instance                         |
| `auth`               | auth   | n/a     | the basic auth configuration of the bin instance    |
| `expire`             | string | n/a     | the default time to live for a paste                |
| `open_discussion`    | bool   | n/a     | the default value of open discussion for a paste    |
| `burn_after_reading` | bool   | n/a     | the default value of burn after reading for a paste |
| `formatter`          | string | n/a     | the formatter for the paste                         |

Auth object:
| key        | type   | default | description             |
|------------|--------|---------|-------------------------|
| `username` | string | n/a     | the basic auth username |
| `password` | string | n/a     | the basic auth password |

Example:
```json
{
	"bin": [
		{
			"name": "demo",
			"host": "demo.example.com",
			"auth": {
				"username": "john.doe",
				"password": "pa$$w0rd"
			},
			"formatter": "markdown",
			"burn_after_reading": false
		}
	],
	"burn_after_reading": true
}
```
