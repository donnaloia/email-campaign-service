
# Email Campaign Service

A modern standalone email campaign service written in Go.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

![example workflow](https://github.com/donnaloia/email-campaign-service/actions/workflows/docker-build-push.yml/badge.svg)


## Features

- REST API endpoints to manage email campaigns (see endpoints below)
- manage email templates, email groups, email reciepents, and campaigns
- designed to plug directly into a distributed or integrated system
- built with performance and scale in mind
- see it deployed in a distributed system here: [SendPulse](https://github.com/donnaloia/sendpulse)


## Tech Stack

**Language:** Go

**Framework:** Echo

**DB:** Postgres

**Deployment:** Docker





## Run Locally
docker-compose spins up the email-campaign-service, along with a postgres instance.
To deploy this project locally run:

```bash
  docker-compose build
  docker-compose up
```


## REST API Reference


#### There is no authentication required for this service, partly because it is designed to sit behind an api gateway and auth will therefor be performed at the gateway level.


#### Get Email Groups

```http
  GET /api/v1/email-groups/<id>
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`    | `string` | **Required**. Uuid of the email group to fetch |

```json
{
    "name": "my first email group",
    "created_at": "2025-01-01 12:00:00",
    "email_addresses": ["id1@example.com", "id2@example.com", "id3@example.com"]
}
```


#### Create User Permission

```http
  POST /api/v1/email-groups/
```


| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`| `json` | **Required**. Name of the email group we are creating|
| `email_addresses`| `json` | **Required**. List of email addresses to add to the group|

```json
{
   "name": "my first email group",
   "email_addresses": ["id1@example.com", "id2@example.com", "id3@example.com"]
}
```


#### Update Email Group

```http
  PUT /api/v1/email-groups/<id>
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Uuid of the email group we are updating|
| `name`| `json` | **Required**. Name of the email group we are updating|
| `email_addresses`| `json` | **Required**. List of email addresses to add to the group|

```json
{
   "name": "my first email group",
   "email_addresses": ["id1@example.com", "id2@example.com", "id3@example.com"]
}
```


#### Get Email Templates

```http
  GET /api/v1/email-templates/<id>
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`    | `string` | **Required**. Uuid of the email template to fetch |

```json
{
    "name": "my first email template",
    "created_at": "2025-01-01 12:00:00",
    "html": "<p>my first email template</p>",
}
```


#### Create Email Template

```http
  POST /api/v1/email-templates/
```


| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`| `json` | **Required**. Name of the email template we are creating|
| `html`| `json` | **Required**. Html of the email template we are creating|

```json
{
   "name": "my first email template",
   "html": "<p>my first email template</p>"
}
```


#### Update Email Template

```http
  PUT /api/v1/email-templates/<id>
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Uuid of the email template we are updating|
| `name`| `json` | **Required**. Name of the email template we are updating|
| `html`| `json` | **Required**. Html of the email template we are updating|

```json
{
   "name": "my first email template",
   "html": "<p>my first email template</p>"
}
```


## Todo

- add more test coverage
- CLI admin tool
