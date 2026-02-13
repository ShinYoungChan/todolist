import 'package:flutter/material.dart';
import '../services/api_service.dart';

class TodoListScreen extends StatefulWidget {
  const TodoListScreen({super.key});

  @override
  State<TodoListScreen> createState() => _TodoListScreenState();
}

class _TodoListScreenState extends State<TodoListScreen> {
  final ApiService _apiService = ApiService();

  void _showAddTodoDialog() {
    // 입력값을 제어할 컨트롤러 두 개 생성
    final TextEditingController _titleController = TextEditingController();
    final TextEditingController _contentController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text("새 할 일 추가"),
        content: Column(
          mainAxisSize: MainAxisSize.min, // 내용물 크기만큼만 다이얼로그 크기 조절
          children: [
            TextField(
              controller: _titleController,
              decoration: const InputDecoration(hintText: "제목을 입력하세요"),
              autofocus: true,
            ),
            const SizedBox(height: 10), // 칸 사이 간격
            TextField(
              controller: _contentController,
              decoration: const InputDecoration(hintText: "설명(내용)을 입력하세요"),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text("취소"),
          ),
          ElevatedButton(
            onPressed: () async {
              if (_titleController.text.trim().isEmpty) return;

              try {
                // ApiService 호출 시 이름 있는 인자로 전달
                await _apiService.createTodo(
                  title: _titleController.text,
                  content: _contentController.text,
                );

                if (mounted) {
                  Navigator.pop(context); // 팝업 닫기
                  setState(() {}); // 화면 새로고침
                }
              } catch (e) {
                print("추가 실패: $e");
              }
            },
            child: const Text("추가"),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text("My Todo List"),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () {
              // TODO: 로그아웃 로직 (토큰 삭제 후 로그인 화면으로 이동)
            },
          ),
        ],
      ),
      // 💡 FutureBuilder: 비동기 데이터(Future)를 UI로 변환해주는 도구
      body: FutureBuilder(
        future: _apiService.getTodos(), // 여기서 getTodos 호출!
        builder: (context, snapshot) {
          // 1. 데이터를 기다리는 중일 때
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }

          // 2. 에러가 발생했을 때 (인터셉터에서 401 등이 났을 때)
          if (snapshot.hasError) {
            return Center(child: Text("에러 발생: ${snapshot.error}"));
          }

          // 3. 데이터가 성공적으로 왔을 때
          if (snapshot.hasData) {
            final List todos =
                snapshot.data!.data['todos']; // Dio response의 data 추출

            if (todos.isEmpty) {
              return const Center(child: Text("할 일이 없어요! 추가해볼까요?"));
            }

            return ListView.builder(
              itemCount: todos.length,
              itemBuilder: (context, index) {
                final todo = todos[index];
                //final bool isDone = todo['status'] ?? false; // 완료 여부 변수화
                return ListTile(
                  leading: Checkbox(
                    value: todo['status'] ?? false,
                    onChanged: (bool? newValue) async {
                      // TODO: 업데이트 API 호출
                      if (newValue == null) return;
                      // 💡 현재 값이 뭔지, 그리고 바꾸려는 값이 뭔지 둘 다 찍어보세요.
                      try {
                        await _apiService.updateTodoState(todo['id'], newValue);
                        if (mounted) {
                          setState(() {
                            todo['status'] = newValue;
                          });
                        }
                      } catch (e) {
                        print("업데이트 에러: $e");
                      }
                    },
                  ),
                  title: Text(
                    todo['title'] ?? '제목 없음',
                    style: TextStyle(
                      decoration: todo['status'] == true
                          ? TextDecoration.lineThrough
                          : null,
                    ),
                  ),
                  trailing: IconButton(
                    icon: const Icon(Icons.delete, color: Colors.red),
                    onPressed: () async {
                      try {
                        await _apiService.deleteTodo(todo['id']);

                        if (mounted) {
                          // 삭제 성공시 화면 새로고침
                          setState(() {});
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text("삭제되었습니다.")),
                          );
                        }
                      } catch (e) {
                        print("삭제 에러: $e");
                      }
                    },
                  ),
                );
              },
            );
          }

          return const Center(child: Text("데이터가 없습니다."));
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showAddTodoDialog,
        child: const Icon(Icons.add),
      ),
    );
  }
}
