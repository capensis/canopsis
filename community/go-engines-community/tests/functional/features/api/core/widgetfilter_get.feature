Feature: Get a widget filter
  I need to be able to get a widget filter
  Only admin should be able to get a widget filter

  Scenario: given all request should return filters
    When I am admin
    When I do GET /api/v4/widget-filters?widget=test-widget-to-filter-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-widgetfilter-to-get-1",
          "title": "test-widgetfilter-to-get-1-title",
          "is_private": false,
          "author": {
            "_id": "nopermsuser",
            "name": "nopermsuser"
          },
          "created": 1605263992,
          "updated": 1605263992,
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-get-1-pattern"
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
                  "value": "test-widgetfilter-to-get-1-pattern"
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
                  "value": "test-widgetfilter-to-get-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-get-5",
          "title": "test-widgetfilter-to-get-5-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992,
          "old_mongo_query": {
            "name": "test-widgetfilter-to-get-5-pattern"
          }
        },
        {
          "_id": "test-widgetfilter-to-get-2",
          "title": "test-widgetfilter-to-get-2-title",
          "is_private": true,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992,
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given public request should return filters
    When I am admin
    When I do GET /api/v4/widget-filters?widget=test-widget-to-filter-get&private=false
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-widgetfilter-to-get-1",
          "title": "test-widgetfilter-to-get-1-title",
          "is_private": false,
          "author": {
            "_id": "nopermsuser",
            "name": "nopermsuser"
          },
          "created": 1605263992,
          "updated": 1605263992,
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-get-1-pattern"
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
                  "value": "test-widgetfilter-to-get-1-pattern"
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
                  "value": "test-widgetfilter-to-get-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-get-5",
          "title": "test-widgetfilter-to-get-5-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992,
          "old_mongo_query": {
            "name": "test-widgetfilter-to-get-5-pattern"
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
    """

  Scenario: given private request should return filters
    When I am admin
    When I do GET /api/v4/widget-filters?widget=test-widget-to-filter-get&private=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-widgetfilter-to-get-2",
          "title": "test-widgetfilter-to-get-2-title",
          "is_private": true,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992,
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given get invalid all request should return error
    When I am admin
    When I do GET /api/v4/widget-filters
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "widget": "Widget is missing."
      }
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/widget-filters
    Then the response code should be 401

  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/widget-filters
    Then the response code should be 403

  Scenario: given get all request and auth user without view permission should not allow access
    When I am admin
    When I do GET /api/v4/widget-filters?widget=test-widget-to-filter-check-access
    Then the response code should be 403

  Scenario: given get all request and not exist widget should not allow access
    When I am admin
    When I do GET /api/v4/widget-filters?widget=test-widget-not-exist
    Then the response code should be 403

  Scenario: given get request should return ok
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-widgetfilter-to-get-1",
      "title": "test-widgetfilter-to-get-1-title",
      "is_private": false,
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-get-1-pattern"
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
              "value": "test-widgetfilter-to-get-1-pattern"
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
              "value": "test-widgetfilter-to-get-1-pattern"
            }
          }
        ]
      ]
    }
    """

  Scenario: given get private filter request should return ok
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-widgetfilter-to-get-2",
      "title": "test-widgetfilter-to-get-2-title",
      "is_private": true,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
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

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/widget-filters/test-widgetfilter-not-exist
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/widget-filters/test-widgetfilter-not-exist
    Then the response code should be 403

  Scenario: given get request and auth user without view permission should not allow access
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-3
    Then the response code should be 403

  Scenario: given get request and another user should return not found
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-4
    Then the response code should be 404

  Scenario: given get not exist filter request should not allow access
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-not-exist
    Then the response code should be 403

  Scenario: given get old filter request should return ok
    When I am admin
    When I do GET /api/v4/widget-filters/test-widgetfilter-to-get-5
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-widgetfilter-to-get-5",
      "title": "test-widgetfilter-to-get-5-title",
      "is_private": false,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992,
      "old_mongo_query": {
        "name": "test-widgetfilter-to-get-5-pattern"
      }
    }
    """
