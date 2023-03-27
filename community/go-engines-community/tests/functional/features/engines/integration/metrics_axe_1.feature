Feature: Metrics should be added on alarm changes
  I need to be able to see metrics.

  @concurrent
  Scenario: given new alarm should add created_alarms metric
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-1",
      "connector_name": "test-connector-name-metrics-axe-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-1",
      "resource": "test-resource-metrics-axe-1",
      "state": 1
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=created_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
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

  @concurrent
  Scenario: given new alarm with auto instruction should add instruction_alarms and non_displayed_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["create"],
      "name": "test-instruction-metrics-axe-2-1-name",
      "description": "test-instruction-metrics-axe-2-1-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-metrics-axe-2-1",
                "test-resource-metrics-axe-2-2"
              ]
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-metrics-axe"
        }
      ],
      "priority": 100
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["create"],
      "name": "test-instruction-metrics-axe-2-2-name",
      "description": "test-instruction-metrics-axe-2-2-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-metrics-axe-2-1",
                "test-resource-metrics-axe-2-2"
              ]
            }
          }
        ]
      ],
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-metrics-axe"
        }
      ],
      "priority": 100
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-metrics-axe-2-1",
                "test-resource-metrics-axe-2-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-1",
        "state": 1
      },
      {
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-2",
        "state": 1
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-2",
        "connector_name": "test-connector-name-metrics-axe-2",
        "component": "test-component-metrics-axe-2",
        "resource": "test-resource-metrics-axe-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-metrics-axe-2-1&with_instructions=true until response code is 200 and response array key "data.0.successful_auto_instructions" contains:
    """json
    [
      "test-instruction-metrics-axe-2-1-name",
      "test-instruction-metrics-axe-2-2-name"
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-metrics-axe-2-2&with_instructions=true until response code is 200 and response array key "data.0.successful_auto_instructions" contains:
    """json
    [
      "test-instruction-metrics-axe-2-1-name",
      "test-instruction-metrics-axe-2-2-name"
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=instruction_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "instruction_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "non_displayed_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given new alarm under pbehavior should add pbehavior_alarms and non_displayed_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-3",
      "connector_name": "test-connector-name-metrics-axe-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-3",
      "resource": "test-resource-metrics-axe-3",
      "state": 0
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-axe-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "pbhenter",
      "connector": "test-connector-metrics-axe-3",
      "connector_name": "test-connector-name-metrics-axe-3",
      "component": "test-component-metrics-axe-3",
      "resource": "test-resource-metrics-axe-3",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-3",
      "connector_name": "test-connector-name-metrics-axe-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-3",
      "resource": "test-resource-metrics-axe-3",
      "state": 1
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=pbehavior_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "pbehavior_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "non_displayed_alarms",
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

  @concurrent
  Scenario: given new alarm and new meta alarm should add correlation_alarms and non_displayed_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarm-metrics-axe-4-1-name",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-metrics-axe-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarm-metrics-axe-4-2-name",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-metrics-axe-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-4",
      "connector_name": "test-connector-name-metrics-axe-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-4",
      "resource": "test-resource-metrics-axe-4",
      "state": 1
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=correlation_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "correlation_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "non_displayed_alarms",
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

  @concurrent
  Scenario: given new alarm and existed meta alarm should add correlation_alarms and non_displayed_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarm-metrics-axe-5-name",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-metrics-axe-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-5-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-metrics-axe-5-1",
                "test-resource-metrics-axe-5-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-5",
      "connector_name": "test-connector-name-metrics-axe-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-5",
      "resource": "test-resource-metrics-axe-5-1",
      "state": 1
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-metrics-axe-5-1&correlation=true until response code is 200 and body contains:
    """json
    {
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
      "connector": "test-connector-metrics-axe-5",
      "connector_name": "test-connector-name-metrics-axe-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-5",
      "resource": "test-resource-metrics-axe-5-2",
      "state": 1
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=correlation_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "correlation_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "non_displayed_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given acked alarm should add ack_alarms and average_ack metrics
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-6-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-6"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-6",
      "connector_name": "test-connector-name-metrics-axe-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-6",
      "resource": "test-resource-metrics-axe-6",
      "state": 1
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-6",
      "connector_name": "test-connector-name-metrics-axe-6",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-metrics-axe-6",
      "resource": "test-resource-metrics-axe-6"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ack_alarms&parameters[]=average_ack&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "average_ack",
          "data": [
            {
              "timestamp": {{ nowDate }}
            }
          ]
        }
      ]
    }
    """
    Then the response key "data.1.data.0.value" should be greater or equal than 1
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-6",
      "connector_name": "test-connector-name-metrics-axe-6",
      "source_type": "resource",
      "event_type": "ackremove",
      "component": "test-component-metrics-axe-6",
      "resource": "test-resource-metrics-axe-6"
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-6",
      "connector_name": "test-connector-name-metrics-axe-6",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-metrics-axe-6",
      "resource": "test-resource-metrics-axe-6"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ack_alarms&parameters[]=average_ack&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "average_ack",
          "data": [
            {
              "timestamp": {{ nowDate }}
            }
          ]
        }
      ]
    }
    """
    Then the response key "data.1.data.0.value" should be greater or equal than 1

  @concurrent
  Scenario: given unacked alarm should add cancel_ack_alarms and ack_active_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-7-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-7"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-7",
      "connector_name": "test-connector-name-metrics-axe-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-7",
      "resource": "test-resource-metrics-axe-7",
      "state": 1
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-7",
      "connector_name": "test-connector-name-metrics-axe-7",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-metrics-axe-7",
      "resource": "test-resource-metrics-axe-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-7",
      "connector_name": "test-connector-name-metrics-axe-7",
      "source_type": "resource",
      "event_type": "ackremove",
      "component": "test-component-metrics-axe-7",
      "resource": "test-resource-metrics-axe-7"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=cancel_ack_alarms&parameters[]=ack_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "cancel_ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "ack_active_alarms",
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-7",
      "connector_name": "test-connector-name-metrics-axe-7",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-metrics-axe-7",
      "resource": "test-resource-metrics-axe-7"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=cancel_ack_alarms&parameters[]=ack_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "cancel_ack_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "ack_active_alarms",
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

  @concurrent
  Scenario: given alarm with ticket should add ticket_active_alarms and without_ticket_active_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-8-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-metrics-axe-8-1",
                "test-resource-metrics-axe-8-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-metrics-axe-8",
        "connector_name": "test-connector-name-metrics-axe-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-8",
        "resource": "test-resource-metrics-axe-8-1",
        "state": 1
      },
      {
        "connector": "test-connector-metrics-axe-8",
        "connector_name": "test-connector-name-metrics-axe-8",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-8",
        "resource": "test-resource-metrics-axe-8-2",
        "state": 1
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-8",
      "connector_name": "test-connector-name-metrics-axe-8",
      "source_type": "resource",
      "event_type": "assocticket",
      "component": "test-component-metrics-axe-8",
      "resource": "test-resource-metrics-axe-8-1",
      "ticket": "testticket"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-8",
      "connector_name": "test-connector-name-metrics-axe-8",
      "source_type": "resource",
      "event_type": "assocticket",
      "component": "test-component-metrics-axe-8",
      "resource": "test-resource-metrics-axe-8-1",
      "ticket": "testticket"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ticket_active_alarms&parameters[]=without_ticket_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "without_ticket_active_alarms",
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

  @concurrent
  Scenario: given resolved alarm should add average_resolve metrics
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-9-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-9"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-9",
      "connector_name": "test-connector-name-metrics-axe-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-9",
      "resource": "test-resource-metrics-axe-9",
      "state": 1
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-9",
      "connector_name": "test-connector-name-metrics-axe-9",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-metrics-axe-9",
      "resource": "test-resource-metrics-axe-9"
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-9",
      "connector_name": "test-connector-name-metrics-axe-9",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-metrics-axe-9",
      "resource": "test-resource-metrics-axe-9"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=average_resolve&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.data.0.value" is greater or equal than 1
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "average_resolve",
          "data": [
            {
              "timestamp": {{ nowDate }}
            }
          ]
        }
      ]
    }
    """
