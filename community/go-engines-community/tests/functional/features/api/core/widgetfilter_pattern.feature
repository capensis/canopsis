Feature: Get a widget filter
  I need to be able to get a widget filter
  Only admin should be able to get a widget filter

  Scenario: given updated or deleted corporate pattern request should return updated filter
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-pattern-1-title",
      "widget": "test-widget-to-filter-edit",
      "is_private": false,
      "corporate_alarm_pattern": "test-pattern-to-filter-pattern-1",
      "corporate_entity_pattern": "test-pattern-to-filter-pattern-2",
      "corporate_pbehavior_pattern": "test-pattern-to-filter-pattern-3"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "corporate_alarm_pattern": "test-pattern-to-filter-pattern-1",
      "corporate_alarm_pattern_title": "test-pattern-to-filter-pattern-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-filter-pattern-2",
      "corporate_entity_pattern_title": "test-pattern-to-filter-pattern-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-2-pattern"
            }
          }
        ]
      ],
      "corporate_pbehavior_pattern": "test-pattern-to-filter-pattern-3",
      "corporate_pbehavior_pattern_title": "test-pattern-to-filter-pattern-3-title",
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-3-pattern"
            }
          }
        ]
      ]
    }
    """
    When I save response filterID={{ .lastResponse._id }}
    When I do PUT /api/v4/patterns/test-pattern-to-filter-pattern-1:
    """json
    {
      "title": "test-pattern-to-filter-pattern-1-title-updated",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-1-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/patterns/test-pattern-to-filter-pattern-2:
    """json
    {
      "title": "test-pattern-to-filter-pattern-2-title-updated",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-2-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/patterns/test-pattern-to-filter-pattern-3:
    """json
    {
      "title": "test-pattern-to-filter-pattern-3-title-updated",
      "type": "pbehavior",
      "is_corporate": true,
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-3-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/widget-filters/{{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "corporate_alarm_pattern": "test-pattern-to-filter-pattern-1",
      "corporate_alarm_pattern_title": "test-pattern-to-filter-pattern-1-title-updated",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-1-pattern-updated"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-filter-pattern-2",
      "corporate_entity_pattern_title": "test-pattern-to-filter-pattern-2-title-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-2-pattern-updated"
            }
          }
        ]
      ],
      "corporate_pbehavior_pattern": "test-pattern-to-filter-pattern-3",
      "corporate_pbehavior_pattern_title": "test-pattern-to-filter-pattern-3-title-updated",
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-3-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    When I do DELETE /api/v4/patterns/test-pattern-to-filter-pattern-1
    Then the response code should be 204
    When I do DELETE /api/v4/patterns/test-pattern-to-filter-pattern-2
    Then the response code should be 204
    When I do DELETE /api/v4/patterns/test-pattern-to-filter-pattern-3
    Then the response code should be 204
    When I do GET /api/v4/widget-filters/{{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-1-pattern-updated"
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
              "value": "test-pattern-to-filter-pattern-2-pattern-updated"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-pattern-3-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response key "corporate_alarm_pattern" should not exist
    Then the response key "corporate_alarm_pattern_title" should not exist
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
    Then the response key "corporate_total_entity_pattern" should not exist
    Then the response key "corporate_total_entity_pattern_title" should not exist
