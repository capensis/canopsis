import AxiosMockAdapter from 'axios-mock-adapter';
import Faker from 'faker';
import { cloneDeep, isFunction, omit } from 'lodash';

import { CANOPSIS_EDITION, EXPORT_STATUSES, VIEW_SCREEN_MODES } from '@/constants';

import request from '@/services/request';

/**
 * @typedef {Object} Module
 * @property {string} name
 * @property {Object.<string, Function | Mock>} [actions]
 * @property {Object} [state]
 * @property {Object.<string, any>} [getters]
 */

const DEFAULT_ENTITY_MODULE_TYPES = {
  FETCH_LIST: 'FETCH_LIST',
  FETCH_LIST_COMPLETED: 'FETCH_LIST_COMPLETED',
  FETCH_LIST_FAILED: 'FETCH_LIST_FAILED',
};

const convertMockedGettersToStore = (getters = {}) =>
  Object.entries(getters).reduce((acc, [getterName, getterOrValue]) => {
    acc[getterName] = isFunction(getterOrValue) ? getterOrValue : () => getterOrValue;

    return acc;
  }, {});

/**
 * Create mocked store module.
 *
 * @param {Module} module
 * @returns {{}}
 */
export const createMockedStoreModule = module => ({
  ...omit(module, ['name']),

  namespaced: true,
  getters: convertMockedGettersToStore(module.getters),
});

const createMockedStorePlugin = (modules) => {
  const store = {
    state: {},
    getters: {},
    actionsMap: {},
    mutationsMap: {},

    commit: vi.fn((type, payload) => {
      const mutation = store.mutationsMap[type];

      return mutation?.(payload);
    }),

    dispatch: vi.fn((type, payload) => {
      const action = store.actionsMap[type];

      return action?.(payload);
    }),

    install(app) {
      app.config.globalProperties.$store = store;
      app.provide('store', store);
    },
  };

  modules.forEach((rawModule) => {
    const { name, actions = {}, mutations = {}, state = {} } = rawModule;
    const module = createMockedStoreModule(rawModule);

    store.state[name] = isFunction(state) ? state() : state;

    Object.entries(module.getters || {}).forEach(([getterName, getter]) => {
      Object.defineProperty(store.getters, `${name}/${getterName}`, {
        configurable: true,
        enumerable: true,
        get: () => getter(store.state[name], store.getters, store.state, store.getters),
      });
    });

    Object.entries(actions).forEach(([actionName, action]) => {
      store.actionsMap[`${name}/${actionName}`] = payload =>
        action(
          {
            commit: store.commit,
            dispatch: store.dispatch,
            getters: store.getters,
            rootGetters: store.getters,
            rootState: store.state,
            state: store.state[name],
          },
          payload,
        );
    });

    Object.entries(mutations).forEach(([mutationName, mutation]) => {
      store.mutationsMap[`${name}/${mutationName}`] = payload =>
        mutation(store.state[name], payload);
    });
  });

  return store;
};

/**
 * Create mocked whole store by special modules
 *
 * @example
 *  createMockedStoreModules({
 *    name: 'info',
 *    getters: {
 *      allowChangeSeverityToInfo: true,
 *      timezone: () => 'Timezone'
 *    },
 *    actions: {
 *      fetchAppInfo: vi.fn()
 *    }
 *  })
 *
 * @param {Module[]} modules
 * @returns {Store}
 */
export const createMockedStoreModules = modules => createMockedStorePlugin(modules);

/**
 * Wrapper for createMockedStoreModule, for mock getters.
 *
 * @param {string} name
 * @param {Object.<string, any>} getters
 * @returns {Store}
 */
export const createMockedStoreGetters = ({ name, ...getters }) =>
  createMockedStoreModules([{ name, getters }]);

/**
 *
 * @param {string} route
 * @param {Object} module
 * @param {Object} schema
 * @param {string} entityType
 * @param {Object[]} entities
 * @param {string[]} entityIds
 * @param {Object} types
 * @return {{ axiosMockAdapter: MockAdapter }}
 */
