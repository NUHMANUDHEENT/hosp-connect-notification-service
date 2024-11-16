# Notification Service

The **Notification Service** handles sending notifications for various events like appointment reminders, payment confirmations, and updates. It consumes messages from Kafka, processes them using Goroutines for efficiency, and sends emails through GORMC.

---

## **Features**

### **Kafka Consumer**
- Listens to Kafka topics for event messages (e.g., payment success, appointment updates).
- Processes and responds in real-time.

### **Concurrent Notifications**
- Uses **Goroutines** to handle multiple notifications concurrently for improved performance.

### **Email Notifications**
- Sends detailed emails to users (e.g., appointment reminders, payment receipts) via the **GORMC** email-sending platform.

---

## **Technology Stack**
- **Backend:** Go (Golang)
- **Messaging Queue:** Kafka
- **Concurrency:** Goroutines
- **Email Platform:** GORMC
- **Database (Optional):** PostgreSQL or MongoDB for storing notification logs.

---

## **How to Run**

### Clone the Repository
```bash
git clone https://github.com/your-username/notification-service.git
cd notification-service
