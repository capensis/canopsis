import { cloneDeep, isFunction, omit } from 'lodash';
import Vuex from 'vuex';
import AxiosMockAdapter from 'axios-mock-adapter';
import Faker from 'faker';

import request from '@/services/request';
import { DEFAULT_ENTITY_MODULE_TYPES } from '@/store/plugins/entities/create-entity-module';

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
  schema,
  entityType,
  entities,
  entityIds,
  types = DEFAULT_ENTITY_MODULE_TYPES,
}) => {
  const { actions, state: initialState, mutations, getters } = module;

  const axiosMockAdapter = new AxiosMockAdapter(request);
  const normalizedData = {
    result: entityIds,
  };
  const responseData = {
    data: entities,
    meta: {
      total_count: Faker.datatype.number(),
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

    const allIds = Faker.datatype.array();

    fetchListCompleted(state, { allIds });

    expect(state).toEqual({ ...state, allIds });
  });

  it('Mutate state after commit FETCH_LIST_FAILED', () => {
    const state = cloneDeep(initialState);

    const fetchListFailed = mutations[types.FETCH_LIST_FAILED];

    fetchListFailed(state);

    expect(state).toEqual({ ...state, pending: false });
  });

  it('Get item by id. Getter: getItemById', () => {
    const item = {
      param: Faker.datatype.string(),
    };
    const getItem = jest.fn(() => item);
    const rootGetters = {
      'entities/getItem': getItem,
    };
    const state = {};

    const id = Faker.datatype.string();

    const data = getters.getItemById(state, getters, {}, rootGetters)(id);

    expect(data).toEqual(item);
    expect(getItem).toHaveBeenCalledWith(entityType, id);
  });

  it('Get items. Getter: items', () => {
    const getList = jest.fn(() => entities);
    const rootGetters = {
      'entities/getList': getList,
    };
    const state = {
      ...initialState,
      allIds: entityIds,
    };

    const data = getters.items(state, getters, {}, rootGetters);

    expect(data).toEqual(entities);
    expect(getList).toHaveBeenCalledWith(entityType, entityIds);
  });

  it('Get pending. Getter: pending', () => {
    const pending = Faker.datatype.boolean();
    const state = { pending };

    const data = getters.pending(state);

    expect(data).toEqual(pending);
  });

  if (getters.meta) {
    it('Get meta. Getter: meta', () => {
      const meta = {
        total_count: Faker.datatype.number(),
      };
      const state = { meta };

      const data = getters.meta(state);

      expect(data).toEqual(meta);
    });
  }

  it('Fetch list. Action: fetchList', async () => {
    const dispatch = jest.fn().mockReturnValue({
      normalizedData,
      data: responseData,
    });
    const commit = jest.fn();

    await actions.fetchList({ dispatch, commit });

    expect(dispatch).toBeCalledWith(
      'entities/fetch',
      {
        route,
        dataPreparer: expect.any(Function),
        params: undefined,
        schema: [schema],
      },
      { root: true },
    );

    expect(commit).toBeCalledWith(
      types.FETCH_LIST_COMPLETED,
      {
        ...responseData,
        allIds: entityIds,
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
    const params = {};
    const dispatch = jest.fn().mockReturnValue({
      normalizedData,
      data: responseData,
    });
    const commit = jest.fn();

    await actions.fetchList({ dispatch, commit }, { params });

    expect(dispatch).toBeCalledWith(
      'entities/fetch',
      {
        route,
        params,
        dataPreparer: expect.any(Function),
        schema: [schema],
      },
      { root: true },
    );

    expect(commit).toBeCalledWith(
      types.FETCH_LIST_COMPLETED,
      {
        ...responseData,
        allIds: entityIds,
      },
    );
  });

  it('Fetch list with error. Action: fetchList', async () => {
    const originalError = console.error;
    console.error = jest.fn();
    const error = new Error(Faker.datatype.string());
    const dispatch = jest.fn().mockRejectedValue(error);
    const commit = jest.fn();

    try {
      await actions.fetchList({ dispatch, commit });
    } catch (err) {
      expect(err).toBe(error);

      expect(commit).toBeCalledWith(types.FETCH_LIST_FAILED);

      expect(console.error).toBeCalledWith(error);

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
  const currentUserPermissionsById = jest.fn()
    .mockReturnValue({});
  const authModule = {
    name: 'auth',
    getters: {
      currentUserPermissionsById,
    },
  };

  afterEach(() => {
    currentUserPermissionsById.mockClear();
  });

  return {
    authModule,
    currentUserPermissionsById,
  };
};

export const createPbehaviorTypesModule = () => {
  const fetchPbehaviorTypesListWithoutStore = jest.fn().mockReturnValue({
    data: [],
  });

  const pbehaviorTypesModule = {
    name: 'pbehaviorTypes',
    actions: {
      fetchListWithoutStore: fetchPbehaviorTypesListWithoutStore,
    },
  };

  return {
    pbehaviorTypesModule,
    fetchPbehaviorTypesListWithoutStore,
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
    .mockReturnValue({ content: {} });
  const updateUserPreference = jest.fn();

  const userPreferenceModule = {
    name: 'userPreference',
    actions: {
      fetchItem: fetchUserPreference,
      update: updateUserPreference,
    },
    getters: {
      getItemByWidgetId: () => getUserPreferenceByWidgetId,
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
  const copyWidget = jest.fn();
  const createWidgetFilter = jest.fn();
  const updateWidgetFilter = jest.fn();
  const removeWidgetFilter = jest.fn();

  afterEach(() => {
    createWidget.mockClear();
    updateWidget.mockClear();
    copyWidget.mockClear();
    createWidgetFilter.mockClear();
    updateWidgetFilter.mockClear();
    removeWidgetFilter.mockClear();
  });

  const widgetModule = {
    name: 'view/widget',
    actions: {
      create: createWidget,
      update: updateWidget,
      copy: copyWidget,
      createWidgetFilter,
      updateWidgetFilter,
      removeWidgetFilter,
    },
  };

  return {
    widgetModule,
    createWidget,
    updateWidget,
    copyWidget,
    createWidgetFilter,
    updateWidgetFilter,
    removeWidgetFilter,
  };
};

export const createServiceModule = () => {
  const fetchEntityInfosKeysWithoutStore = jest.fn().mockReturnValue({
    data: [],
    meta: { total_count: 0 },
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
    serviceModule,
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

  const queryModule = {
    name: 'query',
    getters: {
      getQueryById,
      getQueryNonceById,
    },
    actions: {
      update: updateQuery,
    },
  };

  return {
    getQueryById,
    getQueryNonceById,
    updateQuery,
    queryModule,
  };
};

export const createActiveViewModule = () => {
  const registerEditingOffHandler = jest.fn();
  const unregisterEditingOffHandler = jest.fn();
  const fetchActiveView = jest.fn();

  const activeViewModule = {
    name: 'activeView',
    actions: {
      registerEditingOffHandler,
      unregisterEditingOffHandler,
      fetch: fetchActiveView,
    },
  };

  return {
    registerEditingOffHandler,
    unregisterEditingOffHandler,
    fetchActiveView,
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

  const pbehaviorModule = {
    name: 'pbehavior',
    actions: {
      fetchListByEntityIdWithoutStore: fetchPbehaviorsByEntityIdWithoutStore,
      removeWithoutStore: removePbehavior,
    },
  };

  afterEach(() => {
    removePbehavior.mockClear();
    fetchPbehaviorsByEntityIdWithoutStore.mockClear();
  });

  return {
    removePbehavior,
    fetchPbehaviorsByEntityIdWithoutStore,
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
  const fetchAlarmItemWithoutStore = jest.fn().mockResolvedValue({});

  afterEach(() => {
    fetchAlarmItem.mockClear();
    fetchAlarmItemWithoutStore.mockClear();
  });

  const alarmModule = {
    name: 'alarm',
    actions: {
      fetchItem: fetchAlarmItem,
      fetchItemWithoutStore: fetchAlarmItemWithoutStore,
    },
  };

  return {
    fetchAlarmItem,
    fetchAlarmItemWithoutStore,
    alarmModule,
  };
};

export const createEventModule = () => {
  const createEvent = jest.fn();

  afterEach(() => {
    createEvent.mockClear();
  });

  const eventModule = {
    name: 'event',
    actions: {
      create: createEvent,
    },
  };

  return {
    eventModule,
    createEvent,
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
  const removeAlarmsIntoManualMetaAlarm = jest.fn().mockResolvedValue([]);

  afterEach(() => {
    fetchManualMetaAlarmsListWithoutStore.mockClear();
    createManualMetaAlarm.mockClear();
    addAlarmsIntoManualMetaAlarm.mockClear();
    removeAlarmsIntoManualMetaAlarm.mockClear();
  });

  const manualMetaAlarmModule = {
    name: 'manualMetaAlarm',
    actions: {
      fetchListWithoutStore: fetchManualMetaAlarmsListWithoutStore,
      create: createManualMetaAlarm,
      addAlarms: addAlarmsIntoManualMetaAlarm,
      removeAlarms: removeAlarmsIntoManualMetaAlarm,
    },
  };

  return {
    fetchManualMetaAlarmsListWithoutStore,
    createManualMetaAlarm,
    addAlarmsIntoManualMetaAlarm,
    removeAlarmsIntoManualMetaAlarm,
    manualMetaAlarmModule,
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
