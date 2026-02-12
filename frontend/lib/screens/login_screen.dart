import 'package:flutter/material.dart';
import '../services/api_service.dart';
import '../services/storage_service.dart';
import 'todo_list_screen.dart';

class LoginScreen extends StatefulWidget {
  @override
  _LoginScreenState createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _apiService = ApiService();
  final _storageService = StorageService();

  void _login() async {
    try {
      final response = await _apiService.login(
        _emailController.text,
        _passwordController.text,
      );

      // Go 서버 응답: {"token": "ey... "}
      final token = response.data['token'];

      if (token != null) {
        // 1. 토큰 저장
        await _storageService.saveToken(token);

        // 2. 다음 화면으로 이동 (예: 할 일 목록)
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(const SnackBar(content: Text("로그인 성공!")));
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) => TodoListScreen(),
          ), // 할 일 목록 화면으로!
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text("로그인 실패: $e")));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text("Login")),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            TextField(
              controller: _emailController,
              decoration: const InputDecoration(labelText: "Email"),
            ),
            TextField(
              controller: _passwordController,
              decoration: const InputDecoration(labelText: "Password"),
              obscureText: true,
            ),
            const SizedBox(height: 20),
            ElevatedButton(onPressed: _login, child: const Text("Login")),
          ],
        ),
      ),
    );
  }
}
