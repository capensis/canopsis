Feature: create a pbehavior
  I need to be able to create a pbehavior
  Only admin should be able to create a pbehavior

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/pbehaviors
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without view permission should not allow access
    When I am noperms
    When I do POST /api/v4/pbehaviors
    Then the response code should be 403

  @concurrent
  Scenario: given create request with the custom ID that already exists should cause dup error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "test-pbehavior-to-check-unique",
      "enabled": true,
      "name": "test-pbehavior-to-create",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  @concurrent
  Scenario: given create request with the name that already exists should cause dup error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "name": "test-pbehavior-to-check-unique-name",
      "enabled": true,
      "tstart": 1591172881,
      "tstop": 1591536400,
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "name": "Name is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "reason": "Reason is missing.",
        "tstart": "Start is missing.",
        "type": "Type is missing."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with start > stop should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "tstart": 1592172881,
      "tstop": 1591536400
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with not existed reason should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "reason": "notexist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "reason": "Reason doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with invalid exclude date should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "exdates": [
        {
          "begin": 1592164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "exdates.0.end": "End should be greater than Begin."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "exdates": [
        {
          "begin": 1592164001,
          "end": 1592166001,
          "type": "test-type-not-exist"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "exdates": "Exdates doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with invalid exception should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "exceptions": [
        "test-exception-to-pbh-edit",
        "test-exception-not-exist"
     ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "exceptions": "Exceptions doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with not existed type should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "type": "notexist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with invalid entity pattern should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "entity_pattern": [[]]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-create-pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "corporate_entity_pattern": "test-pattern-to-rule-edit-1",
      "enabled":true,
      "name": "test-pbehavior-to-create-7",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "corporate_entity_pattern": "test-pattern-not-found",
      "enabled":true,
      "name": "test-pbehavior-to-create-7",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with missing stop should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "tstart": 1591172881,
      "type": "test-type-to-pbh-edit-1"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "tstop": "Stop is missing."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with invalid id should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "invalid/id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "invalidid?key=value"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "$invalidid"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "_id": "invalid.id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """

  @concurrent
  Scenario: given invalid create request with invalid color should return error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "color": "notcolor"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "color": "Color is not valid."
      }
    }
    """
