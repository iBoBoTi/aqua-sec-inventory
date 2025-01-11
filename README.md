# Aqua Security Cloud Resource Inventory System

A Golang-based, clean-architecture application to manage customers and their assigned cloud resources, with notifications through RabbitMQ.

## Features

1. **Create Customer**  
    - POST `/customers`
    - RequestBody: 
        ```{
                "name": "John Doe",
                "email": "johndoe@email.com"
        }```   
2. **Get Customer By ID**  
    - GET `/customers/:id`  
3. **Add Cloud Resources**  
    - POST `/customers/:id/resources`
    - RequestBody: 
        ```{
                "resource_names": ["azure_sql_db"]
        }```  
4. **Fetch Cloud Resources**  
    - GET `/customers/:id/resources`  
5. **Update Resource**  
    - PUT `/resources/:id` 
    - RequestBody:
        ```{
                "name": "aws_vpc_main",
                "type": "VPC",
                "region":"us-east-1",
                "customer_id":, 123 //optional
        }```  
6. **Delete Resource**  
   - DELETE `/resources/:id`

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
