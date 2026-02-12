# TODO LIST

## 🛠 Backend (Go Gin Framework) 설정

### 📦 라이브러리 설치
터미널에서 아래 명령어를 실행하여 필요한 패키지를 설치합니다.

```bash
# Gin Web Framework 설치
go get -u github.com/gin-gonic/gin (https://github.com/gin-gonic/gin)

# JWT 인증 라이브러리 설치
go get github.com/golang-jwt/jwt/v5 (https://github.com/golang-jwt/jwt/v5)

# GORM 및 SQLite 드라이버 설치
go get -u gorm.io/gorm
go get -u github.com/glebarez/sqlite (https://github.com/glebarez/sqlite)

# GORM 라이브러리 드라이버 교체
go get -u github.com/glebarez/sqlite

```
```java
import (
    // "gorm.io/driver/sqlite" <-- 이거 지우고
    "github.com/glebarez/sqlite" // <-- 이걸로 교체
    "gorm.io/gorm"
)
```
<h3>Dio 패키지 설치</h3>
2가지 방법

1. pubspec.yaml 파일에서 dependencies 섹션에 dio 추가
2. fluuter pub add dio로 패키지 설치

<h3> 토큰 안전한 저장을 위한 패키지 설치</h3>

```bash
flutter pub add flutter_secure_storage
```

<h3>플러터 실행</h3>
flutter_secure_storage 패키지는 웹에서 추가 설정 없이는 실행 불가

안드로이드 폰 연결 혹은 에뮬레이터 설치해서 동작해야함.

웹 동작으로 확인하고 싶으면 아래와 같은 명령어로 확인

main.go 파일에서도 CORS 관련 추가 설정 필요

```bash
flutter run -d web-server --web-port=5000
```

# 프로젝트 설정 가이드

## 📦 패키지 설치

### 1. Dio (HTTP 클라이언트)
아래 두 가지 방법 중 하나를 선택하여 설치합니다.
* **터미널 명령:** `flutter pub add dio`
* **직접 수정:** `pubspec.yaml`의 `dependencies` 섹션에 `dio: ^최신버전` 추가

### 2. Flutter Secure Storage (보안 저장소)
로그인 토큰 등의 민감 정보를 안전하게 저장하기 위해 사용합니다.
```bash
flutter pub add flutter_secure_storage
```

## 🚀 실행 및 환경 설정

### 📱 플랫폼별 주의사항
* **Mobile (Android/iOS)**
  - 실제 기기 또는 에뮬레이터 연결 시 별도 설정 없이 즉시 사용 가능합니다.
* **Web (Chrome/Edge 등)**
  - `flutter_secure_storage`는 웹 환경에서 데이터를 암호화하여 저장하기 위해 추가적인 설정이 필요할 수 있으며, 환경에 따라 동작이 제한될 수 있습니다.

### 🌐 웹 환경 테스트 및 CORS 대응
웹 브라우저에서 실행 시 발생하는 **CORS(Cross-Origin Resource Sharing)** 이슈를 방지하고 일관된 포트를 유지하기 위해 아래 명령어로 실행하는 것을 권장합니다.

```bash
# 특정 포트(5000)를 지정하여 웹 서버 실행
flutter run -d web-server --web-port=5000
```

### 📂 Go 프로젝트 파일 구조
| 계층 | 역할 |
| :--- | :--- |
| **Route** | 경로(URL) 설정 및 핸들러 연결 |
| **Handler** | HTTP 요청 입구 및 응답 반환 |
| **Service** | 핵심 비즈니스 로직 처리 |
| **Repository** | 데이터베이스(DB) 접근 로직 |
| **Model** | 데이터 구조 정의 (Struct / DB Schema) |

### 🚦 API 응답 상태 코드
| 상태 코드 상수 | 숫자 | 용도 |
| :--- | :---: | :--- |
| `http.StatusOK` | 200 | 조회(GET), 수정(PUT/PATCH) 성공 시 |
| `http.StatusCreated` | 201 | 생성(POST) 성공 시 |
| `http.StatusBadRequest` | 400 | 클라이언트 데이터 오류 (파싱 실패 등) |
| `http.StatusUnauthorized` | 401 | 권한 없음 (로그인 필요) |
| `http.StatusNotFound` | 404 | 데이터 존재하지 않음 |
| `http.StatusInternalServerError` | 500 | 서버 로직/DB 에러 발생 시 |