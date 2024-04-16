# Go 빌드 환경
FROM golang:1.19.4 AS builder

# 작업 디렉토리 설정
WORKDIR /opt

# 필요한 패키지 설치 전 libssl1.0.0 설치
RUN apt-get update && apt-get install -y wget
RUN wget http://security.ubuntu.com/ubuntu/pool/main/o/openssl1.0/libssl1.0.0_1.0.2n-1ubuntu5.8_amd64.deb
RUN dpkg -i libssl1.0.0_1.0.2n-1ubuntu5.8_amd64.deb

# 필요한 패키지 설치 및 librdkafka-dev 설치
RUN apt-get install -y \
    linux-headers-generic \
    build-essential \
    librdkafka-dev

# 의존성 파일들을 컨테이너로 복사
COPY go.mod go.sum ./

# 의존성 다운로드
RUN go mod download

# 소스 코드를 작업 디렉토리로 복사
COPY . .

# CGO_ENABLED를 활성화하고 musl 태그를 추가하여 애플리케이션 빌드
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -tags "musl dynamic" -o workoutstudy_chatting .

# 최종 실행 환경
FROM alpine:latest

WORKDIR /opt

# 필요한 라이브러리 설치
RUN apk --no-cache add ca-certificates librdkafka

# 빌더 스테이지에서 빌드한 실행 파일을 복사
COPY --from=builder /opt/workoutstudy_chatting .

# 실행 권한 부여
RUN chmod +x ./workoutstudy_chatting

# 컨테이너 실행 시 실행할 명령어
CMD ["./workoutstudy_chatting"]
