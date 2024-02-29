# zqz docs

## Intro

## Prepare File

Tell the API about a file you wish to upload. 

Field | Description
----- | -----------
alias | the public name of the file.
type | the content type (if known).
name | the original filename.
size | size in bytes
hash | the sha1 hex digest of the file.

**POST** /prepare

#### Request:
```json
  {
    "alias": "My foo bar file",
    "type": "content/whatever",
    "name": "foobar.jpg",
    "size": 123,
    "hash": "275f6032fb106b0fefa7aab76186a034107f5fd8"
  }
```

#### Responses:
```json
  {
    "bytes_received": 0,
    "hash": "275f6032fb106b0fefa7aab76186a034107f5fd8"
  }
```

## Upload Status

Check if a file with some hash is already uploaded.

**GET** /upload/{hash}

#### Responses:

When the file exists already
```json
  {
    "alias": "Example File",
    "name": "example.jpg",
    "bytes_received": 100,
    "size": 100,
    "hash": "275f6032fb106b0fefa7aab76186a034107f5fd8"
  }
```

When the file does not exist
```json
  {
    "message": "no file with specified hash exists"
  }
```
