Feature: Get a associative table
  I need to be able to get a associative table
  Only admin should be able to get a associative table

  Scenario: given get request should return associative table
    When I am admin
    When I do GET /api/v4/associativetable?name=test-associativetable-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "name": "test-associativetable-to-get",
      "content": [
        {
          "title": "test-associativetable-to-get-content-nested-val-1",
          "names": [
            "test-associativetable-to-get-content-nested-val-2",
            "test-associativetable-to-get-content-nested-val-3"
          ],
          "_id": "test-associativetable-to-get-content-nested-val-4"
        }
      ]
    }
    """

  Scenario: given get not exist request should return empty associative table
    When I am admin
    When I do GET /api/v4/associativetable?name=test-associativetable-not-exist
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "name": "test-associativetable-not-exist",
      "content": null
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/associativetable
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/associativetable
    Then the response code should be 403
