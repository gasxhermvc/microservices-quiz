# cpn-quiz-api-authentication-go
- โปรเจกต์ Golang สำหรับยืนยันตัวตน โดยทำหน้าที่เป็น Web services ตรวจสอบ Access Token ว่าสร้างมาจาก SSO หรือไม่ และหากตรวจสอบ Token ถูกต้องจะสร้าง Access Token สำหรับใช้งานภายใน Application อีกที

# Program Required
* Golang Lastest: https://go.dev/dl/
* Go echo

# How to run for development
* F5 (Run with Debugging) or go run app/main.go local

# How to run for docker
* Execute image-build-dev.sh
* Execute run-image-dev.sh

# Configuration
* .env.yml with development mode
* .env.{environment}.yml without development mode
* config/{environment}.yml
* from db table cq_config