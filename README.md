# Bookify - Microservices Architecture

Bookify is a web-based platform that allows users to create, manage events, and book seats online. 
The project follows a microservices architecture to ensure scalability, maintainability, and high availability.

## Tech Stack

### Backend
- **Golang**: The primary backend language.
- **Gin**: Web framework for building high-performance REST APIs.
- **Gorm**: ORM for handling database operations efficiently.
- **Bleve**: Full-text search engine to support global search functionalities.
- **Kafka / RabbitMQ**: Message brokers for asynchronous communication between services.
- **Redis**: Used for caching and session management.
- **PostgreSQL**: The main relational database for storing structured data.
- **Cronjob (Golang)**: Scheduled tasks for automated system processes.

### Frontend
- **Nuxt.js (Vue.js + TypeScript)**: The main frontend framework for SSR and SPA capabilities.
- **Pinia Persist**: State management with persistence across sessions.
- **IndexedDB (IDB Library)**: Local database storage for offline functionality.
- **$fetch (Nuxt)**: API request handling with built-in SSR support.

### DevOps & Deployment
- **Railway**: Cloud platform for deploying backend services.
- **Docker**: Containerized services for microservices deployment.
- **Kubernetes**: Orchestrating microservices for scalability and fault tolerance.
- **Grafana + Prometheus**: Monitoring and logging for system health tracking.

## Key Features
- **User Authentication & Authorization** (JWT-based access control, refresh tokens).
- **Event Management** (CRUD operations for creating and managing events).
- **Booking System** (Allow users to book, modify, and cancel reservations).
- **Global Search** (Full-text search powered by Bleve for quick data retrieval).
- **Real-time Notifications** (Using WebSockets or Kafka for event updates and system alerts).
- **Optimistic Updates** (Applied only in user table to enhance UX and performance).
- **Google Calendar Integration** (Auto-scheduling events when users book seats).
- **Microservices Communication** (Handled through Kafka or RabbitMQ).
- **Scheduled Jobs** (Automated background tasks via Golang cronjobs).

## Purpose of Libraries
Each library used in this project serves a specific purpose to optimize system performance and maintainability:
- **Bleve**: Enables fast and efficient full-text search across large datasets.
- **Kafka / RabbitMQ**: Facilitates asynchronous messaging for loosely coupled microservices.
- **Pinia Persist**: Ensures persistent state management on the frontend.
- **IndexedDB (IDB Library)**: Provides offline data storage with an ISC-licensed library.
- **Gorm**: Simplifies database interactions while maintaining high performance.
- **Cronjob in Golang**: Handles recurring system tasks like cleanup, reporting, and notifications.

## Testing Strategy
- **Unit Tests**: Table-driven tests for backend services.
- **API Testing**: Direct integration testing rather than using mock controllers.
- **End-to-End Testing**: Ensuring seamless communication between microservices.

## Future Enhancements
- Implementing **event-driven architecture** to further decouple services.
- Improving **data synchronization mechanisms** for better offline support.
- Enhancing **observability** with more detailed logging and distributed tracing.

---
This README provides an overview of Bookify's architecture and key components. If you have any questions or need further details, feel free to reach out!

