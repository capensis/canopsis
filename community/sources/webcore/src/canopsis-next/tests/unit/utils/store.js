import { cloneDeep, isFunction, omit } from 'lodash';
import Vuex from 'vuex';
import AxiosMockAdapter from 'axios-mock-adapter';
import Faker from 'faker';

import { CANOPSIS_EDITION } from '@/constants';

import request from '@/services/request';

import { DEFAULT_ENTITY_MODULE_TYPES } from '@/store/plugins/entities/create-crud-module';

/**
 * @typedef {Object} Module
 * @property {string} name
 * @property {Object.<string, Function | Mock>} [actions]
 * @property {Object} [state]
 * @property {Object.<string, any>} [getters]
 */

const convertMockedGettersToStore = (getters = {}) => Object
  .entries(getters)
  .reduce((acc, [getterName, getterOrValue]) => {
    acc[getterName] = isFunction(getterOrValue)
      ? getterOrValue
      : () => getterOrValue;

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
 *      fetchAppInfo: jest.fn()
 *    }
 *  })
 *
 * @param {Module[]} modules
 * @returns {Store}
 */
export const createMockedStoreModules = modules => new Vuex.Store({
  modules: modules.reduce((acc, module) => {
    acc[module.name] = createMockedStoreModule(module);

    return acc;
  }, {}),
});

/**
 * Wrapper for createMockedStoreModule, for mock getters.
 *
 * @param {string} name
 * @param {Object.<string, any>} getters
 * @returns {Store}
 */
export const createMockedStoreGetters = ({ name, ...getters }) => createMockedStoreModules([{ name, getters }]);

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
    jest.clearAllMocks();
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
    axiosMockAdapter
      .onGet(route)
      .reply(200, response);

    const commit = jest.fn();

    await actions.fetchList({ commit });

    expect(commit).toBeCalledWith(
      types.FETCH_LIST_COMPLETED,
      {
        data: entities,
        meta,
      },
    );
  });

  if (actions.fetchListWithPreviousParams) {
    it('Fetch list with previous params. Action: fetchListWithPreviousParams', async () => {
      const fetchingParams = {
        param: Faker.datatype.string(),
      };
      const state = { fetchingParams };
      const dispatch = jest.fn();

      await actions.fetchListWithPreviousParams({ dispatch, state });

      expect(dispatch).toBeCalledWith('fetchList', { params: fetchingParams });
    });
  }

  it('Fetch list with params. Action: fetchList', async () => {
    const params = { param: 1 };
    const commit = jest.fn();

    axiosMockAdapter
      .onGet(route, params)
      .reply(200, response);

    await actions.fetchList({ commit }, { params });

    expect(commit).toBeCalledWith(
      types.FETCH_LIST_COMPLETED,
      {
        data: entities,
        meta,
      },
    );
  });

  it('Fetch list with error. Action: fetchList', async () => {
    const error = { message: Faker.datatype.string() };

    axiosMockAdapter
      .onGet(route)
      .reply(404, error);

    const originalError = console.error;
    console.error = jest.fn();
    const dispatch = jest.fn().mockRejectedValue(error);
    const commit = jest.fn();

    try {
      await actions.fetchList({ dispatch, commit });
    } catch (err) {
      expect(err.message).toBe(error.message);

      expect(commit).toBeCalledWith(types.FETCH_LIST_FAILED);

      expect(console.error).toBeCalledWith(error);
    } finally {
      console.error = originalError;
    }
  });

  it('Create item. Action: create', async () => {
    const [entity] = entities;

    axiosMockAdapter
      .onPost(route, entity)
      .reply(200);

    await actions.create({}, { data: entity });

    const [entityPostRequest] = axiosMockAdapter.history.post;

    expect(JSON.parse(entityPostRequest.data)).toEqual(entity);
  });

  it('Update item. Action: update', async () => {
    const [entity] = entities;
    const id = Faker.datatype.number();

    axiosMockAdapter
      .onPut(`${route}/${id}`, entity)
      .reply(200);

    await actions.update({}, { id, data: entity });

    const [entityPutRequest] = axiosMockAdapter.history.put;

    expect(JSON.parse(entityPutRequest.data)).toEqual(entity);
  });

  it('Remove item. Action: remove', async () => {
    const id = Faker.datatype.number();

    axiosMockAdapter
      .onDelete(`${route}/${id}`)
      .reply(200);

    await actions.remove({}, { id });

    expect(axiosMockAdapter.history.delete).toHaveLength(1);
  });

  return {
    axiosMockAdapter,
  };
};

