# Go 버전 지정
FROM golang:1.19.4 AS builder

# 작업 디렉토리 설정
WORKDIR /app

# 의존성 파일들을 컨테이너로 복사
COPY go.mod .
COPY go.sum .

# 의존성 다운로드
RUN go mod download

# 소스 코드를 작업 디렉토리로 복사
COPY . .

# 애플리케이션 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o workoutstudy_chatting .

# 실행 이미지
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 빌더 스테이지에서 빌드한 실행 파일을 복사
COPY --from=builder /app/workoutstudy_chatting .

# 컨테이너 실행 시 실행할 명령어
CMD ["./workoutstudy_chatting"]