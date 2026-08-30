cmd: go version
res: go version go1.27.0 windows/amd64

cmd: mongod --version
res: mongod : The term 'mongod' is not recognized as the name of a cmdlet, function, script file, or operable program. Check the spelling of the name, or if a
path was included, verify that the path is correct and try again.
At line:1 char:1
+ mongod --version
+ ~~~~~~
    + CategoryInfo          : ObjectNotFound: (mongod:String) [], CommandNotFoundException
    + FullyQualifiedErrorId : CommandNotFoundException
  
This is expected as we are not running the MongoDB Server (Daemon) installed locally on the machine.


cmd: mongosh --version
res: 2.10.0

We're connected to MongoDB Atlas using the MongoDB Shell (mongosh) version 2.10.0. 
This indicates that we can interact with our MongoDB Atlas cluster through the shell, 
even though the local MongoDB server (mongod) is not running.

url:
https://cloud.mongodb.com/v2/6a8e9989123e7d1554b21677#/overview

Server details:
MONGODB_USERNAME="maheshtallada_db_user"
MONGODB_PASSWORD="<password>"
MONGODB_URI="mongodb+srv://maheshtallada_db_user:<password>@maheshlearningcluster.x2m1d78.mongodb.net"


MaheshLearningCluster
maheshtallada_db_user
<password>

connection string: mongodb+srv://maheshtallada_db_user:<password>@maheshlearningcluster.x2m1d78.mongodb.net/?appName=MaheshLearningCluster


Connect via MongoDB Shell (mongosh):

Connect in mongosh

cmd1:
mongosh "mongodb+srv://maheshtallada_db_user:<password>@maheshlearningcluster.x2m1d78.mongodb.net/"


Or cmd2:
mongosh "mongodb+srv://maheshlearningcluster.x2m1d78.mongodb.net/" --username maheshtallada_db_user


TMG1964mahesh


cmd: mkdir -p day1/go
-p is used to create parent directories as needed. In this case, it will create the directory structure "day1/go" if it does not already exist.

cmd: go mod init day1
res: go: creating new go.mod: module day1
explanation: The command `go mod init day1` initializes a new Go module in the current directory. It creates a `go.mod` file that defines the module's path as "day1". This file is used to manage dependencies and versions for the Go project.

A1:
Q1: What does go mod init create and why?
Ans: <go mod init> creates a go.mod file with module definition , which is used to manage dependencies and versions for the Go project.

Q2: Why must an executable program use package main?
Ans: An executable program must use package main because it is the entry point for the Go program. The main package must contain a main() function, which serves as the starting point of the program's execution.

Q3: What is the difference between <go run .> and <go build>?
Ans: The difference between `go run .` and `go build` is that `go run .` compiles and runs the Go program in a single step, while `go build` compiles the program and produces an executable binary file without running it. `go run .` is useful for quickly testing code, whereas `go build` is used to create a standalone executable for distribution or deployment.

A2:
Trap check:
int - 0
string - ""
bool - false


A3:
Q1 Ans: Go uses capitalized function names to indicate that they are exported and can be accessed from other packages. This is a convention in Go to differentiate between public and private functions. Functions that start with a lowercase letter are unexported and can only be accessed within the same package.
Q2 Ans: (How are Go's multiple return values different from Java?) : In Go, functions can return multiple values simultaneously, which is a unique feature compared to Java. This allows for more flexible and expressive code, especially when dealing with operations that naturally produce multiple related results.
Q3 (What does t.Fatalf do vs t.Error?): In Go's testing framework, `t.Fatalf` is used to log an error message and immediately stop the execution of the test, marking it as failed. On the other hand, `t.Error` logs an error message but allows the test to continue running, which can be useful for reporting multiple errors in a single test case.

B3:
Q1 (What is the default value of hx-swap if you don't set it?): 
Ans: The default value of `hx-swap` is "innerHTML". If you don't set it, the content of the target element will be replaced with the response from the server, effectively updating the inner HTML of that element.
Q2 (What's the difference between innerHTML and outerHTML swaps?): 
Ans: The difference between `innerHTML` and `outerHTML` swaps is that `innerHTML` replaces the content inside the target element, while `outerHTML` replaces the entire target element itself, including its opening and closing tags. Using `innerHTML` allows you to update just the content, whereas `outerHTML` allows you to replace the whole element with new content.
Q3" (Why does HTMX return HTML fragments instead of JSON here?):
Ans: HTMX returns HTML fragments instead of JSON in this context because it is designed to work with server-rendered HTML. By returning HTML fragments, HTMX can directly update parts of the web page without requiring additional client-side rendering logic. This approach simplifies the development process and allows for a more seamless integration of server-side rendering with client-side interactivity, enabling developers to create dynamic web applications with less complexity and overhead compared to using JSON and client-side rendering frameworks.