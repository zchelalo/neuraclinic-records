# Neuraclinic Records Microservice

Go gRPC clinical records service for Neuraclinic.

## Local Setup

Run from `neuraclinic-records`:

```bash
make create-envs
make tls-generate-dev
```

Run shared services from the root `neuraclinic` repository first:

```bash
cd ../neuraclinic
make compose-detached
cd ../neuraclinic-records
make create-network
make compose-build
```

The service listens inside Docker on `:8000` and is exposed on host port `8003`.

`neuraclinic-file-management` and `neuraclinic-rabbitmq` must be running on `neuraclinic-network` before testing attachments end to end.

Relevant env vars for the upload projection flow:

- `RABBITMQ_URL=amqp://guest:guest@neuraclinic-rabbitmq:5672/`
- `RABBITMQ_EXCHANGE=neuraclinic.events`
- `RABBITMQ_ROUTING_KEY=file.record.status_changed.v1`

## Useful Commands

```bash
make proto
make sqlc
make test
make build
make migrate-up
make compose
make compose-down
```
