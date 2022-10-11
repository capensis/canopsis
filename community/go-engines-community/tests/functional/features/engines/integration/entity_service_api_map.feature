Feature: Get a map's state and alarms
  I need to be able to get a map's state and alarms
  Only admin should be able to get a map's state and alarms

  Scenario: given tree of dependencies map should return entities
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-service-api-map-1",
        "connector_name": "test-connector-name-service-api-map-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-1",
        "resource": "test-resource-service-api-map-1-1",
        "state": 2,
        "output": "test-output-service-api-map-1"
      },
      {
        "connector": "test-connector-service-api-map-1",
        "connector_name": "test-connector-name-service-api-map-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-1",
        "resource": "test-resource-service-api-map-1-2",
        "state": 1,
        "output": "test-output-service-api-map-1"
      },
      {
        "connector": "test-connector-service-api-map-1",
        "connector_name": "test-connector-name-service-api-map-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-1",
        "resource": "test-resource-service-api-map-1-3",
        "state": 0,
        "output": "test-output-service-api-map-1"
      },
      {
        "connector": "test-connector-service-api-map-1",
        "connector_name": "test-connector-name-service-api-map-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-1",
        "resource": "test-resource-service-api-map-1-4",
        "state": 1,
        "output": "test-output-service-api-map-1"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-1-1-name",
      "output_template": "test-entityservice-service-api-map-1-1-output",
      "impact_level": 4,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-1-1/test-component-service-api-map-1",
                "test-resource-service-api-map-1-2/test-component-service-api-map-1",
                "test-resource-service-api-map-1-3/test-component-service-api-map-1"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId1={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-1-2-name",
      "output_template": "test-entityservice-service-api-map-1-2-output",
      "impact_level": 3,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-api-map-1-4/test-component-service-api-map-1"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-1-3-name",
      "output_template": "test-entityservice-service-api-map-1-3-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "{{ .serviceId1 }}",
                "{{ .serviceId2 }}"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId3={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-service-api-map-1-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": "{{ .serviceId1 }}",
            "pinned_entities": [
              "test-resource-service-api-map-1-1/test-component-service-api-map-1",
              "test-resource-service-api-map-1-3/test-component-service-api-map-1"
            ]
          },
          {
            "entity": "{{ .serviceId2 }}",
            "pinned_entities": []
          },
          {
            "entity": "{{ .serviceId3 }}",
            "pinned_entities": [
              "{{ .serviceId1 }}",
              "{{ .serviceId2 }}"
            ]
          },
          {
            "entity": "test-resource-service-api-map-1-2/test-component-service-api-map-1",
            "pinned_entities": []
          },
          {
            "entity": "test-resource-service-api-map-1-3/test-component-service-api-map-1",
            "pinned_entities": []
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-1-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId1 }}",
              "name": "test-entityservice-service-api-map-1-1-name",
              "depends_count": 3
            },
            "pinned_entities": [
              {
                "_id": "test-resource-service-api-map-1-1/test-component-service-api-map-1",
                "name": "test-resource-service-api-map-1-1"
              },
              {
                "_id": "test-resource-service-api-map-1-3/test-component-service-api-map-1",
                "name": "test-resource-service-api-map-1-3"
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId2 }}",
              "name": "test-entityservice-service-api-map-1-2-name",
              "depends_count": 1
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "{{ .serviceId3 }}",
              "name": "test-entityservice-service-api-map-1-3-name",
              "depends_count": 2
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId1 }}",
                "name": "test-entityservice-service-api-map-1-1-name"
              },
              {
                "_id": "{{ .serviceId2 }}",
                "name": "test-entityservice-service-api-map-1-2-name"
              }
            ]
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-1-2/test-component-service-api-map-1",
              "name": "test-resource-service-api-map-1-2",
              "depends_count": 0
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-1-3/test-component-service-api-map-1",
              "name": "test-resource-service-api-map-1-3",
              "depends_count": 0
            },
            "pinned_entities": []
          }
        ]
      }
    }
    """
    When I save response mapId={{ .lastResponse._id }}
    When I do GET /api/v4/cat/map-state/{{ .mapId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-1-name",
      "type": "treeofdeps",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId1 }}",
              "name": "test-entityservice-service-api-map-1-1-name",
              "enabled": true,
              "depends_count": 3,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 4,
              "impact_state": 8,
              "state": 2,
              "status": 1,
              "ko_events": 1,
              "ok_events": 0
            },
            "pinned_entities": [
              {
                "_id": "test-resource-service-api-map-1-1/test-component-service-api-map-1",
                "name": "test-resource-service-api-map-1-1",
                "enabled": true,
                "depends_count": 0,
                "type": "resource",
                "category": null,
                "connector": "test-connector-service-api-map-1/test-connector-name-service-api-map-1",
                "component":  "test-component-service-api-map-1",
                "infos": {},
                "impact_level": 1,
                "impact_state": 2,
                "state": 2,
                "status": 1,
                "ko_events": 1,
                "ok_events": 0
              },
              {
                "_id": "test-resource-service-api-map-1-3/test-component-service-api-map-1",
                "name": "test-resource-service-api-map-1-3",
                "enabled": true,
                "depends_count": 0,
                "type": "resource",
                "category": null,
                "connector": "test-connector-service-api-map-1/test-connector-name-service-api-map-1",
                "component":  "test-component-service-api-map-1",
                "infos": {},
                "impact_level": 1,
                "impact_state": 0,
                "state": 0,
                "status": 0,
                "ko_events": 0,
                "ok_events": 1
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId3 }}",
              "name": "test-entityservice-service-api-map-1-3-name",
              "enabled": true,
              "depends_count": 2,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 2,
              "impact_state": 4,
              "state": 2,
              "status": 1,
              "ko_events": 1,
              "ok_events": 0
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId1 }}",
                "name": "test-entityservice-service-api-map-1-1-name",
                "enabled": true,
                "depends_count": 3,
                "type": "service",
                "category": null,
                "infos": {},
                "impact_level": 4,
                "impact_state": 8,
                "state": 2,
                "status": 1,
                "ko_events": 1,
                "ok_events": 0
              },
              {
                "_id": "{{ .serviceId2 }}",
                "name": "test-entityservice-service-api-map-1-2-name",
                "enabled": true,
                "depends_count": 1,
                "type": "service",
                "category": null,
                "infos": {},
                "impact_level": 3,
                "impact_state": 3,
                "state": 1,
                "status": 1,
                "ko_events": 1,
                "ok_events": 0
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId2 }}",
              "name": "test-entityservice-service-api-map-1-2-name",
              "enabled": true,
              "depends_count": 1,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 3,
              "impact_state": 3,
              "state": 1,
              "status": 1,
              "ko_events": 1,
              "ok_events": 0
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-1-2/test-component-service-api-map-1",
              "name": "test-resource-service-api-map-1-2",
              "enabled": true,
              "depends_count": 0,
              "type": "resource",
              "category": null,
              "connector": "test-connector-service-api-map-1/test-connector-name-service-api-map-1",
              "component":  "test-component-service-api-map-1",
              "infos": {},
              "impact_level": 1,
              "impact_state": 1,
              "state": 1,
              "status": 1,
              "ko_events": 1,
              "ok_events": 0
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-1-3/test-component-service-api-map-1",
              "name": "test-resource-service-api-map-1-3",
              "enabled": true,
              "depends_count": 0,
              "type": "resource",
              "category": null,
              "connector": "test-connector-service-api-map-1/test-connector-name-service-api-map-1",
              "component":  "test-component-service-api-map-1",
              "infos": {},
              "impact_level": 1,
              "impact_state": 0,
              "state": 0,
              "status": 0,
              "ko_events": 0,
              "ok_events": 1
            },
            "pinned_entities": []
          }
        ]
      }
    }
    """

  Scenario: given impact chain map should return entities
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-service-api-map-2",
        "connector_name": "test-connector-name-service-api-map-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-2",
        "resource": "test-resource-service-api-map-2-1",
        "state": 2,
        "output": "test-output-service-api-map-2"
      },
      {
        "connector": "test-connector-service-api-map-2",
        "connector_name": "test-connector-name-service-api-map-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-2",
        "resource": "test-resource-service-api-map-2-2",
        "state": 1,
        "output": "test-output-service-api-map-2"
      },
      {
        "connector": "test-connector-service-api-map-2",
        "connector_name": "test-connector-name-service-api-map-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-2",
        "resource": "test-resource-service-api-map-2-3",
        "state": 0,
        "output": "test-output-service-api-map-2"
      },
      {
        "connector": "test-connector-service-api-map-2",
        "connector_name": "test-connector-name-service-api-map-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-2",
        "resource": "test-resource-service-api-map-2-4",
        "state": 1,
        "output": "test-output-service-api-map-2"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-2-1-name",
      "output_template": "test-entityservice-service-api-map-2-1-output",
      "impact_level": 4,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-2-1/test-component-service-api-map-2",
                "test-resource-service-api-map-2-2/test-component-service-api-map-2",
                "test-resource-service-api-map-2-3/test-component-service-api-map-2"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId1={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-2-2-name",
      "output_template": "test-entityservice-service-api-map-2-2-output",
      "impact_level": 3,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-resource-service-api-map-2-4/test-component-service-api-map-2"
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-2-3-name",
      "output_template": "test-entityservice-service-api-map-2-3-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "{{ .serviceId1 }}",
                "{{ .serviceId2 }}"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId3={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-service-api-map-2-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": "{{ .serviceId1 }}",
            "pinned_entities": [
              "{{ .serviceId3 }}"
            ]
          },
          {
            "entity": "{{ .serviceId2 }}",
            "pinned_entities": []
          },
          {
            "entity": "{{ .serviceId3 }}",
            "pinned_entities": []
          },
          {
            "entity": "test-resource-service-api-map-2-2/test-component-service-api-map-2",
            "pinned_entities": [
              "{{ .serviceId1 }}"
            ]
          },
          {
            "entity": "test-resource-service-api-map-2-3/test-component-service-api-map-2",
            "pinned_entities": []
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-2-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId1 }}",
              "name": "test-entityservice-service-api-map-2-1-name",
              "impacts_count": 1
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId3 }}",
                "name": "test-entityservice-service-api-map-2-3-name"
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId2 }}",
              "name": "test-entityservice-service-api-map-2-2-name",
              "impacts_count": 1
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "{{ .serviceId3 }}",
              "name": "test-entityservice-service-api-map-2-3-name",
              "impacts_count": 0
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-2-2/test-component-service-api-map-2",
              "name": "test-resource-service-api-map-2-2",
              "impacts_count": 1
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId1 }}",
                "name": "test-entityservice-service-api-map-2-1-name"
              }
            ]
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-2-3/test-component-service-api-map-2",
              "name": "test-resource-service-api-map-2-3",
              "impacts_count": 1
            },
            "pinned_entities": []
          }
        ]
      }
    }
    """
    When I save response mapId={{ .lastResponse._id }}
    When I do GET /api/v4/cat/map-state/{{ .mapId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-2-name",
      "type": "treeofdeps",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId1 }}",
              "name": "test-entityservice-service-api-map-2-1-name",
              "impacts_count": 1,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 4,
              "impact_state": 8,
              "state": 2,
              "status": 1
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId3 }}",
                "name": "test-entityservice-service-api-map-2-3-name",
                "impacts_count": 0,
                "type": "service",
                "category": null,
                "infos": {},
                "impact_level": 2,
                "impact_state": 4,
                "state": 2,
                "status": 1
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId3 }}",
              "name": "test-entityservice-service-api-map-2-3-name",
              "impacts_count": 0,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 2,
              "impact_state": 4,
              "state": 2,
              "status": 1
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "{{ .serviceId2 }}",
              "name": "test-entityservice-service-api-map-2-2-name",
              "impacts_count": 1,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 3,
              "impact_state": 3,
              "state": 1,
              "status": 1
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-2-2/test-component-service-api-map-2",
              "name": "test-resource-service-api-map-2-2",
              "impacts_count": 1,
              "type": "resource",
              "category": null,
              "connector": "test-connector-service-api-map-2/test-connector-name-service-api-map-2",
              "component":  "test-component-service-api-map-2",
              "infos": {},
              "impact_level": 1,
              "impact_state": 1,
              "state": 1,
              "status": 1
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId1 }}",
                "name": "test-entityservice-service-api-map-2-1-name",
                "impacts_count": 1,
                "type": "service",
                "category": null,
                "infos": {},
                "impact_level": 4,
                "impact_state": 8,
                "state": 2,
                "status": 1
              }
            ]
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-2-3/test-component-service-api-map-2",
              "name": "test-resource-service-api-map-2-3",
              "impacts_count": 1,
              "type": "resource",
              "category": null,
              "connector": "test-connector-service-api-map-2/test-connector-name-service-api-map-2",
              "component":  "test-component-service-api-map-2",
              "infos": {},
              "impact_level": 1,
              "impact_state": 0,
              "state": 0,
              "status": 0
            },
            "pinned_entities": []
          }
        ]
      }
    }
    """

  Scenario: given pinned dependency which is removed from dependencies should not return it
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-service-api-map-3",
        "connector_name": "test-connector-name-service-api-map-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-3",
        "resource": "test-resource-service-api-map-3-1",
        "state": 2,
        "output": "test-output-service-api-map-3"
      },
      {
        "connector": "test-connector-service-api-map-3",
        "connector_name": "test-connector-name-service-api-map-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-3",
        "resource": "test-resource-service-api-map-3-2",
        "state": 1,
        "output": "test-output-service-api-map-3"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-3-1-name",
      "output_template": "test-entityservice-service-api-map-3-1-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-3-1/test-component-service-api-map-3",
                "test-resource-service-api-map-3-2/test-component-service-api-map-3"
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
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-service-api-map-3-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": "{{ .serviceId }}",
            "pinned_entities": [
              "test-resource-service-api-map-3-1/test-component-service-api-map-3",
              "test-resource-service-api-map-3-2/test-component-service-api-map-3"
            ]
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-3-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId }}",
              "name": "test-entityservice-service-api-map-3-1-name",
              "depends_count": 2
            },
            "pinned_entities": [
              {
                "_id": "test-resource-service-api-map-3-1/test-component-service-api-map-3",
                "name": "test-resource-service-api-map-3-1"
              },
              {
                "_id": "test-resource-service-api-map-3-2/test-component-service-api-map-3",
                "name": "test-resource-service-api-map-3-2"
              }
            ]
          }
        ]
      }
    }
    """
    When I save response mapId={{ .lastResponse._id }}
    When I do PUT /api/v4/entityservices/{{ .serviceId }}:
    """json
    {
      "name": "test-entityservice-service-api-map-3-1-name",
      "output_template": "test-entityservice-service-api-map-3-1-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-3-1/test-component-service-api-map-3"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do GET /api/v4/cat/maps/{{ .mapId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-3-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId }}",
              "name": "test-entityservice-service-api-map-3-1-name",
              "depends_count": 1
            },
            "pinned_entities": [
              {
                "_id": "test-resource-service-api-map-3-1/test-component-service-api-map-3",
                "name": "test-resource-service-api-map-3-1"
              }
            ]
          }
        ]
      }
    }
    """
    When I do GET /api/v4/cat/map-state/{{ .mapId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-3-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId }}"
            },
            "pinned_entities": [
              {
                "_id": "test-resource-service-api-map-3-1/test-component-service-api-map-3"
              }
            ]
          }
        ]
      }
    }
    """

  Scenario: given pinned impact which is removed from impact should not return it
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-service-api-map-4",
        "connector_name": "test-connector-name-service-api-map-4",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-4",
        "resource": "test-resource-service-api-map-4-1",
        "state": 2,
        "output": "test-output-service-api-map-4"
      },
      {
        "connector": "test-connector-service-api-map-4",
        "connector_name": "test-connector-name-service-api-map-4",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-4",
        "resource": "test-resource-service-api-map-4-2",
        "state": 1,
        "output": "test-output-service-api-map-4"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-4-1-name",
      "output_template": "test-entityservice-service-api-map-4-1-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-4-1/test-component-service-api-map-4",
                "test-resource-service-api-map-4-2/test-component-service-api-map-4"
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
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-service-api-map-4-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": "test-resource-service-api-map-4-2/test-component-service-api-map-4",
            "pinned_entities": [
              "{{ .serviceId }}"
            ]
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-4-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": {
              "_id": "test-resource-service-api-map-4-2/test-component-service-api-map-4",
              "name": "test-resource-service-api-map-4-2",
              "impacts_count": 1
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId }}",
                "name": "test-entityservice-service-api-map-4-1-name"
              }
            ]
          }
        ]
      }
    }
    """
    When I save response mapId={{ .lastResponse._id }}
    When I do PUT /api/v4/entityservices/{{ .serviceId }}:
    """json
    {
      "name": "test-entityservice-service-api-map-4-1-name",
      "output_template": "test-entityservice-service-api-map-4-1-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-4-1/test-component-service-api-map-4"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 200
    When I wait the end of 2 events processing
    When I do GET /api/v4/cat/maps/{{ .mapId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-4-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": {
              "_id": "test-resource-service-api-map-4-2/test-component-service-api-map-4",
              "name": "test-resource-service-api-map-4-2",
              "impacts_count": 0
            },
            "pinned_entities": []
          }
        ]
      }
    }
    """
    When I do GET /api/v4/cat/map-state/{{ .mapId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-4-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": {
              "_id": "test-resource-service-api-map-4-2/test-component-service-api-map-4"
            },
            "pinned_entities": []
          }
        ]
      }
    }
    """

  Scenario: given tree of dependencies map and expanded entities should return entities
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-service-api-map-5",
        "connector_name": "test-connector-name-service-api-map-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-5",
        "resource": "test-resource-service-api-map-5-1",
        "state": 2,
        "output": "test-output-service-api-map-5"
      },
      {
        "connector": "test-connector-service-api-map-5",
        "connector_name": "test-connector-name-service-api-map-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-5",
        "resource": "test-resource-service-api-map-5-2",
        "state": 1,
        "output": "test-output-service-api-map-5"
      },
      {
        "connector": "test-connector-service-api-map-5",
        "connector_name": "test-connector-name-service-api-map-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-5",
        "resource": "test-resource-service-api-map-5-3",
        "state": 0,
        "output": "test-output-service-api-map-5"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-5-1-name",
      "output_template": "test-entityservice-service-api-map-5-1-output",
      "impact_level": 4,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-5-1/test-component-service-api-map-5",
                "test-resource-service-api-map-5-2/test-component-service-api-map-5"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId1={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-5-2-name",
      "output_template": "test-entityservice-service-api-map-5-2-output",
      "impact_level": 3,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-5-3/test-component-service-api-map-5"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-5-3-name",
      "output_template": "test-entityservice-service-api-map-5-3-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "{{ .serviceId1 }}",
                "{{ .serviceId2 }}"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId3={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-service-api-map-5-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": "{{ .serviceId1 }}",
            "pinned_entities": [
              "test-resource-service-api-map-5-2/test-component-service-api-map-5"
            ]
          },
          {
            "entity": "{{ .serviceId2 }}",
            "pinned_entities": []
          },
          {
            "entity": "{{ .serviceId3 }}",
            "pinned_entities": [
              "{{ .serviceId2 }}"
            ]
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I save response mapId={{ .lastResponse._id }}
    When I do GET /api/v4/cat/map-state/{{ .mapId }}?expanded_entities[]={{ .serviceId1 }}&expanded_entities[]={{ .serviceId2 }}&expanded_entities[]={{ .serviceId3 }}&expanded_entities[]=test-resource-service-api-map-5-2/test-component-service-api-map-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-5-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId1 }}",
              "impact_level": 4,
              "state": 2,
              "impact_state": 8
            },
            "pinned_entities": [
              {
                "_id": "test-resource-service-api-map-5-2/test-component-service-api-map-5",
                "impact_level": 1,
                "state": 1,
                "impact_state": 1
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId3 }}",
              "impact_level": 2,
              "state": 2,
              "impact_state": 4
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId2 }}",
                "impact_level": 3,
                "state": 0,
                "impact_state": 0
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId2 }}",
              "impact_level": 3,
              "state": 0,
              "impact_state": 0
            },
            "pinned_entities": []
          }
        ],
        "expanded_entities": {
          "{{ .serviceId1 }}": [
            {
              "_id": "test-resource-service-api-map-5-1/test-component-service-api-map-5",
              "name": "test-resource-service-api-map-5-1",
              "depends_count": 0,
              "type": "resource",
              "category": null,
              "connector": "test-connector-service-api-map-5/test-connector-name-service-api-map-5",
              "component":  "test-component-service-api-map-5",
              "infos": {},
              "impact_level": 1,
              "impact_state": 2,
              "state": 2,
              "status": 1
            },
            {
              "_id": "test-resource-service-api-map-5-2/test-component-service-api-map-5",
              "name": "test-resource-service-api-map-5-2",
              "depends_count": 0,
              "type": "resource",
              "category": null,
              "connector": "test-connector-service-api-map-5/test-connector-name-service-api-map-5",
              "component":  "test-component-service-api-map-5",
              "infos": {},
              "impact_level": 1,
              "impact_state": 1,
              "state": 1,
              "status": 1
            }
          ],
          "{{ .serviceId2 }}": [
            {
              "_id": "test-resource-service-api-map-5-3/test-component-service-api-map-5",
              "name": "test-resource-service-api-map-5-3",
              "depends_count": 0,
              "type": "resource",
              "category": null,
              "connector": "test-connector-service-api-map-5/test-connector-name-service-api-map-5",
              "component":  "test-component-service-api-map-5",
              "infos": {},
              "impact_level": 1,
              "impact_state": 0,
              "state": 0,
              "status": 0
            }
          ],
          "{{ .serviceId3 }}": [
            {
              "_id": "{{ .serviceId1 }}",
              "name": "test-entityservice-service-api-map-5-1-name",
              "depends_count": 2,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 4,
              "impact_state": 8,
              "state": 2,
              "status": 1
            },
            {
              "_id": "{{ .serviceId2 }}",
              "name": "test-entityservice-service-api-map-5-2-name",
              "depends_count": 1,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 3,
              "impact_state": 0,
              "state": 0,
              "status": 0
            }
          ],
          "test-resource-service-api-map-5-2/test-component-service-api-map-5": []
        }
      }
    }
    """

  Scenario: given impact chain map and expanded entities should return entities
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-service-api-map-6",
        "connector_name": "test-connector-name-service-api-map-6",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-6",
        "resource": "test-resource-service-api-map-6-1",
        "state": 2,
        "output": "test-output-service-api-map-6"
      },
      {
        "connector": "test-connector-service-api-map-6",
        "connector_name": "test-connector-name-service-api-map-6",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-6",
        "resource": "test-resource-service-api-map-6-2",
        "state": 1,
        "output": "test-output-service-api-map-6"
      },
      {
        "connector": "test-connector-service-api-map-6",
        "connector_name": "test-connector-name-service-api-map-6",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-6",
        "resource": "test-resource-service-api-map-6-3",
        "state": 0,
        "output": "test-output-service-api-map-6"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-6-1-name",
      "output_template": "test-entityservice-service-api-map-6-1-output",
      "impact_level": 4,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-6-1/test-component-service-api-map-6",
                "test-resource-service-api-map-6-2/test-component-service-api-map-6"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId1={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-6-2-name",
      "output_template": "test-entityservice-service-api-map-6-2-output",
      "impact_level": 3,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-6-3/test-component-service-api-map-6"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-6-3-name",
      "output_template": "test-entityservice-service-api-map-6-3-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "{{ .serviceId1 }}",
                "{{ .serviceId2 }}"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId3={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-service-api-map-6-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": "test-resource-service-api-map-6-1/test-component-service-api-map-6",
            "pinned_entities": [
              "{{ .serviceId1 }}"
            ]
          },
          {
            "entity": "{{ .serviceId2 }}",
            "pinned_entities": []
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I save response mapId={{ .lastResponse._id }}
    When I do GET /api/v4/cat/map-state/{{ .mapId }}?expanded_entities[]={{ .serviceId2 }}&expanded_entities[]={{ .serviceId3 }}&expanded_entities[]=test-resource-service-api-map-6-1/test-component-service-api-map-6&expanded_entities[]=test-resource-service-api-map-6-3/test-component-service-api-map-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-6-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": {
              "_id": "test-resource-service-api-map-6-1/test-component-service-api-map-6",
              "impact_level": 1,
              "state": 2,
              "impact_state": 2
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId1 }}",
                "impact_level": 4,
                "state": 2,
                "impact_state": 8
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId2 }}",
              "impact_level": 3,
              "state": 0,
              "impact_state": 0
            },
            "pinned_entities": []
          }
        ],
        "expanded_entities": {
          "test-resource-service-api-map-6-1/test-component-service-api-map-6": [
            {
              "_id": "{{ .serviceId1 }}",
              "name": "test-entityservice-service-api-map-6-1-name",
              "impacts_count": 1,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 4,
              "impact_state": 8,
              "state": 2,
              "status": 1
            }
          ],
          "test-resource-service-api-map-6-3/test-component-service-api-map-6": [
            {
              "_id": "{{ .serviceId2 }}",
              "name": "test-entityservice-service-api-map-6-2-name",
              "impacts_count": 1,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 3,
              "impact_state": 0,
              "state": 0,
              "status": 0
            }
          ],
          "{{ .serviceId2 }}": [
            {
              "_id": "{{ .serviceId3 }}",
              "name": "test-entityservice-service-api-map-6-3-name",
              "impacts_count": 0,
              "type": "service",
              "category": null,
              "infos": {},
              "impact_level": 2,
              "impact_state": 4,
              "state": 2,
              "status": 1
            }
          ],
          "{{ .serviceId3 }}": []
        }
      }
    }
    """

  Scenario: given tree of dependencies map and filter should return filtered entities
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-service-api-map-7",
        "connector_name": "test-connector-name-service-api-map-7",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-7",
        "resource": "test-resource-service-api-map-7-1",
        "state": 2,
        "output": "test-output-service-api-map-7"
      },
      {
        "connector": "test-connector-service-api-map-7",
        "connector_name": "test-connector-name-service-api-map-7",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-7",
        "resource": "test-resource-service-api-map-7-2",
        "state": 1,
        "output": "test-output-service-api-map-7"
      },
      {
        "connector": "test-connector-service-api-map-7",
        "connector_name": "test-connector-name-service-api-map-7",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-7",
        "resource": "test-resource-service-api-map-7-3",
        "state": 0,
        "output": "test-output-service-api-map-7"
      },
      {
        "connector": "test-connector-service-api-map-7",
        "connector_name": "test-connector-name-service-api-map-7",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-7",
        "resource": "test-resource-service-api-map-7-4",
        "state": 0,
        "output": "test-output-service-api-map-7"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-7-1-name",
      "output_template": "test-entityservice-service-api-map-7-1-output",
      "impact_level": 4,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-7-1/test-component-service-api-map-7",
                "test-resource-service-api-map-7-2/test-component-service-api-map-7",
                "test-resource-service-api-map-7-3/test-component-service-api-map-7"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId1={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-7-2-name",
      "output_template": "test-entityservice-service-api-map-7-2-output",
      "impact_level": 3,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-7-4/test-component-service-api-map-7"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-7-3-name",
      "output_template": "test-entityservice-service-api-map-7-3-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "{{ .serviceId1 }}",
                "{{ .serviceId2 }}"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId3={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-service-api-map-7-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": "{{ .serviceId1 }}",
            "pinned_entities": [
              "test-resource-service-api-map-7-1/test-component-service-api-map-7",
              "test-resource-service-api-map-7-2/test-component-service-api-map-7"
            ]
          },
          {
            "entity": "{{ .serviceId2 }}",
            "pinned_entities": []
          },
          {
            "entity": "{{ .serviceId3 }}",
            "pinned_entities": [
              "{{ .serviceId2 }}"
            ]
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I save response mapId={{ .lastResponse._id }}
    When I save response mapCreated={{ .lastResponse.created }}
    When I do GET /api/v4/cat/map-state/{{ .mapId }}?filters[]=test-widgetfilter-entity-service-api-map-7-1&expanded_entities[]={{ .serviceId1 }}&expanded_entities[]={{ .serviceId3 }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-7-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps",
        "entities": [
          {
            "entity": {
              "_id": "{{ .serviceId1 }}",
              "depends_count": 2
            },
            "pinned_entities": [
              {
                "_id": "test-resource-service-api-map-7-1/test-component-service-api-map-7",
                "depends_count": 0
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId3 }}",
              "depends_count": 1
            },
            "pinned_entities": []
          }
        ],
        "expanded_entities": {
          "{{ .serviceId1 }}": [
            {
              "_id": "test-resource-service-api-map-7-1/test-component-service-api-map-7",
              "depends_count": 0
            },
            {
              "_id": "test-resource-service-api-map-7-3/test-component-service-api-map-7",
              "depends_count": 0
            }
          ],
          "{{ .serviceId3 }}": [
            {
              "_id": "{{ .serviceId1 }}",
              "depends_count": 2
            }
          ]
        }
      }
    }
    """
    When I do GET /api/v4/cat/map-state/{{ .mapId }}?filters[]=test-widgetfilter-entity-service-api-map-7-2&expanded_entities[]={{ .serviceId1 }}&expanded_entities[]={{ .serviceId3 }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "{{ .mapId }}",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-map-service-api-map-7-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "treeofdeps"
      },
      "created": {{ .mapCreated }},
      "updated": {{ .mapCreated }}
    }
    """

  Scenario: given impact chain map and filter should return filtered entities
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-service-api-map-8",
        "connector_name": "test-connector-name-service-api-map-8",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-8",
        "resource": "test-resource-service-api-map-8-1",
        "state": 2,
        "output": "test-output-service-api-map-8"
      },
      {
        "connector": "test-connector-service-api-map-8",
        "connector_name": "test-connector-name-service-api-map-8",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-8",
        "resource": "test-resource-service-api-map-8-2",
        "state": 1,
        "output": "test-output-service-api-map-8"
      },
      {
        "connector": "test-connector-service-api-map-8",
        "connector_name": "test-connector-name-service-api-map-8",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-service-api-map-8",
        "resource": "test-resource-service-api-map-8-3",
        "state": 0,
        "output": "test-output-service-api-map-8"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-8-1-name",
      "output_template": "test-entityservice-service-api-map-8-1-output",
      "impact_level": 4,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-8-1/test-component-service-api-map-8",
                "test-resource-service-api-map-8-2/test-component-service-api-map-8"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId1={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-8-2-name",
      "output_template": "test-entityservice-service-api-map-8-2-output",
      "impact_level": 3,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "test-resource-service-api-map-8-2/test-component-service-api-map-8",
                "test-resource-service-api-map-8-3/test-component-service-api-map-8"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId2={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-service-api-map-8-3-name",
      "output_template": "test-entityservice-service-api-map-8-3-output",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "is_one_of",
              "value": [
                "{{ .serviceId1 }}",
                "{{ .serviceId2 }}"
              ]
            }
          }
        ]
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I save response serviceId3={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do POST /api/v4/cat/maps:
    """json
    {
      "name": "test-map-service-api-map-8-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": "test-resource-service-api-map-8-1/test-component-service-api-map-8",
            "pinned_entities": []
          },
          {
            "entity": "test-resource-service-api-map-8-2/test-component-service-api-map-8",
            "pinned_entities": [
              "{{ .serviceId1 }}",
              "{{ .serviceId2 }}"
            ]
          },
          {
            "entity": "{{ .serviceId2 }}",
            "pinned_entities": []
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I save response mapId={{ .lastResponse._id }}
    When I do GET /api/v4/cat/map-state/{{ .mapId }}?filters[]=test-widgetfilter-entity-service-api-map-8&expanded_entities[]={{ .serviceId2 }}&expanded_entities[]=test-resource-service-api-map-8-2/test-component-service-api-map-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-map-service-api-map-8-name",
      "type": "treeofdeps",
      "parameters": {
        "type": "impactchain",
        "entities": [
          {
            "entity": {
              "_id": "test-resource-service-api-map-8-1/test-component-service-api-map-8",
              "impacts_count": 0
            },
            "pinned_entities": []
          },
          {
            "entity": {
              "_id": "test-resource-service-api-map-8-2/test-component-service-api-map-8",
              "impacts_count": 1
            },
            "pinned_entities": [
              {
                "_id": "{{ .serviceId2 }}",
                "impacts_count": 1
              }
            ]
          },
          {
            "entity": {
              "_id": "{{ .serviceId2 }}",
              "impacts_count": 1
            },
            "pinned_entities": []
          }
        ],
        "expanded_entities": {
          "test-resource-service-api-map-8-2/test-component-service-api-map-8": [
            {
              "_id": "{{ .serviceId2 }}",
              "impacts_count": 1
            }
          ],
          "{{ .serviceId2 }}": [
            {
              "_id": "{{ .serviceId3 }}",
              "impacts_count": 0
            }
          ]
        }
      }
    }
    """
