# fc-user-service

## description
Golang service that owns user data

# Attention

* Not to use <https://github.com/zenazn/goji>. Checkout the later version <https://goji.io/>. Or look for something else widely used by the community.
* Use <https://github.com/Sirupsen/logrus> for logging for structure logging. The motivation is to easily take advantage of logrus hooks to centralize logs data into one searchable place. With the `correlation_id` on every payload, it should be possible to trace the whole data flow through all micro services.
* Do not share models or business logic code across services.

### User Model

```json
{
  id
  name
  email_address
  phone_number
  identity_card {
    number
  }
  vehicles []{
    type
    license_plate
  }
  c_at
  u_at
}
```
`vehicles.type` can be `car/motobike`

### End Points
`curl -X POST /tenants/:tenant_id/users`
Create a new user. Payload:

```json
{
  correlation_id
  tenant_id
  name (optional)
  email_address (optional)
  phone_number (required)
  identity_card (optional) {
    number (required)
  }
  vehicles (optional) []{
    type (required)
    license_plate (required)
  }
}
```

`curl -X GET /tenants/:tenant_id/users?identity_card.number=123&vehicles.license_plate=abc&limit=1`

Find user(s).

Accept params `id`, `phone_number`, `email_address`, `identity_card.number`, `vehicles.license_plate` as filterers.

Accept params `limit` to limit the results to find.

`curl -X GET /tenants/:tenant_id/users/:user_id`

Find one user with given id

`curl -X PUT /tenants/:tenant_id/users/:user_id`

Replace user with given id with given payload

`curl -X DELETE /tenants/:tenant_id/users/:user_id`

Delete user with the given id

### Events
The service should publish the appropriate message to the message bus when these events happen: `event.users.created`, `event.users.updated`, `event.users.deleted`

The payload should be like:

```json
{
  correlation_id
  tenant_id
  doc { User Model }
}
```
