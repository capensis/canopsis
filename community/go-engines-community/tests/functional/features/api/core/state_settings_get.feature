Feature: get state settings
  I need to be able to get state settings
  Only admin should be able to get state settings

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/state-settings
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/state-settings
    Then the response code should be 403

  @concurrent
  Scenario: given get list request should return ok
    When I am admin
    When I do GET /api/v4/state-settings?search=state-settings-to-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "state-settings-to-get-1",
          "method": "inherited",
          "title": "state-settings-to-get-1-title",
          "enabled": true,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-state-settings-to-get-rule-pattern-1"
                }
              }
            ]
          ],
          "inherited_entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-state-settings-to-get-impacting-pattern-1"
                }
              }
            ]
          ]
        },
        {
          "_id": "state-settings-to-get-2",
          "method": "inherited",
          "title": "state-settings-to-get-2-title",
          "enabled": true,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-state-settings-to-get-rule-pattern-2"
                }
              }
            ]
          ],
          "inherited_entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-state-settings-to-get-impacting-pattern-2"
                }
              }
            ]
          ]
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 2
      }
    }
    """


  @concurrent
  Scenario: given get by id request should return ok
    When I am admin
    When I do GET /api/v4/state-settings/state-settings-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "state-settings-to-get-1",
      "method": "inherited",
      "title": "state-settings-to-get-1-title",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-get-rule-pattern-1"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-get-impacting-pattern-1"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given get by id request should return ok
    When I am admin
    When I do GET /api/v4/state-settings/state-settings-to-get-not found
    Then the response code should be 404
    Then the response body should contain:
    """json
    {
      "error": "Not found"
    }
    """
