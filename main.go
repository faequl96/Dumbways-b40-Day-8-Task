package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

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

	fmt.Println("Server running at localhost port 8000")

	http.ListenAndServe("localhost:8000", route)
}

type Form struct {
	ProjectName    string
	StartDate      string
	EndDate        string
	Description    string
	Node           string
	React          string
	Vue            string
	TypeScript     string
	NodeIcon       string
	ReactIcon      string
	VueIcon        string
	TypeScriptIcon string
	Id             int
	Duration       string
}

var dataForm = []Form{
	{
		ProjectName:    "Dummy Project 1",
		StartDate:      "2022-09-12",
		EndDate:        "2022-09-19",
		Duration:       "1 Weeks",
		Description:    "Description Dummy Project 1",
		Node:           "Node JS",
		React:          "React JS",
		Vue:            "Vue JS",
		TypeScript:     "TypeScript",
		NodeIcon:       "../public/img/node.png",
		ReactIcon:      "../public/img/react.png",
		VueIcon:        "../public/img/vue.png",
		TypeScriptIcon: "../public/img/typescript.png",
		// Id:          0,
	},
	{
		ProjectName:    "Dummy Project 2",
		StartDate:      "2022-09-20",
		EndDate:        "2022-09-25",
		Duration:       "5 Days",
		Description:    "Description Dummy Project 2",
		Node:           "Node JS",
		TypeScript:     "TypeScript",
		NodeIcon:       "../public/img/node.png",
		ReactIcon:      "",
		VueIcon:        "",
		TypeScriptIcon: "../public/img/typescript.png",
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
	node := r.PostForm.Get("nodeJs")
	react := r.PostForm.Get("reactJs")
	vue := r.PostForm.Get("vueJs")
	typeScript := r.PostForm.Get("typeScript")
	nodeIcon := ""
	if node == "Node JS" {
		nodeIcon = "../public/img/node.png"
	}
	reactIcon := ""
	if react == "React JS" {
		reactIcon = "../public/img/react.png"
	}
	vueIcon := ""
	if vue == "Vue JS" {
		vueIcon = "../public/img/vue.png"
	}
	typeScriptIcon := ""
	if typeScript == "TypeScript" {
		typeScriptIcon = "../public/img/typescript.png"
	}
	fmt.Println(startDate)
	layout := "2006-01-02"
	startDateParse, _ := time.Parse(layout, startDate)
	endDateParse, _ := time.Parse(layout, endDate)

	hour := 1
	day := hour * 24
	week := hour * 24 * 7
	month := hour * 24 * 30
	year := hour * 24 * 365

	differHour := endDateParse.Sub(startDateParse).Hours()
	var differHours int = int(differHour)
	// fmt.Println(differHours)
	days := differHours / day
	weeks := differHours / week
	months := differHours / month
	years := differHours / year

	var duration string

	if differHours < week {
		duration = strconv.Itoa(int(days)) + " Days"
	} else if differHours < month {
		duration = strconv.Itoa(int(weeks)) + " Weeks"
	} else if differHours < year {
		duration = strconv.Itoa(int(months)) + " Months"
	} else if differHours > year {
		duration = strconv.Itoa(int(years)) + " Years"
	}

	addNewDataForm := Form{
		ProjectName:    projectName,
		StartDate:      startDate,
		EndDate:        endDate,
		Duration:       duration,
		Description:    description,
		Node:           node,
		React:          react,
		Vue:            vue,
		TypeScript:     typeScript,
		NodeIcon:       nodeIcon,
		ReactIcon:      reactIcon,
		VueIcon:        vueIcon,
		TypeScriptIcon: typeScriptIcon,
		// Id:          len(dataForm),
		// Duration:    time.Now().String(),
	}

	dataForm = append(dataForm, addNewDataForm)

	// fmt.Println(dataForm)

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
				ProjectName:    data.ProjectName,
				StartDate:      data.StartDate,
				EndDate:        data.EndDate,
				Duration:       data.Duration,
				Description:    data.Description,
				Node:           data.Node,
				React:          data.React,
				Vue:            data.Vue,
				TypeScript:     data.TypeScript,
				NodeIcon:       data.NodeIcon,
				ReactIcon:      data.ReactIcon,
				VueIcon:        data.VueIcon,
				TypeScriptIcon: data.TypeScriptIcon,
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
	node := r.PostForm.Get("nodeJs")
	react := r.PostForm.Get("reactJs")
	vue := r.PostForm.Get("vueJs")
	typeScript := r.PostForm.Get("typeScript")
	fmt.Println(node)
	nodeIcon := ""
	if node == "Node JS" {
		nodeIcon = "../public/img/node.png"
	}
	reactIcon := ""
	if react == "React JS" {
		reactIcon = "../public/img/react.png"
	}
	vueIcon := ""
	if vue == "Vue JS" {
		vueIcon = "../public/img/vue.png"
	}
	typeScriptIcon := ""
	if typeScript == "TypeScript" {
		typeScriptIcon = "../public/img/typescript.png"
	}

	layout := "2006-01-02"
	startDateParse, _ := time.Parse(layout, startDate)
	endDateParse, _ := time.Parse(layout, endDate)

	hour := 1
	day := hour * 24
	week := hour * 24 * 7
	month := hour * 24 * 30
	year := hour * 24 * 365

	differHour := endDateParse.Sub(startDateParse).Hours()
	var differHours int = int(differHour)
	// fmt.Println(differHours)
	days := differHours / day
	weeks := differHours / week
	months := differHours / month
	years := differHours / year

	var duration string

	if differHours < week {
		duration = strconv.Itoa(int(days)) + " Days"
	} else if differHours < month {
		duration = strconv.Itoa(int(weeks)) + " Weeks"
	} else if differHours < year {
		duration = strconv.Itoa(int(months)) + " Months"
	} else if differHours > year {
		duration = strconv.Itoa(int(years)) + " Years"
	}

	editDataForm := Form{
		ProjectName:    projectName,
		StartDate:      startDate,
		EndDate:        endDate,
		Duration:       duration,
		Description:    description,
		Node:           node,
		React:          react,
		Vue:            vue,
		TypeScript:     typeScript,
		NodeIcon:       nodeIcon,
		ReactIcon:      reactIcon,
		VueIcon:        vueIcon,
		TypeScriptIcon: typeScriptIcon,
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
