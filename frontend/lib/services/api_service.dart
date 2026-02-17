import 'package:dio/dio.dart';
import 'package:frontend/services/storage_service.dart';

class ApiService {
  // 1. 서버 주소 설정
  // 안드로이드 에뮬레이터: 10.0.2.2 / iOS 시뮬레이터: localhost
  static const String baseUrl = "http://localhost:8080";
  late Dio dio;
  final StorageService _storageService = StorageService(); // 저장소 인스턴스

  ApiService() {
    dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: const Duration(seconds: 5), // 5초 연결 제한
        receiveTimeout: const Duration(seconds: 3),
        contentType: 'application/json',
      ),
    );

    dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) async {
          String? token = await _storageService.getToken();

          if (token != null) {
            options.headers["Authorization"] = "Bearer $token";
          }

          print("요청 경로: ${options.path}"); // 디버깅용 로그
          return handler.next(options); // 다음 단계로 진행
        },
        onError: (DioException e, handler) {
          if (e.response?.statusCode == 401) {
            print("인증이 만료되었습니다. 다시 로그인해주세요.");
          }
          return handler.next(e);
        },
      ),
    );
    // 💡 인터셉터 추가 (나중에 토큰 자동 삽입을 위해 사용)
    dio.interceptors.add(LogInterceptor(responseBody: true, requestBody: true));
  }

  // 로그인 API 예시
  Future<Response> login(String userID, String password) async {
    return await dio.post(
      "/login",
      data: {"user_id": userID, "password": password},
    );
  }

  // 💡 3단계의 하이라이트: 할 일 목록 가져오기 테스트용
  Future<Response> getTodos(String? sortBy, String? filter, String? keyword) async {
    // 이제 여기서는 헤더 설정을 전혀 안 해도 됩니다! 인터셉터가 해주니까요.
    final response = await dio.get("/todos",queryParameters: {
      if(sortBy!=null) "sort":sortBy,
      if(filter!=null)"filter":filter,
      if(keyword!=null&&keyword.isNotEmpty)"keyword":keyword,
    });
    return response;
  }

  /*
  Future<Response> createTodo(Map<String, dynamic> data) async {
    return await dio.post("/todos", data: data);
  }
  */

  Future<Response> createTodo({
    required String title, // 제목은 필수
    String? content, // 내용은 선택 (null 허용)
    DateTime? startDate, // 시작일 선택
  }) async {
    // 서버가 원하는 구조대로 Map 생성
    final Map<String, dynamic> data = {
      "title": title,
      "content": content,
      "start_date": startDate?.toIso8601String(), // 날짜를 문자열로 변환
      "status": false,
    };

    // null인 값은 서버로 보내지 않도록 제거 (선택 사항)
    data.removeWhere((key, value) => value == null);

    return await dio.post("/todos", data: data);
  }

  Future<Response> deleteTodo(int id) async {
    return await dio.delete("/todos/$id");
  }

  Future<Response> updateTodoState(int id, bool status) async {
    return await dio.put("/todos/$id", data: {"id": id, "status": status});
  }

  // api_service.dart 수정
  Future<void> updateTodoDates(int id, {DateTime? startDate, DateTime? dueDate}) async {
    try {
      // 로컬 시간 기준으로 YYYY-MM-DDT00:00:00Z 형식을 수동으로 맞춰줍니다.
      String formatDate(DateTime dt) {
        return "${dt.year.toString().padLeft(4, '0')}-${dt.month.toString().padLeft(2, '0')}-${dt.day.toString().padLeft(2, '0')}T00:00:00Z";
      }
      await dio.put("/todos/$id", data: {
        if (startDate != null) "start_date": formatDate(startDate),
        if (dueDate != null) "due_date": formatDate(dueDate),
      });
    } catch (e) {
      rethrow;
    }
  }
}
