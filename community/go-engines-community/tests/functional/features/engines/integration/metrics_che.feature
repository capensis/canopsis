Feature: Entities should be synchronized in metrics db
  I need to be able to see metrics.

  @concurrent
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "client": "test-client-metrics-che-1",
      "connector": "test-connector-metrics-che-1",
      "connector_name": "test-connector-name-metrics-che-1",
      "component": "test-component-metrics-che-1",
      "resource": "test-resource-metrics-che-1",
      "source_type": "resource"
    }
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filter1ID }}",
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
    }
    """
    When I do POST /api/v4/cat/metrics/alarm until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDateTz }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "client": "test-client-metrics-che-1-updated",
      "connector": "test-connector-metrics-che-1",
      "connector_name": "test-connector-name-metrics-che-1",
      "component": "test-component-metrics-che-1",
      "resource": "test-resource-metrics-che-1",
      "source_type": "resource"
    }
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filter2ID }}",
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
    }
    """
    When I do POST /api/v4/cat/metrics/alarm until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDateTz }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filter1ID }}",
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
    }
    """
    When I do POST /api/v4/cat/metrics/alarm until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDateTz }},
              "value": 0
            }
          ]
        }
      ]
    }
    """

  @concurrent
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 0,
      "client": "test-client-metrics-che-2",
      "connector": "test-connector-metrics-che-2",
      "connector_name": "test-connector-name-metrics-che-2",
      "component": "test-component-metrics-che-2",
      "source_type": "component"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "client": "test-client-metrics-che-2",
      "connector": "test-connector-metrics-che-2",
      "connector_name": "test-connector-name-metrics-che-2",
      "component": "test-component-metrics-che-2",
      "resource": "test-resource-metrics-che-2",
      "source_type": "resource"
    }
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filter1ID }}",
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
    }
    """
    When I do POST /api/v4/cat/metrics/alarm until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDateTz }},
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
      "event_type": "check",
      "state": 0,
      "client": "test-client-metrics-che-2-updated",
      "connector": "test-connector-metrics-che-2",
      "connector_name": "test-connector-name-metrics-che-2",
      "component": "test-component-metrics-che-2",
      "source_type": "component"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-metrics-che-2",
        "connector_name": "test-connector-name-metrics-che-2",
        "component": "test-component-metrics-che-2",
        "source_type": "component"
      },
      {
        "event_type": "entityupdated",
        "connector": "test-connector-metrics-che-2",
        "connector_name": "test-connector-name-metrics-che-2",
        "component": "test-component-metrics-che-2",
        "resource": "test-resource-metrics-che-2",
        "source_type": "resource"
      }
    ]
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filter2ID }}",
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
    }
    """
    When I do POST /api/v4/cat/metrics/alarm until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDateTz }},
              "value": 1
            }
          ]
        }
      ]
    }
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filter1ID }}",
      "sampling": "day",
      "from": {{ nowDateTz }},
      "to": {{ nowDateTz }}
    }
    """
    When I do POST /api/v4/cat/metrics/alarm until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDateTz }},
              "value": 0
            }
          ]
        }
      ]
    }
    """
