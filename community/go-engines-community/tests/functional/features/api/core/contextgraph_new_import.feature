Feature: New import entities
  I need to be able to import entities

  Scenario: given unauthorized import request
    When I do PUT /api/v4/contextgraph-import
    Then the response code should be 401

  Scenario: given import request without permissions
    When I am noperms
    When I do PUT /api/v4/contextgraph-import
    Then the response code should be 403

  Scenario: given set import requests should create new entities and context graph
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-set:
    """json
    [
      {
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
      "impact_level": 1,
      "import_source": "test-new-import-set"
    }
    """
    Then the response key "imported" should be greater than 0
    When I do GET /api/v4/entities?search=test-resource-contextgraph-new-import-1-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
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
          "impact_level": 1,
          "import_source": "test-new-import-set"
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    Then the response key "data.0.imported" should be greater than 0    
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
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-1-2/test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-1-2/test-component-contextgraph-new-import-1",
      "name": "test-resource-contextgraph-new-import-1-2",
      "component": "test-component-contextgraph-new-import-1",
      "enabled": true,
      "infos": {
        "test_info": {
          "description": "description 3",
          "name": "test_info",
          "value": "value 3"
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
    When I do GET /api/v4/entities/context-graph?_id=test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-resource-contextgraph-new-import-1-1/test-component-contextgraph-new-import-1",
      "test-resource-contextgraph-new-import-1-2/test-component-contextgraph-new-import-1"
    ]
    """
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-set:
    """json
    [
      {
        "name": "test-resource-contextgraph-new-import-1-3",
        "type": "resource",
        "infos": {
          "test_info": {
            "description": "description 4",
            "value": "value 4"
          }
        },
        "action": "set",
        "component": "test-component-contextgraph-new-import-1",
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
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-1-3/test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-1-3/test-component-contextgraph-new-import-1",
      "name": "test-resource-contextgraph-new-import-1-3",
      "component": "test-component-contextgraph-new-import-1",
      "enabled": true,
      "infos": {
        "test_info": {
          "description": "description 4",
          "name": "test_info",
          "value": "value 4"
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
    When I do GET /api/v4/entities/context-graph?_id=test-component-contextgraph-new-import-1
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-resource-contextgraph-new-import-1-1/test-component-contextgraph-new-import-1",
      "test-resource-contextgraph-new-import-1-2/test-component-contextgraph-new-import-1",
      "test-resource-contextgraph-new-import-1-3/test-component-contextgraph-new-import-1"
    ]
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-contextgraph-new-import-1-3/test-component-contextgraph-new-import-1
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

  Scenario: given set import request should create component with resource config
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-set:
    """json
    [
      {
        "name": "test-resource-contextgraph-new-import-2-1",
        "type": "resource",
        "infos": {
          "test_info": {
            "description": "description 2",
            "value": "value 2"
          }
        },
        "action": "set",
        "component": "test-component-contextgraph-new-import-2",
        "enabled": true
      },
      {
        "name": "test-resource-contextgraph-new-import-2-2",
        "type": "resource",
        "infos": {
          "test_info": {
            "description": "description 3",
            "value": "value 3"
          }
        },
        "action": "set",
        "component": "test-component-contextgraph-new-import-2",
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
    When I do GET /api/v4/entitybasics?_id=test-component-contextgraph-new-import-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-component-contextgraph-new-import-2",
      "name": "test-component-contextgraph-new-import-2",
      "component": "test-component-contextgraph-new-import-2",
      "enabled": true,
      "type": "component",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-2-1/test-component-contextgraph-new-import-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-2-1/test-component-contextgraph-new-import-2",
      "name": "test-resource-contextgraph-new-import-2-1",
      "component": "test-component-contextgraph-new-import-2",
      "enabled": true,
      "infos": {
        "test_info": {
          "description": "description 2",
          "name": "test_info",
          "value": "value 2"
        }
      },
      "type": "resource",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-2-2/test-component-contextgraph-new-import-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-2-2/test-component-contextgraph-new-import-2",
      "name": "test-resource-contextgraph-new-import-2-2",
      "component": "test-component-contextgraph-new-import-2",
      "enabled": true,
      "infos": {
        "test_info": {
          "description": "description 3",
          "name": "test_info",
          "value": "value 3"
        }
      },
      "type": "resource",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-contextgraph-new-import-2-1/test-component-contextgraph-new-import-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "impact": [
        "test-component-contextgraph-new-import-2"
      ],
      "depends": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-resource-contextgraph-new-import-2-2/test-component-contextgraph-new-import-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "impact": [
        "test-component-contextgraph-new-import-2"
      ],
      "depends": []
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-component-contextgraph-new-import-2
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-resource-contextgraph-new-import-2-1/test-component-contextgraph-new-import-2",
      "test-resource-contextgraph-new-import-2-2/test-component-contextgraph-new-import-2"
    ]
    """

  Scenario: given set import requests should create and update component_infos
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-update:
    """json
    [
      {
        "name": "test-component-contextgraph-new-import-3",
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
        "name": "test-resource-contextgraph-new-import-3-1",
        "type": "resource",
        "infos": {
          "test_info": {
            "description": "description 2",
            "value": "value 2"
          }
        },
        "action": "set",
        "component": "test-component-contextgraph-new-import-3",
        "enabled": true
      },
      {
        "name": "test-resource-contextgraph-new-import-3-2",
        "type": "resource",
        "infos": {
          "test_info": {
            "description": "description 3",
            "value": "value 3"
          }
        },
        "action": "set",
        "component": "test-component-contextgraph-new-import-3",
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
    When I do GET /api/v4/entitybasics?_id=test-component-contextgraph-new-import-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-component-contextgraph-new-import-3",
      "name": "test-component-contextgraph-new-import-3",
      "component": "test-component-contextgraph-new-import-3",
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
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-3-1/test-component-contextgraph-new-import-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-3-1/test-component-contextgraph-new-import-3",
      "name": "test-resource-contextgraph-new-import-3-1",
      "component": "test-component-contextgraph-new-import-3",
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
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-3-2/test-component-contextgraph-new-import-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-3-2/test-component-contextgraph-new-import-3",
      "name": "test-resource-contextgraph-new-import-3-2",
      "component": "test-component-contextgraph-new-import-3",
      "enabled": true,
      "infos": {
        "test_info": {
          "description": "description 3",
          "name": "test_info",
          "value": "value 3"
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
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-update:
    """json
    [
      {
        "_id": "test-component-contextgraph-new-import-3",
        "name": "test-component-contextgraph-new-import-3",
        "type": "component",
        "infos": {
          "test_info": {
            "description": "description 1",
            "value": "new value"
          }
        },
        "action": "set",
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
    When I do GET /api/v4/entitybasics?_id=test-component-contextgraph-new-import-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-component-contextgraph-new-import-3",
      "name": "test-component-contextgraph-new-import-3",
      "component": "test-component-contextgraph-new-import-3",
      "infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "new value"
        }
      },
      "enabled": true,
      "type": "component",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-3-1/test-component-contextgraph-new-import-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-3-1/test-component-contextgraph-new-import-3",
      "name": "test-resource-contextgraph-new-import-3-1",
      "component": "test-component-contextgraph-new-import-3",
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
          "value": "new value"
        }
      },
      "type": "resource",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-3-2/test-component-contextgraph-new-import-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-3-2/test-component-contextgraph-new-import-3",
      "name": "test-resource-contextgraph-new-import-3-2",
      "component": "test-component-contextgraph-new-import-3",
      "enabled": true,
      "infos": {
        "test_info": {
          "description": "description 3",
          "name": "test_info",
          "value": "value 3"
        }
      },
      "component_infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "new value"
        }
      },
      "type": "resource",
      "impact_level": 1
    }
    """
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-update:
    """json
    [
      {
        "name": "test-resource-contextgraph-new-import-3-3",
        "type": "resource",
        "infos": {
          "test_info": {
            "description": "description 4",
            "value": "value 4"
          }
        },
        "action": "set",
        "component": "test-component-contextgraph-new-import-3",
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
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-new-import-3-3/test-component-contextgraph-new-import-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-new-import-3-3/test-component-contextgraph-new-import-3",
      "name": "test-resource-contextgraph-new-import-3-3",
      "component": "test-component-contextgraph-new-import-3",
      "enabled": true,
      "infos": {
        "test_info": {
          "description": "description 4",
          "name": "test_info",
          "value": "value 4"
        }
      },
      "component_infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "new value"
        }
      },
      "type": "resource",
      "impact_level": 1
    }
    """

  Scenario: given set import request without a name should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given enable import should enable entity
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
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
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

  Scenario: given enable import when entity doesn't exist should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "test-entity-contextgraph-new-import-to-set-not-exist",
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given disable import should disable entity
    When I am admin
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-disable-1
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-disable-1",
      "name": "test-entity-contextgraph-new-import-component-to-disable-1",
      "enabled": true,
      "type": "component",
      "impact_level": 1
    }
    """
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "test-entity-contextgraph-new-import-component-to-disable-1",
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-disable-1
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-disable-1",
      "name": "test-entity-contextgraph-new-import-component-to-disable-1",
      "enabled": false,
      "type": "component",
      "impact_level": 1
    }
    """

  Scenario: given disable import should disable component and its resources
    When I am admin
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-disable-2
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-disable-2",
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-disable-2-1/test-entity-contextgraph-new-import-component-to-disable-2
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-resource-to-disable-2-1/test-entity-contextgraph-new-import-component-to-disable-2",
      "enabled": true
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-disable-2-2/test-entity-contextgraph-new-import-component-to-disable-2
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-resource-to-disable-2-2/test-entity-contextgraph-new-import-component-to-disable-2",
      "enabled": true
    }
    """
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "test-entity-contextgraph-new-import-component-to-disable-2",
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-disable-2
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-disable-2",
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-disable-2-1/test-entity-contextgraph-new-import-component-to-disable-2
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-resource-to-disable-2-1/test-entity-contextgraph-new-import-component-to-disable-2",
      "enabled": false
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-disable-2-2/test-entity-contextgraph-new-import-component-to-disable-2
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-resource-to-disable-2-2/test-entity-contextgraph-new-import-component-to-disable-2",
      "enabled": false
    }
    """

  Scenario: given disable import when entity doesn't exist should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given create import request when patterns not for a service should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given create import request with wrong type should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given create import request should calculate context graph for entity service
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "component_SC004C",
        "type": "component",
        "infos": {},
        "action": "set",
        "enabled": true
      },
      {
        "name": "script_new-import_service",
        "component": "component_SC004C",
        "type": "resource",
        "infos": {},
        "action": "set",
        "enabled": true
      },
      {
        "name": "script_new-import_service_2",
        "component": "component_SC004C",
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "done"
    }
    """
    When I do GET /api/v4/entitybasics?_id=component_SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "component_SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "component_SC004C",
      "type": "component"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=component_SC004C
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "script_new-import_service/component_SC004C",
      "script_new-import_service_2/component_SC004C"
    ]
    """
    When I do GET /api/v4/entitybasics?_id=script_new-import_service/component_SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_new-import_service/component_SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_new-import_service",
      "type": "resource"
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-entityservice-context-graph-new-import-1",
      "widget": "test-widget-to-weather-get",
      "is_private": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-service-new-import"
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
          "name": "test-entityservice-service-new-import"
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
    When I save response serviceID={{ (index .lastResponse.data 0)._id }}
    When I do GET /api/v4/entities/context-graph?_id=script_new-import_service/component_SC004C until response code is 200 and response array key "impact" contains:
    """json
    [
      "component_SC004C",
      "{{ .serviceID }}"
    ]
    """
    When I do GET /api/v4/entitybasics?_id=script_new-import_service/component_SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_new-import_service/component_SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_new-import_service",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=script_new-import_service/component_SC004C until response code is 200 and response array key "impact" contains:
    """json
    [
      "component_SC004C",
      "{{ .serviceID }}"
    ]
    """
    When I do GET /api/v4/entitybasics?_id=script_new-import_service_2/component_SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_new-import_service_2/component_SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_new-import_service_2",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=script_new-import_service_2/component_SC004C
    Then the response code should be 200
    Then the response array key "impact" should contain:
    """json
    [
      "component_SC004C",
      "{{ .serviceID }}"
    ]
    """
    When I do GET /api/v4/entityservices/{{ .serviceID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
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
    When I do GET /api/v4/entities/context-graph?_id={{ .serviceID }}
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "script_new-import_service/component_SC004C",
      "script_new-import_service_2/component_SC004C"
    ]
    """

  Scenario: given set import when component is deleted for the created resource should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "test-entity-contextgraph-new-import-resource-to-delete-conflict-1",
        "component": "test-entity-contextgraph-new-import-component-to-delete-conflict-1",
        "type": "resource",
        "action": "set",
        "enabled": true
      },
      {
        "name": "test-entity-contextgraph-new-import-component-to-delete-conflict-1",
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given set import with when component is disabled for the created resource should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "test-entity-contextgraph-new-import-resource-to-disable-conflict-2",
        "component": "test-entity-contextgraph-new-import-component-to-disable-conflict-2",
        "type": "resource",
        "action": "set",
        "enabled": true
      },
      {
        "name": "test-entity-contextgraph-new-import-component-to-disable-conflict-2",
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
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given set import with when component is already disabled for the created resource should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "test-entity-contextgraph-new-import-resource-to-disable-conflict-3",
        "component": "test-entity-contextgraph-new-import-component-to-disable-conflict-3",
        "type": "resource",
        "action": "set",
        "enabled": true
      }
    ]
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: given delete import should delete resource
    When I am admin
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "test-entity-contextgraph-new-import-resource-to-delete-1-1",
        "component": "test-entity-contextgraph-new-import-component-to-delete-1",
        "type": "resource",
        "infos": {
          "new_info": {
            "value": "2"
          }
        },
        "action": "delete",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-delete-1-1/test-entity-contextgraph-new-import-component-to-delete-1
    Then the response code should be 404
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-delete-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-delete-1",
      "name": "test-entity-contextgraph-new-import-component-to-delete-1",
      "enabled": true,
      "infos": {},
      "type": "component",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entity-contextgraph-new-import-component-to-delete-1 until response code is 200 and body contains:
    """json
    {
      "depends": [
        "test-entity-contextgraph-new-import-resource-to-delete-1-2/test-entity-contextgraph-new-import-component-to-delete-1"
      ],
      "impact": []
    }
    """

  Scenario: given delete import should delete component and its resources
    When I am admin
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-delete-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-component-to-delete-2"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-delete-2-1/test-entity-contextgraph-new-import-component-to-delete-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-resource-to-delete-2-1/test-entity-contextgraph-new-import-component-to-delete-2"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-delete-2-2/test-entity-contextgraph-new-import-component-to-delete-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-new-import-resource-to-delete-2-2/test-entity-contextgraph-new-import-component-to-delete-2"
    }
    """
    When I do PUT /api/v4/contextgraph-import?source=test-new-import-source:
    """json
    [
      {
        "name": "test-entity-contextgraph-new-import-component-to-delete-2",
        "type": "component",
        "infos": {
          "new_info": {
            "value": "2"
          }
        },
        "action": "delete",
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
    When I do GET /api/v4/entities?search=test-entity-contextgraph-new-import-resource-to-delete-2
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-component-to-delete-2
    Then the response code should be 404
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-delete-2-1/test-entity-contextgraph-new-import-component-to-delete-2
    Then the response code should be 404
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-new-import-resource-to-delete-2-2/test-entity-contextgraph-new-import-component-to-delete-2
    Then the response code should be 404
