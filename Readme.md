This project contains basic code which has the following endpoints to do CRUD operations with Postgres DB. Listed are the endpoints,
\ - GET, POST => Fetches the list of users in the db
\{id} - GET => Fetches the details of that particular user
\{id} - DELETE => delete the user from db
\{id} - PUT => update the user in the db with the given data


To run this project:
1) Install goSwagger as per this link: https://goswagger.io/install.html
2) Install golang 1.10
3) Install PostgresSQL 9.5, with username and password as "postgres", create a DB named "test", with the following table schema and some test data,

test=# \d+ userinfo
                                         Table "public.userinfo"
 Column |          Type          | Collation | Nullable | Default | Storage  | Stats target | Description
--------+------------------------+-----------+----------+---------+----------+--------------+-------------
 uid    | integer                |           |          |         | plain    |              |
 name   | character varying(100) |           |          |         | extended |              |

Insert some dummy data into the table to play around with the API

4) Once all these are set up, clone this repo using GoLand or any other editor

Run this command inside the /api/swagger folder,

swagger generate server -A user-list -f ./swagger.yml
go install ./cmd/user-list-server/

This will create the server in you GOPATH root, generally ~/go/bin
You can start the server by running that executable in the above path.

Once started you can interact with the API's using curl command

Example commands:
curl -i localhost:45949
curl -X PUT localhost:45949/3 -d "{\"name\":\"sdj\"}" -H 'Content-Type: application/json'
curl -X "DELETE" localhost:39281/2
curl -X POST localhost:39281
curl -i localhost:45949/4

Reference:
https://goswagger.io/
