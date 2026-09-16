import { config } from '@vue/test-utils';
import axios from 'axios';
import { vi } from 'vitest';
import { ref } from 'vue';

const defaultModals = {
  show: vi.fn(),
  hide: vi.fn(),
  minimize: vi.fn(),
  maximize: vi.fn(),
  updateDialogProps: vi.fn(),
  updateModalConfig: vi.fn(),
  dialogPropsMap: {},
  isModalOpenedById: vi.fn(() => false),
};
const defaultPopups = {
  error: vi.fn(),
  success: vi.fn(),
  info: vi.fn(),
};

vi.mock('@/hooks/popups', () => ({
  usePopups: () => globalThis.__TEST_POPUPS__ ?? defaultPopups,
}));

vi.mock('@/hooks/modals', async () => {
  const { computed, toRef } = await vi.importActual('vue');
  const getModals = () => globalThis.__TEST_MODALS__ ?? defaultModals;

  return {
    useModals: () => getModals(),
    useInnerModal: (props = {}) => {
      const modals = getModals();
      const modal = toRef(props, 'modal');

      return {
        modals,
        modal,
        config: computed(() => modal.value?.config ?? {}),
        close: () => modals.hide({ id: props.modal?.id ?? 'test-modal-id' }),
      };
    },
  };
});

vi.mock('@/hooks/store/modules/widget', async () => {
  const { default: request } = await vi.importActual('@/services/request');
  const { API_ROUTES } = await vi.importActual('@/config');
  const dispatch = (type, payload) =>
    typeof globalThis.__TEST_STORE__?.dispatch === 'function'
      ? globalThis.__TEST_STORE__.dispatch(type, payload)
      : undefined;

  return {
    useWidget: () => ({
      fetchWidgetWithoutStore: payload => dispatch('view/widget/fetchItemWithoutStore', payload),
      createWidget: payload => dispatch('view/widget/create', payload),
      cloneWidget: payload => dispatch('view/widget/clone', payload),
      updateWidget: payload => dispatch('view/widget/update', payload),
      removeWidget: payload => dispatch('view/widget/remove', payload),
      updateWidgetGridPositions: payload => dispatch('view/widget/updateGridPositions', payload),
      fetchWidgetFilters: payload =>
        dispatch('view/widget/fetchWidgetFilters', payload) ??
        request.get(API_ROUTES.widget.filters, { params: payload?.params }),
      createWidgetFilter: payload =>
        dispatch('view/widget/createWidgetFilter', payload) ??
        request.post(API_ROUTES.widget.filters, payload?.data),
      updateWidgetFilter: payload =>
        dispatch('view/widget/updateWidgetFilter', payload) ??
        request.put(`${API_ROUTES.widget.filters}/${payload?.id}`, payload?.data),
      removeWidgetFilter: payload =>
        dispatch('view/widget/removeWidgetFilter', payload) ??
        request.delete(`${API_ROUTES.widget.filters}/${payload?.id}`, { data: payload?.data }),
      fetchWidgetFilter: payload =>
        dispatch('view/widget/fetchWidgetFilter', payload) ??
        request.get(`${API_ROUTES.widget.filters}/${payload?.id}`, { params: payload?.params }),
      updateWidgetFiltersPositions: payload =>
        dispatch('view/widget/updateWidgetFiltersPositions', payload) ??
        request.put(API_ROUTES.widget.filterPositions, payload?.data),
    }),
  };
});

