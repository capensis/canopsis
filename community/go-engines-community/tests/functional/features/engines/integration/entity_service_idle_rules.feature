Feature: entity_service idle_rules integration
  Scenario: given service for entity should get idle_since from dependencies
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-idle-since-integration",
      "name": "test-entityservice-idle-since-integration",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-since-integration-resource-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-es-integration-name",
      "description": "test-idle-rule-es-integration-description",
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
              "value": "test-idle-since-integration-resource-1"
            }
          }
        ]
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=test-idle-since-integration-resource-1
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
              "value": "test-entityservice-idle-since-integration"
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
          "name": "test-entityservice-idle-since-integration",
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
    When I do GET /api/v4/weather-services/test-entityservice-idle-since-integration
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-idle-since-integration-resource-1/test-idle-since-integration-component",
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-entityservice-idle-since-integration
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-idle-since-integration-resource-1/test-idle-since-integration-component",
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

  Scenario: given service for entity should get idle_since from depended connector
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-idle-since-integration-2",
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
              "value": "test-idle-since-integration-2-resource-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-es-integration-2-name",
      "description": "test-idle-rule-es-integration-2-description",
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
              "value": "test-idle-since-integration-2-connectorname"
            }
          }
        ]
      ]
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/entities?search=test-idle-since-integration-2-connectorname
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

  Scenario: given entity service should get idle_since from depended service
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-idle-since-integration-3",
      "name": "test-entityservice-idle-since-integration-3",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-since-integration-resource-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-idle-since-integration-4",
      "name": "test-entityservice-idle-since-integration-4",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-idle-since-integration-3"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-es-integration-3-name",
      "description": "test-idle-rule-es-integration-3-description",
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
              "value": "test-idle-since-integration-resource-2"
            }
          }
        ]
      ]
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=test-idle-since-integration-resource-2
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
              "value": "test-entityservice-idle-since-integration-4"
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
          "name": "test-entityservice-idle-since-integration-4",
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
    When I do GET /api/v4/entityservice-dependencies?_id=test-entityservice-idle-since-integration-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-entityservice-idle-since-integration-3",
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

  Scenario: given entity service should update its idle_since from depended resources
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-idle-since-integration-5",
      "name": "test-entityservice-idle-since-integration-5",
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
                "test-idle-since-integration-resource-3",
                "test-idle-since-integration-resource-4"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-es-integration-4-name",
      "description": "test-idle-rule-es-integration-4-description",
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
              "value": "test-idle-since-integration-resource-3"
            }
          }
        ]
      ]
    }
    """
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idle-rule-es-integration-5-name",
      "description": "test-idle-rule-es-integration-5-description",
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
              "value": "test-idle-since-integration-resource-4"
            }
          }
        ]
      ]
    }
    """
    When I wait the end of 2 events processing
    When I wait the next periodical process
    When I wait the next periodical process
    When I do GET /api/v4/entities?search=test-idle-since-integration-resource
    Then the response code should be 200
    When I save response idleSinceFirst={{ (index .lastResponse.data 2).idle_since }}
    When I save response idleSinceSecond={{ (index .lastResponse.data 3).idle_since }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-entityservice-idle-since-integration-5",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-idle-since-integration-5"
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
          "name": "test-entityservice-idle-since-integration-5",
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
      "connector": "test-idle-since-integration-connector",
      "connector_name": "test-idle-since-integration-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-idle-since-integration-component",
      "resource": "test-idle-since-integration-resource-3",
      "state": 2,
      "output": "test-idle-since-integration-resource-3"
    }
    """
    When I wait the end of 2 events processing
    When I wait the next periodical process
    When I wait the next periodical process
    When I do GET /api/v4/weather-services?filters[]={{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-entityservice-idle-since-integration-5",
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