export const testsEntityModule = ({
  route,
  module,
  entities,
  types = DEFAULT_ENTITY_MODULE_TYPES,
}) => {
  const { actions, state: initialState, mutations, getters } = module;

  const axiosMockAdapter = new AxiosMockAdapter(request);
  const meta = {
    total_count: entities.length,
  };
  const response = {
    data: entities,
    meta: {
      total_count: entities.length,
    },
  };

  beforeEach(() => {
    axiosMockAdapter.reset();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('Mutate state after commit FETCH_LIST', () => {
    const state = cloneDeep(initialState);

    const fetchList = mutations[types.FETCH_LIST];

    fetchList(state);

    expect(state).toEqual({ ...state, pending: true });
  });

  it('Mutate state after commit FETCH_LIST with params', () => {
    const state = cloneDeep(initialState);
    const params = {
      param: Faker.datatype.string(),
    };

    const fetchList = mutations[types.FETCH_LIST];

    fetchList(state, { params });

    if (actions.fetchListWithPreviousParams) {
      expect(state).toEqual({ ...state, fetchingParams: params, pending: true });
    } else {
      expect(state).toEqual({ ...state, pending: true });
    }
  });

  it('Mutate state after commit FETCH_LIST_COMPLETED', () => {
    const state = cloneDeep(initialState);

    const fetchListCompleted = mutations[types.FETCH_LIST_COMPLETED];

    fetchListCompleted(state, {
      data: entities,
      meta,
    });

    expect(state).toEqual({ ...state, items: entities, meta });
  });

  it('Mutate state after commit FETCH_LIST_FAILED', () => {
    const state = cloneDeep(initialState);

    const fetchListFailed = mutations[types.FETCH_LIST_FAILED];

    fetchListFailed(state);

    expect(state).toEqual({ ...state, pending: false });
  });

  it('Get items. Getter: items', () => {
    const state = {
      ...initialState,
      items: entities,
    };

    expect(getters.items(state)).toEqual(entities);
  });

  it('Get pending. Getter: pending', () => {
    const pending = Faker.datatype.boolean();
    const state = { pending };

    const data = getters.pending(state);

    expect(data).toEqual(pending);
  });

  it('Get meta. Getter: meta', () => {
    const state = { meta };

    const data = getters.meta(state);

    expect(data).toEqual(meta);
  });

  it('Fetch list. Action: fetchList', async () => {
    axiosMockAdapter.onGet(route).reply(200, response);

    const commit = vi.fn();

    await actions.fetchList({ commit });

    expect(commit).toHaveBeenCalledWith(types.FETCH_LIST_COMPLETED, {
      data: entities,
      meta,
    });
  });

  if (actions.fetchListWithPreviousParams) {
    it('Fetch list with previous params. Action: fetchListWithPreviousParams', async () => {
      const fetchingParams = {
        param: Faker.datatype.string(),
      };
      const state = { fetchingParams };
      const dispatch = vi.fn();

      await actions.fetchListWithPreviousParams({ dispatch, state });

      expect(dispatch).toHaveBeenCalledWith('fetchList', { params: fetchingParams });
    });
  }

  it('Fetch list with params. Action: fetchList', async () => {
    const params = { param: 1 };
    const commit = vi.fn();

    axiosMockAdapter.onGet(route, { params }).reply(200, response);

    await actions.fetchList({ commit }, { params });

    expect(commit).toHaveBeenCalledWith(types.FETCH_LIST_COMPLETED, {
      data: entities,
      meta,
    });
  });

  it('Fetch list with error. Action: fetchList', async () => {
    const error = { message: Faker.datatype.string() };

    axiosMockAdapter.onGet(route).reply(404, error);

    const originalError = console.error;
    console.error = vi.fn();
    const dispatch = vi.fn().mockRejectedValue(error);
    const commit = vi.fn();

    try {
      await actions.fetchList({ dispatch, commit });
    } catch (err) {
      expect(err.message).toBe(error.message);

      expect(commit).toHaveBeenCalledWith(types.FETCH_LIST_FAILED);

      expect(console.error).toHaveBeenCalledWith(error);
    } finally {
      console.error = originalError;
    }
  });

  it('Create item. Action: create', async () => {
    const [entity] = entities;

    axiosMockAdapter.onPost(route, entity).reply(200);

    await actions.create({}, { data: entity });

    const [entityPostRequest] = axiosMockAdapter.history.post;

    expect(JSON.parse(entityPostRequest.data)).toEqual(entity);
  });

  it('Update item. Action: update', async () => {
    const [entity] = entities;
    const id = Faker.datatype.number();

    axiosMockAdapter.onPut(`${route}/${id}`, entity).reply(200);

    await actions.update({}, { id, data: entity });

    const [entityPutRequest] = axiosMockAdapter.history.put;

    expect(JSON.parse(entityPutRequest.data)).toEqual(entity);
  });

  it('Remove item. Action: remove', async () => {
    const id = Faker.datatype.number();

    axiosMockAdapter.onDelete(`${route}/${id}`).reply(200);

    await actions.remove({}, { id });

    expect(axiosMockAdapter.history.delete).toHaveLength(1);
  });

  return {
    axiosMockAdapter,
  };
};

export const createAuthModule = () => {
  const currentUser = vi.fn().mockReturnValue({});
  const currentUserPermissionsById = vi.fn().mockReturnValue({});
  const currentUserViewPermissionsByViewId = vi.fn().mockReturnValue({});
  const login = vi.fn();
  const fetchCurrentUser = vi.fn();

  const authModule = {
    name: 'auth',
    getters: {
      currentUser,
      currentUserPermissionsById,
      currentUserViewPermissionsByViewId,
    },
    actions: {
      login,
      fetchCurrentUser,
    },
  };

  afterEach(() => {
    currentUser.mockClear();
    currentUserPermissionsById.mockClear();
    login.mockClear();
    fetchCurrentUser.mockClear();
  });

  return {
    authModule,
    currentUser,
    currentUserPermissionsById,
    currentUserViewPermissionsByViewId,
    login,
    fetchCurrentUser,
  };
};

/**
 * Mocked `user` store module (entities user) for modals that call `updateCurrentUserTours` / `updateCurrentUser`.
 */
export const createUserModule = () => {
  const updateCurrentUserTours = vi.fn().mockResolvedValue(undefined);
  const updateCurrentUser = vi.fn().mockResolvedValue(undefined);

  const userModule = {
    name: 'user',
    actions: {
      updateCurrentUserTours,
      updateCurrentUser,
    },
  };

  afterEach(() => {
    updateCurrentUserTours.mockClear();
    updateCurrentUser.mockClear();
  });

  return {
    userModule,
    updateCurrentUserTours,
    updateCurrentUser,
  };
};

/**
 * Mocked `llm` store module for AI chat hooks (`fetchListWithoutStore`, `bulkLinkLlmHistory`).
 */
export const createLlmModule = () => {
  const fetchListWithoutStore = vi.fn().mockResolvedValue({ data: [] });
  const bulkLinkLlmHistory = vi.fn().mockResolvedValue({});

  const llmModule = {
    name: 'llm',
    actions: {
      fetchListWithoutStore,
      bulkLinkLlmHistory,
    },
  };

  afterEach(() => {
    fetchListWithoutStore.mockClear();
    bulkLinkLlmHistory.mockClear();
  });

  return {
    llmModule,
    fetchListWithoutStore,
    bulkLinkLlmHistory,
  };
};

export const createPbehaviorTypesModule = () => {
  const fetchPbehaviorTypesListWithoutStore = vi.fn().mockReturnValue({
    data: [],
  });

  const fieldPbehaviorTypes = vi.fn().mockReturnValue([]);

  const fieldPbehaviorTypesPending = vi.fn().mockReturnValue(false);

  const fetchFieldPbehaviorTypes = vi.fn();

  const pbehaviorTypesModule = {
    name: 'pbehaviorTypes',
    getters: {
      fieldItems: fieldPbehaviorTypes,
      fieldPending: fieldPbehaviorTypesPending,
    },
    actions: {
      fetchListWithoutStore: fetchPbehaviorTypesListWithoutStore,
      fetchFieldList: fetchFieldPbehaviorTypes,
    },
  };

  return {
    pbehaviorTypesModule,
    fieldPbehaviorTypes,
    fieldPbehaviorTypesPending,
    fetchPbehaviorTypesListWithoutStore,
    fetchFieldPbehaviorTypes,
  };
};

export const createPbehaviorReasonModule = () => {
  const fetchPbehaviorReasonsListWithoutStore = vi.fn().mockReturnValue({
    meta: {},
    data: [],
  });

  const pbehaviorReasonModule = {
    name: 'pbehaviorReasons',
    actions: {
      fetchListWithoutStore: fetchPbehaviorReasonsListWithoutStore,
    },
  };

  afterEach(() => {
    fetchPbehaviorReasonsListWithoutStore.mockClear();
  });

  return {
    pbehaviorReasonModule,
    fetchPbehaviorReasonsListWithoutStore,
  };
};

export const createUserPreferenceModule = () => {
  const fetchUserPreference = vi.fn();
  const getUserPreferenceByWidgetId = vi.fn().mockReturnValue(() => ({ content: {} }));
  const updateUserPreference = vi.fn();

  const userPreferenceModule = {
    name: 'userPreference',
    actions: {
      fetchItem: fetchUserPreference,
      update: updateUserPreference,
    },
    getters: {
      getItemByWidgetId: getUserPreferenceByWidgetId,
    },
  };

  return {
    fetchUserPreference,
    updateUserPreference,
    userPreferenceModule,
    getUserPreferenceByWidgetId,
  };
};

export const createWidgetModule = () => {
  const createWidget = vi.fn();
  const updateWidget = vi.fn();
  const createWidgetFilter = vi.fn();
  const updateWidgetFilter = vi.fn();
  const removeWidgetFilter = vi.fn();
  const updateGridPositions = vi.fn();
  const updateWidgetFiltersPositions = vi.fn();

  afterEach(() => {
    createWidget.mockClear();
    updateWidget.mockClear();
    createWidgetFilter.mockClear();
    updateWidgetFilter.mockClear();
    removeWidgetFilter.mockClear();
    updateWidgetFiltersPositions.mockClear();
  });

  const widgetModule = {
    name: 'view/widget',
    actions: {
      create: createWidget,
      update: updateWidget,
      createWidgetFilter,
      updateWidgetFilter,
      removeWidgetFilter,
      updateGridPositions,
      updateWidgetFiltersPositions,
    },
  };

  return {
    widgetModule,
    createWidget,
    updateWidget,
    createWidgetFilter,
    updateWidgetFilter,
    removeWidgetFilter,
    updateGridPositions,
    updateWidgetFiltersPositions,
  };
};

export const createViewModule = () => {
  const groups = vi.fn().mockReturnValue([]);
  const pending = vi.fn().mockReturnValue(false);
  const fetchGroupsList = vi.fn();
  const fetchGroupsListWithoutStore = vi.fn();
  const updateViewsPositions = vi.fn();
  const createGroup = vi.fn();
  const updateGroup = vi.fn();
  const removeGroup = vi.fn();
  const copyView = vi.fn();
  const createView = vi.fn();
  const updateView = vi.fn();
  const removeView = vi.fn();

  afterEach(() => {
    fetchGroupsList.mockClear();
    fetchGroupsListWithoutStore.mockClear();
    updateViewsPositions.mockClear();
    createGroup.mockClear();
    updateGroup.mockClear();
    removeGroup.mockClear();
    createView.mockClear();
    copyView.mockClear();
    updateView.mockClear();
    removeView.mockClear();
    groups.mockClear();
    pending.mockClear();
  });

  const viewModule = {
    name: 'view',
    getters: {
      items: groups,
      pending,
    },
    actions: {
      fetchList: fetchGroupsList,
      fetchListWithoutStore: fetchGroupsListWithoutStore,
      updateViewPositions: updateViewsPositions,
      create: createGroup,
      update: updateGroup,
      remove: removeGroup,
      copyView,
      createView,
      updateView,
      removeView,
    },
  };

  return {
    viewModule,

    groups,
    pending,

    fetchGroupsList,
    fetchGroupsListWithoutStore,
    updateViewsPositions,
    createGroup,
    updateGroup,
    removeGroup,
    createView,
    copyView,
    updateView,
    removeView,
  };
};

export const createCopyVarsModule = () => {
  const fetchEventFiltersVarsWithoutStore = vi.fn().mockResolvedValue([]);
  const fetchDynamicInfosVarsWithoutStore = vi.fn().mockResolvedValue([]);

  const copyVarsModule = {
    name: 'copy/vars',
    actions: {
      fetchEventFiltersVarsWithoutStore,
      fetchDynamicInfosVarsWithoutStore,
    },
  };

  afterEach(() => {
    fetchEventFiltersVarsWithoutStore.mockClear();
    fetchDynamicInfosVarsWithoutStore.mockClear();
  });

  return {
    copyVarsModule,
    fetchEventFiltersVarsWithoutStore,
    fetchDynamicInfosVarsWithoutStore,
  };
};

export const createTemplateVarsModule = () => {
  const fetchEntityServicesVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchEventFiltersVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchLinkRulesVarsWithoutStore = vi.fn().mockResolvedValue({ alarm: {}, entity: {} });
  const fetchScenariosVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchWidgetsVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchDeclareTicketRulesVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchDynamicInfosVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchInstructionsVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchJobsVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchMetaAlarmRulesVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchWebhookTokenRulesVarsWithoutStore = vi.fn().mockResolvedValue({});
  const fetchList = vi.fn();
  const items = vi.fn().mockReturnValue([]);
  const pending = vi.fn().mockReturnValue(false);

  const templateVarsModule = {
    name: 'template/vars',
    getters: {
      items,
      pending,
    },
    actions: {
      fetchList,
      fetchEntityServicesVarsWithoutStore,
      fetchEventFiltersVarsWithoutStore,
      fetchLinkRulesVarsWithoutStore,
      fetchScenariosVarsWithoutStore,
      fetchWidgetsVarsWithoutStore,
      fetchDeclareTicketRulesVarsWithoutStore,
      fetchDynamicInfosVarsWithoutStore,
      fetchInstructionsVarsWithoutStore,
      fetchJobsVarsWithoutStore,
      fetchMetaAlarmRulesVarsWithoutStore,
      fetchWebhookTokenRulesVarsWithoutStore,
    },
  };

  afterEach(() => {
    fetchEntityServicesVarsWithoutStore.mockClear();
    fetchEventFiltersVarsWithoutStore.mockClear();
    fetchLinkRulesVarsWithoutStore.mockClear();
    fetchScenariosVarsWithoutStore.mockClear();
    fetchWidgetsVarsWithoutStore.mockClear();
    fetchDeclareTicketRulesVarsWithoutStore.mockClear();
    fetchDynamicInfosVarsWithoutStore.mockClear();
    fetchInstructionsVarsWithoutStore.mockClear();
    fetchJobsVarsWithoutStore.mockClear();
    fetchMetaAlarmRulesVarsWithoutStore.mockClear();
    fetchWebhookTokenRulesVarsWithoutStore.mockClear();
    fetchList.mockClear();
    items.mockClear();
    pending.mockClear();
  });

  return {
    templateVarsModule,
    fetchEntityServicesVarsWithoutStore,
    fetchEventFiltersVarsWithoutStore,
    fetchLinkRulesVarsWithoutStore,
    fetchScenariosVarsWithoutStore,
    fetchWidgetsVarsWithoutStore,
    fetchDeclareTicketRulesVarsWithoutStore,
    fetchDynamicInfosVarsWithoutStore,
    fetchInstructionsVarsWithoutStore,
    fetchJobsVarsWithoutStore,
    fetchMetaAlarmRulesVarsWithoutStore,
    fetchWebhookTokenRulesVarsWithoutStore,
    fetchList,
    items,
    pending,
  };
};

export const createServiceModule = () => {
  const fetchEntityInfosKeysWithoutStore = vi.fn().mockReturnValue({
    data: [],
    meta: { total_count: 0 },
  });
  const fetchDependenciesWithoutStore = vi.fn().mockResolvedValue({
    data: [],
    meta: {},
  });
  const fetchImpactsWithoutStore = vi.fn().mockResolvedValue({
    data: [],
    meta: {},
  });
  const fetchServiceAlarmsWithoutStore = vi.fn();
  const fetchServicesList = vi.fn();
  const getServicesPendingByWidgetId = vi.fn().mockReturnValue(false);
  const getServicesListByWidgetId = vi.fn().mockReturnValue([]);
  const getServicesErrorByWidgetId = vi.fn();

  const serviceModule = {
    name: 'service',
    getters: {
      getPendingByWidgetId: () => getServicesPendingByWidgetId,
      getListByWidgetId: () => getServicesListByWidgetId,
      getErrorByWidgetId: () => getServicesErrorByWidgetId,
    },
    actions: {
      fetchInfosKeysWithoutStore: fetchEntityInfosKeysWithoutStore,
      fetchAlarmsWithoutStore: fetchServiceAlarmsWithoutStore,
      fetchDependenciesWithoutStore,
      fetchImpactsWithoutStore,
      fetchList: fetchServicesList,
    },
  };

  return {
    getServicesPendingByWidgetId,
    getServicesListByWidgetId,
    getServicesErrorByWidgetId,
    fetchEntityInfosKeysWithoutStore,
    fetchServicesList,
    fetchServiceAlarmsWithoutStore,
    fetchDependenciesWithoutStore,
    fetchImpactsWithoutStore,
    serviceModule,
  };
};

export const createEntityModule = () => {
  const fetchStateSettingWithoutStore = vi.fn();
  const fetchEntityInfosLogsListWithoutStore = vi.fn();
  const pause = vi.fn();
  const resume = vi.fn();
  const bulkPause = vi.fn();
  const bulkResume = vi.fn();

  afterEach(() => {
    fetchStateSettingWithoutStore.mockClear();
    fetchEntityInfosLogsListWithoutStore.mockClear();
    pause.mockClear();
    resume.mockClear();
    bulkPause.mockClear();
    bulkResume.mockClear();
  });

  const entityModule = {
    name: 'entity',
    actions: {
      fetchStateSettingWithoutStore,
      fetchEntityInfosLogsListWithoutStore,
      pause,
      resume,
      bulkPause,
      bulkResume,
    },
  };

  return {
    fetchStateSettingWithoutStore,
    fetchEntityInfosLogsListWithoutStore,
    pause,
    resume,
    bulkPause,
    bulkResume,
    entityModule,
  };
};

export const createServiceEntityModule = () => {
  const serviceEntityModule = {
    name: 'service/entity',
    getters: {},
    actions: {},
  };

  return {
    serviceEntityModule,
  };
};

export const createQueryModule = () => {
  const getQueryById = vi.fn().mockReturnValue(() => ({}));
  const getQueryNonceById = vi.fn().mockReturnValue(() => 'nonce');
  const updateQuery = vi.fn();
  const removeQuery = vi.fn();

  const queryModule = {
    name: 'query',
    getters: {
      getQueryById,
      getQueryNonceById,
    },
    actions: {
      update: updateQuery,
      remove: removeQuery,
    },
  };

  return {
    getQueryById,
    getQueryNonceById,
    updateQuery,
    removeQuery,
    queryModule,
  };
};

export const createActiveViewModule = () => {
  const registerEditingOffHandler = vi.fn();
  const unregisterEditingOffHandler = vi.fn();
  const fetchActiveView = vi.fn();
  const toggleEditing = vi.fn();
  const setScreenMode = vi.fn();
  const editing = vi.fn().mockReturnValue(false);
  const item = vi.fn().mockReturnValue({});
  const screenMode = vi.fn().mockReturnValue(VIEW_SCREEN_MODES.default);
  const isKioskScreenMode = vi.fn().mockReturnValue(false);
  const resumePeriodicRefresh = vi.fn();
  const pausePeriodicRefresh = vi.fn();
  const periodicRefreshPaused = vi.fn().mockReturnValue(false);

  const activeViewModule = {
    name: 'activeView',
    getters: {
      editing,
      item,
      screenMode,
      isKioskScreenMode,
      periodicRefreshPaused,
    },
    actions: {
      registerEditingOffHandler,
      unregisterEditingOffHandler,
      toggleEditing,
      setScreenMode,
      fetch: fetchActiveView,
      resumePeriodicRefresh,
      pausePeriodicRefresh,
    },
  };

  afterEach(() => {
    editing.mockClear();
    toggleEditing.mockClear();
    fetchActiveView.mockClear();
    registerEditingOffHandler.mockClear();
    unregisterEditingOffHandler.mockClear();
    resumePeriodicRefresh.mockClear();
    pausePeriodicRefresh.mockClear();
    periodicRefreshPaused.mockClear();
  });

  return {
    item,
    editing,
    registerEditingOffHandler,
    unregisterEditingOffHandler,
    fetchActiveView,
    toggleEditing,
    resumePeriodicRefresh,
    pausePeriodicRefresh,
    periodicRefreshPaused,
    activeViewModule,
  };
};

export const createPbehaviorEntitiesModule = () => {
  const fetchPbehaviorEntitiesListWithoutStore = vi.fn().mockResolvedValue({
    data: [],
    meta: { total_count: 0 },
  });

  const pbehaviorEntitiesModule = {
    name: 'pbehavior/entities',
    actions: {
      fetchListWithoutStore: fetchPbehaviorEntitiesListWithoutStore,
    },
  };

  afterEach(() => {
    fetchPbehaviorEntitiesListWithoutStore.mockClear();
  });

  return {
    fetchPbehaviorEntitiesListWithoutStore,
    pbehaviorEntitiesModule,
  };
};

export const createPbehaviorModule = () => {
  const fetchPbehaviorsByEntityIdWithoutStore = vi.fn().mockResolvedValue([]);
  const removePbehavior = vi.fn();
  const createEntityPbehaviors = vi.fn();
  const removeEntityPbehaviors = vi.fn();

  const pbehaviorModule = {
    name: 'pbehavior',
    actions: {
      fetchListByEntityIdWithoutStore: fetchPbehaviorsByEntityIdWithoutStore,
      removeWithoutStore: removePbehavior,
      bulkCreateEntityPbehaviors: createEntityPbehaviors,
      bulkRemoveEntityPbehaviors: removeEntityPbehaviors,
    },
  };

  afterEach(() => {
    removePbehavior.mockClear();
    fetchPbehaviorsByEntityIdWithoutStore.mockClear();
    createEntityPbehaviors.mockClear();
    removeEntityPbehaviors.mockClear();
  });

  return {
    removePbehavior,
    fetchPbehaviorsByEntityIdWithoutStore,
    createEntityPbehaviors,
    removeEntityPbehaviors,
    pbehaviorModule,
  };
};

export const createPbehaviorCommentModule = () => {
  const createPbehaviorComment = vi.fn();
  const updatePbehaviorComment = vi.fn();
  const removePbehaviorComment = vi.fn();

  const pbehaviorCommentModule = {
    name: 'pbehavior/comment',
    actions: {
      create: createPbehaviorComment,
      update: updatePbehaviorComment,
      remove: removePbehaviorComment,
    },
  };

  afterEach(() => {
    createPbehaviorComment.mockClear();
    updatePbehaviorComment.mockClear();
    removePbehaviorComment.mockClear();
  });

  return {
    createPbehaviorComment,
    updatePbehaviorComment,
    removePbehaviorComment,
    pbehaviorCommentModule,
  };
};

export const createPbehaviorTimespanModule = () => {
  const fetchTimespansListWithoutStore = vi.fn().mockResolvedValue([]);

  afterEach(() => {
    fetchTimespansListWithoutStore.mockClear();
  });

  const pbehaviorTimespanModule = {
    name: 'pbehaviorTimespan',
    actions: {
      fetchListWithoutStore: fetchTimespansListWithoutStore,
    },
  };

  return {
    fetchTimespansListWithoutStore,
    pbehaviorTimespanModule,
  };
};

export const createAlarmModule = () => {
  const fetchAlarmItem = vi.fn();
  const fetchOpenAlarmsListWithoutStore = vi.fn();
  const fetchAlarmItemWithoutStore = vi.fn().mockResolvedValue({});
  const bulkCreateAlarmAckEvent = vi.fn();
  const bulkCreateAlarmAckremoveEvent = vi.fn();
  const bulkCreateAlarmSnoozeEvent = vi.fn();
  const bulkCreateAlarmAssocticketEvent = vi.fn();
  const bulkCreateAlarmTicketremoveEvent = vi.fn();
  const bulkCreateAlarmCommentEvent = vi.fn();
  const bulkCreateAlarmCancelEvent = vi.fn();
  const bulkCreateAlarmUnCancelEvent = vi.fn();
  const bulkCreateAlarmChangestateEvent = vi.fn();
  const addBookmarkToAlarm = vi.fn();
  const removeBookmarkFromAlarm = vi.fn();
  const fetchExecutionsWithoutStore = vi.fn().mockResolvedValue({});

  afterEach(() => {
    fetchAlarmItem.mockClear();
    fetchAlarmItemWithoutStore.mockClear();
    fetchOpenAlarmsListWithoutStore.mockClear();
    bulkCreateAlarmAckEvent.mockClear();
    bulkCreateAlarmAckremoveEvent.mockClear();
    bulkCreateAlarmSnoozeEvent.mockClear();
    bulkCreateAlarmAssocticketEvent.mockClear();
    bulkCreateAlarmTicketremoveEvent.mockClear();
    bulkCreateAlarmCommentEvent.mockClear();
    bulkCreateAlarmCancelEvent.mockClear();
    bulkCreateAlarmUnCancelEvent.mockClear();
    bulkCreateAlarmChangestateEvent.mockClear();
    addBookmarkToAlarm.mockClear();
    removeBookmarkFromAlarm.mockClear();
    fetchExecutionsWithoutStore.mockClear();
  });

  const alarmModule = {
    name: 'alarm',
    actions: {
      fetchItem: fetchAlarmItem,
      fetchItemWithoutStore: fetchAlarmItemWithoutStore,
      fetchOpenAlarmsListWithoutStore,
      bulkCreateAlarmAckEvent,
      bulkCreateAlarmAckremoveEvent,
      bulkCreateAlarmSnoozeEvent,
      bulkCreateAlarmAssocticketEvent,
      bulkCreateAlarmTicketremoveEvent,
      bulkCreateAlarmCommentEvent,
      bulkCreateAlarmCancelEvent,
      bulkCreateAlarmUnCancelEvent,
      bulkCreateAlarmChangestateEvent,
      addBookmarkToAlarm,
      removeBookmarkFromAlarm,
      fetchExecutionsWithoutStore,
    },
  };

  return {
    fetchAlarmItem,
    fetchAlarmItemWithoutStore,
    fetchOpenAlarmsListWithoutStore,
    bulkCreateAlarmAckEvent,
    bulkCreateAlarmAckremoveEvent,
    bulkCreateAlarmSnoozeEvent,
    bulkCreateAlarmAssocticketEvent,
    bulkCreateAlarmTicketremoveEvent,
    bulkCreateAlarmCommentEvent,
    bulkCreateAlarmCancelEvent,
    bulkCreateAlarmUnCancelEvent,
    bulkCreateAlarmChangestateEvent,
    addBookmarkToAlarm,
    removeBookmarkFromAlarm,
    fetchExecutionsWithoutStore,
    alarmModule,
  };
};

export const createAlarmDetailsModule = () => {
  const fetchAlarmDetailsWithoutStore = vi.fn().mockResolvedValue([]);

  afterEach(() => {
    fetchAlarmDetailsWithoutStore.mockClear();
  });

  const alarmDetailsModule = {
    name: 'alarm/details',
    actions: {
      fetchListWithoutStore: fetchAlarmDetailsWithoutStore,
    },
  };

  return {
    fetchAlarmDetailsWithoutStore,
    alarmDetailsModule,
  };
};

export const createWidgetTemplateModule = () => {
  const fetchWidgetTemplatesListWithoutStore = vi.fn().mockReturnValue({
    meta: { total_count: 0 },
    data: [],
  }); // TODO: finish it in the future

  const createWidgetTemplate = vi.fn();
  const updateWidgetTemplate = vi.fn();
  const removeWidgetTemplate = vi.fn();

  const widgetTemplateModule = {
    name: 'widgetTemplate',
    actions: {
      fetchListWithoutStore: fetchWidgetTemplatesListWithoutStore,
      create: createWidgetTemplate,
      update: updateWidgetTemplate,
      remove: removeWidgetTemplate,
    },
  };

  return {
    fetchWidgetTemplatesListWithoutStore,
    createWidgetTemplate,
    updateWidgetTemplate,
    removeWidgetTemplate,
    widgetTemplateModule,
  };
};

export const createAlarmTagModule = () => {
  const alarmTags = vi.fn().mockReturnValue([]);
  const alarmTagsPending = vi.fn().mockReturnValue(false);
  const alarmTagsMeta = vi.fn().mockReturnValue({});

  const fetchAlarmTagsList = vi.fn();
  const createAlarmTag = vi.fn();
  const updateAlarmTag = vi.fn();
  const removeAlarmTag = vi.fn();
  const bulkRemoveAlarmTags = vi.fn();

  afterEach(() => {
    fetchAlarmTagsList.mockClear();
    createAlarmTag.mockClear();
    updateAlarmTag.mockClear();
    removeAlarmTag.mockClear();
    bulkRemoveAlarmTags.mockClear();
  });

  const alarmTagModule = {
    name: 'alarmTag',
    getters: {
      items: alarmTags,
      pending: alarmTagsPending,
      meta: alarmTagsMeta,
    },
    actions: {
      fetchList: fetchAlarmTagsList,
      create: createAlarmTag,
      update: updateAlarmTag,
      remove: removeAlarmTag,
      bulkRemove: bulkRemoveAlarmTags,
    },
  };

  return {
    alarmTagModule,

    alarmTags,
    alarmTagsPending,
    alarmTagsMeta,

    fetchAlarmTagsList,
    createAlarmTag,
    updateAlarmTag,
    removeAlarmTag,
    bulkRemoveAlarmTags,
  };
};

export const createInfosModule = () => {
  const fetchItems = vi.fn();

  const infosModule = {
    name: 'infos',
    getters: {
      alarmInfos: () => [], // TODO: finish it in the future
      alarmInfosRules: () => [],
      entityInfos: () => [],
      pending: () => [],
    },
    actions: {
      fetch: fetchItems,
    },
  };

  return {
    fetchItems,
    infosModule,
  };
};

export const createMetaAlarmModule = () => {
  const removeAlarmsFromMetaAlarm = vi.fn().mockResolvedValue([]);
  const fetchMetaAlarmsListWithoutStore = vi.fn().mockResolvedValue([]);
  const createMetaAlarm = vi.fn().mockResolvedValue([]);
  const addAlarmsIntoMetaAlarm = vi.fn().mockResolvedValue([]);

  afterEach(() => {
    fetchMetaAlarmsListWithoutStore.mockClear();
    createMetaAlarm.mockClear();
    addAlarmsIntoMetaAlarm.mockClear();
    removeAlarmsFromMetaAlarm.mockClear();
  });

  const metaAlarmModule = {
    name: 'metaAlarm',
    actions: {
      removeAlarms: removeAlarmsFromMetaAlarm,
      fetchListWithoutStore: fetchMetaAlarmsListWithoutStore,
      create: createMetaAlarm,
      addAlarms: addAlarmsIntoMetaAlarm,
    },
  };

  return {
    fetchMetaAlarmsListWithoutStore,
    createMetaAlarm,
    addAlarmsIntoMetaAlarm,
    removeAlarmsFromMetaAlarm,
    metaAlarmModule,
  };
};

export const createMetaAlarmRuleModule = () => {
  const checkPatternsAlarmsCount = vi.fn().mockResolvedValue({});

  afterEach(() => {
    checkPatternsAlarmsCount.mockClear();
  });

  const metaAlarmRuleModule = {
    name: 'metaAlarmRule',
    actions: {
      checkPatternsAlarmsCount,
    },
  };

  return {
    checkPatternsAlarmsCount,
    metaAlarmRuleModule,
  };
};

export const createDeclareTicketModule = () => {
  const bulkCreateDeclareTicketExecution = vi.fn().mockResolvedValue([]);
  const fetchAssignedDeclareTicketsWithoutStore = vi.fn().mockResolvedValue({
    by_rules: {},
    by_alarms: {},
  });

  afterEach(() => {
    bulkCreateDeclareTicketExecution.mockClear();
    fetchAssignedDeclareTicketsWithoutStore.mockClear();
  });

  const declareTicketRuleModule = {
    name: 'declareTicketRule',
    actions: {
      bulkCreateDeclareTicketExecution,
      fetchAssignedTicketsWithoutStore: fetchAssignedDeclareTicketsWithoutStore,
    },
  };

  return {
    declareTicketRuleModule,
    bulkCreateDeclareTicketExecution,
    fetchAssignedDeclareTicketsWithoutStore,
  };
};

export const createVectorMetricsModule = () => {
  const getVectorMetricsListByWidgetId = vi.fn().mockReturnValue(() => false);
  const getVectorMetricsPendingByWidgetId = vi.fn().mockReturnValue(() => []);
  const getVectorMetricsMetaByWidgetId = vi.fn().mockReturnValue(() => ({}));
  const fetchVectorMetricsList = vi.fn();

  afterEach(() => {
    getVectorMetricsListByWidgetId.mockClear();
    getVectorMetricsPendingByWidgetId.mockClear();
    getVectorMetricsMetaByWidgetId.mockClear();
    fetchVectorMetricsList.mockClear();
  });

  const vectorMetricsModule = {
    name: 'vectorMetrics',
    getters: {
      getListByWidgetId: getVectorMetricsListByWidgetId,
      getPendingByWidgetId: getVectorMetricsPendingByWidgetId,
      getMetaByWidgetId: getVectorMetricsMetaByWidgetId,
    },
    actions: {
      fetchList: fetchVectorMetricsList,
    },
  };

  return {
    vectorMetricsModule,
    getVectorMetricsListByWidgetId,
    getVectorMetricsPendingByWidgetId,
    getVectorMetricsMetaByWidgetId,
    fetchVectorMetricsList,
  };
};

export const createAggregatedMetricsModule = () => {
  const getAggregatedMetricsListByWidgetId = vi.fn().mockReturnValue(() => false);
  const getAggregatedMetricsPendingByWidgetId = vi.fn().mockReturnValue(() => []);
  const getAggregatedMetricsMetaByWidgetId = vi.fn().mockReturnValue(() => ({}));
  const fetchAggregatedMetricsList = vi.fn();
  const fetchAggregatedMetricsWithoutStore = vi.fn().mockResolvedValue({
    data: [],
  });

  afterEach(() => {
    getAggregatedMetricsListByWidgetId.mockClear();
    getAggregatedMetricsPendingByWidgetId.mockClear();
    getAggregatedMetricsMetaByWidgetId.mockClear();
    fetchAggregatedMetricsList.mockClear();
    fetchAggregatedMetricsWithoutStore.mockClear();
  });

  const aggregatedMetricsModule = {
    name: 'aggregatedMetrics',
    getters: {
      getListByWidgetId: getAggregatedMetricsListByWidgetId,
      getPendingByWidgetId: getAggregatedMetricsPendingByWidgetId,
      getMetaByWidgetId: getAggregatedMetricsMetaByWidgetId,
    },
    actions: {
      fetchList: fetchAggregatedMetricsList,
      fetchListWithoutStore: fetchAggregatedMetricsWithoutStore,
    },
  };

  return {
    aggregatedMetricsModule,
    getAggregatedMetricsListByWidgetId,
    getAggregatedMetricsPendingByWidgetId,
    getAggregatedMetricsMetaByWidgetId,
    fetchAggregatedMetricsList,
    fetchAggregatedMetricsWithoutStore,
  };
};

export const createMetricsModule = () => {
  const fetchExternalMetricsList = vi.fn();
  const fetchAlarmsMetricsWithoutStore = vi.fn().mockResolvedValue({
    data: [],
  });
  const fetchEntityAlarmsMetricsWithoutStore = vi.fn().mockResolvedValue({
    data: [],
  });
  const fetchEntityAggregateMetricsWithoutStore = vi.fn().mockResolvedValue({
    data: [],
  });
  const externalMetrics = vi.fn().mockReturnValue([]);
  const pending = vi.fn().mockReturnValue(false);

  afterEach(() => {
    fetchExternalMetricsList.mockClear();
    fetchAlarmsMetricsWithoutStore.mockClear();
    externalMetrics.mockClear();
    pending.mockClear();
    fetchEntityAlarmsMetricsWithoutStore.mockClear();
    fetchEntityAggregateMetricsWithoutStore.mockClear();
  });

  const metricsModule = {
    name: 'metrics',
    getters: {
      externalMetrics,
      pending,
    },
    actions: {
      fetchExternalMetricsList,
      fetchAlarmsMetricsWithoutStore,
      fetchEntityAlarmsMetricsWithoutStore,
      fetchEntityAggregateMetricsWithoutStore,
    },
  };

  return {
    metricsModule,
    externalMetrics,
    fetchExternalMetricsList,
    fetchAlarmsMetricsWithoutStore,
    fetchEntityAlarmsMetricsWithoutStore,
    fetchEntityAggregateMetricsWithoutStore,
  };
};

export const createPatternModule = () => {
  const checkPatternsEntitiesCount = vi.fn().mockResolvedValue({});
  const checkPatternsAlarmsCount = vi.fn().mockResolvedValue({});

  afterEach(() => {
    checkPatternsEntitiesCount.mockClear();
    checkPatternsAlarmsCount.mockClear();
  });

  const patternModule = {
    name: 'pattern',
    actions: {
      checkPatternsEntitiesCount,
      checkPatternsAlarmsCount,
    },
  };

  return {
    patternModule,
    checkPatternsEntitiesCount,
    checkPatternsAlarmsCount,
  };
};

export const createInfoModule = () => {
  const description = vi.fn().mockReturnValue('');
  const maintenance = vi.fn().mockReturnValue(false);
  const footer = vi.fn().mockReturnValue('');
  const casConfig = vi.fn().mockReturnValue({});
  const samlConfig = vi.fn().mockReturnValue({});
  const oauthConfig = vi.fn().mockReturnValue({});
  const edition = vi.fn().mockReturnValue(CANOPSIS_EDITION.community);
  const version = vi.fn().mockReturnValue();
  const isBasicAuthEnabled = vi.fn().mockReturnValue(false);
  const isCASAuthEnabled = vi.fn().mockReturnValue(false);
  const isSAMLAuthEnabled = vi.fn().mockReturnValue(false);
  const isLDAPAuthEnabled = vi.fn().mockReturnValue(false);
  const isOauthAuthEnabled = vi.fn().mockReturnValue(false);
  const userTimezones = vi.fn().mockReturnValue([]);
  const shownUserTimezone = vi.fn().mockReturnValue(false);

  afterEach(() => {
    maintenance.mockClear();
    description.mockClear();
    footer.mockClear();
    casConfig.mockClear();
    samlConfig.mockClear();
    edition.mockClear();
    version.mockClear();
    isBasicAuthEnabled.mockClear();
    isCASAuthEnabled.mockClear();
    isSAMLAuthEnabled.mockClear();
    isLDAPAuthEnabled.mockClear();
    isOauthAuthEnabled.mockClear();
    userTimezones.mockClear();
    shownUserTimezone.mockClear();
  });

  const infoModule = {
    name: 'info',
    getters: {
      maintenance,
      description,
      footer,
      casConfig,
      samlConfig,
      oauthConfig,
      edition,
      version,
      isBasicAuthEnabled,
      isCASAuthEnabled,
      isSAMLAuthEnabled,
      isLDAPAuthEnabled,
      isOauthAuthEnabled,
      userTimezones,
      shownUserTimezone,
    },
  };

  return {
    infoModule,
    maintenance,
    description,
    footer,
    casConfig,
    samlConfig,
    oauthConfig,
    edition,
    version,
    isBasicAuthEnabled,
    isCASAuthEnabled,
    isSAMLAuthEnabled,
    isLDAPAuthEnabled,
    isOauthAuthEnabled,
    userTimezones,
    shownUserTimezone,
  };
};

export const createNavigationModule = () => {
  const isEditingMode = vi.fn().mockReturnValue(false);
  const toggleEditingMode = vi.fn();

  afterEach(() => {
    isEditingMode.mockClear();
    toggleEditingMode.mockClear();
  });

  const navigationModule = {
    name: 'navigation',
    getters: {
      isEditingMode,
    },
    actions: {
      toggleEditingMode,
    },
  };

  return {
    navigationModule,
    isEditingMode,
    toggleEditingMode,
  };
};

export const createModalsModule = () => {
  const hasMaximizedModal = vi.fn().mockReturnValue(false);

  afterEach(() => {
    hasMaximizedModal.mockClear();
  });

  const modalsModule = {
    name: 'modals',
    getters: {
      hasMaximizedModal,
    },
  };

  return {
    modalsModule,
    hasMaximizedModal,
  };
};

export const createEntitiesModule = () => {
  const registerGetter = vi.fn();
  const unregisterGetter = vi.fn();

  afterEach(() => {
    registerGetter.mockClear();
    unregisterGetter.mockClear();
  });

  const entitiesModule = {
    name: 'entities',
    actions: {
      registerGetter,
      unregisterGetter,
    },
  };

  return {
    entitiesModule,
    registerGetter,
    unregisterGetter,
  };
};

export const createPlaylistModule = () => {
  const playlists = vi.fn().mockReturnValue([]);
  const fetchPlaylistsList = vi.fn();

  afterEach(() => {
    playlists.mockClear();
    fetchPlaylistsList.mockClear();
  });

  const playlistModule = {
    name: 'playlist',
    getters: {
      items: playlists,
    },
    actions: {
      fetchList: fetchPlaylistsList,
    },
  };

  return {
    playlistModule,
    playlists,
    fetchPlaylistsList,
  };
};

export const createAvailabilityModule = () => {
  const exportAvailabilityData = {
    _id: 'export-availability-id',
    status: EXPORT_STATUSES.completed,
  };

  const getAvailabilityPendingByWidgetId = vi.fn().mockReturnValue(false);
  const getAvailabilityListByWidgetId = vi.fn().mockReturnValue([]);
  const getAvailabilityMetaByWidgetId = vi.fn().mockReturnValue({
    total_count: 0,
  });

  const fetchAvailabilityWithoutStore = vi.fn().mockResolvedValue({
    availability: {
      uptime: 10,
      downtime: 20,
      inactive_time: 5,
    },
  });
  const fetchAvailabilityHistoryWithoutStore = vi.fn().mockResolvedValue({
    data: [],
  });
  const fetchAvailabilityList = vi.fn().mockResolvedValue({
    data: [],
  });
  const fetchAvailabilityHistoryExport = vi.fn().mockReturnValue(exportAvailabilityData);

  afterEach(() => {
    fetchAvailabilityList.mockClear();
    fetchAvailabilityWithoutStore.mockClear();
    fetchAvailabilityHistoryWithoutStore.mockClear();
    fetchAvailabilityHistoryExport.mockClear();
  });

  const availabilityModule = {
    name: 'availability',
    getters: {
      getPendingByWidgetId: () => getAvailabilityPendingByWidgetId,
      getListByWidgetId: () => getAvailabilityListByWidgetId,
      getMetaByWidgetId: () => getAvailabilityMetaByWidgetId,
    },
    actions: {
      fetchList: fetchAvailabilityList,
      fetchAvailabilityWithoutStore,
      fetchAvailabilityHistoryWithoutStore,
      fetchAvailabilityHistoryExport,
    },
  };

  return {
    availabilityModule,
    fetchAvailabilityList,
    exportAvailabilityData,
    fetchAvailabilityWithoutStore,
    fetchAvailabilityHistoryWithoutStore,
    fetchAvailabilityHistoryExport,
  };
};

export const createRemediationInstructionModule = () => {
  const fetchRemediationInstructionsListWithoutStore = vi.fn().mockResolvedValue({
    data: [],
    meta: {
      total_count: 0,
    },
  });

  afterEach(() => {
    fetchRemediationInstructionsListWithoutStore.mockClear();
  });

  const remediationInstructionModule = {
    name: 'remediationInstruction',
    actions: {
      fetchListWithoutStore: fetchRemediationInstructionsListWithoutStore,
    },
  };

  return {
    remediationInstructionModule,
    fetchRemediationInstructionsListWithoutStore,
  };
};

export const createEntityInfoPropertyModule = () => {
  const entityInfoProperties = vi.fn().mockReturnValue([]);
  const entityInfoPropertiesWithAlias = vi.fn().mockReturnValue([]);
  const entityInfoPropertiesWithoutAlias = vi.fn().mockReturnValue([]);
  const entityInfoPropertyMeta = vi.fn().mockReturnValue({
    total_count: 0,
  });
  const entityInfoPropertyPending = vi.fn().mockReturnValue(false);

  const fetchEntityInfoPropertiesList = vi.fn().mockResolvedValue({
    data: [],
    meta: {
      total_count: 0,
    },
  });
  const fetchEntityInfoPropertiesListWithoutStore = vi.fn().mockResolvedValue({
    data: [],
    meta: {
      total_count: 0,
    },
  });

  const createEntityInfoProperty = vi.fn();
  const updateEntityInfoProperty = vi.fn();
  const removeEntityInfoProperty = vi.fn();

  afterEach(() => {
    entityInfoProperties.mockClear();
    entityInfoPropertiesWithAlias.mockClear();
    entityInfoPropertiesWithoutAlias.mockClear();
    entityInfoPropertyMeta.mockClear();
    entityInfoPropertyPending.mockClear();
    fetchEntityInfoPropertiesList.mockClear();
    fetchEntityInfoPropertiesListWithoutStore.mockClear();
    createEntityInfoProperty.mockClear();
    updateEntityInfoProperty.mockClear();
    removeEntityInfoProperty.mockClear();
  });

  const entityInfoPropertyModule = {
    name: 'entityInfoProperty',
    getters: {
      items: entityInfoProperties,
      itemsWithAlias: entityInfoPropertiesWithAlias,
      itemsWithoutAlias: entityInfoPropertiesWithoutAlias,
      meta: entityInfoPropertyMeta,
      pending: entityInfoPropertyPending,
    },
    actions: {
      fetchList: fetchEntityInfoPropertiesList,
      create: createEntityInfoProperty,
      update: updateEntityInfoProperty,
      remove: removeEntityInfoProperty,
      fetchListWithoutStore: fetchEntityInfoPropertiesListWithoutStore,
    },
  };

  return {
    entityInfoPropertyModule,
    entityInfoProperties,
    entityInfoPropertiesWithAlias,
    entityInfoPropertiesWithoutAlias,
    entityInfoPropertyMeta,
    entityInfoPropertyPending,
    fetchEntityInfoPropertiesList,
    fetchEntityInfoPropertiesListWithoutStore,
    createEntityInfoProperty,
    updateEntityInfoProperty,
    removeEntityInfoProperty,
  };
};

export const createPbehaviorPatternsModule = () => {
  const runAlarmFiltering = vi.fn().mockResolvedValue();
  const checkPatternsPbehaviorsCount = vi.fn().mockResolvedValue({
    data: 0,
  });

  const pbehaviorPatternsModule = {
    name: 'pbehaviorPatterns',
    actions: {
      runAlarmFiltering,
      checkPatternsPbehaviorsCount,
    },
  };

  afterEach(() => {
    runAlarmFiltering.mockClear();
    checkPatternsPbehaviorsCount.mockClear();
  });

  return {
    runAlarmFiltering,
    checkPatternsPbehaviorsCount,
    pbehaviorPatternsModule,
  };
};

export const createPatternsFieldsModule = () => {
  const fetchFlappingRulePatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchIdleRulePatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchLinkRulePatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchResolveRulePatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchPbehaviorPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchAlarmTagPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchWidgetFilterPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchServicePatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchStateSettingPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchEventFilterPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchScenarioPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchMetaalarmrulePatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchDeclareTicketRulePatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchInstructionPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchKpiFilterPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchDynamicInfosPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });
  const fetchEventRecordPatternFields = vi.fn().mockResolvedValue({
    entity_pattern: [],
    alarm_pattern: [],
    event_pattern: [],
  });

  const patternsFieldsModule = {
    name: 'patternsFields',
    actions: {
      fetchFlappingRulePatternFields,
      fetchIdleRulePatternFields,
      fetchLinkRulePatternFields,
      fetchResolveRulePatternFields,
      fetchPbehaviorPatternFields,
      fetchAlarmTagPatternFields,
      fetchWidgetFilterPatternFields,
      fetchServicePatternFields,
      fetchStateSettingPatternFields,
      fetchEventFilterPatternFields,
      fetchScenarioPatternFields,
      fetchMetaalarmrulePatternFields,
      fetchDeclareTicketRulePatternFields,
      fetchInstructionPatternFields,
      fetchKpiFilterPatternFields,
      fetchDynamicInfosPatternFields,
      fetchEventRecordPatternFields,
    },
  };

  afterEach(() => {
    fetchFlappingRulePatternFields.mockClear();
    fetchIdleRulePatternFields.mockClear();
    fetchLinkRulePatternFields.mockClear();
    fetchResolveRulePatternFields.mockClear();
    fetchPbehaviorPatternFields.mockClear();
    fetchAlarmTagPatternFields.mockClear();
    fetchWidgetFilterPatternFields.mockClear();
    fetchServicePatternFields.mockClear();
    fetchStateSettingPatternFields.mockClear();
    fetchEventFilterPatternFields.mockClear();
    fetchScenarioPatternFields.mockClear();
    fetchMetaalarmrulePatternFields.mockClear();
    fetchDeclareTicketRulePatternFields.mockClear();
    fetchInstructionPatternFields.mockClear();
    fetchKpiFilterPatternFields.mockClear();
    fetchDynamicInfosPatternFields.mockClear();
    fetchEventRecordPatternFields.mockClear();
  });

  return {
    patternsFieldsModule,
    fetchFlappingRulePatternFields,
    fetchIdleRulePatternFields,
    fetchLinkRulePatternFields,
    fetchResolveRulePatternFields,
    fetchPbehaviorPatternFields,
    fetchAlarmTagPatternFields,
    fetchWidgetFilterPatternFields,
    fetchServicePatternFields,
    fetchStateSettingPatternFields,
    fetchEventFilterPatternFields,
    fetchScenarioPatternFields,
    fetchMetaalarmrulePatternFields,
    fetchDeclareTicketRulePatternFields,
    fetchInstructionPatternFields,
    fetchKpiFilterPatternFields,
    fetchDynamicInfosPatternFields,
    fetchEventRecordPatternFields,
  };
};

