package contextgraph_test

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_contextgraph "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/contextgraph"
	mock_metrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/metrics"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestCheckServices(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	dataSets := []struct {
		services       []entityservice.EntityService
		entities       []types.Entity
		expectedResult []types.Entity
		name           string
	}{
		{
			name: "one entity is added to a single service",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Depends: []string{},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                    "id-1",
					Component:             "component-1",
					ImpactedServicesToAdd: []string{"serv-1"},
					ImpactedServices:      []string{"serv-1"},
					Impacts:               []string{"component-1", "serv-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1"},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "one entity is added to multiple services",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID: "serv-1",
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID: "serv-2",
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                    "id-1",
					Component:             "component-1",
					ImpactedServicesToAdd: []string{"serv-1", "serv-2"},
					ImpactedServices:      []string{"serv-1", "serv-2"},
					Impacts:               []string{"component-1", "serv-1", "serv-2"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1"},
					Type:    types.EntityTypeService,
				},
				{
					ID:      "serv-2",
					Depends: []string{"id-1"},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "one entity is added to multiple services impacted services to add/remove should be updated",
			entities: []types.Entity{
				{
					ID:                       "id-1",
					Component:                "component-1",
					Impacts:                  []string{"component-1"},
					ImpactedServicesToAdd:    []string{"serv-4"},
					ImpactedServicesToRemove: []string{"serv-0", "serv-2", "serv-3"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID: "serv-1",
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID: "serv-2",
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-1",
					Component:                "component-1",
					ImpactedServicesToAdd:    []string{"serv-1", "serv-2", "serv-4"},
					ImpactedServicesToRemove: []string{"serv-0", "serv-3"},
					ImpactedServices:         []string{"serv-1", "serv-2"},
					Impacts:                  []string{"component-1", "serv-1", "serv-2"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1"},
					Type:    types.EntityTypeService,
				},
				{
					ID:      "serv-2",
					Depends: []string{"id-1"},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "one entity is removed from a single service",
			entities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					ImpactedServices: []string{"serv-1"},
					Impacts:          []string{"component-1", "serv-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Depends: []string{"id-1"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-2",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-1",
					Component:                "component-1",
					ImpactedServices:         []string{},
					ImpactedServicesToRemove: []string{"serv-1"},
					Impacts:                  []string{"component-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "one entity is removed from multiple services",
			entities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					ImpactedServices: []string{"serv-1", "serv-2"},
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Depends: []string{"id-1"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-2",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Depends: []string{"id-1"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-2",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-1",
					Component:                "component-1",
					ImpactedServices:         []string{},
					ImpactedServicesToRemove: []string{"serv-1", "serv-2"},
					Impacts:                  []string{"component-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{},
					Type:    types.EntityTypeService,
				},
				{
					ID:      "serv-2",
					Depends: []string{},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "one entity is moved from one service to another",
			entities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					ImpactedServices: []string{"serv-1", "serv-2"},
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Depends: []string{"id-1"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Depends: []string{"id-1"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-2",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID: "serv-3",
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-1",
					Component:                "component-1",
					Impacts:                  []string{"component-1", "serv-1", "serv-3"},
					ImpactedServices:         []string{"serv-1", "serv-3"},
					ImpactedServicesToAdd:    []string{"serv-3"},
					ImpactedServicesToRemove: []string{"serv-2"},
				},
				{
					ID:      "serv-2",
					Depends: []string{},
					Type:    types.EntityTypeService,
				},
				{
					ID:      "serv-3",
					Depends: []string{"id-1"},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is added to a single service",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
				{
					ID:        "id-2",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID: "serv-1",
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                    "id-1",
					Component:             "component-1",
					ImpactedServicesToAdd: []string{"serv-1"},
					ImpactedServices:      []string{"serv-1"},
					Impacts:               []string{"component-1", "serv-1"},
				},
				{
					ID:                    "id-2",
					Component:             "component-1",
					ImpactedServicesToAdd: []string{"serv-1"},
					ImpactedServices:      []string{"serv-1"},
					Impacts:               []string{"component-1", "serv-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-2"},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is added to multiple services",
			entities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
				{
					ID:        "id-2",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID: "serv-1",
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID: "serv-2",
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                    "id-1",
					Component:             "component-1",
					ImpactedServicesToAdd: []string{"serv-1", "serv-2"},
					ImpactedServices:      []string{"serv-1", "serv-2"},
					Impacts:               []string{"component-1", "serv-1", "serv-2"},
				},
				{
					ID:                    "id-2",
					Component:             "component-1",
					ImpactedServicesToAdd: []string{"serv-1", "serv-2"},
					ImpactedServices:      []string{"serv-1", "serv-2"},
					Impacts:               []string{"component-1", "serv-1", "serv-2"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-2"},
					Type:    types.EntityTypeService,
				},
				{
					ID:      "serv-2",
					Depends: []string{"id-1", "id-2"},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is removed from a single service",
			entities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					ImpactedServices: []string{"serv-1"},
					Impacts:          []string{"component-1", "serv-1"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					ImpactedServices: []string{"serv-1"},
					Impacts:          []string{"component-1", "serv-1"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Depends: []string{"id-1"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-2",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-1",
					Component:                "component-1",
					Impacts:                  []string{"component-1"},
					ImpactedServices:         []string{},
					ImpactedServicesToRemove: []string{"serv-1"},
				},
				{
					ID:                       "id-2",
					Component:                "component-1",
					Impacts:                  []string{"component-1"},
					ImpactedServices:         []string{},
					ImpactedServicesToRemove: []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is removed from multiple services",
			entities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-1", "serv-2"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-1", "serv-2"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Depends: []string{"id-1", "id-2"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-2",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Depends: []string{"id-1", "id-2"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-2",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-1",
					Component:                "component-1",
					Impacts:                  []string{"component-1"},
					ImpactedServices:         []string{},
					ImpactedServicesToRemove: []string{"serv-1", "serv-2"},
				},
				{
					ID:                       "id-2",
					Component:                "component-1",
					Impacts:                  []string{"component-1"},
					ImpactedServices:         []string{},
					ImpactedServicesToRemove: []string{"serv-1", "serv-2"},
				},
				{
					ID:      "serv-1",
					Depends: []string{},
					Type:    types.EntityTypeService,
				},
				{
					ID:      "serv-2",
					Depends: []string{},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "multiple firstFindCallEntities is moved from one service to another",
			entities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					ImpactedServices: []string{"serv-1", "serv-2"},
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					ImpactedServices: []string{"serv-2", "serv-3"},
					Impacts:          []string{"component-1", "serv-2", "serv-3"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Depends: []string{"id-1"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Depends: []string{"id-1", "id-2"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-2",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-3",
						Depends: []string{"id-2"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-1",
					Component:                "component-1",
					Impacts:                  []string{"component-1", "serv-1", "serv-3"},
					ImpactedServices:         []string{"serv-1", "serv-3"},
					ImpactedServicesToAdd:    []string{"serv-3"},
					ImpactedServicesToRemove: []string{"serv-2"},
				},
				{
					ID:                       "id-2",
					Component:                "component-1",
					Impacts:                  []string{"component-1", "serv-1", "serv-3"},
					ImpactedServices:         []string{"serv-1", "serv-3"},
					ImpactedServicesToAdd:    []string{"serv-1"},
					ImpactedServicesToRemove: []string{"serv-2"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-2"},
					Type:    types.EntityTypeService,
				},
				{
					ID:      "serv-2",
					Depends: []string{},
					Type:    types.EntityTypeService,
				},
				{
					ID:      "serv-3",
					Depends: []string{"id-1", "id-2"},
					Type:    types.EntityTypeService,
				},
			},
		},
		{
			name: "no changes",
			entities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-1", "serv-2"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-1", "serv-2"},
				},
			},
			services: []entityservice.EntityService{
				{
					Entity: types.Entity{
						ID:      "serv-1",
						Depends: []string{"id-1", "id-2"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
				{
					Entity: types.Entity{
						ID:      "serv-2",
						Depends: []string{"id-1", "id-2"},
					},
					EntityPatterns: pattern.EntityPatternList{
						Patterns: []pattern.EntityPattern{
							{
								EntityFields: pattern.EntityFields{
									Component: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "component-1",
											},
										},
									},
								},
							},
						},
						Set:   true,
						Valid: true,
					},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					ImpactedServices: []string{"serv-1", "serv-2"},
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					ImpactedServices: []string{"serv-1", "serv-2"},
					Impacts:          []string{"component-1", "serv-1", "serv-2"},
				},
			},
		},
	}

	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater)

	for _, dataset := range dataSets {
		t.Run(dataset.name, func(t *testing.T) {
			storage.EXPECT().GetAll(gomock.Any()).Return(dataset.services, nil)

			result, err := manager.CheckServices(ctx, dataset.entities)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i].ID < result[j].ID
			})

			for idx := range result {
				sort.Slice(result[idx].ImpactedServices, func(i, j int) bool {
					return result[idx].ImpactedServices[i] < result[idx].ImpactedServices[j]
				})

				sort.Slice(result[idx].ImpactedServicesToAdd, func(i, j int) bool {
					return result[idx].ImpactedServicesToAdd[i] < result[idx].ImpactedServicesToAdd[j]
				})

				sort.Slice(result[idx].ImpactedServicesToRemove, func(i, j int) bool {
					return result[idx].ImpactedServicesToRemove[i] < result[idx].ImpactedServicesToRemove[j]
				})

				sort.Slice(result[idx].Impacts, func(i, j int) bool {
					return result[idx].Impacts[i] < result[idx].Impacts[j]
				})

				sort.Slice(result[idx].Depends, func(i, j int) bool {
					return result[idx].Depends[i] < result[idx].Depends[j]
				})
			}

			for idx := range dataset.expectedResult {
				sort.Slice(dataset.expectedResult[idx].ImpactedServices, func(i, j int) bool {
					return dataset.expectedResult[idx].ImpactedServices[i] < dataset.expectedResult[idx].ImpactedServices[j]
				})

				sort.Slice(dataset.expectedResult[idx].ImpactedServicesToAdd, func(i, j int) bool {
					return dataset.expectedResult[idx].ImpactedServicesToAdd[i] < dataset.expectedResult[idx].ImpactedServicesToAdd[j]
				})

				sort.Slice(dataset.expectedResult[idx].ImpactedServicesToRemove, func(i, j int) bool {
					return dataset.expectedResult[idx].ImpactedServicesToRemove[i] < dataset.expectedResult[idx].ImpactedServicesToRemove[j]
				})

				sort.Slice(dataset.expectedResult[idx].Impacts, func(i, j int) bool {
					return dataset.expectedResult[idx].Impacts[i] < dataset.expectedResult[idx].Impacts[j]
				})

				sort.Slice(dataset.expectedResult[idx].Depends, func(i, j int) bool {
					return dataset.expectedResult[idx].Depends[i] < dataset.expectedResult[idx].Depends[j]
				})
			}

			sort.Slice(dataset.expectedResult, func(i, j int) bool {
				return dataset.expectedResult[i].ID < dataset.expectedResult[j].ID
			})

			if !reflect.DeepEqual(result, dataset.expectedResult) {
				t.Error("result is not expected")
			}
		})
	}
}

func BenchmarkCenterCheckServices(b *testing.B) {
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	collection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	services := make([]entityservice.EntityService, 1000)
	for i := 0; i < 1000; i++ {
		services[i] = entityservice.EntityService{
			Entity: types.Entity{ID: fmt.Sprintf("serv-%d", i)},
			EntityPatterns: pattern.EntityPatternList{
				Patterns: []pattern.EntityPattern{
					{
						EntityFields: pattern.EntityFields{
							Component: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: types.OptionalString{
										Set:   true,
										Value: fmt.Sprintf("component-%d", i),
									},
								},
							},
						},
					},
				},
				Set: true,

				Valid: true,
			},
		}
	}

	entities := make([]types.Entity, 100)
	for i := 0; i < 100; i++ {
		impactedServices := make([]string, 50)
		for j := 0; j < 50; j++ {
			impactedServices[j] = fmt.Sprintf("serv-%d", rand.Intn(1000))
		}

		entities[i] = types.Entity{
			ID:               fmt.Sprintf("id-%d", i),
			Component:        fmt.Sprintf("component-%d", i),
			ImpactedServices: impactedServices,
		}
	}

	storage.EXPECT().GetAll(gomock.Any()).Return(services, nil).AnyTimes()
	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater)

	for i := 0; i < b.N; i++ {
		_, _ = manager.CheckServices(ctx, entities)
	}
}

func TestRecomputeService(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	collection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	dataSets := []struct {
		service                entityservice.EntityService
		firstFindCallEntities  []types.Entity
		secondFindCallEntities []types.Entity
		expectedResult         []types.Entity
		serviceName            string
		name                   string
	}{
		{
			name:        "service deleted",
			service:     entityservice.EntityService{},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:                    "id-3",
					Component:             "component-1",
					Impacts:               []string{"component-1", "serv-0", "serv-1", "serv-2"},
					ImpactedServices:      []string{"serv-0", "serv-1", "serv-2"},
					ImpactedServicesToAdd: []string{"serv-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1"},
					ImpactedServices: []string{},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-2"},
				},
				{
					ID:                    "id-3",
					Component:             "component-1",
					Impacts:               []string{"component-1", "serv-0", "serv-2"},
					ImpactedServices:      []string{"serv-0", "serv-2"},
					ImpactedServicesToAdd: []string{},
				},
			},
		},
		{
			name: "service is disabled",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-2", "id-3"},
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:                    "id-3",
					Component:             "component-1",
					Impacts:               []string{"component-1", "serv-0", "serv-1", "serv-2"},
					ImpactedServices:      []string{"serv-0", "serv-1", "serv-2"},
					ImpactedServicesToAdd: []string{"serv-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1"},
					ImpactedServices: []string{},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-2"},
				},
				{
					ID:                    "id-3",
					Component:             "component-1",
					Impacts:               []string{"component-1", "serv-0", "serv-2"},
					ImpactedServices:      []string{"serv-0", "serv-2"},
					ImpactedServicesToAdd: []string{},
				},
				{
					ID:      "serv-1",
					Depends: []string{},
				},
			},
		},
		{
			name: "New service",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:        "id-1",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
				{
					ID:        "id-2",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
				{
					ID:        "id-3",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:               "id-3",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-2", "id-3"},
				},
			},
		},
		{
			name: "New service(many impacted services)",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-2"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-2"},
				},
				{
					ID:               "id-3",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-2"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-1",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2", "serv-1"},
					ImpactedServices: []string{"serv-1", "serv-0", "serv-2"},
				},
				{
					ID:               "id-2",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2", "serv-1"},
					ImpactedServices: []string{"serv-1", "serv-0", "serv-2"},
				},
				{
					ID:               "id-3",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2", "serv-1"},
					ImpactedServices: []string{"serv-1", "serv-0", "serv-2"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-2", "id-3"},
				},
			},
		},
		{
			name: "Update service, where some depends are not valid anymore",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
					Depends: []string{"id-1", "id-2", "id-3"},
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:               "id-2",
					Component:        "component-2",
					Impacts:          []string{"component-2", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:               "id-3",
					Component:        "component-2",
					Impacts:          []string{"component-2", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
			},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-4",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-2",
					Component:        "component-2",
					Impacts:          []string{"component-2"},
					ImpactedServices: []string{},
				},
				{
					ID:               "id-3",
					Component:        "component-2",
					Impacts:          []string{"component-2"},
					ImpactedServices: []string{},
				},
				{
					ID:               "id-4",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-4"},
				},
			},
		},
		{
			name: "Update service, where some depends are not valid anymore(many impacted services)",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
					Depends: []string{"id-1", "id-2", "id-3"},
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:               "id-2",
					Component:        "component-2",
					Impacts:          []string{"component-2", "serv-0", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:               "id-3",
					Component:        "component-2",
					Impacts:          []string{"component-2", "serv-0", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-1", "serv-2"},
				},
			},
			secondFindCallEntities: []types.Entity{
				{
					ID:               "id-4",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-2"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:               "id-2",
					Component:        "component-2",
					Impacts:          []string{"component-2", "serv-0", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-2"},
				},
				{
					ID:               "id-3",
					Component:        "component-2",
					Impacts:          []string{"component-2", "serv-0", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-2"},
				},
				{
					ID:               "id-4",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-0", "serv-1", "serv-2"},
					ImpactedServices: []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-4"},
				},
			},
		},
		{
			name: "Update service where addedTo is exist, if not valid anymore do not add to RemoveFrom",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
					Depends: []string{"id-1", "id-2", "id-3"},
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:                    "id-2",
					Component:             "component-2",
					Impacts:               []string{"component-2", "serv-1"},
					ImpactedServices:      []string{"serv-1"},
					ImpactedServicesToAdd: []string{"serv-1"},
				},
				{
					ID:                    "id-3",
					Component:             "component-2",
					Impacts:               []string{"component-2", "serv-1"},
					ImpactedServices:      []string{"serv-1"},
					ImpactedServicesToAdd: []string{"serv-1"},
				},
			},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-4",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                    "id-2",
					Component:             "component-2",
					Impacts:               []string{"component-2"},
					ImpactedServices:      []string{},
					ImpactedServicesToAdd: []string{},
				},
				{
					ID:                    "id-3",
					Component:             "component-2",
					Impacts:               []string{"component-2"},
					ImpactedServices:      []string{},
					ImpactedServicesToAdd: []string{},
				},
				{
					ID:               "id-4",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-4"},
				},
			},
		},
		{
			name: "Update service where addedTo is exist, if not valid anymore do not add to RemoveFrom(many impacted services)",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
					Depends: []string{"id-1", "id-2", "id-3"},
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName: "serv-1",
			firstFindCallEntities: []types.Entity{
				{
					ID:                    "id-2",
					Component:             "component-2",
					Impacts:               []string{"component-2", "serv-0", "serv-1", "serv-2"},
					ImpactedServices:      []string{"serv-0", "serv-1", "serv-2"},
					ImpactedServicesToAdd: []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:                    "id-3",
					Component:             "component-2",
					Impacts:               []string{"component-2", "serv-0", "serv-1", "serv-2", "serv-3"},
					ImpactedServices:      []string{"serv-0", "serv-1", "serv-2", "serv-3"},
					ImpactedServicesToAdd: []string{"serv-0", "serv-1", "serv-2"},
				},
			},
			secondFindCallEntities: []types.Entity{
				{
					ID:        "id-4",
					Component: "component-1",
					Impacts:   []string{"component-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                    "id-2",
					Component:             "component-2",
					Impacts:               []string{"component-2", "serv-0", "serv-2"},
					ImpactedServices:      []string{"serv-0", "serv-2"},
					ImpactedServicesToAdd: []string{"serv-0", "serv-2"},
				},
				{
					ID:                    "id-3",
					Component:             "component-2",
					Impacts:               []string{"component-2", "serv-0", "serv-2", "serv-3"},
					ImpactedServices:      []string{"serv-0", "serv-2", "serv-3"},
					ImpactedServicesToAdd: []string{"serv-0", "serv-2"},
				},
				{
					ID:               "id-4",
					Component:        "component-1",
					Impacts:          []string{"component-1", "serv-1"},
					ImpactedServices: []string{"serv-1"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-4"},
				},
			},
		},
		{
			name: "Update service where RemoveTo is exist, if not valid anymore do not add to AddedTo",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
					Depends: []string{"id-1", "id-4"},
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-2",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName:           "serv-1",
			firstFindCallEntities: []types.Entity{},
			secondFindCallEntities: []types.Entity{
				{
					ID:                       "id-2",
					Component:                "component-2",
					Impacts:                  []string{"component-2"},
					ImpactedServices:         []string{},
					ImpactedServicesToRemove: []string{"serv-1"},
				},
				{
					ID:                       "id-3",
					Component:                "component-2",
					Impacts:                  []string{"component-2"},
					ImpactedServices:         []string{},
					ImpactedServicesToRemove: []string{"serv-1"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-2",
					Component:                "component-2",
					Impacts:                  []string{"component-2", "serv-1"},
					ImpactedServices:         []string{"serv-1"},
					ImpactedServicesToRemove: []string{},
				},
				{
					ID:                       "id-3",
					Component:                "component-2",
					Impacts:                  []string{"component-2", "serv-1"},
					ImpactedServices:         []string{"serv-1"},
					ImpactedServicesToRemove: []string{},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-2", "id-3", "id-4"},
				},
			},
		},
		{
			name: "Update service where RemoveTo is exist, if not valid anymore do not add to AddedTo(many services)",
			service: entityservice.EntityService{
				Entity: types.Entity{
					ID:      "serv-1",
					Enabled: true,
					Depends: []string{"id-1", "id-4"},
				},
				EntityPatterns: pattern.EntityPatternList{
					Patterns: []pattern.EntityPattern{
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-1",
										},
									},
								},
							},
						},
						{
							EntityFields: pattern.EntityFields{
								Component: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "component-2",
										},
									},
								},
							},
						},
					},
					Set:   true,
					Valid: true,
				},
			},
			serviceName:           "serv-1",
			firstFindCallEntities: []types.Entity{},
			secondFindCallEntities: []types.Entity{
				{
					ID:                       "id-2",
					Component:                "component-2",
					Impacts:                  []string{"component-2", "serv-3"},
					ImpactedServices:         []string{"serv-3"},
					ImpactedServicesToRemove: []string{"serv-0", "serv-1", "serv-2"},
				},
				{
					ID:                       "id-3",
					Component:                "component-2",
					Impacts:                  []string{"component-2", "serv-3", "serv-4"},
					ImpactedServices:         []string{"serv-3", "serv-4"},
					ImpactedServicesToRemove: []string{"serv-0", "serv-1", "serv-2"},
				},
			},
			expectedResult: []types.Entity{
				{
					ID:                       "id-2",
					Component:                "component-2",
					Impacts:                  []string{"component-2", "serv-1", "serv-3"},
					ImpactedServices:         []string{"serv-1", "serv-3"},
					ImpactedServicesToRemove: []string{"serv-0", "serv-2"},
				},
				{
					ID:                       "id-3",
					Component:                "component-2",
					Impacts:                  []string{"component-2", "serv-1", "serv-3", "serv-4"},
					ImpactedServices:         []string{"serv-1", "serv-3", "serv-4"},
					ImpactedServicesToRemove: []string{"serv-0", "serv-2"},
				},
				{
					ID:      "serv-1",
					Depends: []string{"id-1", "id-2", "id-3", "id-4"},
				},
			},
		},
	}

	call := 0
	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater)
	for _, dataset := range dataSets {
		t.Run(dataset.name, func(t *testing.T) {
			storage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dataset.service, nil)

			call = 0

			cursor := mock_mongo.NewMockCursor(ctrl)
			cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results interface{}) {
				if call == 0 {
					entities := results.(*[]types.Entity)
					*entities = dataset.firstFindCallEntities
				}

				if call == 1 {
					entities := results.(*[]types.Entity)
					*entities = dataset.secondFindCallEntities
				}

				call++
			}).Return(nil).AnyTimes()
			collection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

			_, result, err := manager.RecomputeService(ctx, dataset.serviceName)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(result, func(i, j int) bool {
				return result[i].ID < result[j].ID
			})

			for idx := range result {
				sort.Slice(result[idx].ImpactedServices, func(i, j int) bool {
					return result[idx].ImpactedServices[i] < result[idx].ImpactedServices[j]
				})

				sort.Slice(result[idx].ImpactedServicesToAdd, func(i, j int) bool {
					return result[idx].ImpactedServicesToAdd[i] < result[idx].ImpactedServicesToAdd[j]
				})

				sort.Slice(result[idx].ImpactedServicesToRemove, func(i, j int) bool {
					return result[idx].ImpactedServicesToRemove[i] < result[idx].ImpactedServicesToRemove[j]
				})

				sort.Slice(result[idx].Impacts, func(i, j int) bool {
					return result[idx].Impacts[i] < result[idx].Impacts[j]
				})

				sort.Slice(result[idx].Depends, func(i, j int) bool {
					return result[idx].Depends[i] < result[idx].Depends[j]
				})
			}

			for idx := range dataset.expectedResult {
				sort.Slice(dataset.expectedResult[idx].ImpactedServices, func(i, j int) bool {
					return dataset.expectedResult[idx].ImpactedServices[i] < dataset.expectedResult[idx].ImpactedServices[j]
				})

				sort.Slice(dataset.expectedResult[idx].ImpactedServicesToAdd, func(i, j int) bool {
					return dataset.expectedResult[idx].ImpactedServicesToAdd[i] < dataset.expectedResult[idx].ImpactedServicesToAdd[j]
				})

				sort.Slice(dataset.expectedResult[idx].ImpactedServicesToRemove, func(i, j int) bool {
					return dataset.expectedResult[idx].ImpactedServicesToRemove[i] < dataset.expectedResult[idx].ImpactedServicesToRemove[j]
				})

				sort.Slice(dataset.expectedResult[idx].Impacts, func(i, j int) bool {
					return dataset.expectedResult[idx].Impacts[i] < dataset.expectedResult[idx].Impacts[j]
				})

				sort.Slice(dataset.expectedResult[idx].Depends, func(i, j int) bool {
					return dataset.expectedResult[idx].Depends[i] < dataset.expectedResult[idx].Depends[j]
				})
			}

			sort.Slice(dataset.expectedResult, func(i, j int) bool {
				return dataset.expectedResult[i].ID < dataset.expectedResult[j].ID
			})

			if !reflect.DeepEqual(result, dataset.expectedResult) {
				t.Error("result is not expected")
			}
		})
	}
}

func BenchmarkRecomputeServicesRemoveAll(b *testing.B) {
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	var dependsIds []string
	var entities []types.Entity
	for i := 0; i < 1000; i++ {
		eID := fmt.Sprintf("id-%d", i)
		entities = append(entities, types.Entity{
			ID:               eID,
			Component:        "component-1",
			Impacts:          []string{"component-1", "serv-1"},
			ImpactedServices: []string{"serv-1"},
		})

		dependsIds = append(dependsIds, eID)
	}

	cursor := mock_mongo.NewMockCursor(ctrl)
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results interface{}) {
		ents := results.(*[]types.Entity)
		*ents = append(*ents, entities...)
	}).Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entityservice.EntityService{
		Entity: types.Entity{
			ID:      "serv-1",
			Depends: dependsIds,
		},
		EntityPatterns: pattern.EntityPatternList{
			Patterns: []pattern.EntityPattern{
				{
					EntityFields: pattern.EntityFields{
						Component: pattern.StringPattern{
							StringConditions: pattern.StringConditions{
								Equal: types.OptionalString{
									Set:   true,
									Value: "component-1",
								},
							},
						},
					},
				},
			},
			Set:   true,
			Valid: true,
		},
	}, nil).AnyTimes()

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater)
	for i := 0; i < b.N; i++ {
		_, _, _ = manager.RecomputeService(ctx, "serv-1")
	}
}

