# **File Management API**  
This API allows users to **sign up, log in, upload files, and access files via temporary URLs**.

---

## **How to Run This Project**  

### **1. Clone the Repository**  
Run the following command to clone the repository:  
```bash
git clone https://github.com/naresh2002/8901135580_Backend.git
```

### **2. Install Dependencies**  
Navigate to the project directory and install the required Go dependencies:  
```bash
go get github.com/jackc/pgx/v5
go get github.com/joho/godotenv
go get github.com/gorilla/mux
go get github.com/golang-jwt/jwt/v5
go mod tidy
```

### **3. Set Up PostgreSQL Database**  

#### **Create a New Database**  
- Open PostgreSQL shell (`psql`):  
  - **For Ubuntu**:  
    ```bash
    sudo -u postgres psql
    ```
  - **For Windows/Mac (if psql is configured in PATH)**:  
    ```bash
    psql -U postgres
    ```
- Run the following command in `psql` shell to create a new database:  
  ```sql
  CREATE DATABASE trademarkia;
  ```

#### **Run Database Migrations**  
Execute the `schema.sql` file to create the required tables:  
```bash
psql -U postgres -d trademarkia -h localhost -p 5432 -f schema/schema.sql
```

### **4. Configure Environment Variables**  
Create a `.env` file in the root directory and configure it as follows:  
```plaintext
DB_USER={your_database_username}
DB_PASSWORD={your_database_password}
DB_HOST=localhost
DB_PORT=5432
DB_NAME=trademarkia
JWT_SECRET={your_jwt_secret_key}
```
Replace `{your_database_username}`, `{your_database_password}`, and `{your_jwt_secret_key}` with actual values.

### **5. Start the Project**  
Run the project using:  
```bash
go run main.go
```
Now the API is running at `http://localhost:8000/`.

---

## **Endpoints & Usage**  

### **1. User Signup** [POST]  
Registers a new user in the system.  

**Endpoint:**  
```
POST http://localhost:8000/signup
```
**Request Body (JSON):**
```json
{
    "email": "user@example.com",
    "password": "securepassword"
}
```
**cURL Command:**
```bash
curl -X POST http://localhost:8000/signup \
     -H "Content-Type: application/json" \
     -d '{
          "email": "user@example.com",
          "password": "securepassword"
         }' | jq
```
**Response (JSON):**
```json
{
    "message": "User registered successfully"
}
```

---

### **2. User Login** [GET]  
Logs in the user and provides a **JWT token** for authentication.

**Endpoint:**  
```
GET http://localhost:8000/login
```
**Request Body (JSON):**
```json
{
    "email": "user@example.com",
    "password": "securepassword"
}
```
**cURL Command:**
```bash
curl -X GET http://localhost:8000/login \
     -H "Content-Type: application/json" \
     -d '{
          "email": "user@example.com",
          "password": "securepassword"
         }' | jq
```
**Response (JSON):**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
✅ **Note:** Copy the `token` from the response. It will be used in the `Authorization` header for uploading files.

---

### **3. Upload File** [POST]  
Uploads a file and stores it in the `media/{user_id}/` directory. It also creates an entry in the database.

**Endpoint:**  
```
POST http://localhost:8000/upload
```
**Headers:**
- `Authorization: Bearer <your_token>`

**cURL Command:**
```bash
curl -X POST http://localhost:8000/upload \
     -H "Authorization: Bearer <your_token>" \
     -F "file=@/path/to/your/file.pdf"
```
**Response (JSON):**
```json
{
    "message": "File uploaded successfully",
    "file_id": 1,
    "temporary_url": "http://localhost:8000/file?url=abcd1234xyz"
}
```
✅ **Note:** The `temporary_url` can be used to access the file.

---

### **4. Access File (Public URL)**  
Fetches the file using the generated **temporary URL**.

Enter temprary_url in your browser's search bar  
```plaintext
http://localhost:8000/file?url=abcd1234xyz
```
✅ **Note:** This allows **anyone** to access the file as long as it hasn't expired.

---

## **Authentication Flow**  
1. **Signup:** Create a new user.  
2. **Login:** Get a JWT token.  
3. **Use Token:** Pass the token in `Authorization: Bearer <token>` when uploading a file.  
4. **Access File:** Use the `temporary_url` from the upload response to fetch the file.  

---

### **Project Features**
✅ **Secure File Uploads**: Stores files in user-specific directories (`media/{user_id}/`).  
✅ **JWT-Based Authentication**: Ensures only authenticated users can upload files.  
✅ **Temporary Public URLs**: Allows file access via short-lived public links.  
✅ **PostgreSQL Database**: Uses PostgreSQL to manage users and file metadata.  

---

This setup ensures **secure file storage, access control via JWT, and temporary public URLs** for sharing files. 🚀
