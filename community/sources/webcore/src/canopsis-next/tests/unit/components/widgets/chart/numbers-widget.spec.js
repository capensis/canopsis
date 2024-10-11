import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

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
import {
  AGGREGATE_FUNCTIONS,
  ALARM_METRIC_PARAMETERS,
  QUICK_RANGES,
  SAMPLINGS,
  WIDGET_TYPES,
} from '@/constants';

import NumbersWidget from '@/components/widgets/chart/numbers-widget.vue';

const stubs = {
  'kpi-widget-filters': true,
};

describe('numbers-widget', () => {
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
      metrics: [
        {
          metric: ALARM_METRIC_PARAMETERS.createdAlarms,
          aggregate_func: AGGREGATE_FUNCTIONS.sum,
        },
      ],
    },
  };

  const factory = generateShallowRenderer(NumbersWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(NumbersWidget, {
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
    const parameters = [{
      metric: ALARM_METRIC_PARAMETERS.createdAlarms,
      aggregate_func: AGGREGATE_FUNCTIONS.sum,
    }];

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
          widget_filters: [filter],
          sampling,
          from: 1383865200,
          to: 1386370800,
          parameters,
        },
      },
      undefined,
    );
  });

  test('Renders `numbers-widget` with required props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: {
          _id: 'numbers-widget-id',
          type: WIDGET_TYPES.numbers,
          title: 'Default numbers',
          parameters: {
            default_sampling: SAMPLINGS.day,
            default_time_range: QUICK_RANGES.last7Days.value,
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `numbers-widget` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: {
          _id: 'numbers-widget-id',
          type: WIDGET_TYPES.numbers,
          title: 'Default numbers',
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
