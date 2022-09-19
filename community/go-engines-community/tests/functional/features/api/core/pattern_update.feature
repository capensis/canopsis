Feature: Update a saved pattern
  I need to be able to update a saved pattern
  Only admin should be able to update a saved pattern

  Scenario: given update request should return ok
    When I am noperms
    When I do PUT /api/v4/patterns/test-pattern-to-update-1:
    """json
    {
      "title": "test-pattern-to-update-1-title",
      "type": "alarm",
      "is_corporate": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "title": "test-pattern-to-update-1-title",
      "type": "alarm",
      "is_corporate": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
            }
          }
        ]
      ],
      "created": 1605263992
    }
    """

  Scenario: given update request and another user should return not found
    When I am admin
    When I do PUT /api/v4/patterns/test-pattern-to-update-1:
    """json
    {
      "title": "test-pattern-to-update-1-title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404

  Scenario: given update corporate pattern request and another user should return ok
    When I am admin
    When I do PUT /api/v4/patterns/test-pattern-to-update-2:
    """json
    {
      "title": "test-pattern-to-update-2-title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-2-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-pattern-to-update-2-title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-2-pattern"
            }
          }
        ]
      ]
    }
    """

  Scenario: given update corporate pattern request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/patterns/test-pattern-to-update-2:
    """json
    {
      "title": "test-pattern-to-update-2-title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-2-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/patterns/test-pattern-notexist
    Then the response code should be 401

  Scenario: given update request with not exist id should return not found error
    When I am noperms
    When I do PUT /api/v4/patterns/test-pattern-notexist:
    """json
    {
      "title": "test-pattern-to-update-notexist",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-notexist-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404

  Scenario: given update request with another type should return bad request error
    When I am noperms
    When I do PUT /api/v4/patterns/test-pattern-to-update-1:
    """json
    {
      "title": "test-pattern-to-update-1-title",
      "type": "entity",
      "is_corporate": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
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
      "errors": {"type": "Type cannot be changed"}
    }
    """

  Scenario: given update request with another corporate status should return bad request error
    When I am noperms
    When I do PUT /api/v4/patterns/test-pattern-to-update-2:
    """json
    {
      "title": "test-pattern-to-update-2-title",
      "type": "alarm",
      "is_corporate": false,
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-update-1-pattern"
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
        "is_corporate": "IsCorporate cannot be changed"
      }
    }
    """

  Scenario: given update request with missing fields should return bad request error
    When I am admin
    When I do PUT /api/v4/patterns/test-pattern-to-update-1:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "type": "Type is missing.",
        "is_corporate": "IsCorporate is missing."
      }
    }
    """
