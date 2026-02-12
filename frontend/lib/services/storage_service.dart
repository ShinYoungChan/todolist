import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/* secure_storage 패키지로인해 web에서는 실행 불가능 / 테스트를 위해 map으로 일단 진행... main.go에서도 cors 추가 설정해야함
class StorageService {
  // 보안 저장소 인스턴스 생성
  final _storage = const FlutterSecureStorage();

  // 토큰 저장 (로그인 성공 시 호출)
  Future<void> saveToken(String token) async {
    await _storage.write(key: 'jwt_token', value: token);
  }

  // 토큰 읽기 (API 요청 보낼 때 호출)
  Future<String?> getToken() async {
    return await _storage.read(key: 'jwt_token');
  }

  // 토큰 삭제 (로그아웃 시 호출)
  Future<void> deleteToken() async {
    await _storage.delete(key: 'jwt_token');
  }
}
*/

class StorageService {
  // 실제 저장소 대신 메모리(Map) 변수 사용
  static final Map<String, String> _memoryStorage = {};

  Future<void> saveToken(String token) async {
    _memoryStorage['jwt_token'] = token;
    print("토큰 메모리 저장 완료!");
  }

  Future<String?> getToken() async {
    return _memoryStorage['jwt_token'];
  }

  Future<void> deleteToken() async {
    _memoryStorage.remove('jwt_token');
  }
}