export const createCommentTemplateModule = () => {
  const createCommentTemplate = vi.fn();
  const updateCommentTemplate = vi.fn();
  const removeCommentTemplate = vi.fn();
  const fetchCommentTemplatesListWithoutStore = vi.fn().mockResolvedValue({
    data: [],
    meta: {
      total_count: 0,
    },
  });

  const commentTemplateModule = {
    name: 'commentTemplate',
    actions: {
      fetchListWithoutStore: fetchCommentTemplatesListWithoutStore,
      create: createCommentTemplate,
      update: updateCommentTemplate,
      remove: removeCommentTemplate,
    },
  };

  return {
    createCommentTemplate,
    updateCommentTemplate,
    removeCommentTemplate,
    fetchCommentTemplatesListWithoutStore,
    commentTemplateModule,
  };
};

export const createPatternEntitiesOptimizeModule = () => {
  const optimize = vi.fn().mockResolvedValue({});
  const fetchOptimizeStatus = vi.fn().mockResolvedValue({});
  const update = vi.fn().mockResolvedValue({});
  const remove = vi.fn().mockResolvedValue({});

  const patternEntitiesOptimizeModule = {
    name: 'pattern/entitiesOptimize',
    actions: {
      optimize,
      fetchOptimizeStatus,
      update,
      remove,
    },
  };

  return {
    patternEntitiesOptimizeModule,
    optimize,
    fetchOptimizeStatus,
    update,
    remove,
  };
};
