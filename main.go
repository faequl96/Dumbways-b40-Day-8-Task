package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	route.PathPrefix("/public").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")

	route.HandleFunc("/project", myProject).Methods("GET")
	route.HandleFunc("/project/{index}", myProjectDetail).Methods("GET")
	route.HandleFunc("/form-project", myProjectForm).Methods("GET")
	route.HandleFunc("/add-project", myProjectData).Methods("POST")
	route.HandleFunc("/form-edit-project/{index}", myProjectFormEditProject).Methods("GET")
	route.HandleFunc("/edit-project/{id}", myProjectEdited).Methods("POST")
	route.HandleFunc("/delete-project/{index}", myProjectDelete).Methods("GET")

	route.HandleFunc("/contact", contact).Methods(("GET"))

	fmt.Println("Server running at localhost port 5000")

	http.ListenAndServe("localhost:5000", route)
}

type Form struct {
	ProjectName string
	StartDate   string
	EndDate     string
	Description string
	Id          int
	Duration    string
}

var dataForm = []Form{
	{
		ProjectName: "Project Name",
		StartDate:   "2022-09-12",
		EndDate:     "2022-09-19",
		Description: "Description",
		// Id:          0,
	},
	{
		ProjectName: "Project Name",
		StartDate:   "2022-09-20",
		EndDate:     "2022-09-25",
		Description: "Description",
		// Id:          1,
	},
}

func home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("views/index.html")

	if err != nil {
		panic(err)
	}

	template.Execute(w, nil)
}

// 3
// menampilkan halaman myProject.html
// setelah diredirect oleh func myProjectData dan func myProjectEdited, func ini akan menampilkan halaman myProject.html serta mengisikan data yang telah didapat dari func myProjectData dan mengisikan index untuk route projectdetail, route edit, dan route delete
func myProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/myProject.html")

	response := map[string]interface{}{
		"Projects": dataForm,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		w.Write([]byte("Message: "))
		w.Write([]byte(err.Error()))
	}
	// w.WriteHeader(http.StatusOK)
}

// 1
// menampilkan myProjectForm.html
// ketika menekan tombol create, func ini akan memanggil route /add-project yang berisi func myProjectData
func myProjectForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectForm.html")

	if err == nil {
		tmpl.Execute(w, nil)
	} else {
		panic(err)
	}
}

// 2
// mengisikan arraydata dengan data yang telah diinput di form
// setelah mengisi arraydata kemudian akan meredirect ke route /project yang berisi func myProject
func myProjectData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	projectName := r.PostForm.Get("projectName")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	description := r.PostForm.Get("description")

	addNewDataForm := Form{
		ProjectName: projectName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		// Id:          len(dataForm),
		// Duration:    time.Now().String(),
	}

	dataForm = append(dataForm, addNewDataForm)

	fmt.Println(dataForm)

	http.Redirect(w, r, "/project", http.StatusMovedPermanently)

}

func myProjectDetail(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectDetail.html")

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	// ProjectDetail := Form{
	// 	ProjectName: dataForm[0].ProjectName,
	// 	StartDate:   dataForm[0].StartDate,
	// 	EndDate:     dataForm[0].EndDate,
	// 	Description: dataForm[0].Description,
	// }

	ProjectDetail := Form{}

	for i, data := range dataForm {
		if index == i {
			ProjectDetail = Form{
				ProjectName: data.ProjectName,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
			}
		}
	}

	response := map[string]interface{}{
		"Project": ProjectDetail,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}

}

func myProjectDelete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	// fmt.Println(index)
	// fmt.Println(dataForm[index+1:])

	dataForm = append(dataForm[:index], dataForm[index+1:]...)

	// fmt.Println(dataForm)

	http.Redirect(w, r, "/project", http.StatusFound)
}

// 1
// func ini dipanggil oleh tombol edit project yang ada di masing masing card sesuai indexnya
// menampilkan halaman myProjectFormEditProject.html sekaligus mengisikan value masing-masing field
// ketika menekan tombol save, func ini akan menanamkan Id(query params) pada route yang akan dipanggil yaitum route /edit-project/{{Id}} yang berisi func myProjectEdited
// Id(query params) digunakan untuk mengedit arraydata index yang ditarget
func myProjectFormEditProject(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectFormEditProject.html")

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	ProjectEdit := Form{}

	for i, data := range dataForm {
		if index == i {
			ProjectEdit = Form{
				ProjectName: data.ProjectName,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
				Id:          index,
			}
			fmt.Println(ProjectEdit)
		}
	}

	response := map[string]interface{}{
		"Project": ProjectEdit,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}
}

// 2
// mengedit arraydata yang ditarget dengan data yang telah diinput di form edit
// setelah mengedit arraydata kemudian akan meredirect ke route /project yang berisi func myProject
func myProjectEdited(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	projectName := r.PostForm.Get("projectName")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	description := r.PostForm.Get("description")

	editDataForm := Form{
		ProjectName: projectName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		// Id:          id,
		// Duration:    time.Now().String(),
	}

	dataForm[id] = editDataForm

	fmt.Println(dataForm)

	http.Redirect(w, r, "/project", http.StatusMovedPermanently)

}

func contact(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/contact.html")

	if err == nil {
		tmpl.Execute(w, nil)
	} else {
		panic(err)
	}
}
