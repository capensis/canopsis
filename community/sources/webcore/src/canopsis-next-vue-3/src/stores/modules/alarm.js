import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';

import request, { isRequestCanceled, useRequestCancelling } from '@/services/request';

import { mapIds } from '@/helpers/array';
import { getAlarmIdFromDetailsId } from '@/helpers/entities/alarm/list';

export const useAlarmStore = defineStore('alarm', () => {
  const alarmsById = ref({});
  const widgets = ref({});
  const detailsByWidget = ref({});

  const getItem = computed(() => id => alarmsById.value[id]);
  const getList = computed(() => ids => ids.map(id => alarmsById.value[id]).filter(Boolean));
  const getListByWidgetId = computed(
    () => widgetId => getList.value(widgets.value[widgetId]?.allIds ?? []),
  );
  const getMetaByWidgetId = computed(() => widgetId => widgets.value[widgetId]?.meta ?? {});
  const getPendingByWidgetId = computed(
    () => widgetId => widgets.value[widgetId]?.pending ?? false,
  );
  const getDetailsByWidgetId = computed(() => widgetId => detailsByWidget.value[widgetId] ?? {});
  const getQueriesByWidgetId = computed(
    () => widgetId => Object.values(detailsByWidget.value[widgetId] ?? {}).map(({ query }) => query),
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

  const setAlarms = (alarms = []) => {
    alarmsById.value = alarms.reduce(
      (acc, alarm) => {
        acc[alarm._id] = {
          ...acc[alarm._id],
          ...alarm,
        };

        return acc;
      },
      { ...alarmsById.value },
    );
  };

  const fetchComponentAlarmsListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.componentAlarms, {
      params,
    });
  const fetchResolvedAlarmsListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.resolvedAlarms, {
      params,
    });
  const fetchOpenAlarmsListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.openAlarms, { params });
  const fetchListWithoutStore = ({ params } = {}) => request.get(API_ROUTES.alarms.list, { params });

  const fetchList = ({ widgetId, params, withoutPending } = {}) => useRequestCancelling(async (source) => {
      try {
        if (!withoutPending) {
          setWidgetState(widgetId, { pending: true, fetchingParams: params });
        }

        const { data = [], meta = {} } = await request.get(API_ROUTES.alarms.list, {
          params,
          cancelToken: source.token,
        });

        setAlarms(data);
        setWidgetState(widgetId, { allIds: mapIds(data), meta, pending: false });
      } catch (err) {
        if (isRequestCanceled(err)) {
          throw err;
        }

        console.error(err);

        setWidgetState(widgetId, { pending: false });

        throw err;
      }
    }, `alarms-list-${widgetId}`);

  const fetchItem = async ({ id } = {}) => {
    const alarm = await request.get(`${API_ROUTES.alarms.list}/${id}`);

    setAlarms([alarm]);

    return alarm;
  };

  const fetchItemWithoutStore = ({ id } = {}) => request.get(`${API_ROUTES.alarms.list}/${id}`);
  const fetchAlarmsListExport = ({ params, id } = {}) => request.get(`${API_ROUTES.alarmListExport}/${id}`, {
      params,
    });
  const createAlarmsListExport = ({ data = {} } = {}) => request.post(API_ROUTES.alarmListExport, data);

  const createAlarmEvent = (event, { id, data } = {}) => request.put(`${API_ROUTES.alarms.list}/${id}/${event}`, data);
  const bulkCreateAlarmEvent = (event, { data } = {}) => request.put(`${API_ROUTES.alarms.bulkList}/${event}`, data);

  const addBookmarkToAlarm = ({ id } = {}) => request.put(`${API_ROUTES.alarms.list}/${id}/bookmark`);
  const removeBookmarkFromAlarm = ({ id } = {}) => request.delete(`${API_ROUTES.alarms.list}/${id}/bookmark`);
  const fetchDisplayNamesWithoutStore = ({ params } = {}) => request.get(API_ROUTES.alarmDisplayNames, { params });
  const fetchExecutionsWithoutStore = ({ id, params } = {}) => request.get(`${API_ROUTES.alarmExecutions}/${id}`, { params });
  const updateItemInStore = (alarm) => {
    setAlarms([alarm]);

    return alarm;
  };

  const fetchAlarmLinksWithoutStore = ({ params, id } = {}) => request.get(`${API_ROUTES.alarmLinks}/${id}`, {
      params,
    });

  const fetchAlarmDetailsListWithoutStore = ({ params } = {}) => request.post(API_ROUTES.alarmDetails, params);

  const getAlarmDetailsFromResponse = (response) => {
    const responseData = Array.isArray(response) ? response : response?.data;
    const detailsItem = Array.isArray(responseData) ? responseData[0] : responseData;

    return detailsItem?.data ?? detailsItem ?? {};
  };

  const setDetails = ({ widgetId, id, details }) => {
    detailsByWidget.value = {
      ...detailsByWidget.value,
      [widgetId]: {
        ...detailsByWidget.value[widgetId],
        [id]: {
          ...detailsByWidget.value[widgetId]?.[id],
          data: details,
        },
      },
    };
  };
  const setDetailsQuery = ({ widgetId, id, query }) => {
    detailsByWidget.value = {
      ...detailsByWidget.value,
      [widgetId]: {
        ...detailsByWidget.value[widgetId],
        [id]: {
          ...detailsByWidget.value[widgetId]?.[id],
          query,
        },
      },
    };
  };
  const removeDetailsQuery = ({ widgetId, id }) => {
    const widgetDetails = { ...detailsByWidget.value[widgetId] };

    delete widgetDetails[id];

    detailsByWidget.value = { ...detailsByWidget.value, [widgetId]: widgetDetails };
  };
  const fetchAlarmDetails = async ({ widgetId, id, query } = {}) => {
    setDetailsQuery({ widgetId, id, query });

    const response = await fetchAlarmDetailsListWithoutStore({ params: [query] });
    const details = getAlarmDetailsFromResponse(response);

    setDetails({ widgetId, id, details });

    return details;
  };

  /**
   * Fetch alarm details for all opened expand panels in a widget.
   *
   * @param {Object} payload
   * @param {string} payload.widgetId
   */
  const fetchAlarmDetailsList = async ({ widgetId } = {}) => {
    const queries = getQueriesByWidgetId.value(widgetId).filter(Boolean);

    if (!queries.length) {
      return;
    }

    try {
      const response = await fetchAlarmDetailsListWithoutStore({ params: queries });
      const responseData = Array.isArray(response) ? response : response?.data;
      const detailsItems = Array.isArray(responseData) ? responseData : [];

      detailsItems.forEach((item) => {
        const details = item?.data ?? item ?? {};
        const id = getAlarmIdFromDetailsId(item?._id ?? details?._id, widgetId);

        if (!id) {
          return;
        }

        setDetails({ widgetId, id, details });
      });
    } catch (err) {
      console.error(err);
    }
  };

  /**
   * Update alarm details item in store from a socket payload.
   *
   * @param {Object} payload
   * @param {string} payload.widgetId
   * @param {Object} payload.details
   * @returns {Object}
   */
  const updateDetailsItemInStore = ({ widgetId, details } = {}) => {
    const payload = details?.data ?? details ?? {};
    const id = getAlarmIdFromDetailsId(details?._id ?? payload?._id, widgetId);

    if (!id) {
      return details;
    }

    setDetails({ widgetId, id, details: payload });

    return details;
  };

  return {
    alarmsById,
    widgets,
    detailsByWidget,
    getItem,
    getList,
    getListByWidgetId,
    getMetaByWidgetId,
    getPendingByWidgetId,
    getDetailsByWidgetId,
    getQueriesByWidgetId,
    getFetchingParamsByWidgetId,
    setAlarms,
    fetchComponentAlarmsListWithoutStore,
    fetchResolvedAlarmsListWithoutStore,
    fetchOpenAlarmsListWithoutStore,
    fetchListWithoutStore,
    fetchList,
    fetchItem,
    fetchItemWithoutStore,
    fetchAlarmsListExport,
    createAlarmsListExport,
    createAlarmAckEvent: payload => createAlarmEvent('ack', payload),
    bulkCreateAlarmAckEvent: payload => bulkCreateAlarmEvent('ack', payload),
    createAlarmAckremoveEvent: payload => createAlarmEvent('ackremove', payload),
    bulkCreateAlarmAckremoveEvent: payload => bulkCreateAlarmEvent('ackremove', payload),
    createAlarmSnoozeEvent: payload => createAlarmEvent('snooze', payload),
    bulkCreateAlarmSnoozeEvent: payload => bulkCreateAlarmEvent('snooze', payload),
    createAlarmAssocticketEvent: payload => createAlarmEvent('assocticket', payload),
    bulkCreateAlarmAssocticketEvent: payload => bulkCreateAlarmEvent('assocticket', payload),
    bulkCreateAlarmTicketremoveEvent: payload => bulkCreateAlarmEvent('ticketremove', payload),
    createAlarmCommentEvent: payload => createAlarmEvent('comment', payload),
    bulkCreateAlarmCommentEvent: payload => bulkCreateAlarmEvent('comment', payload),
    createAlarmCancelEvent: payload => createAlarmEvent('cancel', payload),
    bulkCreateAlarmCancelEvent: payload => bulkCreateAlarmEvent('cancel', payload),
    createAlarmUncancelEvent: payload => createAlarmEvent('uncancel', payload),
    bulkCreateAlarmUnCancelEvent: payload => bulkCreateAlarmEvent('uncancel', payload),
    createAlarmChangestateEvent: payload => createAlarmEvent('changestate', payload),
    bulkCreateAlarmChangestateEvent: payload => bulkCreateAlarmEvent('changestate', payload),
    addBookmarkToAlarm,
    removeBookmarkFromAlarm,
    fetchDisplayNamesWithoutStore,
    fetchExecutionsWithoutStore,
    updateItemInStore,
    fetchAlarmLinksWithoutStore,
    fetchAlarmDetailsListWithoutStore,
    fetchAlarmDetails,
    fetchAlarmDetailsList,
    updateDetailsItemInStore,
    setDetails,
    setDetailsQuery,
    removeDetailsQuery,
  };
});
