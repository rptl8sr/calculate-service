# Calculate Service

Service to parse and calculate simple arithmetic expressions.

## Prerequisites
- Installed Go version 1.23 or higher
- curl, Postman, or any similar app to work with HTTP API   

## Installation
1. Clone the repository to your machine `git clone github.com/rptl8sr/calculate-service`
2. Navigate into the project's directory `cd calculate-service`
3. Download all dependencies `go mod download`
4. [MacOS | LINUX] Run app `go run ./cmd/server/main.go`
   
   [Windows] Run app `go run .\cmd\server\main.go`
   
   Also, see the makefile for additional commands.
5. Make http requests
 
## How to use

**Possibilities**
- addition `+`, subtract `-`, multiplicity `*`, divide `/` operations
- any complex nested parentheses with `(` and `)`
- int and float numbers (I hope within the range -1e308..1e308) with `.` as decimal separator ()
- unary minus `-` (regular minus sign) for numbers and parentheses group's

**Environments**

App use next environments with the next default values
- PORT=8080 
- API_VERSION=v1 
- APP_VERSION=v1.0.1 
- APP_NAME=Calculate 
- APP_MODE=production
- LOG_LEVEL=info

But you can make `.env` file in root project's folder to change it.

**Request**

Make POST request to endpoint with the next payload (content-type: application/json):

`{"expression": "arithmetic expression"}`


The result will be an HTTP response with the body (content-type: application/json):

*Successfully*

`{"result": "expression's result as float with defaults to 6 decimal places"}` with `200 OK` HTTP Status Code


*Errors* 

`{"error": "error description"}` with `400 Bad Request`, `422 Unprocessable Entity`, `500 Internal Server Error` HTTP Status Codes


## Examples 


---
**Good requests #1**

Request
`curl -X POST 'localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' -d '{"expression": "2+2"}'`

Response:
`{"result":"4.000000"}`

**Good requests #2**

Request
`curl -X POST 'localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' -d '{"expression": "-2+2--(3+1)"}'`

Body:
`{"result":"4.000000"}`

**Good requests #3**

Request
`curl -X POST 'localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' -d '{"expression": "-(2+3)--(3+1)"}'`

Body:
`{"result":"-1.000000"}`

**Good requests #4**

Request
`curl -X POST 'localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' -d '{"expression": "-(3+1.1342526788908909457457476)*-6.72342534637"}'`

Body:
`{"result":"27.796339"}`

---

---
**Bad requests #1**

Request
`curl -i -X POST 'localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' -d '{"other_field": ""}'`

Headers:
`HTTP/1.1 400 Bad Request
Content-Type: application/json`

Body:
`{"error":"'expression' field is required."}`

**Bad requests #2**

Request
`curl -i -X POST 'localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' -d '{"expression": "(1+2))*3"}'`

Headers:
`HTTP/1.1 422 Unprocessable Entity
Content-Type: text/plain; charset=utf-8`

Body:
`"error":"request error: mismatched parentheses: position 5: )"}`

**Bad requests #3**

Request
`curl -i -X POST 'localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' -d '{"expression": "1+not_a_number"}'`

Headers:
`HTTP/1.1 422 Unprocessable Entity
Content-Type: text/plain; charset=utf-8`

Body:
`{"error":"request error: invalid character: position 2: n"}`

---

## License
1. This project is licensed under the terms of the MIT license. This means that you are free to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software without restriction, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

2. THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE