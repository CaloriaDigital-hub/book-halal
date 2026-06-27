import requests

book_id = "d35247f8-c3c4-4888-bef9-44f7ca86a694"
base_url = "http://localhost:8090"
pages_count = 136

# 1. Сначала удаляем старый кривой индекс
requests.delete(f"http://localhost:8000/admin/books/{book_id}/index")
print("Old index deleted.")

# 2. Формируем запрос на переиндексацию с картинками
pages = []
for page_num in range(1, pages_count + 1):
    # Go сохраняет картинки как page1.jpg или page-1.jpg?
    # В pdftoppm_processor.go было: name = entry.Name(), fmt.Sprintf("/static/books/%s/%s", bookID, name)
    # По умолчанию pdftoppm создает файлы вида: page-001.jpg, page-002.jpg и т.д.
    # Но проще пересоздать книгу через Go API.
    pass

