Lamda Listener:
python3 lambda-handler.py




MinIO Server

CI=true MINIO_ROOT_USER=minio MINO_ROOT_PASSWORD=minio123 MINIO_LAMBDA_WEBHOOK_ENABLE_myfunction=on MINIO_LAMBDA_WEBHOOK_ENDPOINT_myfunction=http://localhost:5000 minio server /tmp/data-lambda-test --address ":22000"

mc alias set local22 http://localhost:22000 minio minio123


mc mb local22/test-bucket 

mc cp 1.txt local22/test-bucket/ 


Consumer/Client
cd lamda-consumer

go run <file>
