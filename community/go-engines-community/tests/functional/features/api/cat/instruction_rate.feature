Feature: rate an instruction
  I need to be able to rate an instruction
  Only admin should be able to rate an instruction

  Scenario: given rate request and admin user should rate instruction which was executed by someone else
    When I am admin
    When I do PUT /api/v4/notification:
    """json
    {
      "instruction": {
        "rate": true,
        "rate_frequency": {
          "value": 1,
          "unit": "s"
        }
      }
    }
    """
    Then the response code should be 200
    Then I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-rate-1&from=1000000000&to=2000000000&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-rate-1",
          "rating": 0,
          "ratable": true,
          "rate_notify": false
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-1/rate:
    """json
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 204
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-1/rate:
    """json
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
       "error": "user has already rated today"
    }
    """
    Then I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-rate-1&from=1000000000&to=2000000000&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-rate-1",
          "rating": 4.5,
          "ratable": false,
          "rate_notify": false
        }
      ]
    }
    """

  Scenario: given rate request and user without instruction create permission should rate instruction which was executed by them
    When I am test-role-to-rate-instruction
    When I do PUT /api/v4/notification:
    """json
    {
      "instruction": {
        "rate": true,
        "rate_frequency": {
          "value": 1,
          "unit": "s"
        }
      }
    }
    """
    Then the response code should be 200
    Then I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-rate-2&from=1000000000&to=2000000000&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-rate-2",
          "rating": 0,
          "ratable": true,
          "rate_notify": true
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-2/rate:
    """json
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 204
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-2/rate:
    """json
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
       "error": "user has already rated today"
    }
    """
    Then I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-rate-2&from=1000000000&to=2000000000&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-rate-2",
          "rating": 4.5,
          "ratable": false,
          "rate_notify": false
        }
      ]
    }
    """

  Scenario: given rate request should be calculated rate properly
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-3/rate:
    """json
    {
      "rating": 1.5,
      "comment": "test"
    }
    """
    Then the response code should be 204
    When I am test-role-to-rate-instruction
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-3/rate:
    """json
    {
      "rating": 3,
      "comment": "test"
    }
    """
    Then the response code should be 204
    Then I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-rate-3&from=1000000000&to=2000000000
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-instruction-to-rate-3",
          "rating": 2.2
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given rate request and admin user should not rate instruction without completed execution
    When I am admin
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-4/rate:
    """json
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "instruction wasn't completed"
    }
    """

  Scenario: given rate request and user without instruction create permission should not rate instruction without completed execution by them
    When I am test-role-to-rate-instruction
    Then I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-rate-5&from=1000000000&to=2000000000
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-5/rate:
    """json
    {
      "rating": 4.5,
      "comment": "test"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "instruction wasn't completed by the user"
    }
    """

  Scenario: given rate request and no auth user should not allow access
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-1/rate
    Then the response code should be 401

  Scenario: given rate request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/instructions/test-instruction-to-rate-1/rate
    Then the response code should be 403
