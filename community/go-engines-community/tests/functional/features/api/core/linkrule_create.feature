Feature: Create an link rule
  I need to be able to create a link rule
  Only admin should be able to create a link rule

  @concurrent
  Scenario: given create alarm request with links should return ok
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-1-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-1-link-1-label",
          "category": "test-link-rule-to-create-1-link-1-category",
          "icon_name": "test-link-rule-to-create-1-link-1-icon",
          "url": "http://test-link-rule-to-create-1-link-1-url.com"
        },
        {
          "label": "test-link-rule-to-create-1-link-2-label",
          "category": "test-link-rule-to-create-1-link-2-category",
          "icon_name": "test-link-rule-to-create-1-link-2-icon",
          "url": "http://test-link-rule-to-create-1-link-2-url.com"
        }
      ],
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{ .Entity.Component }}` }}"
          },
          "regexp": {
            "message": "{{ `{{ .Value.Output }}` }}"
          },
          "collection": "link_mongo_data_regexp"
        }
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-1-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-1-resource"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-1-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-1-link-1-label",
          "category": "test-link-rule-to-create-1-link-1-category",
          "icon_name": "test-link-rule-to-create-1-link-1-icon",
          "url": "http://test-link-rule-to-create-1-link-1-url.com"
        },
        {
          "label": "test-link-rule-to-create-1-link-2-label",
          "category": "test-link-rule-to-create-1-link-2-category",
          "icon_name": "test-link-rule-to-create-1-link-2-icon",
          "url": "http://test-link-rule-to-create-1-link-2-url.com"
        }
      ],
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{ .Entity.Component }}` }}"
          },
          "regexp": {
            "message": "{{ `{{ .Value.Output }}` }}"
          },
          "collection": "link_mongo_data_regexp"
        }
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-1-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-1-resource"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/link-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-1-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-1-link-1-label",
          "category": "test-link-rule-to-create-1-link-1-category",
          "icon_name": "test-link-rule-to-create-1-link-1-icon",
          "url": "http://test-link-rule-to-create-1-link-1-url.com"
        },
        {
          "label": "test-link-rule-to-create-1-link-2-label",
          "category": "test-link-rule-to-create-1-link-2-category",
          "icon_name": "test-link-rule-to-create-1-link-2-icon",
          "url": "http://test-link-rule-to-create-1-link-2-url.com"
        }
      ],
      "external_data": {
        "component": {
          "type": "mongo",
          "sort_by": "status",
          "sort": "asc",
          "select": {
            "customer": "{{ `{{ .Entity.Component }}` }}"
          },
          "regexp": {
            "message": "{{ `{{ .Value.Output }}` }}"
          },
          "collection": "link_mongo_data_regexp"
        }
      },
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-1-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-1-resource"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create alarm request with source code should return ok
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-2-name",
      "type": "alarm",
      "enabled": true,
      "source_code": "function generate(alarms) { return [{\"label\": \"test-link-rule-to-create-2-link-1-label\",\"category\": \"test-link-rule-to-create-2-link-1-category\",\"icon_name\": \"test-link-rule-to-create-2-link-1-icon\",\"url\": \"http://test-link-rule-to-create-2-link-1-url.com\"}] }",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-2-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-2-resource"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-2-name",
      "type": "alarm",
      "enabled": true,
      "source_code": "function generate(alarms) { return [{\"label\": \"test-link-rule-to-create-2-link-1-label\",\"category\": \"test-link-rule-to-create-2-link-1-category\",\"icon_name\": \"test-link-rule-to-create-2-link-1-icon\",\"url\": \"http://test-link-rule-to-create-2-link-1-url.com\"}] }",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-2-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-2-resource"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/link-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-2-name",
      "type": "alarm",
      "enabled": true,
      "source_code": "function generate(alarms) { return [{\"label\": \"test-link-rule-to-create-2-link-1-label\",\"category\": \"test-link-rule-to-create-2-link-1-category\",\"icon_name\": \"test-link-rule-to-create-2-link-1-icon\",\"url\": \"http://test-link-rule-to-create-2-link-1-url.com\"}] }",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-2-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-2-resource"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create entity request with links should return ok
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-3-name",
      "type": "entity",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-3-link-1-label",
          "category": "test-link-rule-to-create-3-link-1-category",
          "icon_name": "test-link-rule-to-create-3-link-1-icon",
          "url": "http://test-link-rule-to-create-3-link-1-url.com"
        },
        {
          "label": "test-link-rule-to-create-3-link-2-label",
          "category": "test-link-rule-to-create-3-link-2-category",
          "icon_name": "test-link-rule-to-create-3-link-2-icon",
          "url": "http://test-link-rule-to-create-3-link-2-url.com"
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-3-resource"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-3-name",
      "type": "entity",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-3-link-1-label",
          "category": "test-link-rule-to-create-3-link-1-category",
          "icon_name": "test-link-rule-to-create-3-link-1-icon",
          "url": "http://test-link-rule-to-create-3-link-1-url.com"
        },
        {
          "label": "test-link-rule-to-create-3-link-2-label",
          "category": "test-link-rule-to-create-3-link-2-category",
          "icon_name": "test-link-rule-to-create-3-link-2-icon",
          "url": "http://test-link-rule-to-create-3-link-2-url.com"
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-3-resource"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/link-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-3-name",
      "type": "entity",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-3-link-1-label",
          "category": "test-link-rule-to-create-3-link-1-category",
          "icon_name": "test-link-rule-to-create-3-link-1-icon",
          "url": "http://test-link-rule-to-create-3-link-1-url.com"
        },
        {
          "label": "test-link-rule-to-create-3-link-2-label",
          "category": "test-link-rule-to-create-3-link-2-category",
          "icon_name": "test-link-rule-to-create-3-link-2-icon",
          "url": "http://test-link-rule-to-create-3-link-2-url.com"
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-3-resource"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request with corporate entity pattern should return success
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-4-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-4-link-1-label",
          "category": "test-link-rule-to-create-4-link-1-category",
          "icon_name": "test-link-rule-to-create-4-link-1-icon",
          "url": "http://test-link-rule-to-create-4-link-1-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-4-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-4-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-4-link-1-label",
          "category": "test-link-rule-to-create-4-link-1-category",
          "icon_name": "test-link-rule-to-create-4-link-1-icon",
          "url": "http://test-link-rule-to-create-4-link-1-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-4-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-2-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request with corporate alarm pattern should return success
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-5-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-5-link-1-label",
          "category": "test-link-rule-to-create-5-link-1-category",
          "icon_name": "test-link-rule-to-create-5-link-1-icon",
          "url": "http://test-link-rule-to-create-5-link-1-url.com"
        }
      ],
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-1-resource"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-5-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-5-link-1-label",
          "category": "test-link-rule-to-create-5-link-1-category",
          "icon_name": "test-link-rule-to-create-5-link-1-icon",
          "url": "http://test-link-rule-to-create-5-link-1-url.com"
        }
      ],
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-1-pattern"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-1-resource"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request with both corporate patterns should return success
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-6-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-6-link-1-label",
          "category": "test-link-rule-to-create-6-link-1-category",
          "icon_name": "test-link-rule-to-create-6-link-1-icon",
          "url": "http://test-link-rule-to-create-6-link-1-url.com"
        }
      ],
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-link-rule-to-create-6-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-6-link-1-label",
          "category": "test-link-rule-to-create-6-link-1-category",
          "icon_name": "test-link-rule-to-create-6-link-1-icon",
          "url": "http://test-link-rule-to-create-6-link-1-url.com"
        }
      ],
      "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
      "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-rule-edit-2-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request with absent alarm corporate pattern should return error
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-7-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-7-link-1-label",
          "category": "test-link-rule-to-create-7-link-1-category",
          "icon_name": "test-link-rule-to-create-7-link-1-icon",
          "url": "http://test-link-rule-to-create-7-link-1-url.com"
        }
      ],
      "corporate_alarm_pattern": "test-pattern-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_alarm_pattern": "CorporateAlarmPattern doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create request with absent alarm corporate pattern should return error
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-7-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-7-link-1-label",
          "category": "test-link-rule-to-create-7-link-1-category",
          "icon_name": "test-link-rule-to-create-7-link-1-icon",
          "url": "http://test-link-rule-to-create-7-link-1-url.com"
        }
      ],
      "corporate_entity_pattern": "test-pattern-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """    
    
  @concurrent
  Scenario: given create request with unacceptable alarm pattern and entity pattern fields for link rules should return error
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-7-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-7-link-1-label",
          "category": "test-link-rule-to-create-7-link-1-category",
          "icon_name": "test-link-rule-to-create-7-link-1-icon",
          "url": "http://test-link-rule-to-create-7-link-1-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-10-pattern"
            }
          },
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-7-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-7-link-1-label",
          "category": "test-link-rule-to-create-7-link-1-category",
          "icon_name": "test-link-rule-to-create-7-link-1-icon",
          "url": "http://test-link-rule-to-create-7-link-1-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-10-pattern"
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-7-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-7-link-1-label",
          "category": "test-link-rule-to-create-7-link-1-category",
          "icon_name": "test-link-rule-to-create-7-link-1-icon",
          "url": "http://test-link-rule-to-create-7-link-1-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-10-pattern"
            }
          },
          {
            "field": "v.linkd",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-7-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-7-link-1-label",
          "category": "test-link-rule-to-create-7-link-1-category",
          "icon_name": "test-link-rule-to-create-7-link-1-icon",
          "url": "http://test-link-rule-to-create-7-link-1-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-10-pattern"
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-7-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-7-link-1-label",
          "category": "test-link-rule-to-create-7-link-1-category",
          "icon_name": "test-link-rule-to-create-7-link-1-icon",
          "url": "http://test-link-rule-to-create-7-link-1-url.com"
        }
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-10-pattern"
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-create-7-name",
      "type": "alarm",
      "enabled": true,
      "links": [
        {
          "label": "test-link-rule-to-create-7-link-1-label",
          "category": "test-link-rule-to-create-7-link-1-category",
          "icon_name": "test-link-rule-to-create-7-link-1-icon",
          "url": "http://test-link-rule-to-create-7-link-1-url.com"
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-link-rule-to-create-10-resource"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """

  @concurrent
  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "type": "Type is missing.",
        "enabled": "Enabled is missing.",
        "links": "Links or SourceCode is required.",
        "source_code": "SourceCode or Links is required."
      }
    }
    """
    When I do POST /api/v4/link-rules:
    """json
    {
      "type": "alarm"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "enabled": "Enabled is missing.",
        "links": "Links or SourceCode is required.",
        "source_code": "SourceCode or Links is required.",
        "alarm_pattern": "AlarmPattern or EntityPattern is required.",
        "entity_pattern": "EntityPattern or AlarmPattern is required."
      }
    }
    """
    When I do POST /api/v4/link-rules:
    """json
    {
      "type": "entity"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "enabled": "Enabled is missing.",
        "links": "Links or SourceCode is required.",
        "source_code": "SourceCode or Links is required.",
        "entity_pattern": "EntityPattern is missing."
      }
    }
    """

  @concurrent
  Scenario: given create request with invalid patterns format should return bad request
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "wrong_field",
            "cond": {
              "type": "eq",
              "value": "ram"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "wrong_field",
            "cond": {
              "type": "eq",
              "value": "ram"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern.",
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/link-rules
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/link-rules
    Then the response code should be 403

  @concurrent
  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/link-rules:
    """json
    {
      "name": "test-link-rule-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """
