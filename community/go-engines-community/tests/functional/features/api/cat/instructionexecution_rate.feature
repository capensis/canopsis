Feature: rate an instruction
  I need to be able to rate an instruction
  Only admin should be able to rate an instruction

  Scenario: Rate an instruction as unauthorized
    When I do PUT /api/v4/cat/executions/execution-for-test-instruction-to-rate-1/rate:
    """
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 401

  Scenario: Rate an instruction without permissions
    When I am noperms
    When I do PUT /api/v4/cat/executions/execution-for-test-instruction-to-rate-1/rate:
    """
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 403

  Scenario: Rate an instruction, but without completed execution
    When I am admin
    When I do PUT /api/v4/cat/executions/execution-for-test-instruction-to-rate-1/rate:
    """
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "instruction wasn't completed by the user"
    }
    """

  Scenario: Rate an instruction, when the user has already rated it
    When I am admin
    When I do PUT /api/v4/cat/executions/execution-for-test-instruction-to-rate-2/rate:
    """
    {
      "rating": 2.5,
      "comment": "test"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "rating already exists"
    }
    """

  Scenario: Rate completed instruction with success
    When I am admin
    When I do PUT /api/v4/cat/executions/execution-for-test-instruction-to-rate-3/rate:
    """
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 204

  Scenario: Rate an instruction with success, check rating calculation
    When I am admin
    When I do PUT /api/v4/cat/executions/execution-for-test-instruction-to-rate-4/rate:
    """
    {
      "rating": 3,
      "comment": "test"
    }
    """
    Then the response code should be 204
    Then I do GET /api/v4/cat/instructions/test-instruction-to-rate-4
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-instruction-to-rate-4",
      "rating": 3,
      "comments": [
        {
          "rating": 5,
          "comment": "test_user_1_comment"
        },
        {
          "rating": 1,
          "comment": "test_user_2_comment"
        },
        {
          "rating": 3,
          "comment": "test"
        }
      ]
    }
    """

  Scenario: Rate failed instruction with success
    When I am admin
    When I do PUT /api/v4/cat/executions/execution-for-test-instruction-to-rate-5/rate:
    """
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 204