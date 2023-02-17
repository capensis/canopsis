Feature: create service entity
  I need to be able to create service entity

  @concurrent
  Scenario: given disabled service entity should not update service context graph
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-second-1-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-second-1-output",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-che-service-second-1-1",
                "test-resource-che-service-second-1-2"
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
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
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
      "connector": "test-connector-che-service-second-1",
      "connector_name": "test-connector-name-che-service-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-second-1",
      "resource": "test-resource-che-service-second-1-1",
      "state": 2,
      "output": "test-output-che-service-second-1"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-che-service-second-1",
        "connector_name": "test-connector-name-che-service-second-1",
        "component": "test-component-che-service-second-1",
        "resource": "test-resource-che-service-second-1-1",
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
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-second-1-1/test-component-che-service-second-1"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-1-1/test-component-che-service-second-1
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-1/test-connector-name-che-service-second-1"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-1",
      "{{ .serviceID }}"
    ]
    """
    When I do PUT /api/v4/entityservices/{{ .serviceID }}:
    """json
    {
      "name": "test-entityservice-che-service-second-1-name",
      "impact_level": 1,
      "output_template": "test-entityservice-che-service-second-1-output",
      "enabled": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-second-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "recomputeentityservice",
      "connector": "service",
      "connector_name": "service",
      "component": "{{ .serviceID }}",
      "source_type": "service"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-second-1",
      "connector_name": "test-connector-name-che-service-second-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-second-1",
      "resource": "test-resource-che-service-second-1-2",
      "state": 1,
      "output": "test-output-che-service-second-1"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-second-1-1/test-component-che-service-second-1"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-1-1/test-component-che-service-second-1
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-1/test-connector-name-che-service-second-1"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-1",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given service with old pattern should update service
    Given I am admin
    When I do PUT /api/v4/entityservices/test-entityservice-che-service-second-2:
    """json
    {
      "name": "test-entityservice-che-service-second-2-name",
      "output_template": "test-entityservice-che-service-second-2-output",
      "impact_level": 1,
      "enabled": true,
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
        "component": "test-entityservice-che-service-second-2",
        "source_type": "service"
      },
      {
        "event_type": "check",
        "connector": "service",
        "connector_name": "service",
        "component": "test-entityservice-che-service-second-2"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-che-service-second-2",
      "connector_name": "test-connector-name-che-service-second-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-second-2",
      "resource": "test-resource-che-service-second-2",
      "state": 2,
      "output": "test-output-che-service-second-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-che-service-second-2",
        "connector_name": "test-connector-name-che-service-second-2",
        "component": "test-component-che-service-second-2",
        "resource": "test-resource-che-service-second-2",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "service",
        "connector_name": "service",
        "component": "test-entityservice-che-service-second-2"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-che-service-second-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-second-2/test-component-che-service-second-2"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-2/test-component-che-service-second-2
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-2/test-connector-name-che-service-second-2"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-2",
      "test-entityservice-che-service-second-2"
    ]
    """

  @concurrent
  Scenario: given service with corporate pattern should update service on pattern update
    Given I am admin
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-che-service-second-3",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-second-3-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response patternID={{ .lastResponse._id }}
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-second-3-name",
      "output_template": "test-entityservice-che-service-second-3-output",
      "impact_level": 1,
      "enabled": true,
      "sli_avail_state": 0,
      "corporate_entity_pattern": "{{ .patternID }}"
    }
    """
    Then the response code should be 201
    Then I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
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
    [
      {
        "connector": "test-connector-che-service-second-3",
        "connector_name": "test-connector-name-che-service-second-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-service-second-3",
        "resource": "test-resource-che-service-second-3-1",
        "state": 2,
        "output": "test-output-che-service-second-3"
      },
      {
        "connector": "test-connector-che-service-second-3",
        "connector_name": "test-connector-name-che-service-second-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-che-service-second-3",
        "resource": "test-resource-che-service-second-3-2",
        "state": 3,
        "output": "test-output-che-service-second-3"
      }
    ]
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-che-service-second-3",
        "connector_name": "test-connector-name-che-service-second-3",
        "component": "test-component-che-service-second-3",
        "resource": "test-resource-che-service-second-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "connector": "test-connector-che-service-second-3",
        "connector_name": "test-connector-name-che-service-second-3",
        "component": "test-component-che-service-second-3",
        "resource": "test-resource-che-service-second-3-2",
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
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-second-3-1/test-component-che-service-second-3"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-3-1/test-component-che-service-second-3
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-3/test-connector-name-che-service-second-3"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-3",
      "{{ .serviceID }}"
    ]
    """
    When I do PUT /api/v4/patterns/{{ .patternID }}:
    """json
    {
      "title": "test-pattern-che-service-second-3",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-che-service-second-3-1",
                "test-resource-che-service-second-3-2"
              ]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
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
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-resource-che-service-second-3-1/test-component-che-service-second-3",
      "test-resource-che-service-second-3-2/test-component-che-service-second-3"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-3-1/test-component-che-service-second-3
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-3/test-connector-name-che-service-second-3"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-3",
      "{{ .serviceID }}"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-3-2/test-component-che-service-second-3
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-3/test-connector-name-che-service-second-3"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-3",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given resource entity and new service entity should update context graph on entity disable or enable
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-second-4",
      "connector_name": "test-connector-name-che-service-second-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-second-4",
      "resource": "test-resource-che-service-second-4",
      "state": 0,
      "output": "test-output-che-service-second-4"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-second-4-name",
      "output_template": "test-entityservice-che-service-second-4-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-che-service-second-4"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceID={{ .lastResponse._id }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
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
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-second-4/test-component-che-service-second-4"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-4/test-component-che-service-second-4
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-4/test-connector-name-che-service-second-4"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-4",
      "{{ .serviceID }}"
    ]
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-che-service-second-4/test-component-che-service-second-4:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "connector": "test-connector-che-service-second-4",
      "connector_name": "test-connector-name-che-service-second-4",
      "component": "test-component-che-service-second-4",
      "resource": "test-resource-che-service-second-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-4/test-component-che-service-second-4
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-second-4/test-connector-name-che-service-second-4"
      ],
      "impact": [
        "test-component-che-service-second-4"
      ]
    }
    """
    When I do PUT /api/v4/entitybasics?_id=test-resource-che-service-second-4/test-component-che-service-second-4:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "connector": "test-connector-che-service-second-4",
      "connector_name": "test-connector-name-che-service-second-4",
      "component": "test-component-che-service-second-4",
      "resource": "test-resource-che-service-second-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-che-service-second-4/test-component-che-service-second-4"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-4/test-component-che-service-second-4
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-4/test-connector-name-che-service-second-4"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-4",
      "{{ .serviceID }}"
    ]
    """

  @concurrent
  Scenario: given resource entity and new service entity should update context graph on entity mass disable or enable
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-second-5",
      "connector_name": "test-connector-name-che-service-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-second-5",
      "resource": "test-resource-che-service-second-5-1",
      "state": 0,
      "output": "test-output-che-service-second-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-service-second-5",
      "connector_name": "test-connector-name-che-service-second-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-service-second-5",
      "resource": "test-resource-che-service-second-5-2",
      "state": 0,
      "output": "test-output-che-service-second-5"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-che-service-second-5-name",
      "output_template": "test-entityservice-che-service-second-5-output",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-che-service-second-5-1",
                "test-resource-che-service-second-5-2"
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
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "connector": "service",
        "connector_name": "service",
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
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-resource-che-service-second-5-1/test-component-che-service-second-5",
      "test-resource-che-service-second-5-2/test-component-che-service-second-5"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-5-1/test-component-che-service-second-5
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-5/test-connector-name-che-service-second-5"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-5",
      "{{ .serviceID }}"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-5-2/test-component-che-service-second-5
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-5/test-connector-name-che-service-second-5"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-5",
      "{{ .serviceID }}"
    ]
    """
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-resource-che-service-second-5-1/test-component-che-service-second-5"
      },
      {
        "_id": "test-resource-che-service-second-5-2/test-component-che-service-second-5"
      }
    ]
    """
    Then the response code should be 207
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-che-service-second-5",
        "connector_name": "test-connector-name-che-service-second-5",
        "component": "test-component-che-service-second-5",
        "resource": "test-resource-che-service-second-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-che-service-second-5",
        "connector_name": "test-connector-name-che-service-second-5",
        "component": "test-component-che-service-second-5",
        "resource": "test-resource-che-service-second-5-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-5-1/test-component-che-service-second-5
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-second-5/test-connector-name-che-service-second-5"
      ],
      "impact": [
        "test-component-che-service-second-5"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-5-2/test-component-che-service-second-5
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-connector-che-service-second-5/test-connector-name-che-service-second-5"
      ],
      "impact": [
        "test-component-che-service-second-5"
      ]
    }
    """
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-resource-che-service-second-5-1/test-component-che-service-second-5"
      },
      {
        "_id": "test-resource-che-service-second-5-2/test-component-che-service-second-5"
      }
    ]
    """
    Then the response code should be 207
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-che-service-second-5",
        "connector_name": "test-connector-name-che-service-second-5",
        "component": "test-component-che-service-second-5",
        "resource": "test-resource-che-service-second-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "entitytoggled",
        "connector": "test-connector-che-service-second-5",
        "connector_name": "test-connector-name-che-service-second-5",
        "component": "test-component-che-service-second-5",
        "resource": "test-resource-che-service-second-5-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-resource-che-service-second-5-1/test-component-che-service-second-5",
      "test-resource-che-service-second-5-2/test-component-che-service-second-5"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-5-1/test-component-che-service-second-5
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-5/test-connector-name-che-service-second-5"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-5",
      "{{ .serviceID }}"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-che-service-second-5-2/test-component-che-service-second-5
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-connector-che-service-second-5/test-connector-name-che-service-second-5"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-component-che-service-second-5",
      "{{ .serviceID }}"
    ]
    """