vi.mock('@/stores/modules/widget', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    useWidgetStore: defineStore('widget', {
      actions: {
        fetchItemWithoutStore: payload => dispatch('view/widget/fetchItemWithoutStore', payload),
        create: payload => dispatch('view/widget/create', payload),
        clone: payload => dispatch('view/widget/clone', payload),
        update: payload => dispatch('view/widget/update', payload),
        remove: payload => dispatch('view/widget/remove', payload),
        updateGridPositions: payload => dispatch('view/widget/updateGridPositions', payload),
        fetchWidgetFilters: payload => dispatch('view/widget/fetchWidgetFilters', payload),
        createWidgetFilter: payload => dispatch('view/widget/createWidgetFilter', payload),
        updateWidgetFilter: payload => dispatch('view/widget/updateWidgetFilter', payload),
        removeWidgetFilter: payload => dispatch('view/widget/removeWidgetFilter', payload),
        fetchWidgetFilter: payload => dispatch('view/widget/fetchWidgetFilter', payload),
        updateWidgetFiltersPositions: payload =>
          dispatch('view/widget/updateWidgetFiltersPositions', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/alarm', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);
  const getter = (type, fallback) => {
    const value = globalThis.__TEST_STORE__?.getters?.[type];

    return value ?? fallback;
  };

  return {
    useAlarmStore: defineStore('alarm', {
      getters: {
        getItem: () => getter('alarm/getItem', () => undefined),
        getList: () => getter('alarm/getList', () => []),
        getListByWidgetId: () => getter('alarm/getListByWidgetId', () => []),
        getMetaByWidgetId: () => getter('alarm/getMetaByWidgetId', () => ({})),
        getPendingByWidgetId: () => getter('alarm/getPendingByWidgetId', () => false),
        getDetailsByWidgetId: () => getter('alarm/getDetailsByWidgetId', () => ({})),
        getQueriesByWidgetId: () => getter('alarm/getQueriesByWidgetId', () => []),
        getFetchingParamsByWidgetId: () =>
          getter('alarm/getFetchingParamsByWidgetId', () => undefined),
      },
      actions: {
        fetchList: payload => dispatch('alarm/fetchList', payload),
        fetchListWithoutStore: payload => dispatch('alarm/fetchListWithoutStore', payload),
        fetchOpenAlarmsListWithoutStore: payload =>
          dispatch('alarm/fetchOpenAlarmsListWithoutStore', payload),
        fetchItem: payload => dispatch('alarm/fetchItem', payload),
        fetchItemWithoutStore: payload => dispatch('alarm/fetchItemWithoutStore', payload),
        fetchAlarmLinksWithoutStore: payload =>
          dispatch('alarm/fetchAlarmLinksWithoutStore', payload),
        fetchAlarmsListExport: payload => dispatch('alarm/fetchAlarmsListExport', payload),
        createAlarmsListExport: payload => dispatch('alarm/createAlarmsListExport', payload),
        bulkCreateAlarmAckEvent: payload => dispatch('alarm/bulkCreateAlarmAckEvent', payload),
        bulkCreateAlarmAckremoveEvent: payload =>
          dispatch('alarm/bulkCreateAlarmAckremoveEvent', payload),
        bulkCreateAlarmSnoozeEvent: payload =>
          dispatch('alarm/bulkCreateAlarmSnoozeEvent', payload),
        bulkCreateAlarmAssocticketEvent: payload =>
          dispatch('alarm/bulkCreateAlarmAssocticketEvent', payload),
        bulkCreateAlarmTicketremoveEvent: payload =>
          dispatch('alarm/bulkCreateAlarmTicketremoveEvent', payload),
        bulkCreateAlarmCommentEvent: payload =>
          dispatch('alarm/bulkCreateAlarmCommentEvent', payload),
        bulkCreateAlarmCancelEvent: payload =>
          dispatch('alarm/bulkCreateAlarmCancelEvent', payload),
        bulkCreateAlarmUnCancelEvent: payload =>
          dispatch('alarm/bulkCreateAlarmUnCancelEvent', payload),
        bulkCreateAlarmChangestateEvent: payload =>
          dispatch('alarm/bulkCreateAlarmChangestateEvent', payload),
        fetchAlarmDetails: payload => dispatch('alarm/details/fetchItem', payload),
        fetchAlarmDetailsList: payload => dispatch('alarm/details/fetchList', payload),
        fetchAlarmDetailsListWithoutStore: payload =>
          dispatch('alarm/details/fetchListWithoutStore', payload),
        removeDetailsQuery: payload => dispatch('alarm/details/removeQuery', payload),
        setDetails: vi.fn(),
        setDetailsQuery: vi.fn(),
        updateDetailsItemInStore: vi.fn(),
      },
    }),
  };
});

vi.mock('@/stores/modules/service', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    useServiceStore: defineStore('service', {
      actions: {
        fetchAlarmsWithoutStore: payload => dispatch('service/fetchAlarmsWithoutStore', payload),
        fetchEntityInfosKeysWithoutStore: payload =>
          dispatch('service/fetchInfosKeysWithoutStore', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/dynamic-info', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    useDynamicInfoStore: defineStore('dynamicInfo', {
      actions: {
        fetchInfosKeysWithoutStore: payload =>
          dispatch('dynamicInfo/fetchInfosKeysWithoutStore', payload) ?? { data: [] },
        fetchListWithoutStore: payload =>
          dispatch('dynamicInfo/fetchListWithoutStore', payload) ?? { data: [] },
      },
    }),
  };
});

vi.mock('@/stores/modules/template-vars', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);
  const fallback = {};

  return {
    useTemplateVarsStore: defineStore('templateVars', {
      state: () => ({
        items: {},
        pending: false,
      }),
      actions: {
        fetchTemplateVarsList: payload => dispatch('template/vars/fetchList', payload) ?? fallback,
        fetchEntityServicesTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchEntityServicesVarsWithoutStore', payload) ?? fallback,
        fetchEventFiltersTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchEventFiltersVarsWithoutStore', payload) ?? fallback,
        fetchScenariosTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchScenariosVarsWithoutStore', payload) ?? fallback,
        fetchLinkRulesTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchLinkRulesVarsWithoutStore', payload) ?? {
            alarm: {},
            entity: {},
          },
        fetchWidgetsTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchWidgetsVarsWithoutStore', payload) ?? fallback,
        fetchDeclareTicketRulesTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchDeclareTicketRulesVarsWithoutStore', payload) ?? fallback,
        fetchDynamicInfosTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchDynamicInfosVarsWithoutStore', payload) ?? fallback,
        fetchInstructionsTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchInstructionsVarsWithoutStore', payload) ?? fallback,
        fetchJobsTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchJobsVarsWithoutStore', payload) ?? fallback,
        fetchMetaAlarmRulesTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchMetaAlarmRulesVarsWithoutStore', payload) ?? fallback,
        fetchWebhookTokenRulesTemplateVarsWithoutStore: payload =>
          dispatch('template/vars/fetchWebhookTokenRulesVarsWithoutStore', payload) ?? fallback,
      },
    }),
  };
});

