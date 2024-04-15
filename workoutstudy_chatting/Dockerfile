# Go 버전 지정
FROM golang:1.19.4 AS builder

# 작업 디렉토리 설정
WORKDIR /

# 필요한 패키지 설치 및 librdkafka 설치
RUN apt-get update
RUN apt-get install -y linux-headers-generic
RUN apt-get install -y build-essential 
RUN apt-get install -y --no-install-recommends librdkafka-dev

# 의존성 파일들을 컨테이너로 복사
COPY go.mod .
COPY go.sum .

# 의존성 다운로드
RUN go mod download

# 소스 코드를 작업 디렉토리로 복사
COPY . .

# CGO_ENABLED를 1로 설정하여 애플리케이션 빌드
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o workoutstudy_chatting .

# 실행 이미지
FROM alpine:latest

WORKDIR /

RUN apk --no-cache add ca-certificates librdkafka

# 빌더 스테이지에서 빌드한 실행 파일을 복사
COPY --from=builder /workoutstudy_chatting .

# 실행 권한 부여
# RUN chmod +x /app/workoutstudy_chatting

# 컨테이너 실행 시 실행할 명령어
CMD [".workoutstudy_chatting"]