Feature: Create a widget filter
  I need to be able to create a widget filter
  Only admin should be able to create a widget filter

  @concurrent
  Scenario: given create public widgetfilter to public widget request should create public widgetfilter
    When I am test-role-to-widget-filter-edit-2
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-1-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-1-pattern"
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
              "value": "test-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-1-pattern"
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
        "_id": "test-user-to-widget-filter-edit-2",
        "name": "test-user-to-widget-filter-edit-2"
      },
      "title": "test-widgetfilter-to-create-1-title",
      "is_user_preference": false,
      "is_private": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-1-pattern"
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
              "value": "test-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/widget-filters/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "test-user-to-widget-filter-edit-2",
        "name": "test-user-to-widget-filter-edit-2"
      },
      "title": "test-widgetfilter-to-create-1-title",
      "is_user_preference": false,
      "is_private": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-1-pattern"
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
              "value": "test-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create private widgetfilter to public widget request should create public widgetfilter with is_user_preference flag
    When I am test-role-to-widget-filter-edit-1
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-4-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-4-pattern"
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
              "value": "test-widgetfilter-to-create-4-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-4-pattern"
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
        "_id": "test-user-to-widget-filter-edit-1",
        "name": "test-user-to-widget-filter-edit-1"
      },
      "title": "test-widgetfilter-to-create-4-title",
      "is_user_preference": true,
      "is_private": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-4-pattern"
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
              "value": "test-widgetfilter-to-create-4-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-4-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/widget-filters/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "test-user-to-widget-filter-edit-1",
        "name": "test-user-to-widget-filter-edit-1"
      },
      "title": "test-widgetfilter-to-create-4-title",
      "is_user_preference": true,
      "is_private": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-4-pattern"
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
              "value": "test-widgetfilter-to-create-4-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-4-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request with corporate patterns should return ok
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-2-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": false,
      "corporate_alarm_pattern": "test-pattern-to-filter-edit-1",
      "corporate_entity_pattern": "test-pattern-to-filter-edit-2",
      "corporate_pbehavior_pattern": "test-pattern-to-filter-edit-3"
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
      "title": "test-widgetfilter-to-create-2-title",
      "is_user_preference": false,
      "corporate_alarm_pattern": "test-pattern-to-filter-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-filter-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-filter-edit-2",
      "corporate_entity_pattern_title": "test-pattern-to-filter-edit-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-2-pattern"
            }
          }
        ]
      ],
      "corporate_pbehavior_pattern": "test-pattern-to-filter-edit-3",
      "corporate_pbehavior_pattern_title": "test-pattern-to-filter-edit-3-title",
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-3-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/widget-filters/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgetfilter-to-create-2-title",
      "is_user_preference": false,
      "corporate_alarm_pattern": "test-pattern-to-filter-edit-1",
      "corporate_alarm_pattern_title": "test-pattern-to-filter-edit-1-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-1-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-filter-edit-2",
      "corporate_entity_pattern_title": "test-pattern-to-filter-edit-2-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-2-pattern"
            }
          }
        ]
      ],
      "corporate_pbehavior_pattern": "test-pattern-to-filter-edit-3",
      "corporate_pbehavior_pattern_title": "test-pattern-to-filter-edit-3-title",
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-3-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request with saved patterns should return ok
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-3-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": true,
      "corporate_alarm_pattern": "test-pattern-to-filter-edit-4",
      "corporate_entity_pattern": "test-pattern-to-filter-edit-5",
      "corporate_pbehavior_pattern": "test-pattern-to-filter-edit-6"
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
      "title": "test-widgetfilter-to-create-3-title",
      "is_user_preference": true,
      "corporate_alarm_pattern": "test-pattern-to-filter-edit-4",
      "corporate_alarm_pattern_title": "test-pattern-to-filter-edit-4-title",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-4-pattern"
            }
          }
        ]
      ],
      "corporate_entity_pattern": "test-pattern-to-filter-edit-5",
      "corporate_entity_pattern_title": "test-pattern-to-filter-edit-5-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-5-pattern"
            }
          }
        ]
      ],
      "corporate_pbehavior_pattern": "test-pattern-to-filter-edit-6",
      "corporate_pbehavior_pattern_title": "test-pattern-to-filter-edit-6-title",
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-filter-edit-6-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/widget-filters
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/widget-filters
    Then the response code should be 403

  @concurrent
  Scenario: given create public request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "widget": "test-widget-to-filter-check-access",
      "title": "test-widgetfilter-to-create-3-title",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-5-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given create private request and auth user without view permission should not allow access
    When I am test-role-to-widget-filter-edit-1
    When I do POST /api/v4/widget-filters:
    """json
    {
      "widget": "test-widget-to-filter-check-access",
      "title": "test-widgetfilter-to-create-3-title",
      "is_user_preference": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-5-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given create request and with not exist widget should not allow access
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "widget": "test-widget-not-exist",
      "title": "test-widgetfilter-to-create-3-title",
      "is_user_preference": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-create-5-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given create request with missing fields should return bad request error
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "widget": "Widget is missing.",
        "is_user_preference": "IsUserPreference is missing.",
        "alarm_pattern": "AlarmPattern is missing.",
        "corporate_alarm_pattern": "CorporateAlarmPattern is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "corporate_entity_pattern": "CorporateEntityPattern is missing.",
        "pbehavior_pattern": "PbehaviorPattern is missing.",
        "corporate_pbehavior_pattern": "CorporatePbehaviorPattern is missing.",
        "weather_service_pattern": "WeatherServicePattern is missing."
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "alarm_pattern": [
        []
      ],
      "entity_pattern": [
        []
      ],
      "pbehavior_pattern": [
        []
      ],
      "weather_service_pattern": [
        []
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern.",
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "pbehavior_pattern": "PbehaviorPattern is invalid pbehavior pattern.",
        "weather_service_pattern": "WeatherServicePattern is invalid weather service pattern."
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-3-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": true,
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-3-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": true,
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
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-3-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": true,
      "corporate_pbehavior_pattern": "test-pattern-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "corporate_pbehavior_pattern": "CorporatePbehaviorPattern doesn't exist."
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-3-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": false,
      "corporate_alarm_pattern": "test-pattern-to-filter-edit-4"
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
  Scenario: given create public widgetfilter to owned private widget request should create private widgetfilter
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-private-widgetfilter-to-create-1-title",
      "widget": "test-private-widget-to-private-filter-create-1",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
      "title": "test-private-widgetfilter-to-create-1-title",
      "is_user_preference": false,
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/widget-filters/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-private-widgetfilter-to-create-1-title",
      "is_user_preference": false,
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create public widgetfilter to not owned private widget request should not allow access
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-private-widgetfilter-to-create-1-title",
      "widget": "test-private-widget-to-private-filter-create-2",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given create private widgetfilter to owned private widget request should create private widget with is_user_preference flag
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-private-widgetfilter-to-create-1-title",
      "widget": "test-private-widget-to-private-filter-create-1",
      "is_user_preference": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
      "title": "test-private-widgetfilter-to-create-1-title",
      "is_user_preference": true,
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/widget-filters/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-private-widgetfilter-to-create-1-title",
      "is_user_preference": true,
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
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
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-1-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create public widgetfilter to owned private widget request
    with api_private_view_groups but without api_view permissions should create private widgetfilter
    When I am test-role-to-private-views-without-view-perm
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-private-widgetfilter-to-create-3-title",
      "widget": "test-private-widget-to-private-filter-create-3",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-3-pattern"
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
              "value": "test-private-widgetfilter-to-create-3-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-3-pattern"
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
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "title": "test-private-widgetfilter-to-create-3-title",
      "is_user_preference": false,
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-3-pattern"
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
              "value": "test-private-widgetfilter-to-create-3-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-3-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/widget-filters/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "title": "test-private-widgetfilter-to-create-3-title",
      "is_user_preference": false,
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-3-pattern"
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
              "value": "test-private-widgetfilter-to-create-3-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-3-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create public widgetfilter to public widget request
    with api_private_view_groups but without api_view permissions should not allow access
    When I am test-role-to-private-views-without-view-perm
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-private-widgetfilter-to-create-3-title",
      "widget": "test-widget-to-filter-edit",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-3-pattern"
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
              "value": "test-private-widgetfilter-to-create-3-pattern"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-create-3-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403
