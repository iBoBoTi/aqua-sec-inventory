# **Aqua Security Cloud Resource Inventory System**

A Golang-based, clean-architecture application for managing customers and their assigned cloud resources, with notification capabilities via RabbitMQ.

---

## **Features**

### **1. Customer Management**
- **Create Customer**  
  **Endpoint:** `POST /customers`  
  **Request Body:**  
  ```json
  {
      "name": "John Doe",
      "email": "johndoe@email.com"
  }
  ```  
  **Response:**  
  ```json
  {
    "data": {
      "id": 1,
      "name": "John Doe",
      "email": "johndoe@email.com"
    }
  }
  ```

- **Get Customer by ID**  
  **Endpoint:** `GET /customers/:id`  
  **Response:**  
  ```json
  {
    "data": {
        "id": 1,
        "name": "John Doe",
        "email": "johndoe@email.com"
    }
  }
  ```

---

### **2. Cloud Resource Management**
- **Add Cloud Resource to Customer**  
  **Endpoint:** `POST /customers/:id/resources`  
  **Request Body:**  
  ```json
  {
      "resource_name": "azure_sql_db"
  }
  ```  
  **Response:**  
  ```json
  {
      "message": "Resources assigned successfully"
  }
  ```

- **Fetch Cloud Resources by Customer**  
  **Endpoint:** `GET /customers/:id/resources`  
  **Response:**  
  ```json
  {
    "data" : [
        {
            "id": 101,
            "name": "azure_sql_db",
            "type": "SQL Database",
            "region": "us-east-1",
            "created_at": "2025-01-11T09:03:22.399082Z",
            "updated_at": "2025-01-11T09:03:22.399082Z"
        }
    ]
  }
  ```

- **Update Resource Information**  
  **Endpoint:** `PUT /resources/:id`  
  **Request Body:**  
  ```json
  {
      "name": "aws_vpc_main",
      "type": "VPC",
      "region": "us-west-2",
  }
  ```  
  **Response:**  
  ```json
  {
    {
        "id": 101,
        "name": "azure_sql_db",
        "type": "SQL Database",
        "region": "us-east-1",
        "created_at": "2025-01-11T09:03:22.399082Z",
        "updated_at": "2025-01-11T09:03:22.399082Z"
    }
  }
  ```

- **Delete Resource**  
  **Endpoint:** `DELETE /resources/:id`  
  **Response:**  
  ```json
  {
      "message": "Resource deleted successfully"
  }
  ```

---

### **3. Notification Service**
- **Get All Notifications**  
  **Endpoint:** `GET /users/:id/notifications`  
  **Response:**  
  ```json
  [
      {
          "id": 1,
          "user_id": 2,
          "message": "New resource added",
          "created": "2025-01-10T10:00:00Z"
      }
  ]
  ```

- **Clear All Notifications**  
  **Endpoint:** `DELETE /users/:id/notifications`  
  **Response:**  
  ```json
  {
      "message": "All notifications cleared"
  }
  ```

- **Clear Single Notification**  
  **Endpoint:** `DELETE /notifications/:id`  
  **Response:**  
  ```json
  {
      "message": "Notification cleared"
  }
  ```
### **4. Notification GRPC Service**
  #### GetAllNotifications
- **Request:**
  ```protobuf
  GetAllNotificationsRequest {
    user_id: 123
  }
  ```
- **Response:**
  ```protobuf
  GetAllNotificationsResponse {
    notifications: [
      {
        id: 1,
        user_id: 123,
        message: "New resource 1 added",
        created_at: "2025-01-12T10:30:00Z"
      },
      {
        id: 2,
        user_id: 123,
        message: "New resource 2 added",
        created_at: "2025-01-11T15:45:00Z"
      }
    ]
  }
  ```

#### ClearSingleNotification
- **Request:**
  ```protobuf
  ClearSingleNotificationRequest {
    notification_id: 1
  }
  ```
- **Response:**
  ```protobuf
  ClearSingleNotificationResponse {
    message: "Notification cleared successfully."
  }
  ```

#### ClearAllNotifications
- **Request:**
  ```protobuf
  ClearAllNotificationsRequest {
    user_id: 123
  }
  ```
- **Response:**
  ```protobuf
  ClearAllNotificationsResponse {
    message: "All notifications cleared successfully."
  }
  ```

---

## **Quick Start**

### **1. Clone the Repository**
```bash
git clone https://github.com/iBoBoTi/aqua-sec-inventory.git
cd aqua-sec-inventory
```

### **2. Build and Run the Application**
Use Docker Compose to build and run the services:
```bash
make build
```

### **3. Run Database Migrations**
Run the database migrations to set up the schema:
```bash
make run-migration
```

### **4. Seed the Database**
Populate the database with predefined cloud resources:
```bash
make seed-db
```

---

## **Automated Testing**

Run unit and integration tests:
```bash
make test
```

---

## **Technologies Used**
- **Programming Language:** Golang
- **Framework:** Gin (for REST APIs)
- **Database:** PostgreSQL
- **Message Queue:** RabbitMQ
- **Containerization:** Docker