# 🧾 Receipt Processor

Welcome to the Receipt Processor project! This project provides a web service for processing receipts and calculating reward points based on various rules. It also includes JWT-based authentication to secure the API endpoints.

## 📜 Table of Contents
- [Features](#-features)
- [Approach](#-approach)
- [Implementation Steps](#️-implementation-steps)
- [Tech Stack](#-tech-stack)
- [Installation and Running the Application](#-installation-and-running-the-application)
- [API Endpoints](#-api-endpoints)
- [Example Usage](#-example-usage)

## ✨ Features
- **Receipt Processing**: Accepts receipt details and processes them to calculate reward points.
- **Point Calculation**: Points are calculated based on rules such as retailer name length, purchase time, and item details.
- **JWT Authentication**: Secures endpoints, allowing only authorized users to access the API.
- **In-Memory Data Storage**: Stores receipts and their points temporarily in memory, ensuring stateless design.
- **Error Handling**: Validates data input with comprehensive error handling and descriptive error messages.

## 🧠 Approach
The project follows a stateless microservice design, using an in-memory store to hold data temporarily. This approach ensures:

- Minimal resource usage without the need for a persistent database.
- Authentication using JWT for secure access.
- Efficient memory management with sync.Mutex for safe, concurrent access to the in-memory store.
- A clean structure with separate modules for handlers, models, and utilities, making it modular and maintainable.

## 🛠️ Implementation Steps
1. **Define Models**: Created data structures for receipts, items, and processed receipts in the models package.
2. **Set up JWT Authentication**: Implemented token generation and validation functions using JWT.
3. **Implement Business Logic**:
   - Developed rules for calculating points based on receipt fields.
   - Included helper functions for each rule to handle points based on item count, date, time, etc.
4. **API Endpoints**:
   - `POST /receipts/process`: Processes receipts, calculates points, and returns a unique receipt ID.
   - `GET /receipts/{id}/points`: Retrieves points for a given receipt ID.
5. **Error Handling**: Validated fields and returned meaningful error messages for invalid requests.
6. **Testing**: Developed test cases to cover different scenarios, including edge cases and invalid inputs.

## 🧰 Tech Stack
- **Language**: Go (Golang)
- **Libraries**:
  - gorilla/mux for routing
  - golang-jwt/jwt for JWT authentication
  - uuid for generating unique receipt IDs

## 🚀 Installation and Running the Application

### Prerequisites
Make sure Go is installed on your system. You can download it from [here](https://golang.org/dl/).

### Steps
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/saurabhag23/receipt-processor.git
   cd receipt-processor
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Generate JWT for Testing**:
   - Run the generate_jwt.go file to generate a JWT token that can be used to access the endpoints.
   ```bash
   go run generate_jwt.go
   ```
   - Copy the generated token from the output.

4. **Start the Server**:
   ```bash
   go run main.go
   ```
   The server will start on http://localhost:8080.

## 📡 API Endpoints

### 1. Process Receipt 🧾
- **URL**: `/receipts/process`
- **Method**: POST
- **Description**: Submits a receipt for processing and returns a unique ID.
- **Headers**:
  - `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Body** (JSON):
  ```json
  {
      "retailer": "Target",
      "purchaseDate": "2022-01-01",
      "purchaseTime": "13:01",
      "items": [
          { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
          { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
      ],
      "total": "18.74"
  }
  ```
- **Response** (JSON):
  ```json
  { "id": "unique-receipt-id" }
  ```

### 2. Get Points 🎯
- **URL**: `/receipts/{id}/points`
- **Method**: GET
- **Description**: Retrieves the points awarded for a specific receipt.
- **Headers**:
  - `Authorization: Bearer <YOUR_JWT_TOKEN>`
- **Response** (JSON):
  ```json
  { "points": 28 }
  ```

## 💡 Example Usage

### Step 1: Generate a JWT Token
Use generate_jwt.go to get a token for testing:
```bash
go run generate_jwt.go
# Output: Generated JWT Token: <YOUR_JWT_TOKEN>
```

### Step 2: Process a Receipt
Submit a receipt for processing:
```bash
curl -X POST http://localhost:8080/receipts/process \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{
         "retailer": "Target",
         "purchaseDate": "2022-01-01",
         "purchaseTime": "13:01",
         "items": [
             { "shortDescription": "Mountain Dew 12PK", "price": "6.49" },
             { "shortDescription": "Emils Cheese Pizza", "price": "12.25" }
         ],
         "total": "18.74"
     }'
# Expected Response: { "id": "<generated_id>" }
```

### Step 3: Retrieve Points for the Receipt
Use the generated id to retrieve points:
```bash
curl -X GET http://localhost:8080/receipts/<generated_id>/points \
     -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
# Expected Response: { "points": <calculated_points> }
```

## 🧪 Testing
Run the application with a variety of test cases to ensure all functionalities work correctly, including:
- Valid and invalid receipts.
- Receipts with various edge cases, such as round totals, odd/even dates, specific times, etc.
- Testing without authorization or with invalid tokens to verify access restrictions.

## 📋 Rules for Point Calculation
Points are calculated based on these rules:
- **Retailer Name**: 1 point per alphanumeric character.
- **Round Dollar Total**: 50 points if the total has no cents.
- **Total is a Multiple of 0.25**: 25 points.
- **Item Count**: 5 points for every two items.
- **Item Description**: If description length is a multiple of 3, award points based on price.
- **Odd Purchase Day**: 6 points if the day is odd.
- **Specific Purchase Time**: 10 points if the time is between 2:00 pm and 4:00 pm.

## ⚠️ Error Handling
The application provides comprehensive error handling with descriptive messages for:
- Missing or incorrectly formatted fields in the receipt.
- Invalid JWT tokens or missing authentication.
- Attempts to retrieve points for non-existent receipt IDs.

## 🤝 Contributions
Contributions are welcome! If you'd like to improve this project, please feel free to fork the repository and submit a pull request.

## 📄 License
This project is open source and available under the MIT License.
