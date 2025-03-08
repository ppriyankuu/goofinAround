# Real-time Data Processing Pipeline

## Overview  

This pipeline can process real-time data streams using Go, with Redis as the message queue. The pipeline consists of the following components:  

- **Data Producer**: Generates and sends data to the Redis queue.  
- **Data Consumer**: Reads data from the Redis queue and processes it.  
- **Data Processor**: Handles the business logic for processing the data.  

## Folder Structure  
```bash
real-time-pipeline/
├── cmd/
│   ├── producer/
│   │   └── main.go
│   └── consumer/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── data_model.go
│   ├── services/
│   │   ├── producer-service.go
│   │   └── consumer-service.go
│   └── utils/
│       └── redis-client.go
├── go.mod
└── go.sum
```