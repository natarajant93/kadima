package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/graphql-go/graphql"
	"database/sql"
 	_ "github.com/go-sql-driver/mysql"
)

var con *sql.DB;

type Employee struct {
	Empno   int 	`json:"EMPNO" bson:"EMPNO"`
	Ename string 	`json:"ENAME"`
	Job string 	`json:"JOB"`
	Mgr int 	`json:"MGR"`
	Salary float64 	`json:"SALARY"`
	Dept Department `json:"DEPT"`
}

type Department struct {
	Deptno int 	`json:"DEPTNO"`
	Dname string 	`json:"DNAME"`
	Loc string 	`json:"LOC"`
}

func get_employee_details(emp_no int) string {

	var employee_det Employee;
	var dept_no int;

	rows, err := con.Query("SELECT * from Employee,Department where EMPNO = ? and Employee.DEPTNO = Department.DEPTNO", emp_no);
	if err != nil {
		return "Error getting employee details from database";
	}
	for rows.Next() {
		err = rows.Scan(&employee_det.Empno, &employee_det.Ename, &employee_det.Job, &employee_det.Mgr, &employee_det.Salary, &employee_det.Dept.Deptno, &dept_no, &employee_det.Dept.Dname, &employee_det.Dept.Loc);
		if err != nil {
			return "Error getting employee details from database";
		}
	}

	employee_json, err := json.Marshal(employee_det)
    if err != nil {
        return err.Error();
    }
    
    return string(employee_json);
}

func get_list_of_employees() string {

	var employee_arr []*Employee;
	var dept_no int;

	rows, err := con.Query("SELECT * from Employee,Department where Employee.DEPTNO = Department.DEPTNO");
	if err != nil {
		return "Error getting employee details from database";
	}
	for rows.Next() {
		employee_det := new(Employee);
		err = rows.Scan(&employee_det.Empno, &employee_det.Ename, &employee_det.Job, &employee_det.Mgr, &employee_det.Salary, &employee_det.Dept.Deptno, &dept_no, &employee_det.Dept.Dname, &employee_det.Dept.Loc);
		if err != nil {
			return "Error getting employee details from database";
		}
		employee_arr = append(employee_arr, employee_det);
	}

	employee_json, err := json.Marshal(employee_arr)
    if err != nil {
        return err.Error();
    }
    
    return string(employee_json);
}

func get_list_of_employees_in_dept(dept_name string) string {

	var employee_arr []*Employee;
	var dept_no int;

	rows, err := con.Query("SELECT * from Employee,Department where Department.DNAME = ? and Employee.DEPTNO = Department.DEPTNO", dept_name);
	if err != nil {
		return "Error getting employee details";
	}
	for rows.Next() {
		employee_det := new(Employee);
		err = rows.Scan(&employee_det.Empno, &employee_det.Ename, &employee_det.Job, &employee_det.Mgr, &employee_det.Salary, &employee_det.Dept.Deptno, &dept_no, &employee_det.Dept.Dname, &employee_det.Dept.Loc);
		if err != nil {
			return "Error getting employee details from database";
		}
		employee_arr = append(employee_arr, employee_det);
	}

	employee_json, err := json.Marshal(employee_arr)
    if err != nil {
        return err.Error();
    }
    
    return string(employee_json);
}

func delete_employee(emp_no int) string {

	stmt, err := con.Prepare("DELETE FROM Employee where EMPNO=?")
    if err != nil {
    	return "Error preparing delete statement"
    }

    _, err = stmt.Exec(emp_no);
    if err != nil {
    	return "Error deleting employee details"
    }
    
    return "";
}

func create_employee(emp_no int, ename string, job string, mgr int, salary float64, dept_no int) string {

	stmt, err := con.Prepare("INSERT INTO Employee values(?,?,?,?,?,?)")
    if err != nil {
    	return "Error preparing insert statement"
    }

    _, err = stmt.Exec(emp_no, ename, job, mgr, salary, dept_no);
    if err != nil {
    	return "Error inserting new employee"
    }
    
    return "";
}

