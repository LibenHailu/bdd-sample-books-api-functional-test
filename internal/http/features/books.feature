Feature: books crud api

  In order to use  books API
  as an API book
  I need to be able to managge books

  # The first example gets an empty books
  Scenario: should get empty books
    When I send "GET" request to "/v1/books"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "books": []
      }
      """

  # The second example gets books
  Scenario: should get books
    Given there are books:
      | id | isbn          | title                                  | author       | created_at                       | updated_at                       |
      | 2  | 9780545069670 | Harry Potter and the Sorcerer's Stone  | J.K. Rowling | 2021-09-12T23:46:26.013163+03:00 | 2021-09-12T23:46:26.013163+03:00 |
      | 3  | 9780439554930 | Harry Potter and the Half-Blood Prince | J.K. Rowling | 2021-09-12T23:47:10.451301+03:00 | 2021-09-12T23:46:26.013163+03:00 |
    When I send "GET" request to "/v1/books"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "books": [
          {
            "id": 2,
            "isbn": "9780545069670",
            "title": "Harry Potter and the Sorcerer's Stone",
            "author": "J.K. Rowling",
            "created_at": "2021-09-12T23:46:26.013Z",
            "updated_at": "2021-09-12T23:46:26.013Z"
          },
          {
            "id": 3,
            "isbn": "9780439554930",
            "title": "Harry Potter and the Half-Blood Prince",
            "author": "J.K. Rowling",
            "created_at": "2021-09-12T23:47:10.451Z",
            "updated_at": "2021-09-12T23:46:26.013Z"
          }
        ]
      }
      """

  #The third  example gets a single book
  Scenario: should get the books when the book is one
    Given there are books:
      | id | isbn          | title                                 | author       | created_at                       | updated_at                       |
      | 1  | 9780545069670 | Harry Potter and the Sorcerer's Stone | J.K. Rowling | 2021-09-12T23:46:26.013163+03:00 | 2021-09-12T23:46:26.013163+03:00 |
    When I send "GET" request to "/v1/books"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "books": [
          {
            "id": 1,
            "isbn": "9780545069670",
            "title": "Harry Potter and the Sorcerer's Stone",
            "author": "J.K. Rowling",
            "created_at": "2021-09-12T23:46:26.013Z",
            "updated_at": "2021-09-12T23:46:26.013Z"
          }
        ]
      }
      """

  # The fourth example results an error when posting an empty books
  Scenario: post an empty book should get an error
    When I send "POST" request to "/v1/books"
    Then the response code should be 400
    And the response should match json:
      """
      {
        "errors": {}
      }
      """

  # The fifth example results an error when posting an empty books
  Scenario: post an missing filed in a book should get an error
    When  I have request json:
      """
      {
        "isbn": "9780439554930",
        "title": "Harry Potter and the Half-Blood Prince"
      }
      """
    And I send "POST" request to "/v1/books"
    Then the response code should be 400
    And the response should match json:
      """
      {
        "errors": {
          "author": "Author is required required"
        }
      }
      """

  # The sixth example results a books creation
  Scenario: should create a book
    When  I have request json:
      """
      {
        "author": "J.K. Rowling",
        "isbn": "9780439554930",
        "title": "Harry Potter and the Half-Blood Prince"
      }
      """
    And I send "POST" request to "/v1/books"
    Then the response code should be 201

  # The seventh example results a books update
  Scenario: should update a book
    Given there are books:
      | id | isbn          | title                                  | author       | created_at                       | updated_at                       |
      | 2  | 9780545069670 | Harry Potter and the Sorcerer's Stone  | J.K. Rowling | 2021-09-12T23:46:26.013163+03:00 | 2021-09-12T23:46:26.013163+03:00 |
      | 3  | 9780439554930 | Harry Potter and the Half-Blood Prince | J.K. Rowling | 2021-09-12T23:47:10.451301+03:00 | 2021-09-12T23:46:26.013163+03:00 |
    When  I have request json:
      """
      {
        "title": "The alchemist"
      }
      """
    And I send "PATCH" request to "/v1/books/2"
    Then the response code should be 200


  # The eighth example results a books update
  Scenario: should delete a book
    Given there are books:
      | id | isbn          | title                                  | author       | created_at                       | updated_at                       |
      | 2  | 9780545069670 | Harry Potter and the Sorcerer's Stone  | J.K. Rowling | 2021-09-12T23:46:26.013163+03:00 | 2021-09-12T23:46:26.013163+03:00 |
      | 3  | 9780439554930 | Harry Potter and the Half-Blood Prince | J.K. Rowling | 2021-09-12T23:47:10.451301+03:00 | 2021-09-12T23:46:26.013163+03:00 |
    When  I send "DELETE" request to "/v1/books/2"
    Then the response code should be 204
