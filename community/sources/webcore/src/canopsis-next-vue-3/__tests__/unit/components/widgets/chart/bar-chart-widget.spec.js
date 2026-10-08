import { randomArrayItem } from '@unit/utils/array';
import { mockDateNow } from '@unit/utils/mock-hooks';
import {
  createActiveViewModule,
  createAlarmModule,
  createAuthModule,
  createMockedStoreModules,
  createQueryModule,
  createServiceModule,
  createUserPreferenceModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS } from '@/constants';

import BarChartWidget from '@/components/widgets/chart/bar-chart-widget.vue';

const metrics = vi.hoisted(() => ({
  createKpiAlarmExport: vi.fn(),
  fetchMetricExport: vi.fn(),
  fetchVectorMetricsList: vi.fn(),
}));

vi.mock('@/stores/modules/metrics', async () => {
  const { defineStore } = await vi.importActual('pinia');

  return {
    useMetricsStore: defineStore('metrics', {
      getters: {
        getVectorMetricsListByWidgetId: () => () => [],
        getVectorMetricsPendingByWidgetId: () => () => false,
        getVectorMetricsMetaByWidgetId: () => () => ({}),
      },
      actions: metrics,
    }),
  };
});

const stubs = {
  'kpi-widget-filters': true,
};

describe('bar-chart-widget', () => {
  const nowTimestamp = 1386435500000;
  mockDateNow(nowTimestamp);

  const widgetId = 'widget-id';

  const { authModule } = createAuthModule();
  const { activeViewModule } = createActiveViewModule();
  const { alarmModule } = createAlarmModule();
  const { userPreferenceModule, fetchUserPreference } = createUserPreferenceModule();
  const { serviceModule } = createServiceModule();
  const { queryModule, updateQuery, getQueryById } = createQueryModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    activeViewModule,
    alarmModule,
    serviceModule,
    queryModule,
  ]);

  const widget = {
    _id: widgetId,
    parameters: {
      default_sampling: SAMPLINGS.month,
      default_time_range: QUICK_RANGES.last7Days.value,
      metrics: [{ metric: ALARM_METRIC_PARAMETERS.createdAlarms }],
    },
  };

  const factory = generateShallowRenderer(BarChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: { timezone: 'Europe/Paris' },
      },
    },
  });

  test('Query updated after mount', async () => {
    factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    expect(fetchUserPreference).toHaveBeenCalledWith(expect.any(Object), { id: widget._id });
    expect(updateQuery).toHaveBeenCalledWith(expect.any(Object), {
      id: widget._id,
      query: {
        search: '',
      },
    });
  });

  test('Vectors metrics fetched with correct query', async () => {
    const filter = Faker.datatype.string();
    const sampling = randomArrayItem(Object.values(SAMPLINGS));
    const parameters = [randomArrayItem(Object.values(ALARM_METRIC_PARAMETERS))];

    getQueryById.mockReturnValue(() => ({
      filter,
      sampling,
      parameters,
      interval: {
        from: QUICK_RANGES.last30Days.start,
        to: QUICK_RANGES.last30Days.stop,
      },
    }));

    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        activeViewModule,
        alarmModule,
        serviceModule,
        queryModule,
      ]),
      propsData: {
        widget,
      },
    });

    await wrapper.vm.fetchList();

    expect(metrics.fetchVectorMetricsList).toHaveBeenCalledWith({
      widgetId: widget._id,
      params: {
        widget_filters: [filter],
        sampling,
        parameters,
        from: 1383865200,
        to: 1386370800,
      },
    });
  });
});
