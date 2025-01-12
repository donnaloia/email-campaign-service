# Email Campaign Service

A modern standalone email campaign service written in Go.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

![example workflow](https://github.com/donnaloia/email-campaign-service/actions/workflows/docker-build-push.yml/badge.svg)


## Features

- REST API endpoints to manage email campaigns (see endpoints below)
- manage email templates, email groups, email reciepents, and campaigns
- designed to plug directly into a distributed or integrated system
- built with performance and scale in mind
- see it deployed in an event-driven, distributed system here: [SendPulse](https://github.com/donnaloia/sendpulse)


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


#### Get Email Address

```http
  GET /api/v1/email-addresses/<id>
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`    | `string` | **Required**. Uuid of the email address to fetch |

```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "organization_id": "123e4567-e89b-12d3-a456-426614174000",
    "created_at": "2025-01-01 12:00:00"
}
```

#### Create Email Address

```http
  POST /api/v1/email-addresses/
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `email`| `string` | **Required**. Email address to create|

```json
{
   "email": "user@example.com"
}
```

#### Update Email Address

```http
  PUT /api/v1/email-addresses/<id>
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Uuid of the email address to update|
| `email`   | `string` | **Required**. New email address|

```json
{
   "email": "newuser@example.com"
}
```

#### Get Email Group

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


#### Create Email Group

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


#### Get Email Template

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
| `name`| `json` | **Required**. Name of the email template we are updating|
| `html`| `json` | **Required**. Html of the email template we are updating|

```json
{
   "name": "my first email template",
   "html": "<p>my first email template</p>"
}
```

#### Get Email Campaign

```http
  GET /api/v1/email-campaigns/<id>
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`    | `string` | **Required**. Uuid of the email group to fetch |

```json
{
    "name": "my first email campaign",
    "created_at": "2025-01-01 12:00:00",
    "email_groups": ["my first email group","my second email group"],
    "email_templates": "my first email template"
}
```


#### Create Email Campaign

```http
  POST /api/v1/email-campaigns/
```


| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`| `json` | **Required**. Name of the email campaign we are creating|
| `email_groups`| `json` | **Required**. List of email groups to add to the campaign|
| `email_templates`| `json` | **Required**. List of email templates to add to the campaign|

```json
{
   "name": "my first email campaign",
   "email_groups": ["my first email group","my second email group"],
   "email_templates": "my first email template"
}
```


#### Update Email Campaign

```http
  PUT /api/v1/email-campaigns/<id>
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`| `json` | **Required**. Name of the email campaign we are updating|
| `email_groups`| `json` | **Required**. List of email groups to add to the campaign|
| `email_templates`| `json` | **Required**. List of email templates to add to the campaign|

```json
{
   "name": "my first email campaign",
   "email_groups": ["my first email group","my second email group"],
   "email_templates": "my first email template"
}
```


## Todo

- add more test coverage
- CLI admin tool
