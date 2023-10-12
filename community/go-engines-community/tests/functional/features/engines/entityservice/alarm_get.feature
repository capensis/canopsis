Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given service should return alarms of dependencies
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-to-alarm-service-get-1-name",
      "output_template": "test-entityservice-to-alarm-service-get-1-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-alarm-service-get-1-1",
                "test-resource-to-alarm-service-get-1-2",
                "test-resource-to-alarm-service-get-1-3"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-to-alarm-service-get-1",
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 3,
        "output": "test-output-to-alarm-service-get-1",
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-to-alarm-service-get-1",
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-to-alarm-service-get-1",
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-4",
        "source_type": "resource"
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-2",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-3",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-4",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/entityservice-alarms/{{ .serviceId }}?sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-1-1/test-component-to-alarm-service-get-1",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-1",
            "connector_name": "test-connector-name-to-alarm-service-get-1",
            "component": "test-component-to-alarm-service-get-1",
            "resource": "test-resource-to-alarm-service-get-1-1"
          }
        },
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-1-2/test-component-to-alarm-service-get-1",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-1",
            "connector_name": "test-connector-name-to-alarm-service-get-1",
            "component": "test-component-to-alarm-service-get-1",
            "resource": "test-resource-to-alarm-service-get-1-2"
          }
        },
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-1-3/test-component-to-alarm-service-get-1",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-1",
            "connector_name": "test-connector-name-to-alarm-service-get-1",
            "component": "test-component-to-alarm-service-get-1",
            "resource": "test-resource-to-alarm-service-get-1-3"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I do GET /api/v4/entityservice-alarms/{{ .serviceId }}?with_service=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "{{ .serviceId }}",
            "depends_count": 3,
            "impacts_count": 0
          },
          "v": {
            "connector": "service",
            "connector_name": "service",
            "component": "{{ .serviceId }}"
          }
        },
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-1-1/test-component-to-alarm-service-get-1",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-1",
            "connector_name": "test-connector-name-to-alarm-service-get-1",
            "component": "test-component-to-alarm-service-get-1",
            "resource": "test-resource-to-alarm-service-get-1-1"
          }
        },
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-1-2/test-component-to-alarm-service-get-1",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-1",
            "connector_name": "test-connector-name-to-alarm-service-get-1",
            "component": "test-component-to-alarm-service-get-1",
            "resource": "test-resource-to-alarm-service-get-1-2"
          }
        },
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-1-3/test-component-to-alarm-service-get-1",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-1",
            "connector_name": "test-connector-name-to-alarm-service-get-1",
            "component": "test-component-to-alarm-service-get-1",
            "resource": "test-resource-to-alarm-service-get-1-3"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

  @concurrent
  Scenario: given service should return fileted alarms of dependencies
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-to-alarm-service-get-2-name",
      "output_template": "test-entityservice-to-alarm-service-get-2-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-alarm-service-get-2-1",
                "test-resource-to-alarm-service-get-2-2",
                "test-resource-to-alarm-service-get-2-3"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-to-alarm-service-get-2",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-1",
        "resource": "test-resource-to-alarm-service-get-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 3,
        "output": "test-output-to-alarm-service-get-2",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-2",
        "resource": "test-resource-to-alarm-service-get-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-to-alarm-service-get-2",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-1",
        "resource": "test-resource-to-alarm-service-get-2-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-to-alarm-service-get-2",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-2",
        "resource": "test-resource-to-alarm-service-get-2-4",
        "source_type": "resource"
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-1",
        "resource": "test-resource-to-alarm-service-get-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-2",
        "resource": "test-resource-to-alarm-service-get-2-2",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-1",
        "resource": "test-resource-to-alarm-service-get-2-3",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-2",
        "resource": "test-resource-to-alarm-service-get-2-4",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-to-alarm-service-get-2-3/test-component-to-alarm-service-get-2-1:
    """json
    {
      "category": "test-category-to-entityservice-alarm-get-2",
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "entityupdated",
      "connector": "test-connector-to-alarm-service-get-2",
      "connector_name": "test-connector-name-to-alarm-service-get-2",
      "component": "test-component-to-alarm-service-get-2-1",
      "resource": "test-resource-to-alarm-service-get-2-3",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-to-alarm-service-get-2-4/test-component-to-alarm-service-get-2-2:
    """json
    {
      "category": "test-category-to-entityservice-alarm-get-2",
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entityupdated",
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "component": "test-component-to-alarm-service-get-2-2",
        "resource": "test-resource-to-alarm-service-get-2-4",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do GET /api/v4/entityservice-alarms/{{ .serviceId }}?search=test-component-to-alarm-service-get-2-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-2-2/test-component-to-alarm-service-get-2-2"
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-2",
            "connector_name": "test-connector-name-to-alarm-service-get-2",
            "component": "test-component-to-alarm-service-get-2-2",
            "resource": "test-resource-to-alarm-service-get-2-2"
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
    When I do GET /api/v4/entityservice-alarms/{{ .serviceId }}?category=test-category-to-entityservice-alarm-get-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-2-3/test-component-to-alarm-service-get-2-1"
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-2",
            "connector_name": "test-connector-name-to-alarm-service-get-2",
            "component": "test-component-to-alarm-service-get-2-1",
            "resource": "test-resource-to-alarm-service-get-2-3"
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
    
  @concurrent
  Scenario: given get services dependencies alarms should return assigned_declare_ticket_rules for dependencies
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-to-alarm-service-get-3-name",
      "output_template": "test-entityservice-to-alarm-service-get-3-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-to-alarm-service-get-3-1",
                "test-resource-to-alarm-service-get-3-2",
                "test-resource-to-alarm-service-get-3-3"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I send an event:
    """json
    [
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-to-alarm-service-get-3",
        "connector": "test-connector-to-alarm-service-get-3",
        "connector_name": "test-connector-name-to-alarm-service-get-3",
        "component": "test-component-to-alarm-service-get-3",
        "resource": "test-resource-to-alarm-service-get-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 3,
        "output": "test-output-to-alarm-service-get-3",
        "connector": "test-connector-to-alarm-service-get-3",
        "connector_name": "test-connector-name-to-alarm-service-get-3",
        "component": "test-component-to-alarm-service-get-3",
        "resource": "test-resource-to-alarm-service-get-3-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-to-alarm-service-get-3",
        "connector": "test-connector-to-alarm-service-get-3",
        "connector_name": "test-connector-name-to-alarm-service-get-3",
        "component": "test-component-to-alarm-service-get-3",
        "resource": "test-resource-to-alarm-service-get-3-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-to-alarm-service-get-3",
        "connector": "test-connector-to-alarm-service-get-3",
        "connector_name": "test-connector-name-to-alarm-service-get-3",
        "component": "test-component-to-alarm-service-get-3",
        "resource": "test-resource-to-alarm-service-get-3-4",
        "source_type": "resource"
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-3",
        "connector_name": "test-connector-name-to-alarm-service-get-3",
        "component": "test-component-to-alarm-service-get-3",
        "resource": "test-resource-to-alarm-service-get-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-3",
        "connector_name": "test-connector-name-to-alarm-service-get-3",
        "component": "test-component-to-alarm-service-get-3",
        "resource": "test-resource-to-alarm-service-get-3-2",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-3",
        "connector_name": "test-connector-name-to-alarm-service-get-3",
        "component": "test-component-to-alarm-service-get-3",
        "resource": "test-resource-to-alarm-service-get-3-3",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-to-alarm-service-get-3",
        "connector_name": "test-connector-name-to-alarm-service-get-3",
        "component": "test-component-to-alarm-service-get-3",
        "resource": "test-resource-to-alarm-service-get-3-4",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-alarm-service-rule-3",
      "system_name": "test-alarm-service-rule-3-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom",
            "empty_response": false
          },
          "stop_on_fail": true
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-alarm-service-get-3-2"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-to-alarm-service-get-3-3"
            }
          }
        ]
      ]
    }
    """
    Then I save response ruleId={{ .lastResponse._id }}
    When I do GET /api/v4/entityservice-alarms/{{ .serviceId }}?sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-3-1/test-component-to-alarm-service-get-3",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-3",
            "connector_name": "test-connector-name-to-alarm-service-get-3",
            "component": "test-component-to-alarm-service-get-3",
            "resource": "test-resource-to-alarm-service-get-3-1"
          }
        },
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-3-2/test-component-to-alarm-service-get-3",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-3",
            "connector_name": "test-connector-name-to-alarm-service-get-3",
            "component": "test-component-to-alarm-service-get-3",
            "resource": "test-resource-to-alarm-service-get-3-2"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleId }}",
              "name": "test-alarm-service-rule-3"
            }
          ]
        },
        {
          "entity": {
            "_id": "test-resource-to-alarm-service-get-3-3/test-component-to-alarm-service-get-3",
            "depends_count": 0,
            "impacts_count": 1
          },
          "v": {
            "connector": "test-connector-to-alarm-service-get-3",
            "connector_name": "test-connector-name-to-alarm-service-get-3",
            "component": "test-component-to-alarm-service-get-3",
            "resource": "test-resource-to-alarm-service-get-3-3"
          },
          "assigned_declare_ticket_rules": [
            {
              "_id": "{{ .ruleId }}",
              "name": "test-alarm-service-rule-3"
            }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    Then the response key "data.0.assigned_declare_ticket_rules" should not exist

  @concurrent
  Scenario: given get dependency alarms unauth request should not allow access
    When I do GET /api/v4/entityservice-alarms/test-entityservice-not-found
    Then the response code should be 401

  @concurrent
  Scenario: given get dependency alarms request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entityservice-alarms/test-entityservice-not-found
    Then the response code should be 403
