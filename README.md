# Infinite Experiment Backend

Infinite Experiment Backend is a Go-based application designed to serve requests for a Discord bot. It provides a hybrid API approach using both REST and gRPC endpoints to handle different use cases. This repository is structured to support fast, iterative local development with hot reloading (via Air) and production-ready builds using Docker multi-stage builds.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Local Development](#local-development)
  - [Production Build](#production-build)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Overview

This project serves as the backend for a Discord bot, managing communication between the bot and services such as a PostgreSQL database, REST endpoints, and gRPC endpoints for real-time communication. The goal is to provide a fast and scalable system that can be easily tested locally and deployed in production.

## Features

- **Hybrid API Design:**  
  Combines REST for simple request-response operations with gRPC for performance-critical or streaming functionalities.

- **Dockerized Development & Production:**  
  Uses a multi-stage Dockerfile to provide separate configurations for local development (with hot reloading via Air) and production builds.

- **Environment Configuration:**  
  Easily switch between local (`.env.local`) and production (`.env.production`) settings.

## Getting Started

### Local Development

To start local development with hot reloading:

1. **Ensure Docker is installed and running.**
2. **Configure your local environment variables:**  
   Create a `.env.local` file in the root directory with variables similar to:
   ```env
   APP_ENV=local
   DEBUG=true
   PORT=8080
3. **Run the application using Docker Compose:**
    `docker-compose up --build`
