Feature: resolve alarm should copy user's bookmarks to a resolved alarm

  @concurrent
  Scenario: given resolve rule should resolve alarm and copy user's bookmarks to a resolved alarm
    Given I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-axe-resolve-bookmarks-1",
      "name": "test-resolve-rule-axe-resolve-bookmarks-1-name",
      "description": "test-resolve-rule-axe-resolve-bookmarks-1-desc",
      "entity_pattern":[
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-resolve-bookmarks-1"
            }
          }
        ]
      ],
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "priority": 1
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-resolve-bookmarks-1",
      "connector_name" : "test-connector-name-axe-resolve-bookmarks-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-resolve-bookmarks-1",
      "resource" : "test-resource-axe-resolve-bookmarks-1",
      "state" : 2,
      "output" : "test-output-axe-resolve-bookmarks-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-resolve-bookmarks-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-resolve-bookmarks-1",
            "connector": "test-connector-axe-resolve-bookmarks-1",
            "connector_name": "test-connector-name-axe-resolve-bookmarks-1",
            "resource": "test-resource-axe-resolve-bookmarks-1"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I save response alarmID={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmID }}/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-axe-resolve-bookmarks-1&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .alarmID }}",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-resolve-bookmarks-1",
      "connector_name" : "test-connector-name-axe-resolve-bookmarks-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-resolve-bookmarks-1",
      "resource" : "test-resource-axe-resolve-bookmarks-1",
      "state" : 0,
      "output" : "test-output-axe-resolve-bookmarks-1"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type" : "resolve_close",
      "connector" : "test-connector-axe-resolve-bookmarks-1",
      "connector_name" : "test-connector-name-axe-resolve-bookmarks-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-resolve-bookmarks-1",
      "resource" : "test-resource-axe-resolve-bookmarks-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-resolve-bookmarks-1&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .alarmID }}",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
