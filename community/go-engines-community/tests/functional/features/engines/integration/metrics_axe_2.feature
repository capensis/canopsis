Feature: Metrics should be added on alarm changes
  I need to be able to see metrics.

  @concurrent
  Scenario: given new alarm with auto instruction and pbehavior should add non_displayed_alarms metrics only once
    Given I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["create"],
      "name": "test-instruction-metrics-axe-second-1-name",
      "description": "test-instruction-metrics-axe-second-1-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-second-1"
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
      "name": "test-filter-metrics-axe-second-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-second-1"
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
      "connector": "test-connector-metrics-axe-second-1",
      "connector_name": "test-connector-name-metrics-axe-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-second-1",
      "resource": "test-resource-metrics-axe-second-1",
      "state": 0
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-axe-second-1",
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
              "value": "test-resource-metrics-axe-second-1"
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
      "connector": "test-connector-metrics-axe-second-1",
      "connector_name": "test-connector-name-metrics-axe-second-1",
      "component": "test-component-metrics-axe-second-1",
      "resource": "test-resource-metrics-axe-second-1",
      "source_type": "resource"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-metrics-axe-second-1",
      "connector_name": "test-connector-name-metrics-axe-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-second-1",
      "resource": "test-resource-metrics-axe-second-1",
      "state": 1
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-metrics-axe-second-1",
        "connector_name": "test-connector-name-metrics-axe-second-1",
        "component": "test-component-metrics-axe-second-1",
        "resource": "test-resource-metrics-axe-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-second-1",
        "connector_name": "test-connector-name-metrics-axe-second-1",
        "component": "test-component-metrics-axe-second-1",
        "resource": "test-resource-metrics-axe-second-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-second-1",
        "connector_name": "test-connector-name-metrics-axe-second-1",
        "component": "test-component-metrics-axe-second-1",
        "resource": "test-resource-metrics-axe-second-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=instruction_alarms&parameters[]=pbehavior_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "instruction_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
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
  Scenario: given new alarm with meta alarm and pbehavior should add non_displayed_alarms metrics only once
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarm-metrics-axe-second-2-name",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-metrics-axe-second-2"
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
      "name": "test-filter-metrics-axe-second-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-second-2"
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
      "connector": "test-connector-metrics-axe-second-2",
      "connector_name": "test-connector-name-metrics-axe-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-second-2",
      "resource": "test-resource-metrics-axe-second-2",
      "state": 0
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-metrics-axe-second-2",
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
              "value": "test-resource-metrics-axe-second-2"
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
      "connector": "test-connector-metrics-axe-second-2",
      "connector_name": "test-connector-name-metrics-axe-second-2",
      "component": "test-component-metrics-axe-second-2",
      "resource": "test-resource-metrics-axe-second-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-second-2",
      "connector_name": "test-connector-name-metrics-axe-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-second-2",
      "resource": "test-resource-metrics-axe-second-2",
      "state": 1
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=correlation_alarms&parameters[]=pbehavior_alarms&parameters[]=non_displayed_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
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
  Scenario: given resolved alarm should decrease active_alarms, ratio_instructions, ratio_tickets, ratio_non_displayed, ack_active_alarms metrics
    Given I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "triggers": ["create"],
      "name": "test-instruction-metrics-axe-second-3-name",
      "description": "test-instruction-metrics-axe-second-3-description",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-second-3-1"
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
      "name": "test-filter-metrics-axe-second-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-metrics-axe-second-3-1",
                "test-resource-metrics-axe-second-3-2"
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
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-1",
        "state": 1
      },
      {
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-2",
        "state": 1
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "trigger",
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-2",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "source_type": "resource",
        "event_type": "assocticket",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-1",
        "ticket": "testticket"
      },
      {
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "source_type": "resource",
        "event_type": "assocticket",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-2",
        "ticket": "testticket"
      },
      {
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "source_type": "resource",
        "event_type": "ack",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-1"
      },
      {
        "connector": "test-connector-metrics-axe-second-3",
        "connector_name": "test-connector-name-metrics-axe-second-3",
        "source_type": "resource",
        "event_type": "ack",
        "component": "test-component-metrics-axe-second-3",
        "resource": "test-resource-metrics-axe-second-3-2"
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=created_alarms&parameters[]=active_alarms&parameters[]=instruction_alarms&parameters[]=ratio_instructions&parameters[]=non_displayed_alarms&parameters[]=ratio_non_displayed&parameters[]=ticket_active_alarms&parameters[]=ratio_tickets&parameters[]=ack_alarms&parameters[]=ack_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "instruction_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "ratio_instructions",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 50
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
        },
        {
          "title": "ratio_non_displayed",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 50
            }
          ]
        },
        {
          "title": "ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 100
            }
          ]
        },
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
          "title": "ack_active_alarms",
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-second-3",
      "connector_name": "test-connector-name-metrics-axe-second-3",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-metrics-axe-second-3",
      "resource": "test-resource-metrics-axe-second-3-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-second-3",
      "connector_name": "test-connector-name-metrics-axe-second-3",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-metrics-axe-second-3",
      "resource": "test-resource-metrics-axe-second-3-1"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=created_alarms&parameters[]=active_alarms&parameters[]=instruction_alarms&parameters[]=ratio_instructions&parameters[]=non_displayed_alarms&parameters[]=ratio_non_displayed&parameters[]=ticket_active_alarms&parameters[]=ratio_tickets&parameters[]=ack_alarms&parameters[]=ack_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "created_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 2
            }
          ]
        },
        {
          "title": "active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "instruction_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 1
            }
          ]
        },
        {
          "title": "ratio_instructions",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
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
        },
        {
          "title": "ratio_non_displayed",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 0
            }
          ]
        },
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
          "title": "ratio_tickets",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 100
            }
          ]
        },
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
  Scenario: given resolved alarm shouldn't decrease ratio_correlation metrics
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarm-metrics-axe-second-4-name",
      "type": "attribute",
      "auto_resolve": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-metrics-axe-second-4"
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
      "name": "test-filter-metrics-axe-second-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-metrics-axe-second-4-1",
                "test-resource-metrics-axe-second-4-2"
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
        "connector": "test-connector-metrics-axe-second-4",
        "connector_name": "test-connector-name-metrics-axe-second-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-second-4",
        "resource": "test-resource-metrics-axe-second-4-1",
        "state": 1
      },
      {
        "connector": "test-connector-metrics-axe-second-4",
        "connector_name": "test-connector-name-metrics-axe-second-4",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-second-4",
        "resource": "test-resource-metrics-axe-second-4-2",
        "state": 1
      }
    ]
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&&parameters[]=correlation_alarms&parameters[]=ratio_correlation&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
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
          "title": "ratio_correlation",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 100
            }
          ]
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-second-4",
      "connector_name": "test-connector-name-metrics-axe-second-4",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-metrics-axe-second-4",
      "resource": "test-resource-metrics-axe-second-4-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-second-4",
      "connector_name": "test-connector-name-metrics-axe-second-4",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-metrics-axe-second-4",
      "resource": "test-resource-metrics-axe-second-4-1"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&&parameters[]=correlation_alarms&parameters[]=ratio_correlation&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
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
          "title": "ratio_correlation",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 100
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given alarm with ticket should add ticket_active_alarms and without_ticket_active_alarms metrics for user
    Given I am admin
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-metrics-axe-second-5-1-name",
      "priority": 100110,
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-metrics-axe-second-5-2"
                }
              }
            ]
          ],
          "type": "assocticket",
          "parameters": {
            "ticket": "test-scenario-metrics-axe-second-5-ticket",
            "output": "test-scenario-metrics-axe-second-5-output",
            "author": "test-author-metrics-axe-second-5"
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-metrics-axe-second-5-2-name",
      "priority": 100111,
      "enabled": true,
      "triggers": ["create"],
      "actions": [
        {
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-resource-metrics-axe-second-5-3"
                }
              }
            ]
          ],
          "type": "webhook",
          "parameters": {
            "request": {
              "method": "POST",
              "url": "{{ .dummyApiURL }}/webhook/request",
              "payload": "{\"ticket_id\":\"testticket\"}"
            },
            "declare_ticket": {
              "ticket_id": "ticket_id"
            }
          },
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declareticket-rule-axe-second-5-name",
      "system_name": "test-declareticket-rule-axe-second-5-system-name",
      "enabled": true,
      "emit_trigger": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-second-5-4"
            }
          }
        ]
      ],
      "webhooks": [
        {
          "request": {
            "method": "POST",
            "url": "{{ .dummyApiURL }}/webhook/request",
            "payload": "{\"ticket_id\":\"testticket\"}"
          },
          "declare_ticket": {
            "ticket_id": "ticket_id"
          },
          "stop_on_fail": true
        }
      ]
    }
    """
    Then the response code should be 201
    Then I save response declareTicketRuleID={{ .lastResponse._id }}
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-axe-second-5-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 100,
      "duration": {
        "value": 1,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-second-5-5"
            }
          }
        ]
      ],
      "operation": {
        "type": "assocticket",
        "parameters": {
          "ticket": "test-scenario-metrics-axe-second-5-ticket",
          "output": "test-scenario-metrics-axe-second-5-output",
          "author": "test-author-metrics-axe-second-5",
          "user": "test-user-metrics-axe-second-5"
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-second-5-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-metrics-axe-second-5-1",
                "test-resource-metrics-axe-second-5-2",
                "test-resource-metrics-axe-second-5-3",
                "test-resource-metrics-axe-second-5-4",
                "test-resource-metrics-axe-second-5-5"
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
        "connector": "test-connector-metrics-axe-second-5",
        "connector_name": "test-connector-name-metrics-axe-second-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-second-5",
        "resource": "test-resource-metrics-axe-second-5-1",
        "state": 1
      },
      {
        "connector": "test-connector-metrics-axe-second-5",
        "connector_name": "test-connector-name-metrics-axe-second-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-second-5",
        "resource": "test-resource-metrics-axe-second-5-2",
        "state": 1
      },
      {
        "connector": "test-connector-metrics-axe-second-5",
        "connector_name": "test-connector-name-metrics-axe-second-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-second-5",
        "resource": "test-resource-metrics-axe-second-5-3",
        "state": 1
      },
      {
        "connector": "test-connector-metrics-axe-second-5",
        "connector_name": "test-connector-name-metrics-axe-second-5",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-metrics-axe-second-5",
        "resource": "test-resource-metrics-axe-second-5-4",
        "state": 1
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-metrics-axe-second-5",
      "connector_name": "test-connector-name-metrics-axe-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-second-5",
      "resource": "test-resource-metrics-axe-second-5-5",
      "state": 1
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-metrics-axe-second-5",
        "connector_name": "test-connector-name-metrics-axe-second-5",
        "component": "test-component-metrics-axe-second-5",
        "resource": "test-resource-metrics-axe-second-5-5",
        "source_type": "resource"
      },
      {
        "event_type": "assocticket",
        "connector": "test-connector-metrics-axe-second-5",
        "connector_name": "test-connector-name-metrics-axe-second-5",
        "component": "test-component-metrics-axe-second-5",
        "resource": "test-resource-metrics-axe-second-5-5",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-second-5",
      "connector_name": "test-connector-name-metrics-axe-second-5",
      "source_type": "resource",
      "event_type": "assocticket",
      "component": "test-component-metrics-axe-second-5",
      "resource": "test-resource-metrics-axe-second-5-1",
      "ticket": "testticket",
      "user_id": "test-user-metrics-axe-second-5",
      "author": "test-user-metrics-axe-second-5",
      "initiator": "user"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-metrics-axe-second-5-4
    Then the response code should be 200
    When I save response alarmID4={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/bulk/declare-ticket-executions:
    """json
    [
      {
        "_id": "{{ .declareTicketRuleID }}",
        "alarms": [
          "{{ .alarmID4 }}"
        ]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200
      }
    ]
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "trigger",
      "connector": "test-connector-metrics-axe-second-5",
      "connector_name": "test-connector-name-metrics-axe-second-5",
      "component": "test-component-metrics-axe-second-5",
      "resource": "test-resource-metrics-axe-second-5-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=ticket_active_alarms&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "title": "ticket_active_alarms",
          "data": [
            {
              "timestamp": {{ nowDate }},
              "value": 5
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/cat/metrics/rating?filter={{ .filterID }}&metric=ticket_active_alarms&criteria=3&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "label": "test-user-metrics-axe-second-5",
          "value": 1
        }
      ]
    }
    """

  @concurrent
  Scenario: given double acked alarm should affect metrics only one time
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-second-6-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-second-6"
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
      "connector": "test-connector-metrics-axe-second-6",
      "connector_name": "test-connector-name-metrics-axe-second-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-metrics-axe-second-6",
      "resource": "test-resource-metrics-axe-second-6",
      "state": 1
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-metrics-axe-second-6",
      "connector_name": "test-connector-name-metrics-axe-second-6",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-metrics-axe-second-6",
      "resource": "test-resource-metrics-axe-second-6"
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
      "connector": "test-connector-metrics-axe-second-6",
      "connector_name": "test-connector-name-metrics-axe-second-6",
      "source_type": "resource",
      "event_type": "ack",
      "component": "test-component-metrics-axe-second-6",
      "resource": "test-resource-metrics-axe-second-6"
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

  Scenario: given acked alarm should add min_ack and max_ack metrics
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-axe-second-7-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-axe-second-7"
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
      "connector" : "test-connector-metrics-axe-second-7",
      "connector_name" : "test-connector-name-metrics-axe-second-7",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-metrics-axe-second-7",
      "resource" : "test-resource-metrics-axe-second-7",
      "state" : 1
    }
    """
    When I wait 2s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-metrics-axe-second-7",
      "connector_name" : "test-connector-name-metrics-axe-second-7",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" : "test-component-metrics-axe-second-7",
      "resource" : "test-resource-metrics-axe-second-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-metrics-axe-second-7",
      "connector_name" : "test-connector-name-metrics-axe-second-7",
      "source_type" : "resource",
      "event_type" : "ackremove",
      "component" : "test-component-metrics-axe-second-7",
      "resource" : "test-resource-metrics-axe-second-7"
    }
    """
    When I wait 4s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-metrics-axe-second-7",
      "connector_name" : "test-connector-name-metrics-axe-second-7",
      "source_type" : "resource",
      "event_type" : "ack",
      "component" : "test-component-metrics-axe-second-7",
      "resource" : "test-resource-metrics-axe-second-7"
    }
    """
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=min_ack&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.data.0.value" is greater or equal than 2
    When I do GET /api/v4/cat/metrics/alarm?filter={{ .filterID }}&parameters[]=max_ack&sampling=day&from={{ nowDate }}&to={{ nowDate }} until response code is 200 and response key "data.0.data.0.value" is greater or equal than 4
