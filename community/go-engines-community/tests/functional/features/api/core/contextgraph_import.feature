Feature: Import entities
  I need to be able to import entities

  Scenario: import unauthorized
    When I do PUT /api/v4/contextgraph/import
    Then the response code should be 401

  Scenario: import without permissions
    When I am noperms
    When I do PUT /api/v4/contextgraph/import
    Then the response code should be 403

  Scenario: import with action create
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
        "json": {
            "cis": [
                {
                    "_id": "SC003C",
                    "name": "SC003C",
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
                    "_id": "script_import/SC003C",
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
                    "_id": "script_import_2/SC003C",
                    "name": "script_import",
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
                        "script_import/SC003C",
                        "script_import_2/SC003C"
                    ],
                    "to": "SC003C",
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
    When I do GET /api/v4/entitybasics?_id=SC003C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "SC003C",
        "name": "SC003C",
        "impact": [],
        "depends": [
            "script_import/SC003C",
            "script_import_2/SC003C"
        ],
        "enable_history": [],
        "infos": {
            "test_info": {
                "description": "description 1",
                "name": "test_info",
                "value": "value 1"
            }
        },
        "measurements": [],
        "enabled": true,
        "type": "component",
        "impact_level": 1
    }
    """
    When I do GET /api/v4/entitybasics?_id=script_import/SC003C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "script_import/SC003C",
        "name": "script_import",
        "impact": ["SC003C"],
        "depends": [],
        "enable_history": [],
        "measurements": [],
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
        "impact": [],
        "depends": [],
        "enable_history": [],
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
        "impact": [],
        "depends": [],
        "enable_history": [],
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
        "impact": [],
        "depends": [],
        "enable_history": [],
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
                    "name": "test-entity-contextgraph-import-to-set-not-exist-name",
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
        "name": "test-entity-contextgraph-import-to-set-not-exist-name",
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
                    "_id": "test-entity-contextgraph-import-resource-to-delete-1",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-delete-1
    Then the response code should be 404
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-delete
    Then the response body should contain:
    """json
    {
        "_id": "test-entity-contextgraph-import-component-to-delete",
        "name": "test-entity-contextgraph-import-component-to-delete-name",
        "impact": [],
        "depends": [
            "test-entity-contextgraph-import-resource-to-delete-2"
        ],
        "enable_history": [],
        "enabled": true,
        "infos": {},
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-delete-2
    Then the response body should contain:
    """json
    {
        "_id": "test-entity-contextgraph-import-resource-to-delete-2",
        "name": "test-entity-contextgraph-import-resource-to-delete-2-name",
        "impact": [],
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "infos": {},
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
        "name": "test-entity-contextgraph-import-component-to-enable-name",
        "impact": [],
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "infos": {},
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
        "name": "test-entity-contextgraph-import-component-to-disable-name",
        "impact": [],
        "depends": [],
        "enable_history": [],
        "enabled": false,
        "infos": {},
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
                        "test-entity-contextgraph-import-resource-to-link-1",
                        "test-entity-contextgraph-import-resource-to-link-2"
                    ],
                    "to": "test-entity-contextgraph-import-component-to-link",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-link
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "test-entity-contextgraph-import-component-to-link",
        "depends": [
            "test-entity-contextgraph-import-resource-to-link-1",
            "test-entity-contextgraph-import-resource-to-link-2"
        ],
        "enable_history": [],
        "enabled": true,
        "impact": [],
        "impact_level": 1,
        "infos": {},
        "measurements": null,
        "name": "test-entity-contextgraph-import-component-to-link-name",
        "description": "test-entity-contextgraph-import-component-to-link-description",
        "type": "component"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-link-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "test-entity-contextgraph-import-resource-to-link-1",
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "impact": [
            "test-entity-contextgraph-import-component-to-link"
        ],
        "impact_level": 1,
        "infos": {
          "old_info": {
            "value": "1"
          }
        },
        "measurements": null,
        "name": "test-entity-contextgraph-import-resource-to-link-1-name",
        "description": "test-entity-contextgraph-import-resource-to-link-1-description",
        "type": "resource"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-link-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "test-entity-contextgraph-import-resource-to-link-2",
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "impact": [
            "test-entity-contextgraph-import-component-to-link"
        ],
        "impact_level": 1,
        "infos": {
          "old_info": {
            "value": "1"
          }
        },
        "measurements": null,
        "name": "test-entity-contextgraph-import-resource-to-link-2-name",
        "description": "test-entity-contextgraph-import-resource-to-link-2-description",
        "type": "resource"
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
                        "test-entity-contextgraph-import-resource-to-unlink-1",
                        "test-entity-contextgraph-import-resource-to-unlink-2"
                    ],
                    "to": "test-entity-contextgraph-import-component-to-unlink",
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
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-component-to-unlink
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "test-entity-contextgraph-import-component-to-unlink",
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "impact": [],
        "impact_level": 1,
        "infos": {},
        "measurements": null,
        "name": "test-entity-contextgraph-import-component-to-unlink-name",
        "description": "test-entity-contextgraph-import-component-to-unlink-description",
        "type": "component"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-unlink-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "test-entity-contextgraph-import-resource-to-unlink-1",
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "impact": [],
        "impact_level": 1,
        "infos": {
          "old_info": {
            "value": "1"
          }
        },
        "measurements": null,
        "name": "test-entity-contextgraph-import-resource-to-unlink-1-name",
        "description": "test-entity-contextgraph-import-resource-to-unlink-1-description",
        "type": "resource"
    }
    """
    When I do GET /api/v4/entitybasics?_id=test-entity-contextgraph-import-resource-to-unlink-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "test-entity-contextgraph-import-resource-to-unlink-2",
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "impact": [],
        "impact_level": 1,
        "infos": {
          "old_info": {
            "value": "1"
          }
        },
        "measurements": null,
        "name": "test-entity-contextgraph-import-resource-to-unlink-2-name",
        "description": "test-entity-contextgraph-import-resource-to-unlink-2-description",
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
                    "_id": "SC003C"
                    "type": "component",
                    "name": "SC003C",
                    "action": "create",
                    "enabled": true,
                    "entity_patterns": [
                        {"name": {"regex_match": "script_import"}}
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
                    "_id": "SC003C"
                    "type": "component",
                    "name": "SC003C",
                    "action": "some",
                    "enabled": true,
                    "entity_patterns": [
                        {"name": {"regex_match": "script_import"}}
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
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "impact": [],
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
                    "name": "test-entityservice-service-import-name",
                    "output_template": "abc",
                    "impact_level": 1,
                    "enabled": true,
                    "entity_patterns": [
                        {"name": {"regex_match": "script_import_service"}}
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
    When I do GET /api/v4/entitybasics?_id=SC004C until response code is 200 and body contains:
    """json
    {
        "_id": "SC004C",
        "depends": [
            "script_import_service/SC004C",
            "script_import_service_2/SC004C"
        ],
        "enable_history": [],
        "enabled": true,
        "impact": [],
        "impact_level": 1,
        "infos": {},
        "name": "SC004C",
        "type": "component"
    }
    """
    When I do GET /api/v4/entitybasics?_id=script_import_service/SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "script_import_service/SC004C",
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "impact": [
            "SC004C",
            "test-entityservice-service-import"
        ],
        "impact_level": 1,
        "infos": {},
        "name": "script_import_service",
        "type": "resource"
    }
    """
    When I do GET /api/v4/entitybasics?_id=script_import_service_2/SC004C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "script_import_service_2/SC004C",
        "depends": [],
        "enable_history": [],
        "enabled": true,
        "impact": [
            "SC004C",
            "test-entityservice-service-import"
        ],
        "impact_level": 1,
        "infos": {},
        "name": "script_import_service_2",
        "type": "resource"
    }
    """
    When I do GET /api/v4/entityservices/test-entityservice-service-import
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entityservice-service-import",
      "depends": [
         "script_import_service/SC004C",
         "script_import_service_2/SC004C"
      ],
      "enabled": true,
      "entity_patterns": [
          {"name": {"regex_match": "script_import_service"}}
      ],
      "impact": [],
      "impact_level": 1,
      "name": "test-entityservice-service-import-name",
      "output_template": "abc",
      "type": "service"
    }
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
        "impact": [],
        "depends": [
            "script_import_2/SC100C"
        ],
        "enable_history": [],
        "infos": {
            "test_info": {
                "description": "description 1",
                "name": "test_info",
                "value": "value 1"
            }
        },
        "measurements": [],
        "enabled": true,
        "type": "component",
        "impact_level": 1
    }
    """
    When I do GET /api/v4/entitybasics?_id=script_import/SC100C
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
        "_id": "script_import/SC100C",
        "name": "script_import",
        "impact": [],
        "depends": [],
        "enable_history": [],
        "measurements": [],
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
        "impact": ["SC100C"],
        "depends": [],
        "enable_history": [],
        "measurements": [],
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

  Scenario: given import with create component and resource should update resource component_infos
    When I am admin
    When I do PUT /api/v4/contextgraph/import?source=test-import-source:
    """json
    {
      "json": {
        "cis": [
          {
            "_id": "test-component-contextgraph-import-28",
            "name": "test-component-contextgraph-import-28",
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
            "_id": "test-resource-contextgraph-import-28/test-component-contextgraph-import-28",
            "name": "test-resource-contextgraph-import-28",
            "type": "resource",
            "action": "create",
            "enabled": true
          }
        ],
        "links": [
          {
            "from": [
              "test-resource-contextgraph-import-28/test-component-contextgraph-import-28"
            ],
            "to": "test-component-contextgraph-import-28",
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
    When I do GET /api/v4/entitybasics?_id=test-resource-contextgraph-import-28/test-component-contextgraph-import-28
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "component": "test-component-contextgraph-import-28",
      "component_infos": {
        "test_info": {
          "description": "description 1",
          "name": "test_info",
          "value": "value 1"
        }
      }
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