func BenchmarkRecomputeServicesAddAll(b *testing.B) {
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	var entities []types.Entity
	for i := 0; i < 1000; i++ {
		eID := fmt.Sprintf("id-%d", i)
		entities = append(entities, types.Entity{
			ID:        eID,
			Component: "component-1",
			Impacts:   []string{"component-1"},
		})
	}

	call := 0
	cursor := mock_mongo.NewMockCursor(ctrl)
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results interface{}) {
		if call == 1 {
			ents := results.(*[]types.Entity)
			*ents = append(*ents, entities...)
		}

		call++
	}).Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entityservice.EntityService{
		Entity: types.Entity{
			ID:      "serv-1",
			Enabled: true,
		},
		EntityPatterns: pattern.EntityPatternList{
			Patterns: []pattern.EntityPattern{
				{
					EntityFields: pattern.EntityFields{
						Component: pattern.StringPattern{
							StringConditions: pattern.StringConditions{
								Equal: types.OptionalString{
									Set:   true,
									Value: "component-1",
								},
							},
						},
					},
				},
			},
			Set:   true,
			Valid: true,
		},
	}, nil).AnyTimes()

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater)
	for i := 0; i < b.N; i++ {
		call = 0
		_, _, _ = manager.RecomputeService(ctx, "serv-1")
	}
}

