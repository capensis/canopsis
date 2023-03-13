Feature: new-import entities
  I need to be able to new-import entities

  @concurrent
  Scenario: given service and new entity by new import should update service
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-new-import-partial-1",
      "name": "test-entityservice-new-import-partial-1-name",
      "output_template": "test-template",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-1"
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
        "component": "test-entityservice-new-import-partial-1"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-1"
      }
    ]
    """
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-1-source:
    """json
    [
      {
        "action": "set",
        "name": "test-component-new-import-partial-1",
        "type": "component",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "component": "test-component-new-import-partial-1"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-component-new-import-partial-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": [
        "{{ .serviceID }}"
      ]
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-new-import-partial-1"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/weather-services/{{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-new-import-partial-1",
          "import_source": "test-new-import-partial-1-source"
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
    Then the response key "data.0.imported" should be greater than 0

  @concurrent
  Scenario: given service and updated entity by new import should update service
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-2",
      "connector_name": "test-connector-name-new-import-partial-2",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-new-import-partial-2",
      "state": 1,
      "output": "test-output-new-import-partial-2"
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-new-import-partial-2",
      "name": "test-entityservice-new-import-partial-2-name",
      "output_template": "test-template",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "infos.test-component-new-import-partial-2-infos-1",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-2-infos-1-value"
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
        "component": "test-entityservice-new-import-partial-2"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-2"
      }
    ]
    """
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-2-source:
    """json
    [
      {
        "action": "set",
        "name": "test-component-new-import-partial-2",
        "type": "component",
        "infos": {
          "test-component-new-import-partial-2-infos-1": {
            "name": "test-component-new-import-partial-2-infos-1",
            "value": "test-component-new-import-partial-2-infos-1-value"
          }
        },
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "entityupdated",
        "component": "test-component-new-import-partial-2"
      },
      {
        "event_type": "activate",
        "component": "test-entityservice-new-import-partial-2"
      }
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-component-new-import-partial-2
    Then the response code should be 200
    Then the response array key "impact" should contain:
    """json
    [
      "test-connector-new-import-partial-2/test-connector-name-new-import-partial-2",
      "{{ .serviceID }}"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-connector-new-import-partial-2/test-connector-name-new-import-partial-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-new-import-partial-2"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-component-new-import-partial-2"
      ],
      "impact": []
    }
    """

  @concurrent
  Scenario: given disabled entity by new import should resolve alarm
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-3",
      "connector_name": "test-connector-name-new-import-partial-3",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-new-import-partial-3",
      "state": 1,
      "output": "test-output-new-import-partial-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-new-import-partial-3&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-new-import-partial-3"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-3-source:
    """json
    [
      {
        "action": "disable",
        "name": "test-component-new-import-partial-3",
        "type": "component",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "entitytoggled",
      "component": "test-component-new-import-partial-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-new-import-partial-3&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given deleted entity by new import should resolve alarm
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-4",
      "connector_name": "test-connector-name-new-import-partial-4",
      "source_type": "component",
      "event_type": "check",
      "component": "test-component-new-import-partial-4",
      "state": 1,
      "output": "test-output-new-import-partial-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-new-import-partial-4&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-new-import-partial-4"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-4-source:
    """json
    [
      {
        "action": "delete",
        "name": "test-component-new-import-partial-4",
        "type": "component",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "resolve_deleted",
      "component": "test-component-new-import-partial-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-new-import-partial-4&opened=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I wait 4s
    Then an entity test-component-new-import-partial-4 should not be in the db

  @concurrent
  Scenario: given deleted component by new import should resolve alarm for resources
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-5",
      "connector_name": "test-connector-name-new-import-partial-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-5",
      "resource": "test-resource-new-import-partial-5-1",
      "state": 1,
      "output": "test-output-new-import-partial-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-5",
      "connector_name": "test-connector-name-new-import-partial-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-5",
      "resource": "test-resource-new-import-partial-5-2",
      "state": 1,
      "output": "test-output-new-import-partial-5"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-5-1&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-new-import-partial-5-1"
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
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-5-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-resource-new-import-partial-5-2"
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-5-source:
    """json
    [
      {
        "action": "delete",
        "name": "test-component-new-import-partial-5",
        "type": "component",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "resolve_deleted",
        "resource": "test-resource-new-import-partial-5-1"
      },
      {
        "event_type": "resolve_deleted",
        "resource": "test-resource-new-import-partial-5-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-5-1&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-5-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I wait 4s
    Then an entity test-component-new-import-partial-5 should not be in the db
    Then an entity test-resource-new-import-partial-5-1/test-component-new-import-partial-5 should not be in the db
    Then an entity test-resource-new-import-partial-5-2/test-component-new-import-partial-5 should not be in the db

  @concurrent
  Scenario: given delete import action should delete component and resources which should update service state
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-new-import-partial-6",
      "name": "test-entityservice-new-import-partial-6-name",
      "output_template": "test-template",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-6-1"
            }
          }
        ],
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-6-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "test-entityservice-new-import-partial-6"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-6"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-6",
      "connector_name": "test-connector-name-new-import-partial-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-6-1",
      "resource": "test-resource-new-import-partial-6-1",
      "state": 1,
      "output": "test-output-import-partial-6"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "resource": "test-resource-new-import-partial-6-1"
      },
      {
        "event_type": "activate",
        "component": "test-entityservice-new-import-partial-6"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-6",
      "connector_name": "test-connector-name-new-import-partial-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-6-1",
      "resource": "test-resource-new-import-partial-6-2",
      "state": 1,
      "output": "test-output-import-partial-6"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "resource": "test-resource-new-import-partial-6-2"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-6"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-6",
      "connector_name": "test-connector-name-new-import-partial-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-6-2",
      "resource": "test-resource-new-import-partial-6-3",
      "state": 3,
      "output": "test-output-import-partial-6"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "resource": "test-resource-new-import-partial-6-3"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-6"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-6",
      "connector_name": "test-connector-name-new-import-partial-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-6-2",
      "resource": "test-resource-new-import-partial-6-4",
      "state": 3,
      "output": "test-output-import-partial-6"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "resource": "test-resource-new-import-partial-6-4"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-6"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-entityservice-new-import-partial-6&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-entityservice-new-import-partial-6",
            "state": {
              "val": 3
            }
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
    When I do GET /api/v4/weather-services/test-entityservice-new-import-partial-6?sort_by=name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-component-new-import-partial-6-2"
        },
        {
          "_id": "test-resource-new-import-partial-6-1/test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-resource-new-import-partial-6-2/test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-resource-new-import-partial-6-3/test-component-new-import-partial-6-2"
        },
        {
          "_id": "test-resource-new-import-partial-6-4/test-component-new-import-partial-6-2"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 6
      }
    }
    """
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-6-source:
    """json
    [
      {
        "action": "delete",
        "name": "test-component-new-import-partial-6-2",
        "type": "component",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "resolve_deleted",
        "resource": "test-resource-new-import-partial-6-3"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-6"
      },
      {
        "event_type": "resolve_deleted",
        "resource": "test-resource-new-import-partial-6-4"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-6"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-entityservice-new-import-partial-6&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-entityservice-new-import-partial-6",
            "state": {
              "val": 1
            }
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
    When I do GET /api/v4/weather-services/test-entityservice-new-import-partial-6?sort_by=name&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-resource-new-import-partial-6-1/test-component-new-import-partial-6-1"
        },
        {
          "_id": "test-resource-new-import-partial-6-2/test-component-new-import-partial-6-1"
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
    When I wait 4s
    Then an entity test-component-new-import-partial-6 should not be in the db
    Then an entity test-resource-new-import-partial-6-3/test-component-new-import-partial-6 should not be in the db
    Then an entity test-resource-new-import-partial-6-4/test-component-new-import-partial-6 should not be in the db

  @concurrent
  Scenario: given delete import action should delete component and resources which should update metaalarm state
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "_id": "test-metaalarm-new-import-partial-7",
      "name": "test-metaalarm-new-import-partial-7-name",
      "type": "complex",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-7-1"
            }
          }
        ],
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-7-2"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 4
      }
    }
    """
    Then the response code should be 201
    When I wait 4s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-7",
      "connector_name": "test-connector-name-new-import-partial-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-7-1",
      "resource": "test-resource-new-import-partial-7-1",
      "state": 1,
      "output": "test-output-import-partial-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-7",
      "connector_name": "test-connector-name-new-import-partial-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-7-1",
      "resource": "test-resource-new-import-partial-7-2",
      "state": 1,
      "output": "test-output-import-partial-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-7",
      "connector_name": "test-connector-name-new-import-partial-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-7-2",
      "resource": "test-resource-new-import-partial-7-3",
      "state": 3,
      "output": "test-output-import-partial-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-7",
      "connector_name": "test-connector-name-new-import-partial-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-7-2",
      "resource": "test-resource-new-import-partial-7-4",
      "state": 3,
      "output": "test-output-import-partial-7"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-7-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-metaalarm-new-import-partial-7-name"
          },
          "v": {
            "state": {
              "val": 3
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-new-import-partial-7-1",
                  "resource": "test-resource-new-import-partial-7-1"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-7-1",
                  "resource": "test-resource-new-import-partial-7-2"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-7-2",
                  "resource": "test-resource-new-import-partial-7-3"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-7-2",
                  "resource": "test-resource-new-import-partial-7-4"
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
        }
      }
    ]
    """
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-7-source:
    """json
    [
      {
        "action": "delete",
        "name": "test-component-new-import-partial-7-2",
        "type": "component",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "resolve_deleted",
        "resource": "test-resource-new-import-partial-7-3"
      },
      {
        "event_type": "resolve_deleted",
        "resource": "test-resource-new-import-partial-7-4"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-7-1&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-metaalarm-new-import-partial-7-name"
          },
          "v": {
            "state": {
              "val": 1
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-new-import-partial-7-1",
                  "resource": "test-resource-new-import-partial-7-1"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-7-1",
                  "resource": "test-resource-new-import-partial-7-2"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """
    When I wait 4s
    Then an entity test-component-new-import-partial-7 should not be in the db
    Then an entity test-resource-new-import-partial-7-3/test-component-new-import-partial-7 should not be in the db
    Then an entity test-resource-new-import-partial-7-4/test-component-new-import-partial-7 should not be in the db

  @concurrent
  Scenario: given delete import action should delete component and resources which should resolve metaalarm
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "_id": "test-metaalarm-new-import-partial-8",
      "name": "test-metaalarm-new-import-partial-8-name",
      "type": "complex",
      "auto_resolve": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-new-import-partial-8-1"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 2
      }
    }
    """
    Then the response code should be 201
    When I wait 4s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-8",
      "connector_name": "test-connector-name-new-import-partial-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-8-1",
      "resource": "test-resource-new-import-partial-8-1",
      "state": 1,
      "output": "test-output-import-partial-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-8",
      "connector_name": "test-connector-name-new-import-partial-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-8-1",
      "resource": "test-resource-new-import-partial-8-2",
      "state": 1,
      "output": "test-output-import-partial-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-8-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-metaalarm-new-import-partial-8-name"
          },
          "v": {
            "state": {
              "val": 1
            }
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "component": "test-component-new-import-partial-8-1",
                  "resource": "test-resource-new-import-partial-8-1"
                }
              },
              {
                "v": {
                  "component": "test-component-new-import-partial-8-1",
                  "resource": "test-resource-new-import-partial-8-2"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 2
            }
          }
        }
      }
    ]
    """
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-8-source:
    """json
    [
      {
        "action": "delete",
        "name": "test-component-new-import-partial-8-1",
        "type": "component",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "resolve_deleted",
        "resource": "test-resource-new-import-partial-8-1"
      },
      {
        "event_type": "resolve_deleted",
        "resource": "test-resource-new-import-partial-8-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-new-import-partial-8-1&correlation=true until response code is 200 and response key "data.0.v.resolved" is greater or equal than 1
    When I wait 4s
    Then an entity test-component-new-import-partial-8 should not be in the db
    Then an entity test-resource-new-import-partial-8-1/test-component-new-import-partial-8 should not be in the db
    Then an entity test-resource-new-import-partial-8-2/test-component-new-import-partial-8 should not be in the db

  @concurrent
  Scenario: given delete import action should delete service and resolve it's alarm
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-new-import-partial-9-1",
      "name": "test-entityservice-new-import-partial-9-1-name",
      "output_template": "test-template",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-new-import-partial-9-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-entityservice-new-import-partial-9-2",
      "name": "test-entityservice-new-import-partial-9-2-name",
      "output_template": "test-template",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-new-import-partial-9-1-name"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-new-import-partial-9-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "test-entityservice-new-import-partial-9-1"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-9-1"
      },
      {
        "event_type": "recomputeentityservice",
        "component": "test-entityservice-new-import-partial-9-2"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-9-2"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-new-import-partial-9",
      "connector_name": "test-connector-name-new-import-partial-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-9",
      "resource": "test-resource-new-import-partial-9-1",
      "state": 3,
      "output": "test-output-import-partial-9"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-new-import-partial-9",
        "connector_name": "test-connector-name-new-import-partial-9",
        "component": "test-component-new-import-partial-9",
        "resource": "test-resource-new-import-partial-9-1",
        "source_type": "resource"
      },
      {
        "event_type": "activate",
        "component": "test-entityservice-new-import-partial-9-1"
      },
      {
        "event_type": "activate",
        "component": "test-entityservice-new-import-partial-9-2"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-new-import-partial-9",
      "connector_name": "test-connector-name-new-import-partial-9",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-new-import-partial-9",
      "resource": "test-resource-new-import-partial-9-2",
      "state": 1,
      "output": "test-output-import-partial-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-entityservice-new-import-partial-9-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-entityservice-new-import-partial-9-2",
            "state": {
              "val": 3
            }
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
    When I do PUT /api/v4/contextgraph-import-partial?source=test-new-import-partial-9-source:
    """json
    [
      {
        "action": "delete",
        "name": "test-entityservice-new-import-partial-9-1-name",
        "type": "service",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
       "status": "done"
    }
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "recomputeentityservice",
        "component": "test-entityservice-new-import-partial-9-1"
      },
      {
        "event_type": "check",
        "component": "test-entityservice-new-import-partial-9-2"
      }
    ]
    """
    When I do GET /api/v4/entityservices/test-entityservice-new-import-partial-9-1
    Then the response code should be 404
    When I do GET /api/v4/entityservice-dependencies?_id=test-entityservice-new-import-partial-9-1
    Then the response code should be 404
    When I do GET /api/v4/entityservice-impacts?_id=test-entityservice-new-import-partial-9-1
    Then the response code should be 404
    When I do GET /api/v4/alarms?search=test-entityservice-new-import-partial-9-2&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-entityservice-new-import-partial-9-2",
            "state": {
              "val": 1
            }
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
    When I wait 4s
    Then an entity test-entityservice-new-import-partial-9-1 should not be in the db
