Pre-req:
Create a mysql database kadima_test with the two tables Employee and Department as specified in the requirement docs.
Change the mysql user name and pasword in line 342

Usage:

Test getEmployee with Get Request: curl -g 'http://localhost:8080/graphql?query={getEmployee(EMPNO:99)}'

Test createEmployee with Get Request: curl -g 'http://localhost:8080/graphql?query={createEmployee(EMPNO:222,ENAME:"Natarajan",JOB:"SDE1",MGR:1,SALARY:334.4,DEPTNO:1)}'

Test updateEmployee with Get Request: curl -g 'http://localhost:8080/graphql?query={updateEmployee(EMPNO:222,ENAME:"Nataraj",JOB:"SDE1",MGR:1,SALARY:333,DEPTNO:1)}'

Test deleteEmployee with Get Request: curl -g 'http://localhost:8080/graphql?query={deleteEmployee(EMPNO:222)}'

Test listOfAllEmployees with Get Request: curl -g 'http://localhost:8080/graphql?query={listOfAllEmployees}'

Test listOfEmployeesInDept with Get Request: curl -g 'http://localhost:8080/graphql?query={listOfEmployeesInDept(DNAME:"Engineering")}'
