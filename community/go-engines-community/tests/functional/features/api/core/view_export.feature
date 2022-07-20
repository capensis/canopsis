Feature: Export views
  I need to be able to export views
  Only admin should be able to export views

  Scenario: given export request should return views
    When I am admin
    When I do POST /api/v4/view-export:
    """json
    {
      "groups": [
        {
          "_id": "test-viewgroup-to-export-1",
          "views": ["test-view-to-export-1", "test-view-to-export-2"]
        },
        {
          "_id": "test-viewgroup-to-export-2"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "groups": [
        {
          "title": "test-viewgroup-to-export-1-title",
          "views": [
            {
              "description": "test-view-to-export-1-description",
              "enabled": true,
              "title": "test-view-to-export-1-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-export-1-tag"
              ],
              "tabs": [
                {
                  "title": "test-tab-to-export-1-title",
                  "widgets": [
                    {
                      "title": "test-widget-to-export-1-title",
                      "type": "test-widget-to-export-1-type",
                      "grid_parameters": {
                        "desktop": {"x": 0, "y": 0}
                      },
                      "parameters": {
                        "mainFilter": "test-widgetfilter-to-export-2",
                        "test-widget-to-view-export-1-parameter-1": {
                          "test-widget-to-view-export-1-parameter-1-subparameter": "test-widget-to-view-export-1-parameter-1-subvalue"
                        },
                        "test-widget-to-view-export-1-parameter-2": [
                          {
                            "test-widget-to-view-export-1-parameter-2-subparameter": "test-widget-to-view-export-1-parameter-2-subvalue"
                          }
                        ]
                      },
                      "filters": [
                        {
                          "_id": "test-widgetfilter-to-export-1",
                          "title": "test-widgetfilter-to-export-1-title",
                          "alarm_pattern": [
                            [
                              {
                                "field": "v.component",
                                "cond": {
                                  "type": "eq",
                                  "value": "test-widgetfilter-to-export-1-pattern"
                                }
                              }
                            ]
                          ]
                        },
                        {
                          "_id": "test-widgetfilter-to-export-2",
                          "title": "test-widgetfilter-to-export-2-title",
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
                      ]
                    }
                  ]
                },
                {
                  "title": "test-tab-to-export-2-title",
                  "widgets": []
                }
              ]
            },
            {
              "description": "test-view-to-export-2-description",
              "enabled": true,
              "title": "test-view-to-export-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-export-2-tag"
              ],
              "tabs": []
            }
          ]
        },
        {
          "title": "test-viewgroup-to-export-2-title",
          "views": []
        }
      ],
      "views": []
    }
    """
    When I do POST /api/v4/view-export:
    """json
    {
      "groups": [
        {
          "_id": "test-viewgroup-to-export-1"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "groups": [
        {
          "title": "test-viewgroup-to-export-1-title",
          "views": []
        }
      ],
      "views": []
    }
    """
    When I do POST /api/v4/view-export:
    """json
    {
      "groups": [
        {
          "_id": "test-viewgroup-to-export-1",
          "views": ["test-view-to-export-2"]
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "groups": [
        {
          "title": "test-viewgroup-to-export-1-title",
          "views": [
            {
              "description": "test-view-to-export-2-description",
              "enabled": true,
              "title": "test-view-to-export-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-export-2-tag"
              ],
              "tabs": []
            }
          ]
        }
      ],
      "views": []
    }
    """
    When I do POST /api/v4/view-export:
    """json
    {
      "groups": [
        {
          "_id": "test-viewgroup-to-export-2"
        }
      ],
      "views": ["test-view-to-export-2"]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "groups": [
        {
          "title": "test-viewgroup-to-export-2-title",
          "views": []
        }
      ],
      "views": [
        {
          "description": "test-view-to-export-2-description",
          "enabled": true,
          "title": "test-view-to-export-2-title",
          "periodic_refresh": {
            "enabled": true,
            "value": 1,
            "unit": "s"
          },
          "tags": [
            "test-view-to-export-2-tag"
          ],
          "tabs": []
        }
      ]
    }
    """

  Scenario: given export request should not return views without access
    When I am admin
    When I do POST /api/v4/view-export:
    """json
    {
      "views": ["test-view-to-check-access"]
    }
    """
    Then the response code should be 403
    When I am admin
    When I do POST /api/v4/view-export:
    """json
    {
      "groups": [
        {
          "_id": "test-viewgroup-to-view-edit",
          "views": ["test-view-to-check-access"]
        }
      ]
    }
    """
    Then the response code should be 403
    When I do POST /api/v4/view-export:
    """json
    {
      "views": ["test-view-not-found"]
    }
    """
    Then the response code should be 403

  Scenario: given get all request and no auth user should not allow access
    When I do POST /api/v4/view-export
    Then the response code should be 401

  Scenario: given get all request and auth user without view permission should not allow access
    When I am noperms
    When I do POST /api/v4/view-export
    Then the response code should be 403

