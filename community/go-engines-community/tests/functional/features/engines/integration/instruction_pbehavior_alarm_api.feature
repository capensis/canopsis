Feature: Get alarms
  I need to be able to get a alarms

  Scenario: given get search request should return assigned instructions for the alarm by pbehavior
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-instruction-pbehavior-alarm-api-1-1",
                "test-resource-instruction-pbehavior-alarm-api-1-2",
                "test-resource-instruction-pbehavior-alarm-api-1-3"
              ]
            }
          }
        ]
      ],
      "description": "test-instruction-instruction-pbehavior-alarm-api-1-1-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-step-1",
          "operations": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-instruction-pbehavior-alarm-api-1-1-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-instruction-pbehavior-alarm-api-1-1-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-instruction-pbehavior-alarm-api-1-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-instruction-pbehavior-alarm-api-1-1",
                "test-resource-instruction-pbehavior-alarm-api-1-2",
                "test-resource-instruction-pbehavior-alarm-api-1-3"
              ]
            }
          }
        ]
      ],
      "active_on_pbh": ["test-maintenance-type-to-engine"],
      "description": "test-instruction-instruction-pbehavior-alarm-api-1-2-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-instruction-pbehavior-alarm-api-1-2-step-1",
          "operations": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-2-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-instruction-pbehavior-alarm-api-1-2-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-instruction-pbehavior-alarm-api-1-2-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID2={{ .lastResponse._id }}
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 0,
      "name": "test-instruction-instruction-pbehavior-alarm-api-1-3-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-instruction-pbehavior-alarm-api-1-1",
                "test-resource-instruction-pbehavior-alarm-api-1-2",
                "test-resource-instruction-pbehavior-alarm-api-1-3"
              ]
            }
          }
        ]
      ],
      "disabled_on_pbh": ["test-maintenance-type-to-engine"],
      "description": "test-instruction-instruction-pbehavior-alarm-api-1-3-description",
      "enabled": true,
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "steps": [
        {
          "name": "test-instruction-instruction-pbehavior-alarm-api-1-3-step-1",
          "operations": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-3-step-1-operation-1",
              "time_to_complete": {"value": 1, "unit":"s"},
              "description": "test-instruction-instruction-pbehavior-alarm-api-1-3-step-1-operation-1-description",
              "jobs": []
            }
          ],
          "stop_on_fail": true,
          "endpoint": "test-instruction-instruction-pbehavior-alarm-api-1-3-step-1-endpoint"
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response instructionID3={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-instruction-pbehavior-alarm-api-1",
      "connector_name": "test-connector-name-instruction-pbehavior-alarm-api-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-instruction-pbehavior-alarm-api-1",
      "resource": "test-resource-instruction-pbehavior-alarm-api-1-1",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-instruction-pbehavior-alarm-api-1",
      "connector_name": "test-connector-name-instruction-pbehavior-alarm-api-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-instruction-pbehavior-alarm-api-1",
      "resource": "test-resource-instruction-pbehavior-alarm-api-1-2",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-instruction-pbehavior-alarm-api-1",
      "connector_name": "test-connector-name-instruction-pbehavior-alarm-api-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-instruction-pbehavior-alarm-api-1",
      "resource": "test-resource-instruction-pbehavior-alarm-api-1-3",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-instruction-pbehavior-alarm-api-1-1",
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
              "value": "test-resource-instruction-pbehavior-alarm-api-1-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-instruction-pbehavior-alarm-api-1-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-pause-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-instruction-pbehavior-alarm-api-1-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?instructions[]={"include":["{{ .instructionID1 }}"]}&with_instructions=true&sort_by=d&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-instruction-pbehavior-alarm-api-1-1/test-component-instruction-pbehavior-alarm-api-1"
          },
          "assigned_instructions": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-name",
              "execution": null
            },
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-3-name",
              "execution": null
            }
          ]
        },
        {
          "entity": {
            "_id": "test-resource-instruction-pbehavior-alarm-api-1-2/test-component-instruction-pbehavior-alarm-api-1"
          },
          "assigned_instructions": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-name",
              "execution": null
            },
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-2-name",
              "execution": null
            }
          ]
        },
        {
          "entity": {
            "_id": "test-resource-instruction-pbehavior-alarm-api-1-3/test-component-instruction-pbehavior-alarm-api-1"
          },
          "assigned_instructions": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-name",
              "execution": null
            },
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-3-name",
              "execution": null
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
    When I do GET /api/v4/alarms?instructions[]={"include":["{{ .instructionID2 }}"]}&with_instructions=true&sort_by=d&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-instruction-pbehavior-alarm-api-1-2/test-component-instruction-pbehavior-alarm-api-1"
          },
          "assigned_instructions": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-name",
              "execution": null
            },
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-2-name",
              "execution": null
            }
          ]
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
    When I do GET /api/v4/alarms?instructions[]={"include":["{{ .instructionID3 }}"]}&with_instructions=true&sort_by=d&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-instruction-pbehavior-alarm-api-1-1/test-component-instruction-pbehavior-alarm-api-1"
          },
          "assigned_instructions": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-name",
              "execution": null
            },
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-3-name",
              "execution": null
            }
          ]
        },
        {
          "entity": {
            "_id": "test-resource-instruction-pbehavior-alarm-api-1-3/test-component-instruction-pbehavior-alarm-api-1"
          },
          "assigned_instructions": [
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-1-name",
              "execution": null
            },
            {
              "name": "test-instruction-instruction-pbehavior-alarm-api-1-3-name",
              "execution": null
            }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: instruction must not run if alarm is under pbh
    When I am admin
    When I do POST /api/v4/cat/instructions:
    """json
    {
      "type": 1,
      "name": "test-instruction-instruction-pbehavior-alarm-api-2-name",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-instruction-pbehavior-alarm-api-2"
            }
          }
        ]
      ],
      "description": "test-instruction-instruction-pbehavior-alarm-api-2-description",
      "enabled": true,
      "disabled_on_pbh": ["test-maintenance-type-to-engine"],
      "timeout_after_execution": {
        "value": 1,
        "unit": "s"
      },
      "jobs": [
        {
          "job": "test-job-to-run-auto-instruction-8"
        }
      ],
      "priority": 30,
      "triggers": ["create"]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-instruction-pbehavior-alarm-api-2",
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
              "value": "test-resource-instruction-pbehavior-alarm-api-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-instruction-pbehavior-alarm-api-2",
      "connector_name": "test-connector-name-instruction-pbehavior-alarm-api-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-instruction-pbehavior-alarm-api-2",
      "resource": "test-resource-instruction-pbehavior-alarm-api-2",
      "state": 1
    }
    """
    When I wait the end of event processing
    When I wait 100ms
    When I do GET /api/v4/alarms?search=test-resource-instruction-pbehavior-alarm-api-2
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
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
          "steps": {
            "data": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "pbhenter"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """
