Feature: Update a widget filter
  I need to be able to update a widget filter
  Only admin should be able to update a widget filter

  @concurrent
  Scenario: given update public filter request should return ok
    When I am test-role-to-widget-filter-edit-2
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-1:
    """json
    {
      "title": "test-widgetfilter-to-update-1-title",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-1-pattern"
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
              "value": "test-widgetfilter-to-update-1-pattern"
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
              "value": "test-widgetfilter-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "test-user-to-widget-filter-edit-2",
        "name": "test-user-to-widget-filter-edit-2"
      },
      "title": "test-widgetfilter-to-update-1-title",
      "is_user_preference": false,
      "is_private": false,
      "created": 1605263992,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-1-pattern"
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
              "value": "test-widgetfilter-to-update-1-pattern"
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
              "value": "test-widgetfilter-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given update private filter request should return ok
    When I am test-role-to-widget-filter-edit-1
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-2:
    """json
    {
      "title": "test-widgetfilter-to-update-2-title",
      "is_user_preference": true,
      "corporate_alarm_pattern": "test-pattern-to-filter-edit-1",
      "corporate_entity_pattern": "test-pattern-to-filter-edit-2",
      "corporate_pbehavior_pattern": "test-pattern-to-filter-edit-3"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "test-user-to-widget-filter-edit-1",
        "name": "test-user-to-widget-filter-edit-1"
      },
      "title": "test-widgetfilter-to-update-2-title",
      "is_user_preference": true,
      "is_private": false,
      "created": 1605263992,
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
  Scenario: given update request for a filter with an old mongo query should return ok
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-6:
    """json
    {
      "title": "test-widgetfilter-to-update-6-title",
      "is_user_preference": false
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgetfilter-to-update-6-title",
      "is_user_preference": false,
      "is_private": false,
      "created": 1605263992,
      "old_mongo_query": {
        "name": "test-widgetfilter-to-update-6-pattern"
      }
    }
    """
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-6:
    """json
    {
      "title": "test-widgetfilter-to-update-6-title",
      "is_user_preference": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-6-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-widgetfilter-to-update-6-title",
      "is_user_preference": false,
      "is_private": false,
      "created": 1605263992,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-6-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response key "old_mongo_query" should not exist

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/widget-filters/test-widgetfilter-not-exist
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/widget-filters/test-widgetfilter-not-exist
    Then the response code should be 403

  @concurrent
  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-3:
    """json
    {
      "title": "test-widgetfilter-to-update-3-title",
      "is_user_preference": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-3-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given update request and another user should return not found
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-4:
    """json
    {
      "title": "test-widgetfilter-to-update-4-title",
      "is_user_preference": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-4-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: Scenario: given update request with another private status should return bad request error
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-5:
    """json
    {
      "title": "test-widgetfilter-to-update-5-title",
      "is_user_preference": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-5-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "is_user_preference": "IsUserPreference cannot be changed"
      }
    }
    """

  @concurrent
  Scenario: given update not exist filter request should not allow access
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-not-exist:
    """json
    {
      "title": "test-widgetfilter-not-exist-title",
      "is_user_preference": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-not-exist-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given update request with missing fields should return bad request error
    When I am admin
    When I do PUT /api/v4/widget-filters/test-widgetfilter-not-exist:
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

  @concurrent
  Scenario: given update owned filter request should return ok
    When I am admin
    When I do PUT /api/v4/widget-filters/test-private-widgetfilter-to-update-1:
    """json
    {
      "title": "test-private-widgetfilter-to-update-1-title",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-update-1-pattern"
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
              "value": "test-private-widgetfilter-to-update-1-pattern"
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
              "value": "test-private-widgetfilter-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-private-widgetfilter-to-update-1-title",
      "is_user_preference": false,
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-update-1-pattern"
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
              "value": "test-private-widgetfilter-to-update-1-pattern"
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
              "value": "test-private-widgetfilter-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given update not owned filter request should not allow access
    When I am admin
    When I do PUT /api/v4/widget-filters/test-private-widgetfilter-to-update-2:
    """json
    {
      "title": "test-private-widgetfilter-to-update-2-title",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-update-1-pattern"
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
              "value": "test-private-widgetfilter-to-update-1-pattern"
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
              "value": "test-private-widgetfilter-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given update owned filter request with api_private_view_groups but without api_view permissions should return ok
    When I am test-role-to-private-views-without-view-perm
    When I do PUT /api/v4/widget-filters/test-private-widgetfilter-to-update-3:
    """json
    {
      "title": "test-private-widgetfilter-to-update-3-title",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-update-3-pattern"
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
              "value": "test-private-widgetfilter-to-update-3-pattern"
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
              "value": "test-private-widgetfilter-to-update-3-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "title": "test-private-widgetfilter-to-update-3-title",
      "is_user_preference": false,
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-private-widgetfilter-to-update-3-pattern"
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
              "value": "test-private-widgetfilter-to-update-3-pattern"
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
              "value": "test-private-widgetfilter-to-update-3-pattern"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given update public filter request with api_private_view_groups
    but without api_view permissions should not allow access
    When I am test-role-to-private-views-without-view-perm
    When I do PUT /api/v4/widget-filters/test-widgetfilter-to-update-1:
    """json
    {
      "title": "test-widgetfilter-to-update-1-title",
      "is_user_preference": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-1-pattern"
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
              "value": "test-widgetfilter-to-update-1-pattern"
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
              "value": "test-widgetfilter-to-update-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403
