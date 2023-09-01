Feature: Update a alarm tag
  I need to be able to update a alarm tag
  Only admin should be able to update a alarm tag

  @concurrent
  Scenario: given update request should return ok
    When I am admin
    When I do PUT /api/v4/alarm-tags/test-alarm-tag-to-update-1:
    """json
    {
      "color": "#AABBCC",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-update-1-pattern"
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
              "value": "test-alarm-tag-to-update-1-pattern"
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
      "_id": "test-alarm-tag-to-update-1",
      "type": 1,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1612139798,
      "value": "test-alarm-tag-to-update-1-value",
      "color": "#AABBCC",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-update-1-pattern"
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
              "value": "test-alarm-tag-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given update external tag request should return ok
    When I am admin
    When I do PUT /api/v4/alarm-tags/test-alarm-tag-to-update-2:
    """json
    {
      "color": "#AABBCC"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-alarm-tag-to-update-2",
      "type": 0,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1612139798,
      "value": "test-alarm-tag-to-update-2-value",
      "color": "#AABBCC"
    }
    """

  @concurrent
  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/alarm-tags/test-alarm-tag-to-update-1:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "color": "Color is missing."
      }
    }
    """
    When I do PUT /api/v4/alarm-tags/test-alarm-tag-to-update-1:
    """json
    {
      "color": "#AABBCC"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern or EntityPattern is required."
      }
    }
    """

  @concurrent
  Scenario: given invalid update external tag request should return errors
    When I am admin
    When I do PUT /api/v4/alarm-tags/test-alarm-tag-to-update-2:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "color": "Color is missing."
      }
    }
    """

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/alarm-tags/notexist
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarm-tags/notexist
    Then the response code should be 403

  @concurrent
  Scenario: given update request with not exist id should return error
    When I am admin
    When I do PUT /api/v4/alarm-tags/notexist:
    """json
    {
      "color": "#AABBCC",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-alarm-tag-to-update-notexist"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404
