Feature: get instruction statistics
  I need to be able to get instruction statistics
  Only admin should be able to get instruction statistics

  Scenario: given get request should return manual instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-changes-get-1/changes
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 1,
              "to": 1
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
          "avg_complete_time": 400,
          "execution_count": 1,
          "successful": 1,
          "modified_on": 1596712203
        },
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 1
            },
            "major": {
              "from": 2,
              "to": 0
            },
            "minor": {
              "from": 0,
              "to": 0
            }
          },
          "ok_alarm_states": 1,
          "avg_complete_time": 275,
          "execution_count": 2,
          "successful": 2,
          "modified_on": 1596712103
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given get request should return auto instruction stats
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-changes-get-2/changes
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 1,
              "to": 1
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
          "avg_complete_time": 400,
          "execution_count": 1,
          "successful": 1,
          "modified_on": 1596712203
        },
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 1
            },
            "major": {
              "from": 2,
              "to": 0
            },
            "minor": {
              "from": 0,
              "to": 0
            }
          },
          "ok_alarm_states": 1,
          "avg_complete_time": 275,
          "execution_count": 2,
          "successful": 2,
          "modified_on": 1596712103
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given instruction without executions should return empty response
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-changes-get-3/changes
    Then the response code should be 200
    Then the response body should be:
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

  Scenario: given get request and user without instruction create permission should return instruction stats
    When I am test-role-to-stats-instruction-get
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-changes-get-1/changes
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "alarm_states": {
            "critical": {
              "from": 1,
              "to": 1
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
          "avg_complete_time": 400,
          "execution_count": 1,
          "successful": 1,
          "modified_on": 1596712203
        },
        {
          "alarm_states": {
            "critical": {
              "from": 0,
              "to": 1
            },
            "major": {
              "from": 2,
              "to": 0
            },
            "minor": {
              "from": 0,
              "to": 0
            }
          },
          "ok_alarm_states": 1,
          "avg_complete_time": 275,
          "execution_count": 2,
          "successful": 2,
          "modified_on": 1596712103
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given get request and and user without instruction create permission should return not found error
    When I am test-role-to-stats-instruction-get
    When I do GET /api/v4/cat/instruction-stats/test-instruction-to-stats-changes-get-2/changes
    Then the response code should be 404

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/instruction-stats/notexist/changes
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/instruction-stats/notexist/changes
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/instruction-stats/notexist/changes
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """
