Feature: Import entities
  I need to be able to import entities

  Scenario: import unauthorized
    When I do PUT /api/v4/contextgraph/import
    Then the response code should be 401

  Scenario: import without permissions
    When I am noperms
    When I do PUT /api/v4/contextgraph/import
    Then the response code should be 403

  Scenario: import unauthorized
    When I do PUT /api/v4/contextgraph/import-partial
    Then the response code should be 401

  Scenario: import without permissions
    When I am noperms
    When I do PUT /api/v4/contextgraph/import-partial
    Then the response code should be 403

  Scenario: import with action create
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1",
            "name": "test-connector-name-contextgraph-import-1",
            "type": "connector",
            "action": "create",
            "enabled": true
          },
          {
            "_id": "test-component-contextgraph-import-1",
            "name": "test-component-contextgraph-import-1",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "create",
            "enabled": true
          },
          {
            "_id": "test-resource-contextgraph-import-1-1/test-component-contextgraph-import-1",
            "name": "test-resource-contextgraph-import-1-1",
            "type": "resource",
            "infos": {
              "test_info": {
                "description": "description 2",
                "value": "value 2"
              }
            },
            "action": "create",
            "enabled": true
          },
          {
            "_id": "test-resource-contextgraph-import-1-2/test-component-contextgraph-import-1",
            "name": "test-resource-contextgraph-import-1-2",
            "type": "resource",
            "infos": {
              "test_info": {
                "description": "description 3",
                "value": "value 3"
              }
            },
            "action": "create",
            "enabled": true
          }
        ],
        "links": [
          {
            "from": ["test-component-contextgraph-import-1"],
            "to": "test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1",
            "action": "create"
          },
          {
            "from": [
              "test-resource-contextgraph-import-1-1/test-component-contextgraph-import-1",
              "test-resource-contextgraph-import-1-2/test-component-contextgraph-import-1"
            ],
            "to": "test-component-contextgraph-import-1",
            "action": "create"
          },
          {
            "from": ["test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1"],
            "to": "test-resource-contextgraph-import-1-1/test-component-contextgraph-import-1",
            "action": "create"
          },
          {
            "from": ["test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1"],
            "to": "test-resource-contextgraph-import-1-2/test-component-contextgraph-import-1",
            "action": "create"
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
    When I do GET /api/v4/entitybasics?_id=test-component-contextgraph-import-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-component-contextgraph-import-1",
      "name": "test-component-contextgraph-import-1",
      "connector": "test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1",
      "component": "test-component-contextgraph-import-1",
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
    When I do GET /api/v4/entities/context-graph?_id=test-component-contextgraph-import-1
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-resource-contextgraph-import-1-1/test-component-contextgraph-import-1",
      "test-resource-contextgraph-import-1-2/test-component-contextgraph-import-1"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1"
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-import-1-1/test-component-contextgraph-import-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-resource-contextgraph-import-1-1/test-component-contextgraph-import-1",
      "name": "test-resource-contextgraph-import-1-1",
      "connector": "test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1",
      "component": "test-component-contextgraph-import-1",
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
      "import_source": "test-import-source"
    }
    """
    Then the response key "imported" should be greater than 0
    When I do GET /api/v4/entities?search=test-resource-contextgraph-import-1-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-resource-contextgraph-import-1-1/test-component-contextgraph-import-1",
          "name": "test-resource-contextgraph-import-1-1",
          "connector": "test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1",
          "component": "test-component-contextgraph-import-1",
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
          "import_source": "test-import-source"
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
    When I do GET /api/v4/entities/context-graph?_id=test-resource-contextgraph-import-1-1/test-component-contextgraph-import-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "impact": [
        "test-component-contextgraph-import-1"
      ],
      "depends": [
        "test-connector-contextgraph-import-1/test-connector-name-contextgraph-import-1"
      ]
    }
    """

  Scenario: import with action update
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-to-update",
            "name": "change name",
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "update",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-to-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-to-update",
      "name": "change name",
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

  Scenario: import with action update, when id doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-to-update-not-exist",
            "name": "change name",
            "type": "component",
            "infos": {
              "new_info": {
                "value": "2"
              }
            },
            "action": "update",
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

  Scenario: import with action set
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-to-set",
            "name": "change name",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-to-set
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-to-set",
      "name": "change name",
      "enabled": true,
      "infos": {
        "old_info": {
          "description": "",
          "name": "old_info",
          "value": "1"
        },
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
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-to-set",
            "name": "change name",
            "type": "component",
            "infos": {
              "old_info": {
                "value": "3"
              },
              "new_new_info": {
                "value": "1"
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-to-set
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-to-set",
      "name": "change name",
      "enabled": true,
      "infos": {
        "old_info": {
          "description": "",
          "name": "old_info",
          "value": "3"
        },
        "new_info": {
          "description": "",
          "name": "new_info",
          "value": "2"
        },
        "new_new_info": {
          "description": "",
          "name": "new_new_info",
          "value": "1"
        }
      },
      "type": "component",
      "impact_level": 1
    }
    """

  Scenario: import with action set, when id doesn't exist, new entity should be created
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-to-set-not-exist",
            "name": "test-entity-contextgraph-import-to-set-not-exist",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-to-set-not-exist
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-to-set-not-exist",
      "name": "test-entity-contextgraph-import-to-set-not-exist",
      "type": "component",
      "infos": {
            "new_info": {
                "value": "2"
            }
        },
      "enabled": true
    }
    """

  Scenario: import with action delete
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-resource-to-delete-1/test-entity-contextgraph-import-component-to-delete",
            "name": "change name",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-delete-1
    Then the response code should be 404
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-delete
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-component-to-delete",
      "name": "test-entity-contextgraph-import-component-to-delete",
      "enabled": true,
      "infos": {},
      "type": "component",
      "impact_level": 1
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entity-contextgraph-import-component-to-delete
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-entity-contextgraph-import-resource-to-delete-2/test-entity-contextgraph-import-component-to-delete"
      ],
      "impact": []
    }
    """
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-component-to-delete",
            "name": "change name",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-delete
    Then the response code should be 404
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-delete-2/test-entity-contextgraph-import-component-to-delete
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-resource-to-delete-2/test-entity-contextgraph-import-component-to-delete",
      "name": "test-entity-contextgraph-import-resource-to-delete-2",
      "enabled": true,
      "type": "resource",
      "impact_level": 1
    }
    """

  Scenario: import with action delete, when id doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-resource-to-delete-not-exist",
            "name": "change name",
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

  Scenario: import with action enable
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-component-to-enable",
            "name": "change name",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-enable
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-component-to-enable",
      "name": "test-entity-contextgraph-import-component-to-enable",
      "enabled": true,
      "type": "component",
      "impact_level": 1
    }
    """

  Scenario: import with action enable, when id doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-component-to-enable-not-exist",
            "name": "change name",
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

  Scenario: import with action disable
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-component-to-disable",
            "name": "change name",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-disable
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-component-to-disable",
      "name": "test-entity-contextgraph-import-component-to-disable",
      "enabled": false,
      "type": "component",
      "impact_level": 1
    }
    """

  Scenario: import with action disable, when id doesn't exist, status should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-entity-contextgraph-import-component-to-disable-not-exist",
            "name": "change name",
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

  Scenario: import with create links action
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "links": [
          {
            "from": [
              "test-entity-contextgraph-import-resource-to-link-1/test-entity-contextgraph-import-component-to-link",
              "test-entity-contextgraph-import-resource-to-link-2/test-entity-contextgraph-import-component-to-link"
            ],
            "to": "test-entity-contextgraph-import-component-to-link",
            "action": "create"
          },
          {
            "from": [
              "test-entity-contextgraph-import-component-to-link"
            ],
            "to": "test-entity-contextgraph-import-connector-to-link",
            "action": "create"
          },
          {
            "from": [
              "test-entity-contextgraph-import-connector-to-link"
            ],
            "to": "test-entity-contextgraph-import-resource-to-link-1/test-entity-contextgraph-import-component-to-link",
            "action": "create"
          },
          {
            "from": [
              "test-entity-contextgraph-import-connector-to-link"
            ],
            "to": "test-entity-contextgraph-import-resource-to-link-2/test-entity-contextgraph-import-component-to-link",
            "action": "create"
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-link
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-component-to-link",
      "enabled": true,
      "connector": "test-entity-contextgraph-import-connector-to-link",
      "impact_level": 1,
      "name": "test-entity-contextgraph-import-component-to-link",
      "component": "test-entity-contextgraph-import-component-to-link",
      "type": "component"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entity-contextgraph-import-component-to-link
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-entity-contextgraph-import-resource-to-link-1/test-entity-contextgraph-import-component-to-link",
      "test-entity-contextgraph-import-resource-to-link-2/test-entity-contextgraph-import-component-to-link"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-entity-contextgraph-import-connector-to-link"
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-link-1/test-entity-contextgraph-import-component-to-link
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-resource-to-link-1/test-entity-contextgraph-import-component-to-link",
      "connector": "test-entity-contextgraph-import-connector-to-link",
      "enabled": true,
      "component": "test-entity-contextgraph-import-component-to-link",
      "impact_level": 1,
      "infos": {
        "old_info": {
          "value": "1"
        }
      },
      "name": "test-entity-contextgraph-import-resource-to-link-1",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entity-contextgraph-import-resource-to-link-1/test-entity-contextgraph-import-component-to-link
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-entity-contextgraph-import-connector-to-link"
      ],
      "impact": [
        "test-entity-contextgraph-import-component-to-link"
      ]
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-link-2/test-entity-contextgraph-import-component-to-link
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-resource-to-link-2/test-entity-contextgraph-import-component-to-link",
      "connector": "test-entity-contextgraph-import-connector-to-link",
      "enabled": true,
      "component": "test-entity-contextgraph-import-component-to-link",
      "impact_level": 1,
      "infos": {
        "old_info": {
          "value": "1"
        }
      },
      "name": "test-entity-contextgraph-import-resource-to-link-2",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entity-contextgraph-import-resource-to-link-2/test-entity-contextgraph-import-component-to-link
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "test-entity-contextgraph-import-connector-to-link"
      ],
      "impact": [
        "test-entity-contextgraph-import-component-to-link"
      ]
    }
    """

  Scenario: import with delete links action
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "links": [
          {
            "from": [
              "test-entity-contextgraph-import-resource-to-unlink-1/test-entity-contextgraph-import-component-to-unlink",
              "test-entity-contextgraph-import-resource-to-unlink-2/test-entity-contextgraph-import-component-to-unlink"
            ],
            "to": "test-entity-contextgraph-import-component-to-unlink",
            "action": "delete"
          },
          {
            "from": [
              "test-entity-contextgraph-import-component-to-unlink"
            ],
            "to": "test-entity-contextgraph-import-connector-to-unlink",
            "action": "delete"
          },
          {
            "from": [
              "test-entity-contextgraph-import-connector-to-unlink"
            ],
            "to": "test-entity-contextgraph-import-resource-to-unlink-1/test-entity-contextgraph-import-component-to-unlink",
            "action": "delete"
          },
          {
            "from": [
              "test-entity-contextgraph-import-connector-to-unlink"
            ],
            "to": "test-entity-contextgraph-import-resource-to-unlink-2/test-entity-contextgraph-import-component-to-unlink",
            "action": "delete"
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-unlink
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-component-to-unlink",
      "enabled": true,
      "impact_level": 1,
      "name": "test-entity-contextgraph-import-component-to-unlink",
      "type": "component"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-unlink-1/test-entity-contextgraph-import-component-to-unlink
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-resource-to-unlink-1/test-entity-contextgraph-import-component-to-unlink",
      "enabled": true,
      "impact_level": 1,
      "infos": {
        "old_info": {
          "value": "1"
        }
      },
      "name": "test-entity-contextgraph-import-resource-to-unlink-1",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-unlink-2/test-entity-contextgraph-import-component-to-unlink
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entity-contextgraph-import-resource-to-unlink-2/test-entity-contextgraph-import-component-to-unlink",
      "enabled": true,
      "impact_level": 1,
      "infos": {
        "old_info": {
          "value": "1"
        }
      },
      "name": "test-entity-contextgraph-import-resource-to-unlink-2",
      "type": "resource"
    }
    """

  Scenario: import with action create, without type, should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "SC003C",
            "name": "SC003C",
            "action": "create",
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

  Scenario: import with action create, without _id, should be failed
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "type": "component",
            "name": "SC003C",
            "action": "create",
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

  Scenario: import with action create, patterns should be only in services
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "SC003C",
            "type": "component",
            "name": "SC003C",
            "action": "create",
            "enabled": true,
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "regexp",
                    "value": "script_import"
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

  Scenario: import with wrong cis action
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "SC003C",
            "type": "component",
            "name": "SC003C",
            "action": "some",
            "enabled": true,
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "regexp",
                    "value": "script_import"
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

  Scenario: import with wrong links action
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "links": [
          {
            "from": [
              "script_import_service/SC004C",
              "script_import_service_2/SC004C"
            ],
            "to": "SC004C",
            "action": "some",
            "infos": {},
            "id": "id_0",
            "properties": []
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

  Scenario: import with not implemented links action
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "links": [
          {
            "from": [
              "script_import_service/SC004C",
              "script_import_service_2/SC004C"
            ],
            "to": "SC004C",
            "action": "update",
            "infos": {},
            "id": "id_0",
            "properties": []
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
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "links": [
          {
            "from": [
              "script_import_service/SC004C",
              "script_import_service_2/SC004C"
            ],
            "to": "SC004C",
            "action": "enable",
            "infos": {},
            "id": "id_0",
            "properties": []
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
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "links": [
          {
            "from": [
              "script_import_service/SC004C",
              "script_import_service_2/SC004C"
            ],
            "to": "SC004C",
            "action": "disable",
            "infos": {},
            "id": "id_0",
            "properties": []
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

  Scenario: import with bad cis
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": "test"
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

  Scenario: import with bad links
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "links": "test"
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

  Scenario: import with empty json
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {}
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
    """json
    {
      "status": "failed"
    }
    """

  Scenario: import with action create, without name, should have a name the same as _id
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "some_id",
            "type": "component",
            "action": "create",
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
    When I do GET /api/v4/entitybasics?_id=some_id
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "some_id",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "some_id",
      "type": "component"
    }
    """

  Scenario: import with action create, context graph should be valid for entity service
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "SC004C",
            "name": "SC004C",
            "type": "component",
            "infos": {},
            "action": "create",
            "enabled": true
          },
          {
            "_id": "script_import_service/SC004C",
            "name": "script_import_service",
            "type": "resource",
            "infos": {},
            "action": "create",
            "enabled": true
          },
          {
            "_id": "script_import_service_2/SC004C",
            "name": "script_import_service_2",
            "type": "resource",
            "infos": {},
            "action": "create",
            "enabled": true
          },
          {
            "_id": "test-entityservice-service-import",
            "name": "test-entityservice-service-import",
            "output_template": "abc",
            "impact_level": 1,
            "enabled": true,
            "entity_pattern": [
              [
                {
                  "field": "name",
                  "cond": {
                    "type": "regexp",
                    "value": "script_import_service"
                  }
                }
              ]
            ],
            "type": "service",
            "action": "create"
          }
        ],
        "links": [
          {
            "from": [
              "script_import_service/SC004C",
              "script_import_service_2/SC004C"
            ],
            "to": "SC004C",
            "action": "create",
            "infos": {},
            "id": "id_0",
            "properties": []
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
    Then the response array key "depends" should contain:
    """json
    [
      "script_import_service/SC004C",
      "script_import_service_2/SC004C"
    ]
    """
    When I do GET /api/v4/entitybasics?_id=script_import_service/SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_import_service/SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_import_service",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=script_import_service/SC004C until response code is 200 and response array key "impact" contains:
    """json
    [
      "SC004C",
      "test-entityservice-service-import"
    ]
    """
    When I do GET /api/v4/entitybasics?_id=script_import_service/SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_import_service/SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_import_service",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=script_import_service/SC004C until response code is 200 and response array key "impact" contains:
    """json
    [
      "SC004C",
      "test-entityservice-service-import"
    ]
    """
    When I do GET /api/v4/entitybasics?_id=script_import_service_2/SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_import_service_2/SC004C",
      "enabled": true,
      "impact_level": 1,
      "infos": {},
      "name": "script_import_service_2",
      "type": "resource"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=script_import_service_2/SC004C
    Then the response code should be 200
    Then the response array key "impact" should contain:
    """json
    [
      "SC004C",
      "test-entityservice-service-import"
    ]
    """
    When I do GET /api/v4/entityservices/test-entityservice-service-import
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entityservice-service-import",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "script_import_service"
            }
          }
        ]
      ],
      "impact_level": 1,
      "name": "test-entityservice-service-import",
      "output_template": "abc",
      "type": "service"
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entityservice-service-import
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "script_import_service/SC004C",
      "script_import_service_2/SC004C"
    ]
    """

  Scenario: import with action create, same links with different actions
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "SC100C",
            "name": "SC100C",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "create",
            "enabled": true
          },
          {
            "_id": "script_import/SC100C",
            "name": "script_import",
            "type": "resource",
            "infos": {
              "test_info": {
                "description": "description 2",
                "value": "value 2"
              }
            },
            "action": "create",
            "enabled": true
          },
          {
            "_id": "script_import_2/SC100C",
            "name": "script_import_2",
            "type": "resource",
            "infos": {
              "test_info": {
                "description": "description 3",
                "value": "value 3"
              }
            },
            "action": "create",
            "enabled": true
          }
        ],
        "links": [
          {
            "from": [
              "script_import/SC100C",
              "script_import_2/SC100C"
            ],
            "to": "SC100C",
            "action": "create",
            "infos": {},
            "id": "id_0",
            "properties": []
          },
          {
            "from": [
              "script_import/SC100C"
            ],
            "to": "SC100C",
            "action": "delete",
            "infos": {},
            "id": "id_0",
            "properties": []
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
    When I do GET /api/v4/entitybasics?_id=SC100C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "SC100C",
      "name": "SC100C",
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
    When I do GET /api/v4/entities/context-graph?_id=SC100C
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "depends": [
        "script_import_2/SC100C"
      ],
      "impact": []
    }
    """
    When I do GET /api/v4/entitybasics?_id=script_import/SC100C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_import/SC100C",
      "name": "script_import",
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
    When I do GET /api/v4/entitybasics?_id=script_import_2/SC100C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "script_import_2/SC100C",
      "name": "script_import_2",
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
    When I do GET /api/v4/entities/context-graph?_id=script_import_2/SC100C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "depends": [],
      "impact": ["SC100C"]
    }
    """

  Scenario: given import with create resource should set resource component_infos
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-component-contextgraph-import-29",
            "name": "test-component-contextgraph-import-29",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "create",
            "enabled": true
          }
        ],
        "links": []
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
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-resource-contextgraph-import-29/test-component-contextgraph-import-29",
            "name": "test-resource-contextgraph-import-29",
            "type": "resource",
            "action": "create",
            "enabled": true
          }
        ],
        "links": [
          {
            "from": [
              "test-resource-contextgraph-import-29/test-component-contextgraph-import-29"
            ],
            "to": "test-component-contextgraph-import-29",
            "action": "create",
            "infos": {},
            "id": "id_0",
            "properties": []
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
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-import-29/test-component-contextgraph-import-29
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "component": "test-component-contextgraph-import-29",
      "component_infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "value 1"
        }
      }
    }
    """

  Scenario: given import with create component should set resource component_infos
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-resource-contextgraph-import-30/test-component-contextgraph-import-30",
            "name": "test-resource-contextgraph-import-30",
            "type": "resource",
            "action": "create",
            "enabled": true
          }
        ],
        "links": []
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
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-component-contextgraph-import-30",
            "name": "test-component-contextgraph-import-30",
            "type": "component",
            "infos": {
              "test_info": {
                "description": "description 1",
                "value": "value 1"
              }
            },
            "action": "create",
            "enabled": true
          }
        ],
        "links": [
          {
            "from": [
              "test-resource-contextgraph-import-30/test-component-contextgraph-import-30"
            ],
            "to": "test-component-contextgraph-import-30",
            "action": "create",
            "infos": {},
            "id": "id_0",
            "properties": []
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
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-import-30/test-component-contextgraph-import-30
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "component": "test-component-contextgraph-import-30",
      "component_infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "value 1"
        }
      }
    }
    """
