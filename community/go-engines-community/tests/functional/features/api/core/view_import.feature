Feature: Import views
  I need to be able to import views
  Only admin should be able to import views

  Scenario: given import request should return views
    When I am admin
    When I do POST /api/v4/view-import:
    """json
    [
      {
        "title": "test-viewgroup-to-import-1-title",
        "views": [
          {
            "title": "test-view-to-import-1-title",
            "description": "test-view-to-import-1-description",
            "enabled": true,
            "periodic_refresh": {
              "enabled": true,
              "value": 1,
              "unit": "s"
            },
            "tags": [
              "test-view-to-import-1-tag"
            ],
            "tabs": [
              {
                "title": "test-tab-to-import-1-title",
                "widgets": [
                  {
                    "title": "test-widget-to-import-1-title",
                    "type": "test-widget-to-import-1-type",
                    "grid_parameters": {
                      "desktop": {"x": 0, "y": 0}
                    },
                    "parameters": {
                      "mainFilter": "test-widgetfilter-to-import-2",
                      "test-widget-to-view-import-1-parameter-1": {
                        "test-widget-to-view-import-1-parameter-1-subparameter": "test-widget-to-view-import-1-parameter-1-subvalue"
                      },
                      "test-widget-to-view-import-1-parameter-2": [
                        {
                          "test-widget-to-view-import-1-parameter-2-subparameter": "test-widget-to-view-import-1-parameter-2-subvalue"
                        }
                      ]
                    },
                    "filters": [
                      {
                        "_id": "test-widgetfilter-to-import-1",
                        "title": "test-widgetfilter-to-import-1-title",
                        "is_private": false,
                        "alarm_pattern": [
                          [
                            {
                              "field": "v.component",
                              "cond": {
                                "type": "eq",
                                "value": "test-widgetfilter-to-import-1-pattern"
                              }
                            }
                          ]
                        ]
                      },
                      {
                        "_id": "test-widgetfilter-to-import-2",
                        "title": "test-widgetfilter-to-import-2-title",
                        "is_private": false,
                        "entity_pattern": [
                          [
                            {
                              "field": "name",
                              "cond": {
                                "type": "eq",
                                "value": "test-widgetfilter-to-import-2-pattern"
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
                "title": "test-tab-to-import-1-2-title"
              }
            ]
          }
        ]
      },
      {
        "title": "test-viewgroup-to-import-1-2-title"
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-import-1&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "test-viewgroup-to-import-1-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "views": [
            {
              "title": "test-view-to-import-1-title",
              "description": "test-view-to-import-1-description",
              "enabled": true,
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-import-1-tag"
              ],
              "author": {
                "_id": "root",
                "name": "root"
              }
            }
          ]
        },
        {
          "title": "test-viewgroup-to-import-1-2-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "views": []
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
    When I save response viewId={{ (index (index .lastResponse.data 0).views 0)._id }}
    When I do GET /api/v4/views/{{ .viewId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-view-to-import-1-title",
      "description": "test-view-to-import-1-description",
      "enabled": true,
      "periodic_refresh": {
        "enabled": true,
        "value": 1,
        "unit": "s"
      },
      "tags": [
        "test-view-to-import-1-tag"
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tabs": [
        {
          "title": "test-tab-to-import-1-title",
          "widgets": [
            {
              "title": "test-widget-to-import-1-title",
              "type": "test-widget-to-import-1-type",
              "grid_parameters": {
                "desktop": {"x": 0, "y": 0}
              },
              "parameters": {
                "test-widget-to-view-import-1-parameter-1": {
                  "test-widget-to-view-import-1-parameter-1-subparameter": "test-widget-to-view-import-1-parameter-1-subvalue"
                },
                "test-widget-to-view-import-1-parameter-2": [
                  {
                    "test-widget-to-view-import-1-parameter-2-subparameter": "test-widget-to-view-import-1-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [
                {
                  "title": "test-widgetfilter-to-import-1-title",
                  "is_private": false,
                  "author": {
                    "_id": "root",
                    "name": "root"
                  },
                  "alarm_pattern": [
                    [
                      {
                        "field": "v.component",
                        "cond": {
                          "type": "eq",
                          "value": "test-widgetfilter-to-import-1-pattern"
                        }
                      }
                    ]
                  ]
                },
                {
                  "title": "test-widgetfilter-to-import-2-title",
                  "is_private": false,
                  "author": {
                    "_id": "root",
                    "name": "root"
                  },
                  "entity_pattern": [
                    [
                      {
                        "field": "name",
                        "cond": {
                          "type": "eq",
                          "value": "test-widgetfilter-to-import-2-pattern"
                        }
                      }
                    ]
                  ]
                }
              ],
              "author": {
                "_id": "root",
                "name": "root"
              }
            }
          ],
          "author": {
            "_id": "root",
            "name": "root"
          }
        },
        {
          "title": "test-tab-to-import-1-2-title",
          "author": {
            "_id": "root",
            "name": "root"
          }
        }
      ]
    }
    """
    Then I save response filterId2={{ (index (index (index .lastResponse.tabs 0).widgets 0).filters 1)._id }}
    Then the response body should contain:
    """json
    {
      "tabs": [
        {
          "widgets": [
            {
              "parameters": {
                "mainFilter": "{{ .filterId2 }}"
              }
            }
          ]
        },
        {}
      ]
    }
    """

  Scenario: given import request with existed groups and models should return views
    When I am admin
    When I do POST /api/v4/view-import:
    """json
    [
      {
        "title": "test-viewgroup-to-import-2-2-title",
        "views": [
          {
            "title": "test-view-to-import-2-3-title"
          }
        ]
      },
      {
        "_id": "test-viewgroup-to-import-2-1",
        "views": [
          {
            "title": "test-view-to-import-2-4-title"
          },
          {
            "_id": "test-view-to-import-2-1"
          },
          {
            "title": "test-view-to-import-2-5-title"
          }
        ]
      },
      {
        "title": "test-viewgroup-to-import-2-3-title",
        "views": [
          {
            "_id": "test-view-to-import-2-2"
          }
        ]
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-import-2&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "title": "test-viewgroup-to-import-2-2-title",
          "views": [
            {
              "title": "test-view-to-import-2-3-title"
            }
          ]
        },
        {
          "_id": "test-viewgroup-to-import-2-1",
          "views": [
            {
              "title": "test-view-to-import-2-4-title"
            },
            {
              "_id": "test-view-to-import-2-1"
            },
            {
              "title": "test-view-to-import-2-5-title"
            }
          ]
        },
        {
          "title": "test-viewgroup-to-import-2-3-title",
          "views": [
            {
              "_id": "test-view-to-import-2-2"
            }
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

  Scenario: given invalid request should return error
    When I am admin
    When I do POST /api/v4/view-import:
    """json
    [
      {
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "0.title": "value is missing"
      }
    }
    """
    When I do POST /api/v4/view-import:
    """json
    [
      {
        "title": "test-viewgroup-to-import-3-title",
        "views": [
          {
          }
        ]
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "0.views.0.title": "value is missing"
      }
    }
    """
    When I do POST /api/v4/view-import:
    """json
    [
      {
        "title": "test-viewgroup-to-import-3-title",
        "views": [
          {
            "title": "test-view-to-import-3-title",
            "tabs": [
              {
              }
            ]
          }
        ]
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "0.views.0.tabs.0.title": "value is missing"
      }
    }
    """

    When I do POST /api/v4/view-import:
    """json
    [
      {
        "title": "test-viewgroup-to-import-3-title",
        "views": [
          {
            "title": "test-view-to-import-3-title",
            "tabs": [
              {
                "title": "test-tab-to-import-3-title",
                "widgets": [
                  {}
                ]
              }
            ]
          }
        ]
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "0.views.0.tabs.0.widgets.0.type": "value is missing"
      }
    }
    """

  Scenario: given import request should not return views without access
    When I am admin
    When I do POST /api/v4/view-import:
    """json
    [
      {
        "title": "test-viewgroup-to-import-3-title",
        "views": [
          {
            "_id": "test-view-to-check-access"
          }
        ]
      }
    ]
    """
    Then the response code should be 403
    When I am admin
    When I do POST /api/v4/view-import:
    """json
    [
      {
        "title": "test-viewgroup-to-import-3-title",
        "views": [
          {
            "_id": "test-view-not-found"
          }
        ]
      }
    ]
    """
    Then the response code should be 403

  Scenario: given get all request and no auth user should not allow access
    When I do POST /api/v4/view-import
    Then the response code should be 401

  Scenario: given get all request and auth user without view permission should not allow access
    When I am noperms
    When I do POST /api/v4/view-import
    Then the response code should be 403

