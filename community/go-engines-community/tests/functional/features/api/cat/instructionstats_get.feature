Feature: get instruction statistics
  I need to be able to get instruction statistics
  Only admin should be able to get instruction statistics

  Scenario: given get all request and admin user should return instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-stats-get&from=1617555600&to=1618246799&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-instruction-to-stats-get-1",
          "alarm_states": {
            "critical": {
              "from": 4,
              "to": 2
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
          "ok_alarm_states": 0,
          "avg_complete_time": 316,
          "execution_count": 8,
          "successful": 7,
          "last_executed_on": 1618394399,
          "last_modified": 1617555600,
          "created": 1617555600,
          "name": "test-instruction-to-stats-get-1-name",
          "rating": 3.2,
          "type": 0,
          "has_executions": true
        },
        {
          "_id": "test-instruction-to-stats-get-2",
          "alarm_states": {
            "critical": {
              "from": 4,
              "to": 2
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
          "ok_alarm_states": 0,
          "avg_complete_time": 316,
          "execution_count": 8,
          "successful": 7,
          "last_executed_on": 1618394399,
          "last_modified": 1617555600,
          "created": 1617555600,
          "name": "test-instruction-to-stats-get-2-name",
          "rating": 0,
          "type": 1,
          "has_executions": true
        },
        {
          "_id": "test-instruction-to-stats-get-3",
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
          "last_executed_on": null,
          "last_modified": 1617555600,
          "created": 1617555600,
          "name": "test-instruction-to-stats-get-3-name",
          "rating": 0,
          "type": 1,
          "has_executions": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-stats-get&from=1616950800&to=1618419599&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-instruction-to-stats-get-1",
          "alarm_states": {
            "critical": {
              "from": 5,
              "to": 3
            },
            "major": {
              "from": 4,
              "to": 2
            },
            "minor": {
              "from": 1,
              "to": 3
            }
          },
          "ok_alarm_states": 2,
          "avg_complete_time": 316,
          "execution_count": 11,
          "successful": 10,
          "last_executed_on": 1618394399,
          "last_modified": 1617555600,
          "created": 1617555600,
          "name": "test-instruction-to-stats-get-1-name",
          "rating": 3.2,
          "type": 0,
          "has_executions": true
        },
        {
          "_id": "test-instruction-to-stats-get-2",
          "alarm_states": {
            "critical": {
              "from": 5,
              "to": 3
            },
            "major": {
              "from": 4,
              "to": 2
            },
            "minor": {
              "from": 1,
              "to": 3
            }
          },
          "ok_alarm_states": 2,
          "avg_complete_time": 316,
          "execution_count": 11,
          "successful": 10,
          "last_executed_on": 1618394399,
          "last_modified": 1617555600,
          "created": 1617555600,
          "name": "test-instruction-to-stats-get-2-name",
          "rating": 0,
          "type": 1,
          "has_executions": true
        },
        {
          "_id": "test-instruction-to-stats-get-3",
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
          "last_modified": 1617555600,
          "created": 1617555600,
          "name": "test-instruction-to-stats-get-3-name",
          "rating": 0,
          "type": 1,
          "has_executions": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request and user without instruction create permission should return instruction stats
    When I am test-role-to-stats-instruction-get
    When I do GET /api/v4/cat/instruction-stats?search=test-instruction-to-stats-get&from=1617555600&to=1618246799&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-instruction-to-stats-get-1",
          "alarm_states": {
            "critical": {
              "from": 4,
              "to": 2
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
          "ok_alarm_states": 0,
          "avg_complete_time": 316,
          "execution_count": 8,
          "successful": 7,
          "last_executed_on": 1618394399,
          "last_modified": 1617555600,
          "created": 1617555600,
          "name": "test-instruction-to-stats-get-1-name",
          "rating": 3.2,
          "type": 0,
          "has_executions": true
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

  Scenario: given sort all request should return sorted instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats?sort_by=name&sort=desc&search=test-instruction-to-stats-get&from=1617555600&to=1618246799
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-instruction-to-stats-get-3"
        },
        {
          "_id": "test-instruction-to-stats-get-2"
        },
        {
          "_id": "test-instruction-to-stats-get-1"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/instruction-stats
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/instruction-stats
    Then the response code should be 403
