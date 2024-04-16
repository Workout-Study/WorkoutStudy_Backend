# Go 빌드 환경
FROM golang:1.19.4 AS builder

# 작업 디렉토리 설정
WORKDIR /opt

# librdkafka의 최신 버전을 설치하기 위해 Confluent 리포지토리 추가
RUN apt-get update && apt-get install -y software-properties-common \
    && add-apt-repository "deb [arch=amd64] http://packages.confluent.io/deb/7.0 stable main" \
    && wget -qO - http://packages.confluent.io/deb/7.0/archive.key | apt-key add - \
    && apt-get update

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
