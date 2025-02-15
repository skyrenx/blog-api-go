# **blog-api-go**

`blog-api-go` is an AWS Lambda-based Go module designed to serve as a backend API. The module handles API Gateway traffic and supports both local development and production deployment.

---

## **Requirements**
- **Go 1.x**
- **AWS SAM CLI** (for local SAM testing)
- **Docker** (for local SAM testing)
- **AWS CLI** (for deployment to test/prod)

---

## **Local Development**
1. Ensure you have the **AWS SAM CLI** installed.
2. Execute the `run-sam.sh` script to start a local HTTP server that simulates API Gateway traffic:
   ```bash
   ./run-sam.sh
   ```
   This will run the Lambda locally using `sam local start-api`, allowing you to test your API endpoints at:
   ```
   http://localhost:3000
   ```

---

## **Deployment Instructions**

### **Test (TBD)**
Deployment to the test environment will be documented in future updates.

### **Production Deployment**
1. Use the `deploy.sh` script to deploy the module to the production AWS Lambda:
   ```bash
   ./deploy.sh
   ```
2. This script will:
   - Build the Go module and package it into a `.zip` file.
   - Upload it to the production Lambda using the AWS CLI.
