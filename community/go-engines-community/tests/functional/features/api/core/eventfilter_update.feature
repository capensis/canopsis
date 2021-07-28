Feature: Update an eventfilter
  I need to be able to update an eventfilter
  Only admin should be able to update an eventfilter

  Scenario: given update request should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update:
    """
    {
      "description": "drop filter",
      "type": "drop",
      "patterns": [
        {
          "resource": "test-eventfilter-to-update-pattern-updated"
        }
      ],
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update",
      "author": "root",
      "description": "drop filter",
      "type": "drop",
      "patterns": [
        {
          "resource": "test-eventfilter-to-update-pattern-updated"
        }
      ],
      "priority": 1,
      "enabled": true,
      "created": 1608635535
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-not-found:
    """
    {
      "description": "drop filter",
      "type": "drop",
      "patterns": [
        {
          "resource": "test-eventfilter-to-update-pattern-updated"
        }
      ],
      "priority": 1,
      "enabled": true
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
    
  Scenario: given PUT change_entity rule requests should return error, because of empty config
    Given I am admin
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
      "description": "update change_entity",
      "type": "change_entity",
      "patterns": [
        {
          "resource": "never be used change entity update test"
        }
      ],
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
      "description": "update change_entity",
      "type": "change_entity",
      "patterns": [
        {
          "resource": "never be used change entity update test"
        }
      ],
      "config": {
        "component": "",
        "connector": "",
        "resource": "",
        "connector_name": ""
      },
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
      "description": "update change_entity",
      "type": "change_entity",
      "patterns": [
        {
          "resource": "never be used change entity update test"
        }
      ],
      "config": {},
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """
