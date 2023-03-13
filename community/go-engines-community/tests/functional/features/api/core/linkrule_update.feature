Feature: Update an link rule
  I need to be able to update an link rule
  Only admin should be able to update an link rule

  @concurrent
  Scenario: given update request should update link rule
    When I am admin
    Then I do PUT /api/v4/link-rules/test-link-rule-to-update:
    """json
    {
      "name": "test-link-rule-to-update-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-update-link-1-label",
          "category": "test-link-rule-to-update-link-1-category",
          "icon_name": "test-link-rule-to-update-link-1-icon",
          "url": "http://test-link-rule-to-update-link-1-url.com",
          "single": true
        },
        {
          "label": "test-link-rule-to-update-link-2-label",
          "category": "test-link-rule-to-update-link-2-category",
          "icon_name": "test-link-rule-to-update-link-2-icon",
          "url": "http://test-link-rule-to-update-link-2-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-update-pattern"
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
              "value": "test-link-rule-to-update-resource"
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
      "created": 1605263992,
      "name": "test-link-rule-to-update-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-update-link-1-label",
          "category": "test-link-rule-to-update-link-1-category",
          "icon_name": "test-link-rule-to-update-link-1-icon",
          "url": "http://test-link-rule-to-update-link-1-url.com",
          "single": true
        },
        {
          "label": "test-link-rule-to-update-link-2-label",
          "category": "test-link-rule-to-update-link-2-category",
          "icon_name": "test-link-rule-to-update-link-2-icon",
          "url": "http://test-link-rule-to-update-link-2-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-update-pattern"
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
              "value": "test-link-rule-to-update-resource"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/link-rules/test-link-rule-to-update
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/link-rules/test-link-rule-to-update
    Then the response code should be 403

  @concurrent
  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/link-rules/test-link-rule-not-found:
    """json
    {
      "name": "test-link-rule-not-found-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-not-found-link-1-label",
          "category": "test-link-rule-not-found-link-1-category",
          "icon_name": "test-link-rule-not-found-link-1-icon",
          "url": "http://test-link-rule-not-found-link-1-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-not-found-pattern"
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
              "value": "test-link-rule-not-found-resource"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given update request with missing fields should return bad request
    When I am admin
    Then I do PUT /api/v4/link-rules/test-link-rule-to-update:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "type": "Type is missing.",
        "enabled": "Enabled is missing.",
        "links": "Links or SourceCode is required.",
        "source_code": "SourceCode or Links is required."
      }
    }
    """
    Then I do PUT /api/v4/link-rules/test-link-rule-to-update:
    """json
    {
      "type": "alarm"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "enabled": "Enabled is missing.",
        "links": "Links or SourceCode is required.",
        "source_code": "SourceCode or Links is required.",
        "alarm_pattern": "AlarmPattern or EntityPattern is required.",
        "entity_pattern": "EntityPattern or AlarmPattern is required."
      }
    }
    """
    Then I do PUT /api/v4/link-rules/test-link-rule-to-update:
    """json
    {
      "type": "entity"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "enabled": "Enabled is missing.",
        "links": "Links or SourceCode is required.",
        "source_code": "SourceCode or Links is required.",
        "entity_pattern": "EntityPattern is missing."
      }
    }
    """

  @concurrent
  Scenario: given create request with already exists id and name should return error
    When I am admin
    When I do PUT /api/v4/link-rules/test-link-rule-to-update-1:
    """json
    {
      "name": "test-link-rule-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """
