# zqz.ca
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/a75e339ad8a045949fd2d032103f71cd)](https://app.codacy.com/gh/zqz/upl/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![codebeat badge](https://codebeat.co/badges/d0afb6f1-2490-4ec1-b4a4-213f88800fe4)](https://codebeat.co/projects/github-com-zqz-upl-master)
[![Maintainability](https://api.codeclimate.com/v1/badges/74bcc076dbf4d07c141d/maintainability)](https://codeclimate.com/github/zqz/upl/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/zqz/upl)](https://goreportcard.com/report/github.com/zqz/upl)
![Deploy to Production](https://github.com/zqz/upl/workflows/Deploy%20to%20Production/badge.svg)

### Development

#### Configuration
The application is configured via environment variables. These can be specified
in a local file named `.env` in the root of the working directory.

ENV|Example|Desc
-|-|-
`DATABASE_URL`|`postgres://user:password@host:5432/database?args`|Postgres connection URL
`PORT`|`3000`|The port the webserver runs on
`FILES_PATH`|`/var/opt/files`|Path to store uploaded files

#### Running the app
```
go run cmd/backend/main.go
```

#### Database Schema
If the database schema changes you will need to run `sqlboiler psql` to regenerate the models package.

#### Development Tools

Install Go (1.22+) and then install package dependencies
```
go mod download
```

Install [SQLBoiler](https://github.com/volatiletech/sqlboiler?tab=readme-ov-file#download) and the PSQL driver
```
go install github.com/volatiletech/sqlboiler/v4@latest
go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
```

