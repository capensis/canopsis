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
  createVectorMetricsModule,
} from '@unit/utils/store';
import { randomArrayItem } from '@unit/utils/array';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS, WIDGET_TYPES } from '@/constants';

import BarChartWidget from '@/components/widgets/chart/bar-chart-widget.vue';

const stubs = {
  'chart-widget-filters': true,
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
  const { vectorMetricsModule, fetchVectorMetricsList } = createVectorMetricsModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    activeViewModule,
    alarmModule,
    serviceModule,
    queryModule,
    vectorMetricsModule,
  ]);

  const widget = {
    _id: widgetId,
    parameters: {
      default_sampling: SAMPLINGS.month,
      default_time_range: QUICK_RANGES.last7Days.value,
      metrics: [
        { metric: ALARM_METRIC_PARAMETERS.createdAlarms },
      ],
    },
  };

  const factory = generateShallowRenderer(BarChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(BarChartWidget, {
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
    const parameters = [randomArrayItem(Object.values(ALARM_METRIC_PARAMETERS))];

    getQueryById.mockReturnValueOnce(() => ({
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
        vectorMetricsModule,
      ]),
      propsData: {
        widget,
      },
    });

    await wrapper.vm.fetchList();

    expect(fetchVectorMetricsList).toBeCalledWith(
      expect.any(Object),
      {
        widgetId: widget._id,
        params: {
          widget_filters: [filter],
          sampling,
          parameters,
          from: 1383843500,
          to: 1386435500,
        },
      },
      undefined,
    );
  });

  test('Renders `bar-chart-widget` with required props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: {
          _id: 'bar-chart-widget-id',
          type: WIDGET_TYPES.barChart,
          title: 'Default bar chart',
          parameters: {
            default_sampling: SAMPLINGS.day,
            default_time_range: QUICK_RANGES.last7Days.value,
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `bar-chart-widget` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: {
          _id: 'bar-chart-widget-id',
          type: WIDGET_TYPES.barChart,
          title: 'Default bar chart',
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
