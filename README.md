This project implements a server with Go and Gin for recording and querying one's work hours.

For authorization purposes the email, password hash and token key are currently being read from a .env file in the root directory of the project, which can be certainly improved upon from a security perspective. The data is saved to an SQLite database file.