func BenchmarkRecomputeServicesMixed(b *testing.B) {
	ctx := context.Background()
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	var dependsIds []string
	var entitiesToRemove []types.Entity
	for i := 0; i < 500; i++ {
		eID := fmt.Sprintf("id-%d", i)
		entitiesToRemove = append(entitiesToRemove, types.Entity{
			ID:               eID,
			Component:        "component-1",
			Impacts:          []string{"component-1", "serv-1"},
			ImpactedServices: []string{"serv-1"},
		})

		dependsIds = append(dependsIds, eID)
	}

	var entitiesToAdd []types.Entity
	for i := 500; i < 1000; i++ {
		eID := fmt.Sprintf("id-%d", i)
		entitiesToAdd = append(entitiesToAdd, types.Entity{
			ID:        eID,
			Component: "component-1",
			Impacts:   []string{"component-1"},
		})
	}

	call := 0
	cursor := mock_mongo.NewMockCursor(ctrl)
	cursor.EXPECT().All(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, results interface{}) {
		if call == 0 {
			ents := results.(*[]types.Entity)
			*ents = append(*ents, entitiesToRemove...)
		}

		if call == 1 {
			ents := results.(*[]types.Entity)
			*ents = append(*ents, entitiesToAdd...)
		}

		call++
	}).Return(nil).AnyTimes()

	collection := mock_mongo.NewMockDbCollection(ctrl)
	collection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(cursor, nil).AnyTimes()

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(collection).AnyTimes()

	adapter := entity.NewAdapter(dbClient)
	storage := mock_contextgraph.NewMockEntityServiceStorage(ctrl)
	storage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entityservice.EntityService{
		Entity: types.Entity{
			ID:      "serv-1",
			Enabled: true,
			Depends: dependsIds,
		},
		EntityPatterns: pattern.EntityPatternList{
			Patterns: []pattern.EntityPattern{
				{
					EntityFields: pattern.EntityFields{
						Component: pattern.StringPattern{
							StringConditions: pattern.StringConditions{
								Equal: types.OptionalString{
									Set:   true,
									Value: "component-1",
								},
							},
						},
					},
				},
			},
			Set:   true,
			Valid: true,
		},
	}, nil).AnyTimes()

	metaUpdater := mock_metrics.NewMockMetaUpdater(ctrl)

	manager := contextgraph.NewManager(adapter, dbClient, storage, metaUpdater)
	for i := 0; i < b.N; i++ {
		call = 0
		_, _, _ = manager.RecomputeService(ctx, "serv-1")
	}
}
