import 'package:dio/dio.dart';

class ApiService {
  // 1. 서버 주소 설정
  // 안드로이드 에뮬레이터: 10.0.2.2 / iOS 시뮬레이터: localhost
  static const String baseUrl = "http://localhost:8080"; 

  late Dio dio;

  ApiService() {
    dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 5), // 5초 연결 제한
      receiveTimeout: const Duration(seconds: 3),
      contentType: 'application/json',
    ));

    // 💡 인터셉터 추가 (나중에 토큰 자동 삽입을 위해 사용)
    dio.interceptors.add(LogInterceptor(responseBody: true, requestBody: true));
  }

  // 로그인 API 예시
  Future<Response> login(String email, String password) async {
    return await dio.post("/login", data: {
      "email": email,
      "password": password,
    });
  }
}