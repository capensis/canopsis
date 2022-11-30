Feature: New import entities
  I need to be able to import entities

  Scenario: import unauthorized
    When I do PUT /api/v4/contextgraph/new-import
    Then the response code should be 401

  Scenario: import without permissions
    When I am noperms
    When I do PUT /api/v4/contextgraph/new-import
    Then the response code should be 403

  Scenario: import with action set should create new entities
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-set:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-component-contextgraph-new-import-1",
            "name": "test-component-contextgraph-new-import-1",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "set",
            "enabled": true
          },
          {
            "_id": "test-resource-contextgraph-new-import-1-1/test-component-contextgraph-new-import-1",
            "name": "test-resource-contextgraph-new-import-1-1",
            "type": "resource",
            "infos": {
              "test_info": {
                "description": "description 2",
                "value": "value 2"
              }
            },
            "action": "set",
            "component": "test-component-contextgraph-new-import-1",
            "enabled": true
          },
          {
            "_id": "test-resource-contextgraph-new-import-1-2/test-component-contextgraph-new-import-1",
            "name": "test-resource-contextgraph-new-import-1-2",
            "type": "resource",
            "infos": {
              "test_info": {
                "description": "description 3",
                "value": "value 3"
              }
            },
            "action": "set",
            "component": "test-component-contextgraph-new-import-1",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-component-contextgraph-new-import-1",
      "name": "test-component-contextgraph-new-import-1",
      "component": "test-component-contextgraph-new-import-1",
      "infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "value 1"
        }
      },
      "enabled": true,
      "type": "component",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-resource-contextgraph-new-import-1-1/test-component-contextgraph-new-import-1",
        "test-resource-contextgraph-new-import-1-2/test-component-contextgraph-new-import-1"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-1-1/test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-1-1/test-component-contextgraph-new-import-1",
      "name": "test-resource-contextgraph-new-import-1-1",
      "component": "test-component-contextgraph-new-import-1",
      "enabled": true,
      "infos": {
        "test_info": {
          "description": "description 2",
          "name": "test_info",
          "value": "value 2"
        }
      },
      "component_infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "value 1"
        }
      },
      "type": "resource",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-contextgraph-new-import-1-1/test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "impact": [
        "test-component-contextgraph-new-import-1"
      ],
      "depends": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-contextgraph-new-import-1-2/test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "impact": [
        "test-component-contextgraph-new-import-1"
      ],
      "depends": []
    }
    """

  Scenario: import with action update should update entity
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-update:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-new-import-to-update",
            "name": "test-entity-contextgraph-new-import-to-update",
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "set",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-to-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-to-update",
      "name": "test-entity-contextgraph-new-import-to-update",
      "enabled": true,
      "infos": {
        "new_info": {
          "description": "",
          "name": "new_info",
          "value": "2"
        }
      },
      "type": "component",
      "impact_level": 1
    }
    """

  Scenario: new import with action set, when name doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "set",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: new-import with action enable
    When I am admin
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-enable
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-enable",
      "name": "test-entity-contextgraph-new-import-component-to-enable",
      "enabled": false,
      "type": "component",
      "impact_level": 1
    }
    """
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-entity-contextgraph-new-import-component-to-enable",
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "enable",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-enable
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-enable",
      "name": "test-entity-contextgraph-new-import-component-to-enable",
      "enabled": true,
      "type": "component",
      "impact_level": 1
    }
    """

  Scenario: new-import with action enable, when name doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "enable",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: new-import with action enable, when entity doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-entity-contextgraph-import-to-set-not-exist",
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "enable",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: new-import with action disable
    When I am admin
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-disable
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-disable",
      "name": "test-entity-contextgraph-new-import-component-to-disable",
      "enabled": true,
      "type": "component",
      "impact_level": 1
    }
    """
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-entity-contextgraph-new-import-component-to-disable",
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "disable",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-disable
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-disable",
      "name": "test-entity-contextgraph-new-import-component-to-disable",
      "enabled": false,
      "type": "component",
      "impact_level": 1
    }
    """

  Scenario: new-import with action disable, when name doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "disable",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: new-import with action disable, when entity doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-entity-contextgraph-import-component-to-disable-not-exist",
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "disable",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: new-import with action set, patterns should be only in services
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "SC003C",
            "type": "component",
            "name": "SC003C",
            "action": "set",
            "enabled": true,
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "regexp",
                    "value": "script_new_import"
                  }
                }
              ]
            ]
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: new-import with action set, wrong action
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "SC003C",
            "type": "service",
            "name": "SC003C",
            "action": "some",
            "enabled": true,
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "regexp",
                    "value": "script_new_import"
                  }
                }
              ]
            ]
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: new-import with action create, context graph should be valid for entity service
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "SC004C",
            "type": "component",
            "infos": {},
            "action": "set",
            "enabled": true
          },
          {
            "name": "script_new-import_service",
            "component": "SC004C",
            "type": "resource",
            "infos": {},
            "action": "set",
            "enabled": true
          },
          {
            "name": "script_new-import_service_2",
            "component": "SC004C",
            "type": "resource",
            "infos": {},
            "action": "set",
            "enabled": true
          },
          {
            "name": "test-entityservice-service-new-import",
            "output_template": "abc",
            "impact_level": 1,
            "enabled": true,
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "regexp",
                    "value": "script_new-import_service"
                  }
                }
              ]
            ],
            "type": "service",
            "action": "set"
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "SC004C",
      "type": "component"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=SC004C
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "script_new-import_service/SC004C",
        "script_new-import_service_2/SC004C"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entitybasics?_id=script_new-import_service/SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_new-import_service/SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_new-import_service",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=script_new-import_service/SC004C until response code is 200 and body is:
    """json
    {
      "depends": [],
      "impact": [
        "SC004C",
        "test-entityservice-service-new-import"
      ]
    }
    """
    When I do GET /api/v4/entitybasics?_id=script_new-import_service/SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_new-import_service/SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_new-import_service",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=script_new-import_service/SC004C until response code is 200 and body is:
    """json
    {
      "depends": [],
      "impact": [
        "SC004C",
        "test-entityservice-service-new-import"
      ]
    }
    """
    When I do GET /api/v4/entitybasics?_id=script_new-import_service_2/SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_new-import_service_2/SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_new-import_service_2",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=script_new-import_service_2/SC004C
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [],
      "impact": [
        "SC004C",
        "test-entityservice-service-new-import"
      ]
    }
    """
    When I do GET /api/v4/entityservices/test-entityservice-service-new-import
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entityservice-service-new-import",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "script_new-import_service"
            }
          }
        ]
      ],
      "impact_level": 1,
      "name": "test-entityservice-service-new-import",
      "output_template": "abc",
      "type": "service"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-service-new-import
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "script_new-import_service/SC004C",
        "script_new-import_service_2/SC004C"
      ],
      "impact": []
    }
    """

  Scenario: given new-import with create resource should set resource component_infos
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-component-contextgraph-new-import-29",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "set",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-resource-contextgraph-new-import-29",
            "component": "test-component-contextgraph-new-import-29",
            "type": "resource",
            "action": "set",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-29/test-component-contextgraph-new-import-29
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "component": "test-component-contextgraph-new-import-29",
      "component_infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "value 1"
        }
      }
    }
    """

  Scenario: given new-import with create component should set resource component_infos
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-resource-contextgraph-new-import-30",
            "component": "test-component-contextgraph-new-import-30",
            "type": "resource",
            "action": "set",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-component-contextgraph-new-import-30",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "set",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-30/test-component-contextgraph-new-import-30
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "component": "test-component-contextgraph-new-import-30",
      "component_infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "value 1"
        }
      }
    }
    """

  Scenario: given new-import with create resource and delete component should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-entity-contextgraph-new-import-resource-to-delete-conflict-1",
            "component": "test-entity-contextgraph-new-import-component-to-delete-conflict",
            "type": "resource",
            "action": "set",
            "enabled": true
          },
          {
            "name": "test-entity-contextgraph-new-import-component-to-delete-conflict",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "delete",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given new-import with create resource and disable component should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-entity-contextgraph-new-import-resource-to-disable-conflict-1",
            "component": "test-entity-contextgraph-new-import-component-to-disable-conflict-1",
            "type": "resource",
            "action": "set",
            "enabled": true
          },
          {
            "name": "test-entity-contextgraph-new-import-component-to-disable-conflict-1",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "disable",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given new-import with create resource when component is disabled
    When I am admin
    When I do PUT /api/v4/contextgraph/new-import?source=test-new-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "name": "test-entity-contextgraph-new-import-resource-to-disable-conflict-2",
            "component": "test-entity-contextgraph-new-import-component-to-disable-conflict-2",
            "type": "resource",
            "action": "set",
            "enabled": true
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """
