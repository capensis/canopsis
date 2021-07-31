Feature: Create an baggot rule
  I need to be able to create a baggot rule
  Only admin should be able to create a baggot rule

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/baggot-rules:
    """json
    {
      "_id": "test-baggot-rule-to-create-1",
      "description": "test create 1",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-baggot-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns":[
        {
          "name": "test-baggot-rule-to-create-1-pattern"
        }
      ],
      "duration": {
        "seconds": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-baggot-rule-to-create-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test create 1",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-baggot-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns":[
        {
          "name": "test-baggot-rule-to-create-1-pattern"
        }
      ],
      "duration": {
        "seconds": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    When I do GET /api/v4/baggot-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-baggot-rule-to-create-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test create 1",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-baggot-rule-to-create-1-pattern"
          }
        }
      ],
      "entity_patterns":[
        {
          "name": "test-baggot-rule-to-create-1-pattern"
        }
      ],
      "duration": {
        "seconds": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """

  Scenario: given create request should update priority of next rules
    When I am admin
    When I do POST /api/v4/baggot-rules:
    """json
    {
      "_id": "test-baggot-rule-to-create-2",
      "description": "test create 2",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-baggot-rule-to-create-2-pattern"
          }
        }
      ],
      "entity_patterns":[
        {
          "name": "test-baggot-rule-to-create-2-pattern"
        }
      ],
      "duration": {
        "seconds": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/baggot-rules:
    """json
    {
      "_id": "test-baggot-rule-to-create-3",
      "description": "test create 3",
      "alarm_patterns": [
        {
          "v": {
            "component": "test-baggot-rule-to-create-3-pattern"
          }
        }
      ],
      "entity_patterns":[
        {
          "name": "test-baggot-rule-to-create-2-pattern"
        }
      ],
      "duration": {
        "seconds": 10,
        "unit": "s"
      },
      "priority": 5
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/baggot-rules/test-baggot-rule-to-create-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "priority": 6
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/baggot-rules:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing.",
        "duration.seconds": "Seconds is missing.",
        "duration.unit": "Unit is missing.",
        "priority": "Priority is missing."
      }
    }
    """

  Scenario: given create request with wrong alarm patterns format should return bad request
    When I am admin
    When I do POST /api/v4/baggot-rules:
    """json
    {
      "alarm_patterns": [
        {
          "v": {
            "component_name": "ram"
          }
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_patterns":"Invalid alarm patterns."
      }
    }
    """

  Scenario: given create request with wrong entity patterns format should return bad request
    When I am admin
    When I do POST /api/v4/baggot-rules:
    """json
    {
      "entity_patterns": [{"test": "test"}]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_patterns":"Invalid entity patterns."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/baggot-rules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/baggot-rules
    Then the response code should be 403

  Scenario: given create request with already exists id should return error
    When I am admin
    When I do POST /api/v4/baggot-rules:
    """json
    {
      "_id": "test-baggot-rule-to-check-unique"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """
