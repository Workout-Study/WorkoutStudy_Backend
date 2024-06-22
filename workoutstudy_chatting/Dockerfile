# Go 빌드 환경
FROM golang:1.19.4 AS builder

# 작업 디렉토리 설정
WORKDIR /opt

# 의존성 파일들을 컨테이너로 복사
COPY go.mod go.sum ./

# 의존성 다운로드
RUN go mod download

# 소스 코드를 작업 디렉토리로 복사
COPY . .

# CGO_ENABLED를 비활성화하여 애플리케이션 빌드 (kafka-go는 CGO가 필요하지 않음)
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o workoutstudy_chatting .

# 최종 실행 환경
FROM alpine:latest

WORKDIR /opt

# 필요한 라이브러리 설치 (ca-certificates, curl, jq 추가)
RUN apk --no-cache add ca-certificates curl jq

# 빌더 스테이지에서 빌드한 실행 파일을 복사
COPY --from=builder /opt/workoutstudy_chatting .
# docs 디렉토리를 복사
COPY --from=builder /opt/docs /opt/docs  

# 실행 권한 부여
RUN chmod +x ./workoutstudy_chatting

# 컨테이너 실행 시 실행할 명령어
CMD ["./workoutstudy_chatting"]
