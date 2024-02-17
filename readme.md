# zqz.ca
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/ff2f2665c8864a0aa3b7334088eef646)](https://www.codacy.com/app/fish/upl?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=zqz/upl&amp;utm_campaign=Badge_Grade)
[![codebeat badge](https://codebeat.co/badges/c2cc652d-7230-4eea-a642-976064865d2d)](https://codebeat.co/projects/github-com-zqz-upl-master)
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