vi.mock('@/stores/modules/meta-alarm', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    useMetaAlarmStore: defineStore('metaAlarm', {
      actions: {
        create: payload => dispatch('metaAlarm/create', payload),
        addAlarms: payload => dispatch('metaAlarm/addAlarms', payload),
        removeAlarms: payload => dispatch('metaAlarm/removeAlarms', payload),
        fetchListWithoutStore: payload => dispatch('metaAlarm/fetchListWithoutStore', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/declare-ticket-rule', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    useDeclareTicketRuleStore: defineStore('declareTicketRule', {
      actions: {
        bulkCreateDeclareTicketExecution: payload =>
          dispatch('declareTicketRule/bulkCreateDeclareTicketExecution', payload),
        fetchAssignedTicketsWithoutStore: payload =>
          dispatch('declareTicketRule/fetchAssignedTicketsWithoutStore', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/pbehavior-type', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    usePbehaviorTypeStore: defineStore('pbehaviorType', {
      actions: {
        fetchListWithoutStore: payload => dispatch('pbehaviorTypes/fetchListWithoutStore', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/pbehavior-reason', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    usePbehaviorReasonStore: defineStore('pbehaviorReason', {
      actions: {
        fetchListWithoutStore: payload =>
          dispatch('pbehaviorReasons/fetchListWithoutStore', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/pbehavior', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    usePbehaviorStore: defineStore('pbehavior', {
      actions: {
        bulkCreateEntityPbehaviors: payload =>
          dispatch('pbehavior/bulkCreateEntityPbehaviors', payload),
        bulkRemoveEntityPbehaviors: payload =>
          dispatch('pbehavior/bulkRemoveEntityPbehaviors', payload),
        fetchListByEntityIdWithoutStore: payload =>
          dispatch('pbehavior/fetchListByEntityIdWithoutStore', payload),
        removeWithoutStore: payload => dispatch('pbehavior/removeWithoutStore', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/entity-info-property', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);
  const getter = (type, fallback) => {
    const value = globalThis.__TEST_STORE__?.getters?.[type];

    return value ?? fallback;
  };

  return {
    useEntityInfoPropertyStore: defineStore('entityInfoProperty', {
      getters: {
        items: () => getter('entityInfoProperty/items', []),
        itemsWithAlias: () => getter('entityInfoProperty/itemsWithAlias', []),
        itemsWithoutAlias: () => getter('entityInfoProperty/itemsWithoutAlias', []),
        meta: () => getter('entityInfoProperty/meta', {}),
        pending: () => getter('entityInfoProperty/pending', false),
      },
      actions: {
        fetchList: payload => dispatch('entityInfoProperty/fetchList', payload),
        fetchListWithoutStore: payload =>
          dispatch('entityInfoProperty/fetchListWithoutStore', payload),
        create: payload => dispatch('entityInfoProperty/create', payload),
        update: payload => dispatch('entityInfoProperty/update', payload),
        remove: payload => dispatch('entityInfoProperty/remove', payload),
        bulkRemove: payload => dispatch('entityInfoProperty/bulkRemove', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/associative-table', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);
  const getter = (type, fallback) => {
    const value = globalThis.__TEST_STORE__?.getters?.[type];

    return value ?? fallback;
  };

  return {
    useAssociativeTableStore: defineStore('associativeTable', {
      getters: {
        getItem: () => getter('associativeTable/getItem', () => undefined),
      },
      actions: {
        fetch: payload => dispatch('associativeTable/fetch', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/map', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const { default: request } = await vi.importActual('@/services/request');
  const { API_ROUTES } = await vi.importActual('@/config');
  const dispatch = (type, payload) =>
    typeof globalThis.__TEST_STORE__?.dispatch === 'function'
      ? globalThis.__TEST_STORE__.dispatch(type, payload)
      : undefined;
  const getter = (type, fallback) => {
    const value = globalThis.__TEST_STORE__?.getters?.[type];

    return value ?? fallback;
  };

  return {
    useMapStore: defineStore('map', {
      state: () => ({
        fetchingParams: undefined,
        list: [],
        listMeta: {},
        listPending: false,
      }),
      getters: {
        items: state => getter('map/items', state.list),
        meta: state => getter('map/meta', state.listMeta),
        pending: state => getter('map/pending', state.listPending),
      },
      actions: {
        async fetchList(payload = {}) {
          const dispatched = dispatch('map/fetchList', payload);

          if (dispatched !== undefined) {
            return dispatched;
          }

          this.listPending = true;

          try {
            const response = await request.get(API_ROUTES.maps, { params: payload.params });

            this.list = response.data;
            this.listMeta = response.meta;

            return response;
          } catch (err) {
            console.error(err);

            throw err;
          } finally {
            this.listPending = false;
          }
        },
        fetchListWithPreviousParams(payload) {
          return this.fetchList(payload);
        },
        fetchListWithoutStore(payload = {}) {
          return (
            dispatch('map/fetchListWithoutStore', payload) ??
            request.get(API_ROUTES.maps, { params: payload.params })
          );
        },
        fetchItemWithoutStore(payload = {}) {
          return (
            dispatch('map/fetchItemWithoutStore', payload) ??
            request.get(`${API_ROUTES.maps}/${payload.id}`, { params: payload.params })
          );
        },
        fetchItemStateWithoutStore(payload = {}) {
          return (
            dispatch('map/fetchItemStateWithoutStore', payload) ??
            request.get(`${API_ROUTES.mapState}/${payload.id}`, { params: payload.params })
          );
        },
        create(payload = {}) {
          return dispatch('map/create', payload) ?? request.post(API_ROUTES.maps, payload.data);
        },
        update(payload = {}) {
          return (
            dispatch('map/update', payload) ??
            request.put(`${API_ROUTES.maps}/${encodeURIComponent(payload.id)}`, payload.data)
          );
        },
        remove(payload = {}) {
          return (
            dispatch('map/remove', payload) ??
            request.delete(`${API_ROUTES.maps}/${encodeURIComponent(payload.id)}`)
          );
        },
        bulkRemove(payload = {}) {
          return (
            dispatch('map/bulkRemove', payload) ??
            request.delete(API_ROUTES.bulkMaps, { data: payload.data })
          );
        },
      },
    }),
  };
});

vi.mock('@/stores/modules/query', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);
  const getter = (type, fallback) => {
    const value = globalThis.__TEST_STORE__?.getters?.[type];

    return value ?? fallback;
  };

  return {
    useQueryStore: defineStore('query', {
      getters: {
        getQueryById: () => getter('query/getQueryById', () => ({})),
        getQueryNonceById: () => getter('query/getQueryNonceById', () => 0),
      },
      actions: {
        update: payload => dispatch('query/update', payload),
        merge: payload => dispatch('query/merge', payload),
        remove: payload => dispatch('query/remove', payload),
        forceUpdate: payload => dispatch('query/forceUpdate', payload),
        updateLocked: payload => dispatch('query/updateLocked', payload),
        removeLocked: payload => dispatch('query/removeLocked', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/user-preference', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);
  const getter = (type, fallback) => {
    const value = globalThis.__TEST_STORE__?.getters?.[type];

    return value ?? fallback;
  };

  return {
    useUserPreferenceStore: defineStore('userPreference', {
      getters: {
        items: () => [],
        pending: () => false,
        getItemByWidgetId: () =>
          getter('userPreference/getItemByWidgetId', () => ({ content: {} })),
      },
      actions: {
        fetchItem: payload => dispatch('userPreference/fetchItem', payload),
        fetchItemWithoutStore: payload => dispatch('userPreference/fetchItemWithoutStore', payload),
        update: payload => dispatch('userPreference/update', payload),
        updateLocal: payload => dispatch('userPreference/updateLocal', payload),
      },
    }),
  };
});

const createTestStoreGetter = (type, fallback) => {
  const value = globalThis.__TEST_STORE__?.getters?.[type];

  if (typeof value === 'function') {
    const state = globalThis.__TEST_STORE__?.state?.view ?? {};

    return value(state, globalThis.__TEST_STORE__?.getters, globalThis.__TEST_STORE__?.state, globalThis.__TEST_STORE__?.getters);
  }

  return value ?? fallback;
};

vi.mock('@/stores/modules/view', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    useViewStore: defineStore('view', {
      getters: {
        items: () => createTestStoreGetter('view/items', []),
        pending: () => createTestStoreGetter('view/pending', false),
        getGroupById: () => id => createTestStoreGetter('view/items', []).find?.(group => group._id === id),
        getViewById: () => () => createTestStoreGetter('view/getViewById', undefined),
        getViewTabById: () => () => createTestStoreGetter('view/getViewTabById', undefined),
      },
      actions: {
        fetchList: payload => dispatch('view/fetchList', payload),
        fetchListWithoutStore: payload => dispatch('view/fetchListWithoutStore', payload),
        create: payload => dispatch('view/create', payload),
        update: payload => dispatch('view/update', payload),
        remove: payload => dispatch('view/remove', payload),
        createPrivateGroup: payload => dispatch('view/createPrivateGroup', payload),
        updatePrivateGroup: payload => dispatch('view/updatePrivateGroup', payload),
        removePrivateGroup: payload => dispatch('view/removePrivateGroup', payload),
        createView: payload => dispatch('view/createView', payload),
        updateView: payload => dispatch('view/updateView', payload),
        updateWithoutStoreView: payload => dispatch('view/updateWithoutStoreView', payload),
        updateViewPositions: payload => dispatch('view/updateViewPositions', payload),
        removeView: payload => dispatch('view/removeView', payload),
        copyView: payload => dispatch('view/copyView', payload),
        exportViewWithoutStore: payload => dispatch('view/exportViewWithoutStore', payload),
        importViewWithoutStore: payload => dispatch('view/importViewWithoutStore', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/navigation', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    useNavigationStore: defineStore('navigation', {
      getters: {
        isEditingMode: () => createTestStoreGetter('navigation/isEditingMode', false),
      },
      actions: {
        toggleEditingMode: payload => dispatch('navigation/toggleEditingMode', payload),
      },
    }),
  };
});

vi.mock('@/stores/modules/auth', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);
  const getter = (type, fallback) => {
    const value = globalThis.__TEST_STORE__?.getters?.[type];

    return value ?? fallback;
  };

  return {
    useAuthStore: defineStore('auth', {
      getters: {
        isLoggedIn: () => true,
        currentUser: () => getter('auth/currentUser', {}),
        pending: () => false,
        currentUserPermissionsById: () => getter('auth/currentUserPermissionsById', {}),
        currentUserViewPermissionsByViewId: () =>
          getter('auth/currentUserViewPermissionsByViewId', {}),
      },
      actions: {
        login: payload => dispatch('auth/login', payload),
        applyAccessToken: payload => dispatch('auth/applyAccessToken', payload),
        logout: payload => dispatch('auth/logout', payload),
        fetchCurrentUser: payload => dispatch('auth/fetchCurrentUser', payload),
        filesAccess: payload => dispatch('auth/filesAccess', payload),
      },
    }),
  };
});

vi.mock('@/hooks/store/modules/active-view', () => ({
  useActiveView: () => ({
    view: ref(globalThis.__TEST_STORE__?.getters?.['activeView/item'] ?? {}),
    activeViewPending: ref(false),
    activeViewScreenMode: ref(globalThis.__TEST_STORE__?.getters?.['activeView/screenMode']),
    activeViewEditing: ref(globalThis.__TEST_STORE__?.getters?.['activeView/editing'] ?? false),
    activeViewEditingProcess: ref(false),
    activeViewPeriodicRefreshPaused: ref(false),
    activeViewIsKioskScreenMode: ref(
      globalThis.__TEST_STORE__?.getters?.['activeView/isKioskScreenMode'] ?? false,
    ),
    toggleEditing: payload =>
      globalThis.__TEST_STORE__?.dispatch('activeView/toggleEditing', payload),
    registerEditingOffHandler: payload =>
      globalThis.__TEST_STORE__?.dispatch('activeView/registerEditingOffHandler', payload),
    unregisterEditingOffHandler: payload =>
      globalThis.__TEST_STORE__?.dispatch('activeView/unregisterEditingOffHandler', payload),
    fetchActiveView: payload => globalThis.__TEST_STORE__?.dispatch('activeView/fetch', payload),
    clearActiveView: payload => globalThis.__TEST_STORE__?.dispatch('activeView/clear', payload),
    resumePeriodicRefresh: payload =>
      globalThis.__TEST_STORE__?.dispatch('activeView/resumePeriodicRefresh', payload),
    pausePeriodicRefresh: payload =>
      globalThis.__TEST_STORE__?.dispatch('activeView/pausePeriodicRefresh', payload),
    setActiveViewScreenMode: payload =>
      globalThis.__TEST_STORE__?.dispatch('activeView/setScreenMode', payload),
  }),
}));

vi.mock('@/stores/modules/active-view', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const getter = (type, fallback) => {
    const value = globalThis.__TEST_STORE__?.getters?.[type];

    return value ?? fallback;
  };
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);

  return {
    useActiveViewStore: defineStore('activeView', {
      getters: {
        item: () => getter('activeView/item', {}),
        pending: () => false,
        editing: () => getter('activeView/editing', false),
        editingProcess: () => false,
        screenMode: () => getter('activeView/screenMode', undefined),
        periodicRefreshPaused: () => getter('activeView/periodicRefreshPaused', false),
        isKioskScreenMode: () => getter('activeView/isKioskScreenMode', false),
      },
      actions: {
        registerEditingOffHandler: payload =>
          dispatch('activeView/registerEditingOffHandler', payload),
        unregisterEditingOffHandler: payload =>
          dispatch('activeView/unregisterEditingOffHandler', payload),
        toggleEditing: payload => dispatch('activeView/toggleEditing', payload),
        fetch: payload => dispatch('activeView/fetch', payload),
        clear: payload => dispatch('activeView/clear', payload),
        resumePeriodicRefresh: payload => dispatch('activeView/resumePeriodicRefresh', payload),
        pausePeriodicRefresh: payload => dispatch('activeView/pausePeriodicRefresh', payload),
        setScreenMode: payload => dispatch('activeView/setScreenMode', payload),
      },
    }),
  };
});

vi.mock('@/hooks/store/modules/user-preference', async () => {
  const { default: request } = await vi.importActual('@/services/request');
  const { API_ROUTES } = await vi.importActual('@/config');

  return {
    useUserPreference: () => ({
      userPreferences: ref([]),
      userPreferencesByWidget: ref({}),
      userPreferencesPending: ref(false),
      getUserPreferenceByWidgetId: ref((id) => {
        const getter = globalThis.__TEST_STORE__?.getters?.['userPreference/getItemByWidgetId'];

        return typeof getter === 'function' ? getter(id) : { content: {} };
      }),
      fetchUserPreference: payload =>
        typeof globalThis.__TEST_STORE__?.dispatch === 'function'
          ? globalThis.__TEST_STORE__.dispatch('userPreference/fetchItem', payload)
          : request.get(`${API_ROUTES.userPreferences}/${payload?.id}`, {
              params: payload?.params,
            }),
      fetchUserPreferenceWithoutStore: payload =>
        typeof globalThis.__TEST_STORE__?.dispatch === 'function'
          ? globalThis.__TEST_STORE__.dispatch('userPreference/fetchItemWithoutStore', payload)
          : request.get(`${API_ROUTES.userPreferences}/${payload?.id}`, {
              params: payload?.params,
            }),
      updateUserPreference: payload =>
        typeof globalThis.__TEST_STORE__?.dispatch === 'function'
          ? globalThis.__TEST_STORE__.dispatch('userPreference/update', payload)
          : request.put(`${API_ROUTES.userPreferences}/${payload?.id}`, payload?.data),
      updateLocalUserPreference: vi.fn(),
    }),
  };
});

vi.mock('@/hooks/sidebar', () => ({
  useSidebar: () => ({
    moduleName: 'sidebar',
    show: payload =>
      globalThis.__TEST_SIDEBAR__
        ? globalThis.__TEST_SIDEBAR__.show(payload)
        : globalThis.__TEST_STORE__?.dispatch('sidebar/show', payload),
    hide: payload =>
      globalThis.__TEST_SIDEBAR__
        ? globalThis.__TEST_SIDEBAR__.hide(payload)
        : globalThis.__TEST_STORE__?.dispatch('sidebar/hide', payload),
    updateConfig: payload =>
      globalThis.__TEST_SIDEBAR__
        ? globalThis.__TEST_SIDEBAR__.updateConfig(payload)
        : globalThis.__TEST_STORE__?.dispatch('sidebar/updateConfig', payload),
  }),
}));

vi.mock('@/hooks/store/modules/patterns-fields', () => {
  const fields = ref({
    alarm_pattern: [],
    entity_pattern: [],
    event_pattern: [],
    pbehavior_pattern: [],
    weather_service_pattern: [],
  });

  return {
    usePatternsFields: () => ({
      fetchPatternsFields: vi.fn().mockResolvedValue({}),
      fetchAlarmTagPatternFields: vi.fn().mockResolvedValue({}),
      fetchDeclareTicketRulePatternFields: vi.fn().mockResolvedValue({}),
      fetchDynamicInfosPatternFields: vi.fn().mockResolvedValue({}),
      fetchEventFilterPatternFields: vi.fn().mockResolvedValue({}),
      fetchEventRecordPatternFields: vi.fn().mockResolvedValue({}),
      fetchFlappingRulePatternFields: vi.fn().mockResolvedValue({}),
      fetchIdleRulePatternFields: vi.fn().mockResolvedValue({}),
      fetchInstructionPatternFields: vi.fn().mockResolvedValue({}),
      fetchKpiFilterPatternFields: vi.fn().mockResolvedValue({}),
      fetchLinkRulePatternFields: vi.fn().mockResolvedValue({}),
      fetchMetaalarmrulePatternFields: vi.fn().mockResolvedValue({}),
      fetchPbehaviorPatternFields: vi.fn().mockResolvedValue({}),
      fetchResolveRulePatternFields: vi.fn().mockResolvedValue({}),
      fetchScenarioPatternFields: vi.fn().mockResolvedValue({}),
      fetchServicePatternFields: vi.fn().mockResolvedValue({}),
      fetchStateSettingPatternFields: vi.fn().mockResolvedValue({}),
      fetchWidgetFilterPatternFields: vi.fn().mockResolvedValue({}),
      patternsFields: fields,
      patternsFieldsPending: ref(false),
    }),
    usePatternsFieldsFetching: () => ({
      alarmAttributes: ref([]),
      entityAttributes: ref([]),
      eventAttributes: ref([]),
      fetchPatternsFields: vi.fn().mockResolvedValue({}),
      pbehaviorAttributes: ref([]),
      pending: ref(false),
      weatherServiceAttributes: ref([]),
    }),
  };
});

globalThis.jest = vi;

axios.defaults.adapter = async axiosConfig => ({
  config: axiosConfig,
  data: {
    data: [],
    meta: {
      page: 1,
      per_page: 10,
      total_count: 0,
    },
  },
  headers: {},
  request: {},
  status: 200,
  statusText: 'OK',
});

globalThis.IntersectionObserver = vi.fn(function IntersectionObserver() {
  return {
    disconnect: vi.fn(),
    observe: vi.fn(),
    unobserve: vi.fn(),
  };
});

globalThis.ResizeObserver = vi.fn(function ResizeObserver() {
  return {
    disconnect: vi.fn(),
    observe: vi.fn(),
    unobserve: vi.fn(),
  };
});

Element.prototype.scrollIntoView = vi.fn();
Element.prototype.scrollTo = vi.fn();

HTMLCanvasElement.prototype.getContext = vi.fn(() => ({
  arc: vi.fn(),
  beginPath: vi.fn(),
  clearRect: vi.fn(),
  closePath: vi.fn(),
  drawImage: vi.fn(),
  fill: vi.fn(),
  fillRect: vi.fn(),
  fillText: vi.fn(),
  getImageData: vi.fn(() => ({ data: [] })),
  lineTo: vi.fn(),
  measureText: vi.fn(() => ({ width: 0 })),
  moveTo: vi.fn(),
  putImageData: vi.fn(),
  restore: vi.fn(),
  save: vi.fn(),
  scale: vi.fn(),
  setLineDash: vi.fn(),
  stroke: vi.fn(),
  translate: vi.fn(),
}));
HTMLCanvasElement.prototype.toBlob = vi.fn(callback => callback(new Blob()));

globalThis.visualViewport = {
  addEventListener: vi.fn(),
  height: 768,
  offsetLeft: 0,
  offsetTop: 0,
  removeEventListener: vi.fn(),
  scale: 1,
  width: 1024,
};

config.global.mocks = {
  ...config.global.mocks,
  $system: {
    timezone: 'UTC',
    theme: 'canopsis',
  },
  $t: key => key,
  $tc: key => key,
  $te: () => true,
};

config.global.provide = {
  ...config.global.provide,
  $system: {
    timezone: 'UTC',
    theme: 'canopsis',
  },
  $clickOutside: {
    register: vi.fn(),
    unregister: vi.fn(),
  },
};

config.global.config = {
  ...config.global.config,
  compilerOptions: {
    ...(config.global.config ?? {}).compilerOptions,
    isCustomElement: tag => tag.endsWith('-stub'),
  },
};

expect.extend({
  toEmitInput(wrapper, value) {
    const emitted = wrapper.emitted?.() ?? {};
    const events = [
      ...(emitted.input ?? []),
      ...(emitted.change ?? []),
      ...Object.entries(emitted)
        .filter(([eventName]) => eventName.startsWith('update:'))
        .flatMap(([, eventCalls]) => eventCalls),
    ];
    const formValues = globalThis.__TEST_FORM__?.values?.value;
    const pass =
      events.some(args => this.equals(args?.[0], value)) || this.equals(formValues, value);

    return {
      pass,
      message: () =>
        `expected wrapper ${pass ? 'not ' : ''}to emit input/update:value/update:modelValue with ${this.utils.printExpected(value)}\nReceived: ${this.utils.printReceived(events.map(args => args?.[0]))}`,
    };
  },

  toHaveBeenEmit(wrapper, eventName) {
    const emitted = wrapper.emitted?.() ?? {};
    const pass = Boolean(emitted[eventName]?.length);

    return {
      pass,
      message: () => `expected wrapper ${pass ? 'not ' : ''}to emit ${eventName}`,
    };
  },

  toEmit(wrapper, eventName, value) {
    const emitted = wrapper.emitted?.() ?? {};
    const events = emitted[eventName] ?? [];
    const pass =
      arguments.length >= 3
        ? events.some(args => this.equals(args?.[0], value))
        : Boolean(events.length);

    return {
      pass,
      message: () => `expected wrapper ${pass ? 'not ' : ''}to emit ${eventName}`,
    };
  },

  toHaveBeenEmitWith(wrapper, eventName, value) {
    const emitted = wrapper.emitted?.() ?? {};
    const events = emitted[eventName] ?? [];
    const pass = events.some(args => this.equals(args?.[0], value));

    return {
      pass,
      message: () =>
        `expected wrapper ${pass ? 'not ' : ''}to emit ${eventName} with ${this.utils.printExpected(value)}`,
    };
  },

  toBeDispatchedWith(mock, payload) {
    const calls = mock?.mock?.calls ?? [];
    const pass = calls.some(args => args.some(arg => this.equals(arg, payload)));

    return {
      pass,
      message: () =>
        `expected dispatch mock ${pass ? 'not ' : ''}to be called with ${this.utils.printExpected(payload)}`,
    };
  },
});
