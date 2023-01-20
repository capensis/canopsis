Feature: update events count
  I need to be able to get a events count of alarms
  
  @concurrent
  Scenario: given alarm should update events count
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-events-count-1",
      "connector_name" : "test-connector-name-axe-events-count-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-events-count-1",
      "resource" : "test-resource-axe-events-count-1",
      "state" : 2,
      "output" : "test-output-axe-events-count-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-events-count-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "events_count": 1
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-events-count-1",
      "connector_name" : "test-connector-name-axe-events-count-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-events-count-1",
      "resource" : "test-resource-axe-events-count-1",
      "state" : 2,
      "output" : "test-output-axe-events-count-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-events-count-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "events_count": 2
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
