Feature: entity_service idle_rules integration

  @concurrent
  Scenario: given service for entity should get idle_since from dependencies
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-idle-since-integration-1",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-idle-since-integration-1"
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
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-idle-since-integration-1-name",
      "description": "test-idle-rule-idle-since-integration-1-description",
      "type": "entity",
      "enabled": true,
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
              "value": "test-resource-idle-since-integration-1"
            }
          }
        ]
      ]
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-idle-since-integration-1",
        "connector_name": "test-connector-name-idle-since-integration-1",
        "component": "test-component-idle-since-integration-1",
        "resource": "test-resource-idle-since-integration-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do GET /api/v4/entities?search=test-resource-idle-since-integration-1
    Then the response code should be 200
    When I save response idleSince={{ (index .lastResponse.data 0).idle_since }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-entityservice-idle-since-integration-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-idle-since-integration-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-entityservice-idle-since-integration-1",
          "idle_since": {{ .idleSince }}
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
    When I do GET /api/v4/weather-services/{{ .serviceId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-idle-since-integration-1/test-component-idle-since-integration-1",
          "idle_since": {{ .idleSince }}
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
    When I do GET /api/v4/entityservice-dependencies?_id={{ .serviceId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-idle-since-integration-1/test-component-idle-since-integration-1",
          "idle_since": {{ .idleSince }}
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
  Scenario: given service for entity should get idle_since from depended connector
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-idle-since-integration-2",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-idle-since-integration-2"
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
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-idle-since-integration-2-name",
      "description": "test-idle-rule-idle-since-integration-2-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
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
              "value": "test-connector-name-idle-since-integration-2"
            }
          }
        ]
      ]
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "activate",
      "connector": "test-connector-idle-since-integration-2",
      "connector_name": "test-connector-name-idle-since-integration-2",
      "source_type": "connector"
    }
    """
    When I do GET /api/v4/entities?search=test-connector-name-idle-since-integration-2
    Then the response code should be 200
    When I save response idleSince={{ (index .lastResponse.data 0).idle_since }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-entityservice-idle-since-integration-2",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-idle-since-integration-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-entityservice-idle-since-integration-2",
          "idle_since": {{ .idleSince }}
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
  Scenario: given entity service should get idle_since from depended service
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-idle-since-integration-3-1",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-idle-since-integration-3"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId1={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId1 }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-idle-since-integration-3-2",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-idle-since-integration-3-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId2={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId2 }}",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId2 }}",
        "source_type": "service"
      }
    ]
    """
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-idle-since-integration-3-name",
      "description": "test-idle-rule-idle-since-integration-3-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
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
              "value": "test-resource-idle-since-integration-3"
            }
          }
        ]
      ]
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-idle-since-integration-3",
        "connector_name": "test-connector-name-idle-since-integration-3",
        "component": "test-component-idle-since-integration-3",
        "resource": "test-resource-idle-since-integration-3",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId1 }}"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId2 }}"
      }
    ]
    """
    When I do GET /api/v4/entities?search=test-resource-idle-since-integration-3
    Then the response code should be 200
    When I save response idleSince={{ (index .lastResponse.data 0).idle_since }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-entityservice-idle-since-integration-3",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-idle-since-integration-3-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/weather-services?filters[]={{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-entityservice-idle-since-integration-3-2",
          "idle_since": {{ .idleSince }}
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
    When I do GET /api/v4/entityservice-dependencies?_id={{ .serviceId2 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .serviceId1 }}",
          "idle_since": {{ .idleSince }}
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
  Scenario: given entity service should update its idle_since from depended resources
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-idle-since-integration-4",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-idle-since-integration-4-1",
                "test-resource-idle-since-integration-4-2"
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
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-idle-since-integration-4-1-name",
      "description": "test-idle-rule-idle-since-integration-4-1-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
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
              "value": "test-resource-idle-since-integration-4-1"
            }
          }
        ]
      ]
    }
    """
    When I save response ruleID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-idle-since-integration-4",
        "connector_name": "test-connector-name-idle-since-integration-4",
        "component": "test-component-idle-since-integration-4",
        "resource": "test-resource-idle-since-integration-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-idle-since-integration-4-2-name",
      "description": "test-idle-rule-idle-since-integration-4-2-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
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
              "value": "test-resource-idle-since-integration-4-2"
            }
          }
        ]
      ]
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-idle-since-integration-4",
        "connector_name": "test-connector-name-idle-since-integration-4",
        "component": "test-component-idle-since-integration-4",
        "resource": "test-resource-idle-since-integration-4-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I wait the next periodical process
    When I wait the next periodical process
    When I do GET /api/v4/entities?search=test-resource-idle-since-integration-4
    Then the response code should be 200
    When I save response idleSinceFirst={{ (index .lastResponse.data 0).idle_since }}
    When I save response idleSinceSecond={{ (index .lastResponse.data 1).idle_since }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-entityservice-idle-since-integration-4",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-idle-since-integration-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/weather-services?filters[]={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-entityservice-idle-since-integration-4",
          "idle_since": {{ .idleSinceFirst }}
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
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-idle-since-integration-4",
      "connector": "test-connector-idle-since-integration-4",
      "connector_name": "test-connector-name-idle-since-integration-4",
      "component": "test-component-idle-since-integration-4",
      "resource": "test-resource-idle-since-integration-4-1",
      "source_type": "resource"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "check",
        "connector": "test-connector-idle-since-integration-4",
        "connector_name": "test-connector-name-idle-since-integration-4",
        "component": "test-component-idle-since-integration-4",
        "resource": "test-resource-idle-since-integration-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "{{ .serviceId }}"
      }
    ]
    """
    When I wait the next periodical process
    When I wait the next periodical process
    When I do GET /api/v4/weather-services?filters[]={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-entityservice-idle-since-integration-4",
          "idle_since": {{ .idleSinceSecond }}
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
