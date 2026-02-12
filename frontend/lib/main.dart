import 'package:flutter/material.dart';
import 'screens/login_screen.dart'; // 로그인 화면 불러오기

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Todo App',
      debugShowCheckedModeBanner: false, // 오른쪽 상단 DEBUG 띠 제거
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true, // 최신 디자인 스타일 적용
      ),
      // 앱의 첫 페이지를 LoginScreen으로 설정!
      home: LoginScreen(), 
    );
  }
}