export const createAuthModule = () => {
  const currentUser = jest.fn()
    .mockReturnValue({});
  const currentUserPermissionsById = jest.fn()
    .mockReturnValue({});
  const login = jest.fn();
  const fetchCurrentUser = jest.fn();

  const authModule = {
    name: 'auth',
    getters: {
      currentUser,
      currentUserPermissionsById,
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
    login,
    fetchCurrentUser,
  };
};

export const createPbehaviorTypesModule = () => {
  const fetchPbehaviorTypesListWithoutStore = jest.fn().mockReturnValue({
    data: [],
  });

  const fieldPbehaviorTypes = jest.fn()
    .mockReturnValue([]);

  const fieldPbehaviorTypesPending = jest.fn()
    .mockReturnValue(false);

  const fetchFieldPbehaviorTypes = jest.fn();

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
  const fetchPbehaviorReasonsListWithoutStore = jest.fn().mockReturnValue({
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
  const fetchUserPreference = jest.fn();
  const getUserPreferenceByWidgetId = jest.fn()
    .mockReturnValue(() => ({ content: {} }));
  const updateUserPreference = jest.fn();

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
  const createWidget = jest.fn();
  const updateWidget = jest.fn();
  const createWidgetFilter = jest.fn();
  const updateWidgetFilter = jest.fn();
  const removeWidgetFilter = jest.fn();
  const updateGridPositions = jest.fn();

  afterEach(() => {
    createWidget.mockClear();
    updateWidget.mockClear();
    createWidgetFilter.mockClear();
    updateWidgetFilter.mockClear();
    removeWidgetFilter.mockClear();
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
  };
};

export const createViewModule = () => {
  const groups = jest.fn().mockReturnValue([]);
  const pending = jest.fn().mockReturnValue(false);
  const fetchGroupsList = jest.fn();
  const fetchGroupsListWithoutStore = jest.fn();
  const updateViewsPositions = jest.fn();
  const createGroup = jest.fn();
  const updateGroup = jest.fn();
  const removeGroup = jest.fn();
  const copyView = jest.fn();
  const createView = jest.fn();
  const updateView = jest.fn();
  const removeView = jest.fn();

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

export const createServiceModule = () => {
  const fetchEntityInfosKeysWithoutStore = jest.fn().mockReturnValue({
    data: [],
    meta: { total_count: 0 },
  });
  const fetchDependenciesWithoutStore = jest.fn()
    .mockResolvedValue({
      data: [],
      meta: {},
    });
  const fetchImpactsWithoutStore = jest.fn()
    .mockResolvedValue({
      data: [],
      meta: {},
    });
  const fetchServiceAlarmsWithoutStore = jest.fn();
  const fetchServicesList = jest.fn();
  const getServicesPendingByWidgetId = jest.fn().mockReturnValue(false);
  const getServicesListByWidgetId = jest.fn().mockReturnValue([]);
  const getServicesErrorByWidgetId = jest.fn();

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
  const fetchStateSettingWithoutStore = jest.fn();

  const entityModule = {
    name: 'entity',
    actions: {
      fetchStateSettingWithoutStore,
    },
  };

  return {
    fetchStateSettingWithoutStore,
    entityModule,
  };
};

export const createServiceEntityModule = () => {
  const serviceEntityModule = {
    name: 'service/entity',
    getters: {
    },
    actions: {
    },
  };

  return {
    serviceEntityModule,
  };
};

export const createQueryModule = () => {
  const getQueryById = jest.fn().mockReturnValue(() => ({}));
  const getQueryNonceById = jest.fn().mockReturnValue(() => 'nonce');
  const updateQuery = jest.fn();
  const removeQuery = jest.fn();

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
  const registerEditingOffHandler = jest.fn();
  const unregisterEditingOffHandler = jest.fn();
  const fetchActiveView = jest.fn();
  const toggleEditing = jest.fn();
  const editing = jest.fn().mockReturnValue(false);

  const activeViewModule = {
    name: 'activeView',
    getters: {
      editing,
    },
    actions: {
      registerEditingOffHandler,
      unregisterEditingOffHandler,
      toggleEditing,
      fetch: fetchActiveView,
    },
  };

  afterEach(() => {
    editing.mockClear();
    toggleEditing.mockClear();
    fetchActiveView.mockClear();
    registerEditingOffHandler.mockClear();
    unregisterEditingOffHandler.mockClear();
  });

  return {
    editing,
    registerEditingOffHandler,
    unregisterEditingOffHandler,
    fetchActiveView,
    toggleEditing,
    activeViewModule,
  };
};

export const createPbehaviorEntitiesModule = () => {
  const fetchPbehaviorEntitiesListWithoutStore = jest.fn().mockResolvedValue({
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
  const fetchPbehaviorsByEntityIdWithoutStore = jest.fn().mockResolvedValue([]);
  const removePbehavior = jest.fn();
  const createEntityPbehaviors = jest.fn();
  const removeEntityPbehaviors = jest.fn();

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

export const createPbehaviorTimespanModule = () => {
  const fetchTimespansListWithoutStore = jest.fn().mockResolvedValue([]);

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
  const fetchAlarmItem = jest.fn();
  const fetchOpenAlarmsListWithoutStore = jest.fn();
  const fetchAlarmItemWithoutStore = jest.fn().mockResolvedValue({});
  const bulkCreateAlarmAckEvent = jest.fn();
  const bulkCreateAlarmAckremoveEvent = jest.fn();
  const bulkCreateAlarmSnoozeEvent = jest.fn();
  const bulkCreateAlarmAssocticketEvent = jest.fn();
  const bulkCreateAlarmCommentEvent = jest.fn();
  const bulkCreateAlarmCancelEvent = jest.fn();
  const bulkCreateAlarmUnCancelEvent = jest.fn();
  const bulkCreateAlarmChangestateEvent = jest.fn();
  const addBookmarkToAlarm = jest.fn();
  const removeBookmarkFromAlarm = jest.fn();
  const fetchExecutionsWithoutStore = jest.fn().mockResolvedValue({});

  afterEach(() => {
    fetchAlarmItem.mockClear();
    fetchAlarmItemWithoutStore.mockClear();
    fetchOpenAlarmsListWithoutStore.mockClear();
    bulkCreateAlarmAckEvent.mockClear();
    bulkCreateAlarmAckremoveEvent.mockClear();
    bulkCreateAlarmSnoozeEvent.mockClear();
    bulkCreateAlarmAssocticketEvent.mockClear();
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
  const fetchAlarmDetailsWithoutStore = jest.fn().mockResolvedValue([]);

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
  const fetchWidgetTemplatesListWithoutStore = jest.fn()
    .mockReturnValue({
      meta: { total_count: 0 },
      data: [],
    }); // TODO: finish it in the future

  const createWidgetTemplate = jest.fn();
  const updateWidgetTemplate = jest.fn();
  const removeWidgetTemplate = jest.fn();

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
  const alarmTags = jest.fn().mockReturnValue([]);
  const alarmTagsPending = jest.fn().mockReturnValue(false);
  const alarmTagsMeta = jest.fn().mockReturnValue({});

  const fetchAlarmTagsList = jest.fn();
  const createAlarmTag = jest.fn();
  const updateAlarmTag = jest.fn();
  const removeAlarmTag = jest.fn();
  const bulkRemoveAlarmTags = jest.fn();

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
  const fetchItems = jest.fn();

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

export const createManualMetaAlarmModule = () => {
  const fetchManualMetaAlarmsListWithoutStore = jest.fn().mockResolvedValue([]);
  const createManualMetaAlarm = jest.fn().mockResolvedValue([]);
  const addAlarmsIntoManualMetaAlarm = jest.fn().mockResolvedValue([]);
  const removeAlarmsFromManualMetaAlarm = jest.fn().mockResolvedValue([]);

  afterEach(() => {
    fetchManualMetaAlarmsListWithoutStore.mockClear();
    createManualMetaAlarm.mockClear();
    addAlarmsIntoManualMetaAlarm.mockClear();
    removeAlarmsFromManualMetaAlarm.mockClear();
  });

  const manualMetaAlarmModule = {
    name: 'manualMetaAlarm',
    actions: {
      fetchListWithoutStore: fetchManualMetaAlarmsListWithoutStore,
      create: createManualMetaAlarm,
      addAlarms: addAlarmsIntoManualMetaAlarm,
      removeAlarms: removeAlarmsFromManualMetaAlarm,
    },
  };

  return {
    fetchManualMetaAlarmsListWithoutStore,
    createManualMetaAlarm,
    addAlarmsIntoManualMetaAlarm,
    removeAlarmsFromManualMetaAlarm,
    manualMetaAlarmModule,
  };
};

export const createMetaAlarmModule = () => {
  const removeAlarmsFromMetaAlarm = jest.fn().mockResolvedValue([]);

  afterEach(() => {
    removeAlarmsFromMetaAlarm.mockClear();
  });

  const metaAlarmModule = {
    name: 'metaAlarm',
    actions: {
      removeAlarms: removeAlarmsFromMetaAlarm,
    },
  };

  return {
    removeAlarmsFromMetaAlarm,
    metaAlarmModule,
  };
};

export const createDeclareTicketModule = () => {
  const bulkCreateDeclareTicketExecution = jest.fn().mockResolvedValue([]);
  const fetchAssignedDeclareTicketsWithoutStore = jest.fn().mockResolvedValue({
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
  const getVectorMetricsListByWidgetId = jest.fn().mockReturnValue(() => false);
  const getVectorMetricsPendingByWidgetId = jest.fn().mockReturnValue(() => []);
  const getVectorMetricsMetaByWidgetId = jest.fn().mockReturnValue(() => ({}));
  const fetchVectorMetricsList = jest.fn();

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
  const getAggregatedMetricsListByWidgetId = jest.fn().mockReturnValue(() => false);
  const getAggregatedMetricsPendingByWidgetId = jest.fn().mockReturnValue(() => []);
  const getAggregatedMetricsMetaByWidgetId = jest.fn().mockReturnValue(() => ({}));
  const fetchAggregatedMetricsList = jest.fn();
  const fetchAggregatedMetricsWithoutStore = jest.fn().mockResolvedValue({
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
  const fetchExternalMetricsList = jest.fn();
  const fetchAlarmsMetricsWithoutStore = jest.fn().mockResolvedValue({
    data: [],
  });
  const fetchEntityAlarmsMetricsWithoutStore = jest.fn().mockResolvedValue({
    data: [],
  });
  const fetchEntityAggregateMetricsWithoutStore = jest.fn().mockResolvedValue({
    data: [],
  });
  const externalMetrics = jest.fn().mockReturnValue([]);
  const pending = jest.fn().mockReturnValue(false);

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
  const checkPatternsEntitiesCount = jest.fn().mockResolvedValue({});
  const checkPatternsAlarmsCount = jest.fn().mockResolvedValue({});

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
  const description = jest.fn().mockReturnValue('');
  const maintenance = jest.fn().mockReturnValue(false);
  const footer = jest.fn().mockReturnValue('');
  const casConfig = jest.fn().mockReturnValue({});
  const samlConfig = jest.fn().mockReturnValue({});
  const oauthConfig = jest.fn().mockReturnValue({});
  const edition = jest.fn().mockReturnValue(CANOPSIS_EDITION.community);
  const version = jest.fn().mockReturnValue();
  const isBasicAuthEnabled = jest.fn().mockReturnValue(false);
  const isCASAuthEnabled = jest.fn().mockReturnValue(false);
  const isSAMLAuthEnabled = jest.fn().mockReturnValue(false);
  const isLDAPAuthEnabled = jest.fn().mockReturnValue(false);
  const isOauthAuthEnabled = jest.fn().mockReturnValue(false);

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
  };
};

export const createNavigationModule = () => {
  const isEditingMode = jest.fn().mockReturnValue(false);
  const toggleEditingMode = jest.fn();

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
  const hasMaximizedModal = jest.fn().mockReturnValue(false);

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
  const registerGetter = jest.fn();
  const unregisterGetter = jest.fn();

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
  const playlists = jest.fn().mockReturnValue([]);
  const fetchPlaylistsList = jest.fn();

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
