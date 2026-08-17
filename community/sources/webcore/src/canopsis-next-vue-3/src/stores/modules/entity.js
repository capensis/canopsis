import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';

import request, { isRequestCanceled, useRequestCancelling } from '@/services/request';

import { mapIds } from '@/helpers/array';

export const useEntityStore = defineStore('entity', () => {
  const entitiesById = ref({});
  const widgets = ref({});

  const getList = computed(() => ids => ids.map(id => entitiesById.value[id]).filter(Boolean));
  const getListByWidgetId = computed(
    () => widgetId => getList.value(widgets.value[widgetId]?.allIds ?? []),
  );
  const getMetaByWidgetId = computed(() => widgetId => widgets.value[widgetId]?.meta ?? {});
  const getPendingByWidgetId = computed(
    () => widgetId => widgets.value[widgetId]?.pending ?? false,
  );
  const getFetchingParamsByWidgetId = computed(
    () => widgetId => widgets.value[widgetId]?.fetchingParams,
  );

  const setWidgetState = (widgetId, patch) => {
    widgets.value = {
      ...widgets.value,
      [widgetId]: {
        ...widgets.value[widgetId],
        ...patch,
      },
    };
  };

  const setEntities = (entities = []) => {
    entitiesById.value = entities.reduce(
      (acc, entity) => {
        acc[entity._id] = {
          ...acc[entity._id],
          ...entity,
        };

        return acc;
      },
      { ...entitiesById.value },
    );
  };

  const fetchContextEntitiesListWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entity, { params });
  const fetchListWithoutStore = fetchContextEntitiesListWithoutStore;
  const fetchEntityInfosLogsListWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityInfosLog, { params });
  const archiveDisabledEntitiesData = ({ data } = {}) =>
    request.post(`${API_ROUTES.entity}/archive-disabled`, data);
  const archiveUnlinkedEntitiesData = ({ data } = {}) =>
    request.post(`${API_ROUTES.entity}/archive-unlinked`, data);
  const cleanArchivedEntitiesData = () => request.post(`${API_ROUTES.entity}/clean-archived`);
  const checkStateSetting = ({ data } = {}) =>
    request.post(API_ROUTES.entityCheckStateSetting, data);
  const fetchStateSettingWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.entityStateSetting, { params });

  const fetchList = ({ widgetId, params } = {}) =>
    useRequestCancelling(async (source) => {
      try {
        setWidgetState(widgetId, { pending: true, fetchingParams: params });

        const { data = [], meta = {} } = await request.get(API_ROUTES.entity, {
          params,
          cancelToken: source.token,
        });

        setEntities(data);
        setWidgetState(widgetId, { allIds: mapIds(data), meta, pending: false });
      } catch (err) {
        if (isRequestCanceled(err)) {
          throw err;
        }

        console.error(err);

        setWidgetState(widgetId, { pending: false });

        throw err;
      }
    }, `entities-list-${widgetId}`);

  const fetchListWithPreviousParams = () => {
    const widgetsIds = Object.keys(widgets.value);

    return Promise.all(widgetsIds.map(widgetId => fetchList({
      widgetId,
      params: getFetchingParamsByWidgetId.value(widgetId),
    })));
  };

  const create = ({ data, id } = {}) =>
    request.post(API_ROUTES.entityBasics, data, {
      params: { _id: id },
    });
  const update = ({ data, id } = {}) =>
    request.put(API_ROUTES.entityBasics, data, {
      params: { _id: id },
    });
  const remove = ({ id } = {}) =>
    request.delete(API_ROUTES.entityBasics, {
      params: { _id: id },
    });

  const fetchBasicEntityWithoutStore = ({ id } = {}) =>
    request.get(API_ROUTES.entityBasics, {
      params: { _id: id },
    });
  const fetchContextGraphWithoutStore = ({ id } = {}) =>
    request.get(API_ROUTES.entityContextGraph, {
      params: { _id: id },
    });

  const createContextExport = ({ data = {} } = {}) =>
    request.post(API_ROUTES.contextExport, data);
  const fetchContextExport = ({ params, id } = {}) =>
    request.get(`${API_ROUTES.contextExport}/${id}`, { params });

  const bulkEnable = ({ data } = {}) => request.put(API_ROUTES.bulkEntitiesEnable, data);
  const bulkDisable = ({ data } = {}) => request.put(API_ROUTES.bulkEntitiesDisable, data);

  return {
    getListByWidgetId,
    getMetaByWidgetId,
    getPendingByWidgetId,
    getFetchingParamsByWidgetId,
    fetchContextEntitiesListWithoutStore,
    fetchListWithoutStore,
    fetchEntityInfosLogsListWithoutStore,
    archiveDisabledEntitiesData,
    archiveUnlinkedEntitiesData,
    cleanArchivedEntitiesData,
    checkStateSetting,
    fetchStateSettingWithoutStore,
    fetchList,
    fetchListWithPreviousParams,
    create,
    update,
    remove,
    fetchBasicEntityWithoutStore,
    fetchContextGraphWithoutStore,
    createContextExport,
    fetchContextExport,
    bulkEnable,
    bulkDisable,
  };
});
