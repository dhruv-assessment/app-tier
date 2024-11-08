# Multi-Tier Face Recognition Application

This project is a cloud-based, multi-tier face recognition application that processes images of people, identifies faces, and returns their names with minimal latency. The system is designed for scalability and performance, with components distributed across different cloud infrastructure services.

Given an image, the application uses a pre-trained machine learning model (`facenet_pytorch`) to recognize the person in the image. The architecture is composed of three main layers: the Web Tier, App Tier, and a Custom Load Balancer, each with distinct functionalities to ensure smooth and efficient processing.

## Architecture Design

![Face Recognition Architecture](assets/assessment.png)

*Figure: System Architecture of the Multi-Tier Face Recognition Application*

### Architecture Overview

The application is divided into three primary components, each hosted in a separate repository:

1. **Web Tier**: Handles HTTP requests, forwards images to the App Tier, and returns results to clients.
2. **App Tier**: Processes images using face recognition and returns the name of the person in the image.
3. **Custom Load Balancer**: Manages scaling of App Tier instances on AWS based on request volume.

### Key Components and Flow

- **Client**: The client sends an image to the Web Tier, which acts as the entry point for processing.
- **Web Tier**: Receives images, forwards them to the App Tier, and returns the recognized results to the client. It also connects to an InfluxDB instance for dummy data population.
- **Amazon SQS**: Used to manage queues between the Web Tier and App Tier, facilitating scalable message passing.
- **App Tier**: Uses a pre-trained face recognition model (`facenet_pytorch`) to process images and return results.
- **Amazon S3**: Stores input images and output results as needed for additional data persistence.
- **Custom Load Balancer**: Monitors incoming request load and dynamically scales App Tier instances on AWS to maintain performance.

### Key Features

- **Distributed Processing**: Each tier operates independently, allowing scalable processing and easy maintenance.
- **Custom Load Balancing**: Manages the App Tier instances dynamically, without relying on AWS's auto-scaling services.
- **Dockerized Deployment**: Each component is containerized, simplifying deployment and ensuring environment consistency.
- **Zero-Downtime CI/CD Pipeline**: Ensures seamless deployment without affecting the application's availability.

# App Tier

## Description
The App Tier is the core processing unit responsible for face recognition. It receives images from the Web Tier, uses a pre-trained ML model (`facenet_pytorch`) to identify faces, and returns the recognized names to the Web Tier.

## Features
- **Face Recognition**: Processes images using `facenet_pytorch` to identify faces.
- **Efficient Image Processing**: Optimized for low latency and accurate recognition.
- **Scalable**: Designed to work with the Custom Load Balancer to manage instance scaling on AWS.

## Tech Stack
- **Programming Language**: Go, Python
- **Machine Learning Model**: `facenet_pytorch`
- **Dependencies**: AWS SDK

## Setup

1. **Export AWS Environment Variables or create .env file**
    ```bash
     export AWS_ACCESS_KEY_ID=your_aws_access_key
     export AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key
     export AWS_DEFAULT_REGION=us-east-1
     export AWS_IN_BUCKET_NAME=your_input_bucket_name
     export AWS_OUT_BUCKET_NAME=your_output_bucket_name
     export AWS_REQ_URL=your_request_url
     export AWS_RESP_URL=your_response_url
    ```
2. **Run the app-tier**
    ```bash
    go run main.go
    ```
