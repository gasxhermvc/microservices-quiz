# microservices-quiz
สำหรับ Quiz CPN ทดสอบจัดทำโดย gasxhermvc
email: dev.awesome.th@gmail.com
github: gasxhermvc

# How to test
1. docker-compose up -d
2. เปิด Program postman และนำเข้า Collection & Environment จาก Folder postman
3. ทำการใช้ Enviroment CPN.Dev
4. ทำการเรียก API token จาก Collection keycloak
5. ทำการเรียก API token จาก Collection authentication-service
6. เมื่อเรียกได้รับ token แล้วคุณสามารถเลือกใช้งานได้ระหว่าง mailer-service หรือ file-manage-service

# วิธีการเช็ค Queue ใน Redis
* ทำการ Access หา Container redis ด้วยคำสั่ง redis-cli -a Wh5a2N7CJve36HRpSZAukKMQDPLdmYG9tTcbxrz4
* คำสั่งในการตรวจสอบ Message queue ทั้งหมดใน Queue collection lrange send-email 0 -1

# เพื่อยกระดับการป้องกัน Email Service นี่คือสิ่งที่ควรทำ
* 1. ใช้บริการเว็บ CDN ที่มีคุณสมบัติช่วย Security, Cache ,DDoS Protection เพื่อลง Traffic ที่อาจะเกิดขึ้นโดยตรงกับเว็บเรา เช่น Cloudflare เป็นต้น
* 2. กำหนด Rate limit
* 3. กำหนด White list IP
* 4. กำหนด Permission และ API Key ที่จะอนุญาตเสมอ
* 5. กำหนด Limit size ของ Request
* 6. กำหนดนามสกุลไฟล์ที่อนุญาตให้อัปโหลดเสมอ
* 7. กำหนด Limit จำนวนไฟล์ และขนาด File เสมอ

# เพื่อยกระดับ Performance โดยต่อยอดจาก Solution ปัจจุบัน
* Note: หากคุณต้องการยกระดับ Performance ของ Mailer service คุณสามารถลด Process บางอย่างที่เกิดขึ้น เช่นการอัปโหลดไฟล์ไปที่ File manage service เพื่อจัดเก็บไฟล์ ซึ่งกระบวนการนี้สามารถลดลง โดยการ Config การ Mount path file ของ Docker ให้สามารถมองเห็น Folder จัดเก็บไฟล์ของ File manage service ซึ่งจะสามาถช่วยลด Letency ระหว่างการเชื่อมต่อลงได้
### หมายเหตุ: ทั้งนี้ทั้งนั้นขึ้นอยู่กับการตกลงกัน และการยอมรับได้ เพราะการลด Process บางอย่างลงอาจนำมาซึ่ง Security ที่ลดลงเช่นกัน

# ส่วนของ URL และ API Gateway
* คุณสามารถ Config kong ให้ใช้งานร่วมกับ Nginx เพื่อใช้งานคุณสมบัติ Reverse proxy สำหรับแปลงตัวเองเป็น Domain public โดยสมบูรณ์สำหรับให้ Client สามารถเข้าถึง SSO Keycloak ได้
### หมายเหตุ: ปัจจุบัน Keycloak จะใช้งาน Port 8080 ซึ่งไม่ได้มีการ Allow ให้เข้าถึงจากภายนอกการใช้ Reverse proxy เพื่อให้เข้าถึง Keycloak แบบ Public Domain สามารถทำได้สะดวก และเป็นวิธีที่ง่าย

# ส่วนสุดท้าย
* สำหรับการนำเทคโนโลยี Microservices และ Container เข้าใช้งาน การจัดการ Container, ข้อมูล และการ Monitor เป็นสิ่งสำคัญที่ขาดไม่ได้ ปัจจุบันมี Solution มากมายที่ใช้ เช่น Elaticsearch+Logstash+Kibana สำหรับจัดการ Logs, Grafana+Prometheus สำหรับดู Performance และ Portianer สำหรับจัดการ Container เป็นต้น เพื่อให้การดูแล Microservices ทำได้ง่ายดายขึ้น