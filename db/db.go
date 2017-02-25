package db

import (
	"database/sql"
	"html/template"
	"log"
	"time"

	"github.com/djung460/Tasks/types"

	_ "github.com/mattn/go-sqlite3"
	md "github.com/shurcooL/github_flavored_markdown"
)

var db *sql.DB
var taskStatus map[string]int
var err error
var status string

func init() {
	db, err = sql.Open("sqlite3", "./tasks.db")
	taskStatus = map[string]int{"COMPLETE": 1, "PENDING": 2, "DELETED": 3}
	if err != nil {
		log.Fatal(err)
	}
}

// Close database
func Close() {
	db.Close()
}

// GetPendingTasks returns the pending tasks of the username passed as the argument
func GetPendingTasks(username string) (types.Context, error) {
	log.Println("getting tasks for ", status)
	var tasks []types.Task
	var categories []types.Category
	var task types.Task

	var TaskCreated time.Time
	var context types.Context
	var rows *sql.Rows

	basicSQL := "select t.id, title, content, created_date, priority, case when c.name is null then 'NA' else c.name end from task t, status s, user u left outer join  category c on c.id=t.cat_id where u.username=? and s.id=t.task_status_id and u.id=t.user_id and s.status='PENDING' and t.hide!=1 order by t.created_date asc"

	rows, err := db.Query(basicSQL, "david")
	if err != nil {
		log.Println("something went wrong in fetching tasks")
		return types.Context{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Title, &task.Content, &TaskCreated, &task.Priority, &task.Category)
		if err != nil {
			log.Fatal(err)
		}
		task.HTMLContent = template.HTML(md.Markdown([]byte(task.Content)))
		TaskCreated = TaskCreated.Local()
		task.Created = TaskCreated.Format("Jan 2 2006")
		tasks = append(tasks, task)
	}

	categories = GetCategories(username)

	context = types.Context{Tasks: tasks, CSRFToken: "supersecret", Categories: categories}
	return context, nil
}

//GetCategories will return the list of categories to be
//rendered in the template
func GetCategories(username string) []types.Category {
	stmt := "select name from category c, user u where u.id = c.user_id and u.username=?"
	rows, _ := db.Query(stmt, username)
	var categories []types.Category
	var category types.Category

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&category.Name)
		if err != nil {
			log.Println(err)
		}
		categories = append(categories, category)
	}
	return categories
}

//taskQuery encapsulates running multiple queries which don't do much things
func taskQuery(sql string, args ...interface{}) error {
	log.Print("inside task query")
	SQL, err := db.Prepare(sql)
	if err != nil {
		log.Println("error in preparing query")
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		log.Println("error starting transaction")
		return err
	}
	_, err = tx.Stmt(SQL).Exec(args...)
	if err != nil {
		log.Println("taskQuery: ", err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Commit successful")
	}
	return err
}

//CompleteTask  is used to mark tasks as complete
func CompleteTask(username string, id int) error {
	err := taskQuery("update task set task_status_id=?, finish_date=datetime(),last_modified_at=datetime() where id=? and user_id=(select id from user where username=?) ", taskStatus["COMPLETE"], id, username)
	return err
}

// AddTask adds a task to the database of the currently logged in user
// AddTask is used to add the task in the database
func AddTask(username string, task types.Task) error {
	log.Println("AddTask: started function")
	var err error

	userID, err := GetUserID(username)
	if err != nil && (task.Title != "" || task.Content != "") {
		return err
	}

	err = taskQuery("insert into task(title, content, priority, task_status_id, created_date, last_modified_at, user_id,hide) values(?,?,?,?,datetime(), datetime(),?,?)", task.Title, task.Content, task.Priority, taskStatus["PENDING"], userID, task.Hidden)
	return err
}
