Feature: Copy a view
  I need to be able to copy a view
  Only admin should be able to copy a view

  @concurrent
  Scenario: given copy public view to public group request should create a public view
    When I am admin
    When I do POST /api/v4/view-copy/test-view-to-copy-1:
    """json
    {
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "group": "test-viewgroup-to-view-copy-1",
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
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
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": false,
      "group": {
        "_id": "test-viewgroup-to-view-copy-1",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "created": 1611229670,
        "title": "test-viewgroup-to-view-copy-1-title",
        "updated": 1611229670
      },
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
      ]
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": false,
      "group": {
        "_id": "test-viewgroup-to-view-copy-1",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "created": 1611229670,
        "title": "test-viewgroup-to-view-copy-1-title",
        "updated": 1611229670
      },
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
      ],
      "tabs": [
        {
          "title": "test-tab-to-view-copy-1-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "is_private": false,
          "widgets": [
            {
              "title": "test-widget-to-view-copy-1-title",
              "type": "test-widget-to-view-copy-1-type",
              "is_private": false,
              "grid_parameters": {
                "desktop": {"x": 0, "y": 0}
              },
              "parameters": {
                "test-widget-to-view-copy-1-parameter-1": {
                  "test-widget-to-view-copy-1-parameter-1-subparameter": "test-widget-to-view-copy-1-parameter-1-subvalue"
                },
                "test-widget-to-view-copy-1-parameter-2": [
                  {
                    "test-widget-to-view-copy-1-parameter-2-subparameter": "test-widget-to-view-copy-1-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [
                {
                  "title": "test-widgetfilter-to-view-copy-1-title",
                  "is_private": false,
                  "is_user_preference": false,
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
                          "value": "test-widgetfilter-to-view-copy-1-pattern"
                        }
                      }
                    ]
                  ]
                },
                {
                  "title": "test-widgetfilter-to-view-copy-2-title",
                  "is_private": false,
                  "is_user_preference": false,
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
                          "value": "test-widgetfilter-to-view-copy-2-pattern"
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
          ]
        }
      ]
    }
    """
    Then the response key "_id" should not be "test-view-to-copy-1"
    Then the response key "tabs.0._id" should not be "test-tab-to-view-copy-1"
    Then the response key "tabs.0.widgets.0._id" should not be "test-widget-to-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.0._id" should not be "test-widgetfilter-to-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.1._id" should not be "test-widgetfilter-to-view-copy-2"
    Then the response key "tabs.0.widgets.0.parameters.mainFilter" should not be "test-widgetfilter-to-view-copy-1"
    Then I save response viewId={{ .lastResponse._id }}
    When I do GET /api/v4/permissions?search={{ .viewId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .viewId }}",
          "name": "{{ .viewId }}",
          "description": "Rights on view : test-view-to-copy-1-title-updated",
          "view": {
            "_id": "{{ .viewId }}",
            "title": "test-view-to-copy-1-title-updated"
          },
          "view_group": {
            "_id": "test-viewgroup-to-view-copy-1",
            "title": "test-viewgroup-to-view-copy-1-title"
          },
          "type": "RW"
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
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-view-copy-1&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "views": [
            {
              "_id": "test-view-to-copy-2",
              "is_private": false
            },
            {
              "title": "test-view-to-copy-1-title-updated",
              "is_private": false
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given copy request with not found group should return error
    When I am admin
    When I do POST /api/v4/view-copy/test-view-to-copy-3:
    """json
    {
      "description": "test-view-to-copy-3-description-updated",
      "enabled": true,
      "group": "test-viewgroup-to-view-copy-not-found",
      "title": "test-view-to-copy-3-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-3-tag-updated"
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "group": "Group doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given invalid copy request should return errors
    When I am admin
    When I do POST /api/v4/view-copy/test-view-to-copy-1:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "group": "Group is missing.",
        "title": "Title is missing."
      }
    }
    """

  @concurrent
  Scenario: given copy request and no auth user should not allow access
    When I do POST /api/v4/view-copy/test-view-to-copy-1
    Then the response code should be 401

  @concurrent
  Scenario: given copy request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/view-copy/test-view-to-copy-1
    Then the response code should be 403

  @concurrent
  Scenario: given copy request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/view-copy/test-view-to-copy-4
    Then the response code should be 403

  @concurrent
  Scenario: given copy request with not exist tab should return not found error
    When I am admin
    When I do POST /api/v4/view-copy/test-view-not-found
    Then the response code should be 404

  @concurrent
  Scenario: given copy owned private view request to the same owned private viewgroup should create a private view
    When I am admin
    When I do POST /api/v4/view-copy/test-private-view-to-copy-1:
    """json
    {
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "group": "test-private-viewgroup-to-copy-1",
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
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
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-copy-1",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-copy-1-title"
      },
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
      ]
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-copy-1",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-copy-1-title"
      },
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
      ],
      "tabs": [
        {
          "title": "test-private-tab-to-private-view-copy-1-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "is_private": true,
          "widgets": [
            {
              "title": "test-private-widget-to-private-view-copy-1-title",
              "type": "test-private-widget-to-private-view-copy-1-type",
              "is_private": true,
              "grid_parameters": {
                "desktop": {"x": 0, "y": 0}
              },
              "parameters": {
                "test-private-widget-to-private-view-copy-1-parameter-1": {
                  "test-private-widget-to-private-view-copy-1-parameter-1-subparameter": "test-private-widget-to-private-view-copy-1-parameter-1-subvalue"
                },
                "test-private-widget-to-private-view-copy-1-parameter-2": [
                  {
                    "test-private-widget-to-private-view-copy-1-parameter-2-subparameter": "test-private-widget-to-private-view-copy-1-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [
                {
                  "title": "test-private-widgetfilter-to-private-view-copy-1-title",
                  "is_private": true,
                  "is_user_preference": false,
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
                          "value": "test-private-widgetfilter-to-private-view-copy-1-pattern"
                        }
                      }
                    ]
                  ]
                },
                {
                  "title": "test-private-widgetfilter-to-private-view-copy-2-title",
                  "is_private": true,
                  "is_user_preference": false,
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
                          "value": "test-private-widgetfilter-to-private-view-copy-2-pattern"
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
          ]
        }
      ]
    }
    """
    Then the response key "_id" should not be "test-private-view-to-copy-1"
    Then the response key "tabs.0._id" should not be "test-private-tab-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0._id" should not be "test-private-widget-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.0._id" should not be "test-private-widgetfilter-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.1._id" should not be "test-private-widgetfilter-to-private-view-copy-2"
    Then the response key "tabs.0.widgets.0.parameters.mainFilter" should not be "test-private-widgetfilter-to-private-view-copy-1"    
    Then I save response viewId={{ .lastResponse._id }}
    When I do GET /api/v4/permissions?search={{ .viewId }}
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
    When I do GET /api/v4/view-groups?search=test-private-viewgroup-to-copy-1&with_views=true&with_private=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "views": [
            {
              "_id": "test-private-view-to-copy-1",
              "is_private": true
            },
            {
              "_id": "{{ .viewId }}",
              "is_private": true
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given copy owned private view request to another owned private viewgroup should create a private view
    When I am admin
    When I do POST /api/v4/view-copy/test-private-view-to-copy-1:
    """json
    {
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "group": "test-private-viewgroup-to-copy-2",
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
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
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-copy-2",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-copy-2-title"
      },
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
      ]
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-copy-2",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-copy-2-title"
      },
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
      ],
      "tabs": [
        {
          "title": "test-private-tab-to-private-view-copy-1-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "is_private": true,
          "widgets": [
            {
              "title": "test-private-widget-to-private-view-copy-1-title",
              "type": "test-private-widget-to-private-view-copy-1-type",
              "is_private": true,
              "grid_parameters": {
                "desktop": {"x": 0, "y": 0}
              },
              "parameters": {
                "test-private-widget-to-private-view-copy-1-parameter-1": {
                  "test-private-widget-to-private-view-copy-1-parameter-1-subparameter": "test-private-widget-to-private-view-copy-1-parameter-1-subvalue"
                },
                "test-private-widget-to-private-view-copy-1-parameter-2": [
                  {
                    "test-private-widget-to-private-view-copy-1-parameter-2-subparameter": "test-private-widget-to-private-view-copy-1-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [
                {
                  "title": "test-private-widgetfilter-to-private-view-copy-1-title",
                  "is_private": true,
                  "is_user_preference": false,
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
                          "value": "test-private-widgetfilter-to-private-view-copy-1-pattern"
                        }
                      }
                    ]
                  ]
                },
                {
                  "title": "test-private-widgetfilter-to-private-view-copy-2-title",
                  "is_private": true,
                  "is_user_preference": false,
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
                          "value": "test-private-widgetfilter-to-private-view-copy-2-pattern"
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
          ]
        }
      ]
    }
    """
    Then the response key "_id" should not be "test-private-view-to-copy-1"
    Then the response key "tabs.0._id" should not be "test-private-tab-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0._id" should not be "test-private-widget-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.0._id" should not be "test-private-widgetfilter-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.1._id" should not be "test-private-widgetfilter-to-private-view-copy-2"
    Then the response key "tabs.0.widgets.0.parameters.mainFilter" should not be "test-private-widgetfilter-to-private-view-copy-1"
    Then I save response viewId={{ .lastResponse._id }}
    When I do GET /api/v4/permissions?search={{ .viewId }}
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
    When I do GET /api/v4/view-groups?search=test-private-viewgroup-to-copy-2&with_views=true&with_private=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "views": [
            {
              "_id": "{{ .viewId }}",
              "is_private": true
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given copy owned private view request to not owned private viewgroup should return error
    When I am admin
    When I do POST /api/v4/view-copy/test-private-view-to-copy-1:
    """json
    {
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "group": "test-private-viewgroup-to-copy-3",
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "group": "Group is private."
      }
    }
    """

  @concurrent
  Scenario: given copy owned private view request to public viewgroup should create a public view
    When I am admin
    When I do POST /api/v4/view-copy/test-private-view-to-copy-1:
    """json
    {
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "group": "test-viewgroup-to-view-copy-2",
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
      ]
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-private-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": false,
      "group": {
        "_id": "test-viewgroup-to-view-copy-2",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "created": 1611229670,
        "title": "test-viewgroup-to-view-copy-2-title",
        "updated": 1611229670
      },
      "title": "test-private-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-copy-1-tag-updated"
      ],
      "tabs": [
        {
          "title": "test-private-tab-to-private-view-copy-1-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "is_private": false,
          "widgets": [
            {
              "title": "test-private-widget-to-private-view-copy-1-title",
              "type": "test-private-widget-to-private-view-copy-1-type",
              "is_private": false,
              "grid_parameters": {
                "desktop": {"x": 0, "y": 0}
              },
              "parameters": {
                "test-private-widget-to-private-view-copy-1-parameter-1": {
                  "test-private-widget-to-private-view-copy-1-parameter-1-subparameter": "test-private-widget-to-private-view-copy-1-parameter-1-subvalue"
                },
                "test-private-widget-to-private-view-copy-1-parameter-2": [
                  {
                    "test-private-widget-to-private-view-copy-1-parameter-2-subparameter": "test-private-widget-to-private-view-copy-1-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [
                {
                  "title": "test-private-widgetfilter-to-private-view-copy-1-title",
                  "is_private": false,
                  "is_user_preference":false,
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
                          "value": "test-private-widgetfilter-to-private-view-copy-1-pattern"
                        }
                      }
                    ]
                  ]
                },
                {
                  "title": "test-private-widgetfilter-to-private-view-copy-2-title",
                  "is_private": false,
                  "is_user_preference":false,
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
                          "value": "test-private-widgetfilter-to-private-view-copy-2-pattern"
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
          ]
        }
      ]
    }
    """
    Then the response key "_id" should not be "test-private-view-to-copy-1"
    Then the response key "tabs.0._id" should not be "test-private-tab-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0._id" should not be "test-private-widget-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.0._id" should not be "test-private-widgetfilter-to-private-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.1._id" should not be "test-private-widgetfilter-to-private-view-copy-2"
    Then the response key "tabs.0.widgets.0.parameters.mainFilter" should not be "test-private-widgetfilter-to-private-view-copy-1"
    Then I save response viewId={{ .lastResponse._id }}
    When I do GET /api/v4/permissions?search={{ .viewId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .viewId }}",
          "name": "{{ .viewId }}",
          "description": "Rights on view : test-private-view-to-copy-1-title-updated",
          "view": {
            "_id": "{{ .viewId }}",
            "title": "test-private-view-to-copy-1-title-updated"
          },
          "view_group": {
            "_id": "test-viewgroup-to-view-copy-2",
            "title": "test-viewgroup-to-view-copy-2-title"
          },
          "type": "RW"
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
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-view-copy-2&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "views": [
            {
              "title": "test-private-view-to-copy-1-title-updated",
              "is_private": false
            }
          ]
        }
      ]
    }
    """

  @concurrent
  Scenario: given copy public view to owned private group request should create a private view
    When I am admin
    When I do POST /api/v4/view-copy/test-view-to-copy-1:
    """json
    {
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "group": "test-private-viewgroup-to-copy-4",
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
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
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-copy-4",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-copy-4-title"
      },
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
      ]
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-copy-4",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-copy-4-title"
      },
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
      ],
      "tabs": [
        {
          "title": "test-tab-to-view-copy-1-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "is_private": true,
          "widgets": [
            {
              "title": "test-widget-to-view-copy-1-title",
              "type": "test-widget-to-view-copy-1-type",
              "is_private": true,
              "grid_parameters": {
                "desktop": {"x": 0, "y": 0}
              },
              "parameters": {
                "test-widget-to-view-copy-1-parameter-1": {
                  "test-widget-to-view-copy-1-parameter-1-subparameter": "test-widget-to-view-copy-1-parameter-1-subvalue"
                },
                "test-widget-to-view-copy-1-parameter-2": [
                  {
                    "test-widget-to-view-copy-1-parameter-2-subparameter": "test-widget-to-view-copy-1-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [
                {
                  "title": "test-widgetfilter-to-view-copy-1-title",
                  "is_private": true,
                  "is_user_preference": false,
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
                          "value": "test-widgetfilter-to-view-copy-1-pattern"
                        }
                      }
                    ]
                  ]
                },
                {
                  "title": "test-widgetfilter-to-view-copy-2-title",
                  "is_private": true,
                  "is_user_preference": false,
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
                          "value": "test-widgetfilter-to-view-copy-2-pattern"
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
          ]
        }
      ]
    }
    """
    Then the response key "_id" should not be "test-view-to-copy-1"
    Then the response key "tabs.0._id" should not be "test-tab-to-view-copy-1"
    Then the response key "tabs.0.widgets.0._id" should not be "test-widget-to-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.0._id" should not be "test-widgetfilter-to-view-copy-1"
    Then the response key "tabs.0.widgets.0.filters.1._id" should not be "test-widgetfilter-to-view-copy-2"
    Then the response key "tabs.0.widgets.0.parameters.mainFilter" should not be "test-widgetfilter-to-view-copy-1"
    Then I save response viewId={{ .lastResponse._id }}
    When I do GET /api/v4/permissions?search={{ .viewId }}
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
    When I do GET /api/v4/view-groups?search=test-private-viewgroup-to-copy-4&with_views=true&with_private=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "views": [
            {
              "title": "test-view-to-copy-1-title-updated",
              "is_private": true
            }
          ]
        }
      ]
    }
    """
