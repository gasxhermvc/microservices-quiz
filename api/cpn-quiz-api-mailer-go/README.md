# cpn-quiz-api-mailer-go
- โปรเจกต์ Golang สำหรับทำหน้าที่อัปโหลดไฟล์ไปเก็บ cpn-quiz-api-file-manage-go และนำ Email เก็บลง Redis Queue เพื่อให้ Schedule cpn-quiz-schedule-messenger-go ทำการ Dequeue และนำไปส่ง Email

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