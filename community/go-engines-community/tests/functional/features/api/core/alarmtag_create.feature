Feature: Create a alarm tag
  I need to be able to create a alarm tag
  Only admin should be able to create a alarm tag

  @concurrent
  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-alarm-tag-to-create-1",
      "color": "#AABBCC",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-create-1-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "type": 1,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "value": "test-alarm-tag-to-create-1",
      "color": "#AABBCC",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-create-1-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/alarm-tags/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": 1,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "value": "test-alarm-tag-to-create-1",
      "color": "#AABBCC",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-create-1-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request with both corporate patterns should return success
    When I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-alarm-tag-to-create-2",
      "color": "#AABBCC",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "type": 1,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "value": "test-alarm-tag-to-create-2",
      "color": "#AABBCC",
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-2-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/alarm-tags
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/alarm-tags
    Then the response code should be 403

  @concurrent
  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "value": "Value is missing.",
        "color": "Color is missing.",
        "entity_pattern": "EntityPattern or AlarmPattern is required.",
        "alarm_pattern": "AlarmPattern or EntityPattern is required."
      }
    }
    """

  @concurrent
  Scenario: given create request with invalid patterns format should return bad request
    When I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "wrong_field",
            "cond": {
              "type": "eq",
              "value": "ram"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "wrong_field",
            "cond": {
              "type": "eq",
              "value": "ram"
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
        "alarm_pattern": "AlarmPattern is invalid alarm pattern.",
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """    

  @concurrent
  Scenario: given create request with already exists value should return error
    When I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-alarm-tag-to-check-unique",
      "color": "#AABBCC",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-create"
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
        "value": "Value already exists."
      }
    }
    """
