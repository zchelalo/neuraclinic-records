# Neuraclinic Records Microservice

Go gRPC clinical records service for Neuraclinic.

## Local Setup

Run from `neuraclinic-records`:

```bash
make create-envs
make tls-generate-dev
make compose-build
```

The service listens inside Docker on `:8000` and is exposed on host port `8003`.

`AttachmentService` calls the external `FileManagementService` over gRPC. Start that service on the shared `neuraclinic-network` before testing attachments end to end.

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

