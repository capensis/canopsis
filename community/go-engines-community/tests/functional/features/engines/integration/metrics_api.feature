Feature: Entities and users should be synchronized in metrics db
  I need to be able to see metrics.

  @concurrent
  Scenario: given updated entity should get metrics by updated entity
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-api-1-1-name",
      "entity_pattern": [
        [
          {
            "field": "infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-metrics-api-1"
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
      "name": "test-filter-metrics-api-1-2-name",
      "entity_pattern": [
        [
          {
            "field": "infos.client",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-client-metrics-api-1-updated"
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
              "value": "test-resource-metrics-api-1"
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
      "description": "test-eventfilter-metrics-api-1-description",
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
      "client": "test-client-metrics-api-1",
      "connector": "test-connector-metrics-api-1",
      "connector_name": "test-connector-name-metrics-api-1",
      "component": "test-component-metrics-api-1",
      "resource": "test-resource-metrics-api-1",
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-metrics-api-1/test-component-metrics-api-1:
    """json
    {
      "enabled": true,
      "impact_level": 3,
      "sli_avail_state": 1,
      "infos": [
        {
          "name": "client",
          "value": "test-client-metrics-api-1-updated"
        }
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entityupdated",
      "component": "test-component-metrics-api-1",
      "resource": "test-resource-metrics-api-1",
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
  Scenario: given updated user should get metrics by updated user
    Given I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-metrics-api-2-name",
      "email": "test-user-metrics-api-2@canopsis.net",
      "roles": ["admin"],
      "enable": true,
      "password": "test-password"
    }
    """
    Then the response code should be 201
    When I save response userID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-api-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-api-2"
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
      "event_type": "check",
      "state": 1,
      "connector": "test-connector-metrics-api-2",
      "connector_name": "test-connector-name-metrics-api-2",
      "component": "test-component-metrics-api-2",
      "resource": "test-resource-metrics-api-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ack",
      "user_id": "{{ .userID }}",
      "initiator": "user",
      "connector": "test-connector-metrics-api-2",
      "connector_name": "test-connector-name-metrics-api-2",
      "component": "test-component-metrics-api-2",
      "resource": "test-resource-metrics-api-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/cat/metrics/rating?filter={{ .filterID }}&metric=ack_alarms&criteria=3&from={{ nowDateTz }}&to={{ nowDateTz }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "label": "test-user-metrics-api-2-name",
          "value": 1
        }
      ]
    }
    """
    When I do PUT /api/v4/users/{{ .userID }}:
    """json
    {
      "name": "test-user-metrics-api-2-name-updated",
      "email": "test-user-metrics-api-2@canopsis.net",
      "roles": ["admin"],
      "enable": true
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/metrics/rating?filter={{ .filterID }}&metric=ack_alarms&criteria=3&from={{ nowDateTz }}&to={{ nowDateTz }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "label": "test-user-metrics-api-2-name-updated",
          "value": 1
        }
      ]
    }
    """

  @concurrent
  Scenario: given deleted user should not get metrics by deleted user
    Given I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-metrics-api-3-name",
      "email": "test-user-metrics-api-3@canopsis.net",
      "roles": ["admin"],
      "enable": true,
      "password": "test-password"
    }
    """
    Then the response code should be 201
    When I save response userID={{ .lastResponse._id }}
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-api-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-api-3"
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
      "event_type": "check",
      "state": 1,
      "connector": "test-connector-metrics-api-3",
      "connector_name": "test-connector-name-metrics-api-3",
      "component": "test-component-metrics-api-3",
      "resource": "test-resource-metrics-api-3",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ack",
      "user_id": "{{ .userID }}",
      "initiator": "user",
      "connector": "test-connector-metrics-api-3",
      "connector_name": "test-connector-name-metrics-api-3",
      "component": "test-component-metrics-api-3",
      "resource": "test-resource-metrics-api-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/cat/metrics/rating?filter={{ .filterID }}&metric=ack_alarms&criteria=3&from={{ nowDateTz }}&to={{ nowDateTz }} until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "label": "test-user-metrics-api-3-name",
          "value": 1
        }
      ]
    }
    """
    When I do DELETE /api/v4/users/{{ .userID }}
    Then the response code should be 204
    When I do GET /api/v4/cat/metrics/rating?filter={{ .filterID }}&metric=ack_alarms&criteria=3&from={{ nowDateTz }}&to={{ nowDateTz }} until response code is 200 and body contains:
    """json
    {
      "data": []
    }
    """

  @concurrent
  Scenario: given created service should get metrics by created entity
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-api-4-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-metrics-api-4-name"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-metrics-api-4-name",
      "output_template": "test-entityservice-metrics-api-4-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-api-4"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 1,
      "connector": "test-connector-metrics-api-4",
      "connector_name": "test-connector-name-metrics-api-4",
      "component": "test-component-metrics-api-4",
      "resource": "test-resource-metrics-api-4",
      "source_type": "resource"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-metrics-api-4",
        "connector_name": "test-connector-name-metrics-api-4",
        "component": "test-component-metrics-api-4",
        "resource": "test-resource-metrics-api-4",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filterID }}",
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

  @concurrent
  Scenario: given updated service should get metrics by updated entity
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-api-5-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-metrics-api-5-name"
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
      "name": "test-filter-metrics-api-5-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-metrics-api-5-name-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filter2ID={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-metrics-api-5-name",
      "output_template": "test-entityservice-metrics-api-5-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-api-5"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 1,
      "connector": "test-connector-metrics-api-5",
      "connector_name": "test-connector-name-metrics-api-5",
      "component": "test-component-metrics-api-5",
      "resource": "test-resource-metrics-api-5",
      "source_type": "resource"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-metrics-api-5",
        "connector_name": "test-connector-name-metrics-api-5",
        "component": "test-component-metrics-api-5",
        "resource": "test-resource-metrics-api-5",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
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
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-metrics-api-5-name-updated",
      "output_template": "test-entityservice-metrics-api-5-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-api-5"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entityupdated",
      "component": "{{ .serviceID }}"
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
  Scenario: given deleted service should get metrics by deleted entity
    Given I am admin
    When I do POST /api/v4/cat/kpi-filters:
    """json
    {
      "name": "test-filter-metrics-api-6-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-metrics-api-6-name"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-metrics-api-6-name",
      "output_template": "test-entityservice-metrics-api-6-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-metrics-api-6"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "{{ .serviceID }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 1,
      "connector": "test-connector-metrics-api-6",
      "connector_name": "test-connector-name-metrics-api-6",
      "component": "test-component-metrics-api-6",
      "resource": "test-resource-metrics-api-6",
      "source_type": "resource"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-metrics-api-6",
        "connector_name": "test-connector-name-metrics-api-6",
        "component": "test-component-metrics-api-6",
        "resource": "test-resource-metrics-api-6",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceID }}"
      }
    ]
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filterID }}",
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
    When I do DELETE /api/v4/entityservices/{{ .serviceID }}
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "recomputeentityservice",
      "component": "{{ .serviceID }}",
      "source_type": "service"
    }
    """
    When I save request:
    """json
    {
      "parameters": [
        {"metric": "created_alarms"}
      ],
      "filter": "{{ .filterID }}",
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
