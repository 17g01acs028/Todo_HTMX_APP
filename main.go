package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"

	"github.com/google/uuid"
)

type Todos struct {
	Title string
	ID    string
}

var todoList = []Todos{
	{ID: "1", Title: "Hit the gym"},
	{ID: "2", Title: "Pay bills"},
	{ID: "3", Title: "Meet George"},
	{ID: "4", Title: "Buy eggs"},
	{ID: "5", Title: "Read a book"},
	{ID: "6", Title: "Organize office"},
	{ID: "7", Title: "Beef Meet"},
	{ID: "12763dc3-46e6-4abe-aa9a-6be1e35c8584", Title: "Code and Chil"},
}

func todo(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println("Error parsing body:", err)
		return
	}

	// Extract the title parameter from the map
	id := uuid.New()
	fmt.Println(id.String())
	title := values.Get("title")

	newTodos := []Todos{
		{ID: id.String(), Title: title},
	}

	todoList = append(todoList, newTodos...)

	//htmlStr := fmt.Sprintf(`<li id="todo-%s">%s <button  class="deletebtn" hx-delete="/todos/delete/%s" hx-target="#todo-%s" hx-swap="outerHTML" hx-confirm="Are you sure you want to delete this TODO?">Delete</button></li>`, id.String(), title, id.String(), id.String())
	htmlStr := fmt.Sprintf(`<li id="todo-%s">%s<div class="w3-bar" style="width:fit-content;">
	<button class="w3-bar-item w3-button w3-black" hx-get="/todos/view/%s" hx-target="body" hx-swap="innerHTML" hx-push-url="/todos/view/%s">View</button>
	<button class="w3-bar-item w3-button w3-teal" hx-get="/todos/update/view/%s" hx-target="body" hx-swap="innerHTML" hx-push-url="/todos/update/view/%s">Edit</button>
	<button class="w3-bar-item w3-button w3-red" hx-delete="/todos/delete/%s" hx-target="#todo-%s" hx-swap="outerHTML" hx-confirm="Are you sure you want to delete this TODO?">Delete</button>
  </div> </li>`, id, title, id, id, id, id, id, id)
	temp, _ := template.New("t").Parse(htmlStr)
	temp.Execute(w, nil)

}

func todos(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("static/index.html"))

	tmp.Execute(w, todoList)
}

func delTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/delete/"):]

	for i, todo := range todoList {
		if todo.ID == idStr {
			todoList = append(todoList[:i], todoList[i+1:]...)
			temp, _ := template.New("t").Parse("")
			temp.Execute(w, nil)
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}
func updateTodo(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		fmt.Println("Error parsing body:", err)
		return
	}

	// Extract the title parameter from the map
	title := values.Get("title")
	idStr := r.URL.Path[len("/todos/update/"):]

	getUpdatedTodoListByID := func(todoList []Todos, id string, newTitle string) []Todos {
		for i := range todoList {
			if todoList[i].ID == id {
				todoList[i].Title = newTitle
				break
			}
		}
		return todoList
	}

	updatedTodoList := getUpdatedTodoListByID(todoList, idStr, title)
	tmp := template.Must(template.ParseFiles("static/index.html"))

	tmp.Execute(w, updatedTodoList)
}

func getUpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/update/view/"):]

	getTodoSliceByID := func(id string) []Todos {
		var todoSlice []Todos
		for _, todo := range todoList {
			if todo.ID == id {
				todoSlice = append(todoSlice, todo)
			}
		}
		return todoSlice
	}

	foundTodoSlice := getTodoSliceByID(idStr)
	if len(foundTodoSlice) > 0 {
		temp := template.Must(template.ParseFiles("static/edit.html"))
		temp.Execute(w, foundTodoSlice)
	} else {
		fmt.Println("Todo not found with ID:", idStr)
	}
}

func getViewTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/todos/view/"):]

	getTodoSliceByID := func(id string) []Todos {
		var todoSlice []Todos
		for _, todo := range todoList {
			if todo.ID == id {
				todoSlice = append(todoSlice, todo)
			}
		}
		return todoSlice
	}

	foundTodoSlice := getTodoSliceByID(idStr)
	if len(foundTodoSlice) > 0 {
		temp := template.Must(template.ParseFiles("static/view.html"))
		temp.Execute(w, foundTodoSlice)
	} else {
		fmt.Println("Todo not found with ID:", idStr)
	}
}
func main() {
	//Load css to the template
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))

	//Routes handlers
	http.HandleFunc("/", todos)
	http.HandleFunc("/todo", todo)
	http.HandleFunc("/todos/delete/", delTodo)
	http.HandleFunc("/todos/update/view/", getUpdateTodo)
	http.HandleFunc("/todos/update/", updateTodo)
	http.HandleFunc("/todos/view/", getViewTodo)

	// Start server on port 8000
	fmt.Println("Started server started on port 9000")
	err := http.ListenAndServe(":9000", nil)

	if err != nil {
		panic(err)
	}
}
