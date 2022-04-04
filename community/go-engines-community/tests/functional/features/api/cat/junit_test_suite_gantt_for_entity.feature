Feature: get test suite's gantt intervals for entity
  I need to be able to get test suite's gantt intervals for entity
  Only admin should be able to get test suite's gantt intervals for entity

  Scenario: given get request for component should return gantt
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-entity-gantt?_id=test-suite-to-get-entity-gantt-1-name.test-report-to-get-entity-gantt.xml
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "avg_status": 0,
          "avg_to": 0,
          "avg_time": 0,
          "from": 0,
          "_id": "test-case-to-get-entity-gantt-1-1",
          "message": "test-case-to-get-entity-gantt-1-1-msg",
          "name": "test-case-to-get-entity-gantt-1-1-name",
          "status": 0,
          "time": 1,
          "to": 1
        },
        {
          "avg_status": 0,
          "avg_to": 0,
          "avg_time": 0,
          "from": 1,
          "_id": "test-case-to-get-entity-gantt-1-2",
          "message": "test-case-to-get-entity-gantt-1-2-msg",
          "name": "test-case-to-get-entity-gantt-1-2-name",
          "status": 1,
          "time": 2,
          "to": 3
        },
        {
          "avg_status": 0,
          "avg_to": 0,
          "avg_time": 0,
          "from": 3,
          "_id": "test-case-to-get-entity-gantt-1-3",
          "message": "test-case-to-get-entity-gantt-1-3-msg",
          "name": "test-case-to-get-entity-gantt-1-3-name",
          "status": 2,
          "time": 3,
          "to": 6
        },
        {
          "avg_status": 0,
          "avg_to": 0,
          "avg_time": 0,
          "from": 6,
          "_id": "test-case-to-get-entity-gantt-1-4",
          "message": "test-case-to-get-entity-gantt-1-4-msg",
          "name": "test-case-to-get-entity-gantt-1-4-name",
          "status": 3,
          "time": 4,
          "to": 10
        }
      ],
      "meta": {
        "time": 10,
        "avg_time": 0,
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

  Scenario: given get request for resource should return gantt
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-entity-gantt?_id=test-case-to-get-entity-gantt-1-1-name/test-suite-to-get-entity-gantt-1-name.test-report-to-get-entity-gantt.xml
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "avg_status": 0,
          "avg_to": 0,
          "avg_time": 0,
          "from": 0,
          "_id": "test-case-to-get-entity-gantt-1-1",
          "message": "test-case-to-get-entity-gantt-1-1-msg",
          "name": "test-case-to-get-entity-gantt-1-1-name",
          "status": 0,
          "time": 1,
          "to": 1
        },
        {
          "avg_status": 0,
          "avg_to": 0,
          "avg_time": 0,
          "from": 1,
          "_id": "test-case-to-get-entity-gantt-1-2",
          "message": "test-case-to-get-entity-gantt-1-2-msg",
          "name": "test-case-to-get-entity-gantt-1-2-name",
          "status": 1,
          "time": 2,
          "to": 3
        },
        {
          "avg_status": 0,
          "avg_to": 0,
          "avg_time": 0,
          "from": 3,
          "_id": "test-case-to-get-entity-gantt-1-3",
          "message": "test-case-to-get-entity-gantt-1-3-msg",
          "name": "test-case-to-get-entity-gantt-1-3-name",
          "status": 2,
          "time": 3,
          "to": 6
        },
        {
          "avg_status": 0,
          "avg_to": 0,
          "avg_time": 0,
          "from": 6,
          "_id": "test-case-to-get-entity-gantt-1-4",
          "message": "test-case-to-get-entity-gantt-1-4-msg",
          "name": "test-case-to-get-entity-gantt-1-4-name",
          "status": 3,
          "time": 4,
          "to": 10
        }
      ],
      "meta": {
        "time": 10,
        "avg_time": 0,
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/junit/test-suites-entity-gantt
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/junit/test-suites-entity-gantt
    Then the response code should be 403

  Scenario: given invalid get request should return error
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-entity-gantt?_id=notexist
    Then the response code should be 404

  Scenario: given invalid get request should return error
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-entity-gantt?_id=test-junit-entity-to-get-entity-gantt-not-linked-to-test-suite
    Then the response code should be 404
