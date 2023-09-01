Feature: update a pbehavior
  I need to be able to update a pbehavior
  Only admin should be able to update a pbehavior

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without view permission should not allow access
    When I am noperms
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 403

  @concurrent
  Scenario: given update request with the name that already exists should cause dup error
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-1:
    """json
    {
      "name": "test-pbehavior-to-check-unique-name",
      "enabled": true,
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-update-1-pattern"
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
  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-1:
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
  Scenario: given no exist pbehavior id should return error
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-not-exist:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-not-exist",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-not-exist"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404
