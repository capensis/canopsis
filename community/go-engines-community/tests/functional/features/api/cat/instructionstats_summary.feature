Feature: get instruction statistics
  I need to be able to get instruction statistics
  Only admin should be able to get instruction statistics

  Scenario: given request should return manual instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-summary-get-1/summary
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-instruction-to-stats-summary-get-1",
      "alarm_states": {
        "critical": {
          "from": 4,
          "to": 1
        },
        "major": {
          "from": 2,
          "to": 2
        },
        "minor": {
          "from": 1,
          "to": 3
        }
      },
      "ok_alarm_states": 1,
      "avg_complete_time": 316,
      "execution_count": 7,
      "successful": 7,
      "last_executed_on": 1618280210,
      "last_modified": 1596712303,
      "created": 1596712203,
      "name": "test-instruction-to-stats-summary-get-1-name",
      "type": 0
    }
    """

  Scenario: given request should return auto instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-summary-get-2/summary
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-instruction-to-stats-summary-get-2",
      "alarm_states": {
        "critical": {
          "from": 4,
          "to": 1
        },
        "major": {
          "from": 2,
          "to": 2
        },
        "minor": {
          "from": 1,
          "to": 3
        }
      },
      "ok_alarm_states": 1,
      "avg_complete_time": 316,
      "execution_count": 7,
      "successful": 7,
      "last_executed_on": 1618280210,
      "last_modified": 1596712303,
      "created": 1596712203,
      "name": "test-instruction-to-stats-summary-get-2-name",
      "type": 1
    }
    """

  Scenario: given request should return empty instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-summary-get-3/summary
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-instruction-to-stats-summary-get-3",
      "alarm_states": {
        "critical": {
          "from": 0,
          "to": 0
        },
        "major": {
          "from": 0,
          "to": 0
        },
        "minor": {
          "from": 0,
          "to": 0
        }
      },
      "ok_alarm_states": 0,
      "avg_complete_time": 0,
      "execution_count": 0,
      "successful": 0,
      "last_executed_on": null,
      "last_modified": 1596712203,
      "created": 1596712203,
      "name": "test-instruction-to-stats-summary-get-3-name",
      "type": 1
    }
    """

  Scenario: given request and user without instruction create permission should return instruction stats
    When I am test-role-to-stats-instruction-get
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-summary-get-1/summary
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-instruction-to-stats-summary-get-1",
      "alarm_states": {
        "critical": {
          "from": 4,
          "to": 1
        },
        "major": {
          "from": 2,
          "to": 2
        },
        "minor": {
          "from": 1,
          "to": 3
        }
      },
      "ok_alarm_states": 1,
      "avg_complete_time": 316,
      "execution_count": 7,
      "successful": 7,
      "last_executed_on": 1618280210,
      "last_modified": 1596712303,
      "created": 1596712203,
      "name": "test-instruction-to-stats-summary-get-1-name",
      "type": 0
    }
    """

  Scenario: given request and user without instruction create permission should return not found error
    When I am test-role-to-stats-instruction-get
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-summary-get-2/summary
    Then the response code should be 404

  Scenario: given request and no auth user should not allow access
    When I do GET /api/v4/cat/instruction-stats/notexist/summary
    Then the response code should be 401

  Scenario: given request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/instruction-stats/notexist/summary
    Then the response code should be 403

  Scenario: given request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/notexist/summary
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given request should return manual instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-summary-get-4/summary
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-instruction-to-stats-summary-get-4",
      "alarm_states": {
        "critical": {
          "from": 2,
          "to": 1
        },
        "major": {
          "from": 2,
          "to": 2
        },
        "minor": {
          "from": 1,
          "to": 1
        }
      },
      "ok_alarm_states": 1,
      "avg_complete_time": 401,
      "execution_count": 6,
      "successful": 5,
      "last_executed_on": 1596712209,
      "last_modified": 1596712204,
      "created": 1596712203,
      "name": "test-instruction-to-stats-summary-get-4-name",
      "type": 0
    }
    """
