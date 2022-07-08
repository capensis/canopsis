Feature: Entities should be synchronized in metrics db
  I need to be able to see metrics.

  Scenario: given updated entity should get metrics by updated entity
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-che-1-1-name",
      "entity_pattern": [
        [
          {
            "field": "infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-metrics-che-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filter1ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-che-1-2-name",
      "entity_pattern": [
        [
          {
            "field": "infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-metrics-che-1-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filter2ID={{ .lastResponse._id }}
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-che-1"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "client",
            "description": "Client",
            "value": "{{ `{{ .Event.ExtraInfos.client }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority": 2,
      "description": "test-eventfilter-metrics-che-1-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter2ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter1ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
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
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-che-2-1-name",
      "entity_pattern": [
        [
          {
            "field": "component_infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-metrics-che-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filter1ID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-che-2-2-name",
      "entity_pattern": [
        [
          {
            "field": "component_infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-metrics-che-2-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filter2ID={{ .lastResponse._id }}
    When I do POST /api/v4/eventfilter/rules:
    """json
    {
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-metrics-che-2"
            }
          },
          {
            "field": "source_type",
            "cond": {
              "type": "eq",
              "value": "component"
            }
          },
          {
            "field": "event_type",
            "cond": {
              "type": "eq",
              "value": "check"
            }
          }
        ]
      ],
      "config": {
        "actions": [
          {
            "type": "set_entity_info_from_template",
            "name": "client",
            "description": "Client",
            "value": "{{ `{{ .Event.ExtraInfos.client }}` }}"
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      },
      "priority": 2,
      "description": "test-eventfilter-metrics-che-2-description",
      "enabled": true
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter2ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
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
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filter1ID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
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
