import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import {
  createActiveViewModule,
  createAlarmModule,
  createAuthModule,
  createMockedStoreModules,
  createQueryModule,
  createServiceModule,
  createUserPreferenceModule,
  createAggregatedMetricsModule,
} from '@unit/utils/store';
import { randomArrayItem } from '@unit/utils/array';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { AGGREGATE_FUNCTIONS, ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS, WIDGET_TYPES } from '@/constants';

import PieChartWidget from '@/components/widgets/chart/pie-chart-widget.vue';

const stubs = {
  'chart-widget-filters': true,
};

describe('pie-chart-widget', () => {
  const nowTimestamp = 1386435500000;
  mockDateNow(nowTimestamp);

  const widgetId = 'widget-id';

  const { authModule } = createAuthModule();
  const { activeViewModule } = createActiveViewModule();
  const { alarmModule } = createAlarmModule();
  const { userPreferenceModule, fetchUserPreference } = createUserPreferenceModule();
  const { serviceModule } = createServiceModule();
  const { queryModule, updateQuery, getQueryById } = createQueryModule();
  const { aggregatedMetricsModule, fetchAggregatedMetricsList } = createAggregatedMetricsModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    activeViewModule,
    alarmModule,
    serviceModule,
    queryModule,
    aggregatedMetricsModule,
  ]);

  const widget = {
    _id: widgetId,
    parameters: {
      default_sampling: SAMPLINGS.month,
      default_time_range: QUICK_RANGES.last7Days.value,
      aggregate_func: AGGREGATE_FUNCTIONS.avg,
      metrics: [
        {
          metric: ALARM_METRIC_PARAMETERS.createdAlarms,
        },
      ],
    },
  };

  const factory = generateShallowRenderer(PieChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(PieChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
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

    expect(fetchUserPreference).toBeCalledWith(
      expect.any(Object),
      { id: widget._id },
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          search: '',
        },
      },
      undefined,
    );
  });

  test('Vectors metrics fetched with correct query', async () => {
    const filter = Faker.datatype.string();
    const sampling = randomArrayItem(Object.values(SAMPLINGS));

    getQueryById.mockReturnValueOnce(() => ({
      filter,
      sampling,
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
        aggregatedMetricsModule,
      ]),
      propsData: {
        widget,
      },
    });

    await wrapper.vm.fetchList();

    expect(fetchAggregatedMetricsList).toBeCalledWith(
      expect.any(Object),
      {
        widgetId: widget._id,
        params: {
          filter,
          sampling,
          from: 1383843500,
          to: 1386435500,
          parameters: [{
            metric: ALARM_METRIC_PARAMETERS.createdAlarms,
            aggregate_func: AGGREGATE_FUNCTIONS.avg,
          }],
        },
      },
      undefined,
    );
  });

  test('Renders `pie-chart-widget` with required props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: {
          _id: 'pie-chart-widget-id',
          type: WIDGET_TYPES.pieChart,
          title: 'Default pie chart',
          parameters: {
            default_sampling: SAMPLINGS.day,
            default_time_range: QUICK_RANGES.last7Days.value,
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pie-chart-widget` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: {
          _id: 'pie-chart-widget-id',
          type: WIDGET_TYPES.pieChart,
          title: 'Default pie chart',
          parameters: {
            default_sampling: SAMPLINGS.month,
            default_time_range: QUICK_RANGES.last7Days.value,
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
