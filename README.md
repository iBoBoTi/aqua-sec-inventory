# Aqua Security Cloud Resource Inventory System

A Golang-based, clean-architecture application to manage customers and their assigned cloud resources, with notifications through RabbitMQ.

## Features

1. **Create Customer**  
   - POST `/customers`  
2. **Get Customer By ID**  
   - GET `/customers/:id`  
3. **Add Cloud Resources**  
   - POST `/customers/:customer_id/resources`  
4. **Fetch Cloud Resources**  
   - GET `/customers/:customer_id/resources`  
5. **Update Resource**  
   - PUT `/resources/:resource_id`  
6. **Delete Resource**  
   - DELETE `/resources/:resource_id`  
7. **Notifications**  
   - **REST**:  
     - GET `/notifications/:user_id` (Get all)  
     - DELETE `/notifications/:user_id` (Clear all)  
     - DELETE `/notifications/:user_id/:notification_id` (Clear single)  
   - **gRPC (optional)**: see `internal/transport/grpc/notification_service.go` for a stub.

## Quick Start

1. **Clone** the github repository:
   ```bash
   git clone https://github.com/iBoBoTi/aqua-sec-inventory.git
   cd aqua-sec-inventory
```

2. Build the application to run using docker with the make command:
    `make build`

3. In a seperate terminal Run your database migration:
    `make run-migration`

4. Seed the available Cloud Resources:
    `make seed-db`

## Automated Testing
1. Run the make command to test:
    `make test`
