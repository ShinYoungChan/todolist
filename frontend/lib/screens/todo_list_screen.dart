import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import '../services/api_service.dart';

class TodoListScreen extends StatefulWidget {
  const TodoListScreen({super.key});

  @override
  State<TodoListScreen> createState() => _TodoListScreenState();
}

class _TodoListScreenState extends State<TodoListScreen> {
  final ApiService _apiService = ApiService();
  String _currentSort = "created_at";
  String _filter = "all";
  String _keyword = "";

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

  void _showDatePicker(BuildContext context, Map todo, int todoId, String field, DateTime initialDate) {
    DateTime selectedDate = initialDate;

    showCupertinoModalPopup(
      context: context,
      builder: (_) => Container(
        height: 300,
        color: Colors.white,
        child: Column(
          children: [
            // 상단 완료 버튼 바
            Container(
              height: 50,
              color: const Color(0xfff8f8f8),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  CupertinoButton(
                    child: const Text('완료'),
                    onPressed: () async {
                      // 검증 로직 추가
                      if (field == 'due') {
                        DateTime start = DateTime.parse(todo['start_date']);
                        if (selectedDate.isBefore(start)) {
                            // 경고창 띄우기
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text("마감일은 시작일보다 빠를 수 없습니다."))
                            );
                            return; // 함수 종료 (서버에 안 보냄)
                        }
                      }
                      // API 호출: 선택된 날짜를 ISO8601 포맷 문자열로 전송
                      await _apiService.updateTodoDates(
                        todoId,
                        startDate: field == 'start' ? selectedDate : null, // DateTime 객체 그대로 전달
                        dueDate: field == 'due' ? selectedDate : null,
                      );
                      Navigator.pop(context);
                      setState(() {}); // 리스트 새로고침
                    },
                  ),
                ],
              ),
            ),
            // 피커 본체
            Expanded(
              child: CupertinoDatePicker(
                mode: CupertinoDatePickerMode.date,
                initialDateTime: initialDate,
                onDateTimeChanged: (DateTime newDate) {
                  selectedDate = newDate; // 휠을 돌릴 때마다 변수에 저장
                },
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFilterChip(String label, String value) {
    return ChoiceChip(
      label: Text(label, style: const TextStyle(fontSize: 12)),
      selected: _filter == value,
      showCheckmark: false,
      selectedColor: Colors.blueAccent.withOpacity(0.2), // 필터는 정렬과 색상을 다르게 하면 구분하기 좋습니다.
      onSelected: (bool selected) {
        if (selected) {
          setState(() {
            _filter = value;
          });
        }
      },
      labelPadding: const EdgeInsets.symmetric(horizontal: 4),
      visualDensity: const VisualDensity(horizontal: -2, vertical: -2),
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
              // TODO: 로그아웃 로직
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // 1. 상단 정렬 버튼 영역
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
            child: TextField(
              decoration: InputDecoration(
                hintText: "할 일 검색...",
                prefixIcon: const Icon(Icons.search),
                contentPadding: const EdgeInsets.symmetric(vertical: 0),
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
                fillColor: Colors.grey[100],
                filled: true,
              ),
              onChanged: (value) {
                setState(() {
                  _keyword = value; // 키워드 저장 후 리스트 새로고침
                });
              },
            ),
          ),

          // 💡 2. 필터 칩 영역 (전체 / 진행중 / 완료)
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                _buildFilterChip("전체", "all"),
                const SizedBox(width: 8),
                _buildFilterChip("진행중", "pending"),
                const SizedBox(width: 8),
                _buildFilterChip("완료", "completed"),
              ],
            ),
          ),

          // 💡 3. 기존 정렬 버튼 영역 (한 줄 띄우기 위해 Divider 추가 가능)
          const Divider(height: 1),
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                _buildSortChip("최신순", "created_at"),
                const SizedBox(width: 8),
                _buildSortChip("시작일순", "start_date"),
                const SizedBox(width: 8),
                _buildSortChip("마감일순", "due_date"),
              ],
            ),
          ),

          // 2. 리스트 영역 (Expanded로 감싸야 Column 안에서 정상 작동합니다)
          Expanded(
            child: FutureBuilder(
              future: _apiService.getTodos(_currentSort, _filter, _keyword),
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                }

                if (snapshot.hasError) {
                  return Center(child: Text("에러 발생: ${snapshot.error}"));
                }

                if (snapshot.hasData) {
                  final List todos = snapshot.data!.data['todos'];

                  if (todos.isEmpty) {
                    return const Center(child: Text("할 일이 없어요! 추가해볼까요?"));
                  }

                  return ListView.builder(
                    itemCount: todos.length,
                    itemBuilder: (context, index) {
                      final todo = todos[index];
                      final bool isDone = todo['status'] ?? false;

                      return Opacity(
                        opacity: isDone ? 0.5 : 1.0,
                        child: Card(
                          margin: const EdgeInsets.symmetric(
                            vertical: 6,
                            horizontal: 16,
                          ),
                          elevation: isDone ? 0 : 2,
                          child: ListTile(
                            leading: CupertinoSwitch(
                              value: isDone,
                              activeColor: CupertinoColors.activeGreen,
                              onChanged: (bool newValue) async {
                                await _apiService.updateTodoState(
                                  todo['id'],
                                  newValue,
                                );
                                if (mounted) {
                                  setState(() {
                                    todo['status'] = newValue;
                                  });
                                }
                              },
                            ),
                            title: Text(
                              todo['title'] ?? '',
                              style: TextStyle(
                                fontWeight: isDone
                                    ? FontWeight.normal
                                    : FontWeight.bold,
                                color: isDone ? Colors.grey : Colors.black87,
                              ),
                            ),
                            subtitle: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const SizedBox(height: 4),
                                // 1. 할 일 내용
                                Text(
                                  todo['content'] ?? '내용 없음',
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                  style: const TextStyle(
                                    fontSize: 13,
                                    color: Colors.black54,
                                  ),
                                ),
                                const SizedBox(height: 6),
                                // 2. 날짜 영역 (아이콘 + 시작일 ~ 마감일)
                                Row(
                                  children: [
                                    const Icon(Icons.calendar_month, size: 14, color: Colors.grey),
                                    const SizedBox(width: 4),
                                    // 시작일 클릭
                                    GestureDetector(
                                      onTap: () => _showDatePicker(
                                        context, 
                                        todo,
                                        todo['id'], 
                                        'start', 
                                        DateTime.parse(todo['start_date']),
                                      ),
                                      child: Text(
                                        todo['start_date']?.substring(0, 10) ?? '시작일',
                                        style: const TextStyle(fontSize: 11, color: Colors.blue, decoration: TextDecoration.underline),
                                      ),
                                    ),
                                    const Text(" ~ ", style: TextStyle(fontSize: 11)),
                                    // 마감일 클릭
                                    GestureDetector(
                                      behavior: HitTestBehavior.opaque, // 👈 이거 추가! 빈 공간을 눌러도 반응하게 합니다.
                                      onTap: () => _showDatePicker(
                                        context, 
                                        todo,
                                        todo['id'], 
                                        'due', 
                                        DateTime.parse(todo['due_date'] ?? DateTime.now().toString()),
                                      ),
                                      child: Padding(
                                        padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2), // 👈 터치 영역 확보
                                        child: Text(
                                          todo['due_date']?.substring(0, 10) ?? '마감일',
                                          style: const TextStyle(fontSize: 11, color: Colors.redAccent),
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ],
                            ),
                            trailing: IconButton(
                              icon: const Icon(
                                Icons.delete_outline,
                                color: Colors.redAccent,
                              ),
                              onPressed: () async {
                                try {
                                  await _apiService.deleteTodo(todo['id']);
                                  if (mounted) {
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
                          ),
                        ),
                      );
                    },
                  );
                }
                return const Center(child: Text("데이터가 없습니다."));
              },
            ),
          ),
        ], // Column의 children 끝
      ), // Column 끝
      floatingActionButton: FloatingActionButton(
        onPressed: _showAddTodoDialog,
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildSortChip(String label, String value) {
    return ChoiceChip(
      label: Text(label, style: const TextStyle(fontSize: 12)), // 텍스트 크기 살짝 줄임
      selected: _currentSort == value,
      showCheckmark: false, // 💡 1. 체크 표시(V) 안 보이게 설정
      selectedColor: const Color.fromARGB(255, 128, 128, 128), // 선택됐을 때 색상
      onSelected: (bool selected) {
        if (selected) {
          setState(() {
            _currentSort = value;
          });
        }
      },
      // 💡 2. 칩 내부의 여백을 줄여서 짤림 방지
      labelPadding: const EdgeInsets.symmetric(horizontal: 4),
      visualDensity: const VisualDensity(horizontal: -2, vertical: -2),
    );
  }
}
