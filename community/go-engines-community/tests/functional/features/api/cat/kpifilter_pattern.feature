Feature: Update a kpi filter
  I need to be able to update a kpi filter

  Scenario: given updated or deleted corporate pattern request should return updated kpi filter
    When I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-kpi-filter-to-pattern-1-name",
      "corporate_entity_pattern": "test-pattern-to-kpi-filter-pattern-1"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "corporate_entity_pattern": "test-pattern-to-kpi-filter-pattern-1",
      "corporate_entity_pattern_title": "test-pattern-to-kpi-filter-pattern-1-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-kpi-filter-pattern-1-pattern"
            }
          }
        ]
      ]
    }
    """
    When I save response filterID={{ .lastResponse._id }}
    When I do PUT /api/v4/patterns/test-pattern-to-kpi-filter-pattern-1:
    """json
    {
      "title": "test-pattern-to-kpi-filter-pattern-1-title-updated",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-kpi-filter-pattern-1-pattern-updated"
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
    Then the response code should be 200
    When I do GET /api/v4/cat/kpi-filters/{{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "corporate_entity_pattern": "test-pattern-to-kpi-filter-pattern-1",
      "corporate_entity_pattern_title": "test-pattern-to-kpi-filter-pattern-1-title-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-kpi-filter-pattern-1-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    When I do DELETE /api/v4/patterns/test-pattern-to-kpi-filter-pattern-1
    Then the response code should be 204
    When I do GET /api/v4/cat/kpi-filters/{{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-kpi-filter-pattern-1-pattern-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
