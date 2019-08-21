package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"

    _ "github.com/go-sql-driver/mysql"
)

type Usuario struct {
    Id    int
    Nombre  string
    Apellidos string
    Dni string
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "username"
    dbPass := "password"
    dbName := "dbName"
    host := "xxxxxxx.com"
    port := "3306"
	
    db, err := sql.Open(dbDriver, dbUser + ":" + dbPass +"@tcp("+ host +":"+ port +")/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM Usuario ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    emp := Usuario{}
    res := []Usuario{}
    for selDB.Next() {
        var id int
        var nombre, apellidos, dni string
        err = selDB.Scan(&id, &nombre, &apellidos, &dni)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Nombre = nombre
        emp.Apellidos = apellidos
        emp.Dni = dni
        res = append(res, emp)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Usuario WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Usuario{}
    for selDB.Next() {
        var id int
        var nombre, apellidos, dni string
        err = selDB.Scan(&id, &nombre, &apellidos, &dni)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Nombre = nombre
        emp.Apellidos = apellidos
        emp.Dni = dni
    }
    tmpl.ExecuteTemplate(w, "Show", emp)
    defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Usuario WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Usuario{}
    for selDB.Next() {
        var id int
        var nombre, apellidos, dni string
        err = selDB.Scan(&id, &nombre, &apellidos, &dni)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Nombre = nombre
        emp.Apellidos = apellidos
        emp.Dni = dni
    }
    tmpl.ExecuteTemplate(w, "Edit", emp)
    defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        nombre := r.FormValue("nombre")
        apellidos := r.FormValue("apellidos")
        dni := r.FormValue("dni")
        insForm, err := db.Prepare("INSERT INTO Usuario(nombre, apellidos, dni) VALUES(?,?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(nombre, apellidos, dni)
        log.Println("INSERT: Nombre: " + nombre + " | Apellidos: " + apellidos + " | Dni: " + dni)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        nombre := r.FormValue("nombre")
        apellidos := r.FormValue("apellidos")
        dni := r.FormValue("dni")
        id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE Usuario SET nombre=?, apellidos=?, dni=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(nombre, apellidos, dni, id)
        log.Println("UPDATE: Nombre: " + nombre + " | Apellidos: " + apellidos + " | Dni: " + dni)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Usuario WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("BORRAR")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
    log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", Index)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/new", New)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}
