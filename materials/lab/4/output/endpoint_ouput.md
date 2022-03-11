# Required Output

I used Thunder Client to generate these requests (similar to PostMan).

## Request: GET: http://localhost:8080/
* **Response**: Welcome to my assignment page!
* Status: 200 OK
* Size: 30 Bytes
* Time: 10 ms

## Request: GET: http://localhost:8080/api-status
* **Response**: API is up and running
* Status: 200 OK
* Size: 21 Bytes
* Time: 5 ms

## Request: POST: http://localhost:8080/assignment
* **Response**: Empty, as expected
* Status: 201 Created
* Size: 0 Bytes
* Time: 6 ms
* Note: This needs to use the "Form-encode" body type or you'll have a bad time. 

## Request: GET: http://localhost:8080/assignments
* **Response**: {
  "assignments": [
    {
      "id": "Mike1A",
      "Title": "Lab 4 ",
      "desc": "Some lab this guy made yesterday?",
      "points": 20
    },
    {
      "id": "Clay1A",
      "Title": "Lab 1",
      "desc": "Clay's first lab",
      "points": 25
    }
  ]
}
* Status: 200 OK
* Size: 104 Bytes
* Time: 5 ms

 ## Request: GET: http://localhost:8080/assignment/Clay1A
 * **Response**: {
  "id": "Clay1A",
  "Title": "Lab 1",
  "desc": "Clay's first lab",
  "points": 25
}
* Status: 200 OK
* Size: 70 Bytes
* Time: 1 ms

## Request: DELETE: http://localhost:8080/assignment/Clay1A
* **Response**: {
  "status": "Success"
}
* Status: 200 OK
* Size: 20 Bytes
* Time: 5 ms