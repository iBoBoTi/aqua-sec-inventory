# Aqua Security Cloud Resource Inventory System

A Golang-based, clean-architecture application to manage customers and their assigned cloud resources, with notifications through RabbitMQ.

## Features

1. **Create Customer**  
    - POST `/customers`
    - RequestBody: 
        <pre>
        {
            "name": "John Doe",
            "email": "johndoe@email.com"
        }
        </pre>  
2. **Get Customer By ID**  
    - GET `/customers/:id`  
3. **Add Cloud Resource**  
    - POST `/customers/:id/resources`
    - RequestBody:
        <pre>{"resource_name": "azure_sql_db"}</pre>
4. **Fetch Cloud Resources**  
    - GET `/customers/:id/resources`  
5. **Update Resource**  
    - PUT `/resources/:id` 
    - RequestBody:
        <pre>
        {
            "name": "aws_vpc_main",
            "type": "VPC",
            "region":"us-east-1",
            "customer_id":, 123 //optional
        }
        </pre>  
6. **Delete Resource**  
   - DELETE `/resources/:id`
7. **Notifications**  
   - **REST**:  
     - GET `/notifications/:user_id` (Get all)  
     - DELETE `/notifications/:user_id` (Clear all)  
     - DELETE `/notifications/:user_id/:notification_id` (Clear single)  

## Quick Start

1. **Clone** the github repository:
    <pre>bash
    git clone https://github.com/iBoBoTi/aqua-sec-inventory.git
    cd aqua-sec-inventory</pre>

2. Build the application to run using docker with the make command:
    `make build`

3. In a seperate terminal Run your database migration:
    `make run-migration`

4. Seed the available Cloud Resources:
    `make seed-db`

## Automated Testing
1. Run the make command to test:
    `make test`