func update_employee(emp_no int, ename string, job string, mgr int, salary float64, dept_no int) string {

	stmt, err := con.Prepare("UPDATE Employee SET ENAME=?, JOB=?, MGR=?, SALARY=?, DEPTNO=? where EMPNO=?")
    if err != nil {
    	return "Error preparing update statement"
    }

    _, err = stmt.Exec(ename, job, mgr, salary, dept_no, emp_no);
    if err != nil {
    	return "Error updating employee details"
    }
    
    return "";
}

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"getEmployee": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"EMPNO": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					emp_no, isOK := p.Args["EMPNO"].(int)
					if isOK {
						result := get_employee_details(emp_no)
						return result, nil
					}
					return nil, nil
				},
			},
			"deleteEmployee": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"EMPNO": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					emp_no, isOK := p.Args["EMPNO"].(int)
					if isOK {
						result := delete_employee(emp_no)
						if result == "" {
							return "Employee record deleted", nil
						}
						return result, nil;
					}
					return nil, nil
				},
			},
			"listOfEmployeesInDept": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"DNAME": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					d_name, isOK := p.Args["DNAME"].(string)
					if isOK {
						result := get_list_of_employees_in_dept(d_name);
						return result, nil
					}
					return nil, nil
				},
			},
			"listOfAllEmployees": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					result := get_list_of_employees();
					return result, nil
				},
			},
			"createEmployee": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"EMPNO": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"ENAME": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"JOB": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"MGR": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"SALARY": &graphql.ArgumentConfig{
						Type: graphql.Float,
					},
					"DEPTNO": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					emp_no, isOK := p.Args["EMPNO"].(int)
					if !isOK {
						return "Error in argument", nil
					}
					ename, isOK := p.Args["ENAME"].(string)
					if !isOK {
						return "Error in argument", nil
					}
					job, isOK := p.Args["JOB"].(string)
					if !isOK {
						return "Error in argument", nil
					}
					mgr, isOK := p.Args["MGR"].(int)
					if !isOK {
						return "Error in argument", nil
					}
					salary, isOK := p.Args["SALARY"].(float64)
					if !isOK {
						return "Error in argument", nil
					}
					dept_no, isOK := p.Args["DEPTNO"].(int)
					if !isOK {
						return "Error in argument", nil
					}
					err := create_employee(emp_no, ename, job, mgr, salary, dept_no)
					if err == "" {
						return "Employee created", nil
					}
					return err, nil
				},
			},
			"updateEmployee": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"EMPNO": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"ENAME": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"JOB": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"MGR": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"SALARY": &graphql.ArgumentConfig{
						Type: graphql.Float,
					},
					"DEPTNO": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					emp_no, isOK := p.Args["EMPNO"].(int)
					if !isOK {
						return "Error in argument", nil
					}
					ename, isOK := p.Args["ENAME"].(string)
					if !isOK {
						return "Error in argument", nil
					}
					job, isOK := p.Args["JOB"].(string)
					if !isOK {
						return "Error in argument", nil
					}
					mgr, isOK := p.Args["MGR"].(int)
					if !isOK {
						return "Error in argument", nil
					}
					salary, isOK := p.Args["SALARY"].(float64)
					if !isOK {
						return "Error in argument", nil
					}
					dept_no, isOK := p.Args["DEPTNO"].(int)
					if !isOK {
						return "Error in argument", nil
					}
					err := update_employee(emp_no, ename, job, mgr, salary, dept_no)
					if err == "" {
						return "Employee details updated", nil
					}
					return err, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {

	con, _ = sql.Open("mysql", "db_user_name:db_password@/database_name")
	defer con.Close()
    
    http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query()["query"][0], schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")

	fmt.Println(`Test getEmployee with Get Request: curl -g 'http://localhost:8080/graphql?query={getEmployee(EMPNO:99)}'`)
	fmt.Println(`Test createEmployee with Get Request: curl -g 'http://localhost:8080/graphql?query={createEmployee(EMPNO:222,ENAME:"Natarajan",JOB:"SDE1",MGR:1,SALARY:334.4,DEPTNO:1)}'`)
	fmt.Println(`Test updateEmployee with Get Request: curl -g 'http://localhost:8080/graphql?query={updateEmployee(EMPNO:222,ENAME:"Nataraj",JOB:"SDE1",MGR:1,SALARY:333,DEPTNO:1)}'`)
	fmt.Println(`Test deleteEmployee with Get Request: curl -g 'http://localhost:8080/graphql?query={deleteEmployee(EMPNO:222)}'`)
	fmt.Println(`Test listOfAllEmployees with Get Request: curl -g 'http://localhost:8080/graphql?query={listOfAllEmployees}'`)
	fmt.Println(`Test listOfEmployeesInDept with Get Request: curl -g 'http://localhost:8080/graphql?query={listOfEmployeesInDept(DNAME:"Engineering")}'`)

	http.ListenAndServe(":8080", nil)
}
