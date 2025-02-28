# RE Partners Challenge

This is a RESTful API service that calculates and manages package sizes for orders. The API helps determine the most efficient way to pack items using available package sizes.

## Technologies Used

- **Go (Golang)**: Main programming language
- **Echo Framework**: High performance, minimalist Go web framework
- **Zerolog**: Zero allocation JSON logger
- **Swagger**: API documentation
- **Docker**: Containerization
- **In-memory Cache**: For storing package sizes

## API Endpoints

### Order Endpoints

#### Create Order
- **URL**: `/order`
- **Method**: `POST`
- **Description**: Calculate the best packing options for an order
- **Request Body**:
```json
{
    "items": 1000
}
```
- **Success Response**: HTTP 200
```json
{
    "items": 1000,
    "order_packs": [
        {
            "pack_size": 1000,
            "pack_quantity": 1
        }
    ]
}
```

### Pack Endpoints

#### Update Package Sizes
- **URL**: `/pack/sizes`
- **Method**: `PUT`
- **Description**: Update available package sizes for order packing
- **Request Body**:
```json
{
    "sizes": [250, 500, 1000, 2000, 5000]
}
```
- **Success Response**: HTTP 200
```json
"success"
```

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory with the following variables:

```env
APP_PORT=3000
LOG_LEVEL=info
```

## Running Locally

1. Make sure you have Go 1.24 or later installed
2. Clone the repository
3. Install dependencies:
```bash
go mod download
```
4. Create and configure the `.env` file
5. Run the application:
```bash
go run main.go
```

## Running with Docker

1. Make sure you have Docker and Docker Compose installed
2. Clone the repository
3. Create and configure the `.env` file
4. Build and run the containers:
```bash
docker-compose up --build
```

The API will be available at `http://localhost:3000`

## Running Tests

To run the unit tests:

```bash
go test ./... -v
```

## API Documentation

The API documentation is available through Swagger UI at `http://localhost:3000/swagger/`

## Project Structure

```
.
├── config/         # Configuration management
├── docs/          # Swagger documentation
├── internal/      # Internal packages
│   ├── api/       # API handlers and routes
│   ├── clients/   # External clients (cache, etc)
│   ├── constants/ # Constants and default values
│   ├── models/    # Data models
│   └── services/  # Business logic
├── web/           # Web interface files
├── main.go        # Application entry point
├── Dockerfile     # Docker configuration
└── docker-compose.yml