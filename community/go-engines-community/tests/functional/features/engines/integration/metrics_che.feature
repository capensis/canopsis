Feature: Entities should be synchronized in metrics db
  I need to be able to see metrics.

  Scenario: given updated entity should get metrics by updated entity
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-che-1-1-name",
      "entity_patterns": [
        {
          "infos": {
            "client": {
              "value": "test-client-metrics-che-1"
            }
          }
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filter1ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-che-1-2-name",
      "entity_patterns": [
        {
          "infos": {
            "client": {
              "value": "test-client-metrics-che-1-updated"
            }
          }
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filter2ID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-metrics-che-1",
      "connector_name": "test-connector-name-metrics-che-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-che-1",
      "resource": "test-resource-metrics-che-1",
      "state": 1,
      "client": "test-client-metrics-che-1"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter1ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-metrics-che-1",
      "connector_name": "test-connector-name-metrics-che-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-che-1",
      "resource": "test-resource-metrics-che-1",
      "state": 1,
      "client": "test-client-metrics-che-1-updated"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter2ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter1ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  Scenario: given updated component should get metrics by updated resource
    Given I am admin
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-che-2-1-name",
      "entity_patterns": [
        {
          "component_infos": {
            "client": {
              "value": "test-client-metrics-che-2"
            }
          }
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filter1ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/filters:
    """json
    {
      "name": "test-filter-metrics-che-2-2-name",
      "entity_patterns": [
        {
          "component_infos": {
            "client": {
              "value": "test-client-metrics-che-2-updated"
            }
          }
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response filter2ID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-metrics-che-2",
      "connector_name": "test-connector-name-metrics-che-2",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-metrics-che-2",
      "state": 0,
      "client": "test-client-metrics-che-2"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-metrics-che-2",
      "connector_name": "test-connector-name-metrics-che-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-che-2",
      "resource": "test-resource-metrics-che-2",
      "state": 1,
      "client": "test-client-metrics-che-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter1ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-metrics-che-2",
      "connector_name": "test-connector-name-metrics-che-2",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-metrics-che-2",
      "state": 0,
      "client": "test-client-metrics-che-2-updated"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter2ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter1ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        }
      ]
    }
    """
