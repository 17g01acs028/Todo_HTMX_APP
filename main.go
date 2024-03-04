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
	{ID: "todo-12763dc3-46e6-4abe-aa9a-6be1e35c8584", Title: "Code and Chil"},
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

	htmlStr := fmt.Sprintf(`<li id="todo-%s">%s <button  class="deletebtn" hx-delete="/todos/delete/%s" hx-target="#todo-%s" hx-swap="outerHTML" hx-confirm="Are you sure you want to delete this TODO?">Delete</button></li>`, id.String(), title, id.String(), id.String())

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

func main() {
	//Load css to the template
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))

	//Routes handlers
	http.HandleFunc("/", todos)
	http.HandleFunc("/todo", todo)
	http.HandleFunc("/todos/delete/", delTodo)

	// Start server on port 8000
	fmt.Println("Started server started on port 8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		panic(err)
	}
}
