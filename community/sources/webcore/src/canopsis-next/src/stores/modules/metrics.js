import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

export const useMetricsStore = defineStore('metrics', () => {
  const widgetMetrics = ref({});
  const groupMetrics = ref({});
  const vectorMetrics = ref({});
  const availabilityMetrics = ref({});

  const getListByWidgetId = computed(() => widgetId => widgetMetrics.value[widgetId]?.items ?? []);
  const getPendingByWidgetId = computed(
    () => widgetId => widgetMetrics.value[widgetId]?.pending ?? false,
  );
  const getGroupMetricsListByWidgetId = computed(
    () => widgetId => groupMetrics.value[widgetId]?.items ?? [],
  );
  const getGroupMetricsPendingByWidgetId = computed(
    () => widgetId => groupMetrics.value[widgetId]?.pending ?? false,
  );
  const getGroupMetricsMetaByWidgetId = computed(
    () => widgetId => groupMetrics.value[widgetId]?.meta ?? {},
  );
  const getVectorMetricsListByWidgetId = computed(
    () => widgetId => vectorMetrics.value[widgetId]?.items ?? [],
  );
  const getVectorMetricsPendingByWidgetId = computed(
    () => widgetId => vectorMetrics.value[widgetId]?.pending ?? false,
  );
  const getVectorMetricsMetaByWidgetId = computed(
    () => widgetId => vectorMetrics.value[widgetId]?.meta ?? {},
  );
  const getAvailabilityListByWidgetId = computed(
    () => widgetId => availabilityMetrics.value[widgetId]?.items ?? [],
  );
  const getAvailabilityPendingByWidgetId = computed(
    () => widgetId => availabilityMetrics.value[widgetId]?.pending ?? false,
  );
  const getAvailabilityMetaByWidgetId = computed(
    () => widgetId => availabilityMetrics.value[widgetId]?.meta ?? {},
  );

  const setWidgetMetrics = (widgetId, patch) => {
    widgetMetrics.value = {
      ...widgetMetrics.value,
      [widgetId]: {
        ...widgetMetrics.value[widgetId],
        ...patch,
      },
    };
  };

  const patchWidgetBucket = (bucket, widgetId, patch) => ({
    ...bucket,
    [widgetId]: {
      ...bucket[widgetId],
      ...patch,
    },
  });

  const setGroupMetrics = (widgetId, patch) => {
    groupMetrics.value = patchWidgetBucket(groupMetrics.value, widgetId, patch);
  };

  const setVectorMetrics = (widgetId, patch) => {
    vectorMetrics.value = patchWidgetBucket(vectorMetrics.value, widgetId, patch);
  };

  const setAvailabilityMetrics = (widgetId, patch) => {
    availabilityMetrics.value = patchWidgetBucket(availabilityMetrics.value, widgetId, patch);
  };

  const createBucketSetter = bucketName =>
    ({
      group: setGroupMetrics,
      vector: setVectorMetrics,
      availability: setAvailabilityMetrics,
    })[bucketName];

  const setBucketByName = (bucketName, widgetId, patch) => {
    const setBucket = createBucketSetter(bucketName);

    setBucket(widgetId, patch);
  };

  const prepareAvailabilityInfos = infos =>
    infos
      ? Object.entries(infos).reduce((acc, [name, value]) => {
          acc[name] = { value, name };

          return acc;
        }, {})
      : {};

  const prepareAvailabilitiesResponse = response =>
    response?.data?.map(availability => ({
      ...availability,
      entity: {
        ...availability.entity,
        infos: prepareAvailabilityInfos(availability.entity?.infos),
        component_infos: prepareAvailabilityInfos(availability.entity?.component_infos),
      },
    })) ?? [];

  const fetchWidgetList = async (
    bucketName,
    { widgetId, params } = {},
    {
      route,
      method = 'post',
      dataPreparer = response => response?.data ?? [],
      metaPreparer = response => response?.meta ?? {},
    },
  ) => {
    try {
      setBucketByName(bucketName, widgetId, { pending: true, error: null });

      const response =
        method === 'get' ? await request.get(route, { params }) : await request.post(route, params);

      const items = dataPreparer(response);
      const meta = metaPreparer(response);

      setBucketByName(bucketName, widgetId, { widgetId, items, meta, pending: false });

      return response;
    } catch (err) {
      console.error(err);

      setBucketByName(bucketName, widgetId, { widgetId, pending: false });

      return undefined;
    }
  };

  const fetchSliMetricsWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.metrics.sli, { params });
  const fetchRatingMetricsWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.metrics.rating, { params });
  const fetchAlarmsMetricsWithoutStore = ({ params } = {}) =>
    request.post(API_ROUTES.metrics.alarm, params);
  const createKpiAlarmExport = ({ data } = {}) =>
    request.post(API_ROUTES.metrics.exportAlarm, data);
  const createKpiAlarmAggregateExport = ({ data } = {}) =>
    request.post(API_ROUTES.metrics.exportAggregate, data);
  const createRemediationExport = ({ data } = {}) =>
    request.post(API_ROUTES.metrics.exportRemediation, null, { params: data });
  const createKpiRatingExport = ({ data } = {}) =>
    request.post(API_ROUTES.metrics.exportRating, null, {
      params: data,
    });
  const createKpiSliExport = ({ data } = {}) =>
    request.post(API_ROUTES.metrics.exportSli, null, { params: data });
  const fetchMetricExport = ({ id } = {}) =>
    request.get(`${API_ROUTES.metrics.exportMetric}/${id}`);
  const fetchExternalMetricsWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.metrics.perfDataMetrics, {
      params,
    });
  const fetchEntityAlarmsMetricsWithoutStore = ({ params } = {}) =>
    request.post(API_ROUTES.metrics.entityAlarmMetrics, params);
  const fetchEntityAggregateMetricsWithoutStore = ({ params } = {}) =>
    request.post(API_ROUTES.metrics.entityAggregateMetrics, params);

  const fetchAggregatedMetricsWithoutStore = async ({
    params: { with_history: withHistory, prev_from: prevFrom, prev_to: prevTo, ...params } = {},
  } = {}) => {
    const { data: metrics = [] } = await request.post(API_ROUTES.metrics.aggregate, params);

    if (!withHistory) {
      return { data: metrics };
    }

    try {
      const previousInterval = {
        from: prevFrom ?? params.from - (params.to - params.from),
        to: prevTo ?? params.from,
      };
      const { data: previousMetrics = [] } = await request.post(API_ROUTES.metrics.aggregate, {
        ...params,
        ...previousInterval,
      });

      return {
        data: metrics.map((metric, index) => ({
          ...metric,
          previous_metric: previousMetrics[index]?.value,
          previous_interval: previousInterval,
        })),
      };
    } catch (err) {
      console.error(err);

      return { data: metrics };
    }
  };

  const fetchAggregatedMetrics = async ({ widgetId, params = {} } = {}) => {
    setWidgetMetrics(widgetId, { pending: true, error: null });

    const { data = [] } = await fetchAggregatedMetricsWithoutStore({ params });

    setWidgetMetrics(widgetId, { items: data, pending: false });

    return data;
  };

  const fetchGroupMetricsList = payload =>
    fetchWidgetList('group', payload, {
      route: API_ROUTES.metrics.group,
    });
  const fetchGroupMetrics = fetchGroupMetricsList;
  const createGroupMetricsExport = ({ data } = {}) =>
    request.post(API_ROUTES.metrics.exportGroup, data);
  const fetchVectorMetricsList = payload =>
    fetchWidgetList('vector', payload, {
      route: API_ROUTES.metrics.alarm,
    });
  const fetchVectorMetrics = fetchVectorMetricsList;
  const fetchAvailabilityList = payload =>
    fetchWidgetList('availability', payload, {
      route: API_ROUTES.metrics.availability,
      method: 'get',
      dataPreparer: prepareAvailabilitiesResponse,
    });
  const fetchAvailability = fetchAvailabilityList;
  const fetchAvailabilityWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.metrics.entityAggregateAvailability, { params });
  const fetchAvailabilityHistoryWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.metrics.entityAvailability, { params });
  const fetchAvailabilityHistoryExport = ({ params } = {}) =>
    request.get(API_ROUTES.metrics.exportAvailabilityByEntity, { params });
  const createAvailabilityExport = ({ data } = {}) =>
    request.post(API_ROUTES.metrics.exportAvailability, data);
  const fetchAvailabilityExport = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.metrics.exportMetric}/${id}`, { params });

  return {
    widgetMetrics,
    groupMetrics,
    vectorMetrics,
    availabilityMetrics,
    getListByWidgetId,
    getPendingByWidgetId,
    getGroupMetricsListByWidgetId,
    getGroupMetricsPendingByWidgetId,
    getGroupMetricsMetaByWidgetId,
    getVectorMetricsListByWidgetId,
    getVectorMetricsPendingByWidgetId,
    getVectorMetricsMetaByWidgetId,
    getAvailabilityListByWidgetId,
    getAvailabilityPendingByWidgetId,
    getAvailabilityMetaByWidgetId,
    fetchSliMetricsWithoutStore,
    fetchRatingMetricsWithoutStore,
    fetchAlarmsMetricsWithoutStore,
    createKpiAlarmExport,
    createKpiAlarmAggregateExport,
    createRemediationExport,
    createKpiRatingExport,
    createKpiSliExport,
    fetchMetricExport,
    fetchExternalMetricsWithoutStore,
    fetchEntityAlarmsMetricsWithoutStore,
    fetchEntityAggregateMetricsWithoutStore,
    fetchAggregatedMetrics,
    fetchAggregatedMetricsWithoutStore,
    fetchGroupMetricsList,
    fetchGroupMetrics,
    createGroupMetricsExport,
    fetchVectorMetricsList,
    fetchVectorMetrics,
    fetchAvailabilityList,
    fetchAvailability,
    fetchAvailabilityWithoutStore,
    fetchAvailabilityHistoryWithoutStore,
    fetchAvailabilityHistoryExport,
    createAvailabilityExport,
    fetchAvailabilityExport,
  };
});
