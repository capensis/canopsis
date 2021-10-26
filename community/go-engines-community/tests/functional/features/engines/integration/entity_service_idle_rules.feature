Feature: entity_service idle_rules integration
  Scenario: given service for entity should get idle_since from dependencies
    When I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "_id": "test-entityservice-idle-since-integration",
      "name": "test-entityservice-idle-since-integration",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
            "name": "test-idle-since-integration-resource-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idle-rule-es-integration-name",
      "description": "test-idle-rule-es-integration-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "seconds": 1,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "name": "test-idle-since-integration-resource-1"
        }
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=test-idle-since-integration-resource-1
    Then the response code should be 200
    When I save response idleSince={{ (index .lastResponse.data 0).idle_since }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I do GET /api/v4/weather-services?filter={"name":"test-entityservice-idle-since-integration"}
    Then the response code should be 200
    Then the response body should contain:
    """
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
    """
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
    """
    {
      "data": [
        {
          "entity": {
            "_id": "test-idle-since-integration-resource-1/test-idle-since-integration-component",
            "idle_since": {{ .idleSince }}
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

  Scenario: given service for entity should get idle_since from depended connector
    When I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "_id": "test-entityservice-idle-since-integration-2",
      "name": "test-entityservice-idle-since-integration-2",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": "test-idle-since-integration-2-resource-1"
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idle-rule-es-integration-2-name",
      "description": "test-idle-rule-es-integration-2-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "seconds": 1,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "name": "test-idle-since-integration-2-connectorname"
        }
      ]
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/entities?search=test-idle-since-integration-2-connectorname
    Then the response code should be 200
    When I save response idleSince={{ (index .lastResponse.data 0).idle_since }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I do GET /api/v4/weather-services?filter={"name":"test-entityservice-idle-since-integration-2"}
    Then the response code should be 200
    Then the response body should contain:
    """
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
    """
    {
      "_id": "test-entityservice-idle-since-integration-3",
      "name": "test-entityservice-idle-since-integration-3",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": "test-idle-since-integration-resource-2"
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """
    {
      "_id": "test-entityservice-idle-since-integration-4",
      "name": "test-entityservice-idle-since-integration-4",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": "test-entityservice-idle-since-integration-3"
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idle-rule-es-integration-3-name",
      "description": "test-idle-rule-es-integration-3-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "seconds": 1,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "name": "test-idle-since-integration-resource-2"
        }
      ]
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=test-idle-since-integration-resource-2
    Then the response code should be 200
    When I save response idleSince={{ (index .lastResponse.data 0).idle_since }}
    When I wait the next periodical process
    When I wait the next periodical process
    When I do GET /api/v4/weather-services?filter={"name":"test-entityservice-idle-since-integration-4"}
    Then the response code should be 200
    Then the response body should contain:
    """
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
    """
    {
      "data": [
        {
          "entity": {
            "_id": "test-entityservice-idle-since-integration-3",
            "idle_since": {{ .idleSince }}
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

  Scenario: given entity service should update its idle_since from depended resources
    When I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "_id": "test-entityservice-idle-since-integration-5",
      "name": "test-entityservice-idle-since-integration-5",
      "output_template": "123",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": "test-idle-since-integration-resource-3"
        },
        {
          "name": "test-idle-since-integration-resource-4"
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idle-rule-es-integration-4-name",
      "description": "test-idle-rule-es-integration-4-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "seconds": 1,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "name": "test-idle-since-integration-resource-3"
        }
      ]
    }
    """
    When I save response ruleID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204
    When I do POST /api/v4/idle-rules:
    """
    {
      "name": "test-idle-rule-es-integration-5-name",
      "description": "test-idle-rule-es-integration-5-description",
      "type": "entity",
      "enabled": true,
      "priority": 21,
      "duration": {
        "seconds": 1,
        "unit": "s"
      },
      "entity_patterns": [
        {
          "name": "test-idle-since-integration-resource-4"
        }
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
    When I do GET /api/v4/weather-services?filter={"name":"test-entityservice-idle-since-integration-5"}
    Then the response code should be 200
    Then the response body should contain:
    """
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
    """
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
    When I do GET /api/v4/weather-services?filter={"name":"test-entityservice-idle-since-integration-5"}
    Then the response code should be 200
    Then the response body should contain:
    """
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
