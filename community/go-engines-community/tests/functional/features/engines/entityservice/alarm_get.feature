Feature: Get alarms
  I need to be able to get a alarms

  Scenario: given service should return alarms of dependencies
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-to-alarm-service-get-1-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-1",
        "state": 1,
        "output": "test-output-to-alarm-service-get-1"
      },
      {
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-2",
        "state": 3,
        "output": "test-output-to-alarm-service-get-1"
      },
      {
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-3",
        "state": 2,
        "output": "test-output-to-alarm-service-get-1"
      },
      {
        "connector": "test-connector-to-alarm-service-get-1",
        "connector_name": "test-connector-name-to-alarm-service-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-service-get-1",
        "resource": "test-resource-to-alarm-service-get-1-4",
        "state": 2,
        "output": "test-output-to-alarm-service-get-1"
      }
    ]
    """
    When I wait the end of 7 events processing
    When I do GET /api/v4/entityservice-alarms/{{ .serviceID }}?sort_by=v.resource&sort=asc
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
    When I do GET /api/v4/entityservice-alarms/{{ .serviceID }}?with_service=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "{{ .serviceID }}",
            "depends_count": 3,
            "impacts_count": 0
          },
          "v": {
            "connector": "service",
            "connector_name": "service",
            "component": "{{ .serviceID }}"
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

  Scenario: given service should return fileted alarms of dependencies
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-to-alarm-service-get-2-name",
      "output_template": "All: {{ `{{.All}}` }}; Alarms: {{ `{{.Alarms}}` }}; Acknowledged: {{ `{{.Acknowledged}}` }}; NotAcknowledged: {{ `{{.NotAcknowledged}}` }}; StateCritical: {{ `{{.State.Critical}}` }}; StateMajor: {{ `{{.State.Major}}` }}; StateMinor: {{ `{{.State.Minor}}` }}; StateInfo: {{ `{{.State.Info}}` }}; Pbehaviors: {{ `{{.PbehaviorCounters}}` }};",
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
    When I save response serviceID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-service-get-2-1",
        "resource": "test-resource-to-alarm-service-get-2-1",
        "state": 1,
        "output": "test-output-to-alarm-service-get-2"
      },
      {
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-service-get-2-2",
        "resource": "test-resource-to-alarm-service-get-2-2",
        "state": 3,
        "output": "test-output-to-alarm-service-get-2"
      },
      {
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-service-get-2-1",
        "resource": "test-resource-to-alarm-service-get-2-3",
        "state": 2,
        "output": "test-output-to-alarm-service-get-2"
      },
      {
        "connector": "test-connector-to-alarm-service-get-2",
        "connector_name": "test-connector-name-to-alarm-service-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-to-alarm-service-get-2-2",
        "resource": "test-resource-to-alarm-service-get-2-4",
        "state": 2,
        "output": "test-output-to-alarm-service-get-2"
      }
    ]
    """
    When I wait the end of 7 events processing
    When I do PUT /api/v4/entitybasics?_id=test-resource-to-alarm-service-get-2-3/test-component-to-alarm-service-get-2-1:
    """json
    {
      "category": "test-category-to-entityservice-alarm-get-2",
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0,
      "impact": [
        "test-component-to-alarm-service-get-2-1"
      ],
      "depends": [
        "test-connector-to-alarm-service-get-2/test-connector-name-to-alarm-service-get-2"
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do PUT /api/v4/entitybasics?_id=test-resource-to-alarm-service-get-2-4/test-component-to-alarm-service-get-2-2:
    """json
    {
      "category": "test-category-to-entityservice-alarm-get-2",
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0,
      "impact": [
        "test-component-to-alarm-service-get-2-2"
      ],
      "depends": [
        "test-connector-to-alarm-service-get-2/test-connector-name-to-alarm-service-get-2"
      ]
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do GET /api/v4/entityservice-alarms/{{ .serviceID }}?search=test-component-to-alarm-service-get-2-2
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
    When I do GET /api/v4/entityservice-alarms/{{ .serviceID }}?category=test-category-to-entityservice-alarm-get-2
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

  Scenario: given get dependency alarms unauth request should not allow access
    When I do GET /api/v4/entityservice-alarms/test-entityservice-not-found
    Then the response code should be 401

  Scenario: given get dependency alarms request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entityservice-alarms/test-entityservice-not-found
    Then the response code should be 403
