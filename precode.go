package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта

// обрабочик получения всех задач методом get
func GetTasks(w http.ResponseWriter, r *http.Request) {
	//сериализация данных из мапы задач
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//записываем тип контента в заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	//записываем статус 200 ОК
	w.WriteHeader(http.StatusOK)
	//записываем json в ответ
	w.Write(resp)
}

// обрабочик получения конкретной задачи методом get
func GetTask(w http.ResponseWriter, r *http.Request) {
	//выявляем ключ задачи для мапы из эндпоинта
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задачи с таким ID нет", http.StatusBadRequest)
		return
	}

	//сериализация данных задачи
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//записываем тип контента в заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	//записываем статус 200 ОК
	w.WriteHeader(http.StatusOK)
	//записываем json в ответ
	w.Write(resp)
}

// обрабочик добавления новой задачи методом post
func PostTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	var buf bytes.Buffer

	//считываем данные для добавления из тела запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//десереализуем данные из запроса в переменную newTask
	if err = json.Unmarshal(buf.Bytes(), &newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//проверяем нет ли уже задачи с таким ID в списке
	_, ok := tasks[newTask.ID]
	if ok {
		http.Error(w, "Задачи с таким ID уже есть в списке", http.StatusBadRequest)
		return
	}

	//добавляем новую задачу из переменной newTask в мапу tasks
	tasks[newTask.ID] = newTask

	//записываем тип контента в заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	//записываем статус 201 Created
	w.WriteHeader(http.StatusCreated)
}

// обрабочик удаления задачи методом delite
func DelTask(w http.ResponseWriter, r *http.Request) {
	//выявляем ключ задачи для мапы из эндпоинта
	id := chi.URLParam(r, "id")

	//проверка наличия задачи в списке
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задачи с таким ID нет", http.StatusBadRequest)
		return
	}

	//удаление задачи из списка задач
	delete(tasks, id)

	//записываем тип контента в заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	//записываем статус 200 ОК
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", GetTasks)
	r.Post("/tasks", PostTask)
	r.Get("/tasks/{id}", GetTask)
	r.Delete("/tasks/{id}", DelTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
