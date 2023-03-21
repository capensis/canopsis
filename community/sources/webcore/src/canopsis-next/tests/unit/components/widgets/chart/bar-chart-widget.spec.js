import { generateRenderer } from '@unit/utils/vue';
import {
  createActiveViewModule,
  createAlarmModule, createAuthModule, createMockedStoreModules, createQueryModule,
  createServiceModule,
  createUserPreferenceModule, createVectorMetricsModule,
} from '@unit/utils/store';
import { QUICK_RANGES, SAMPLINGS, WIDGET_TYPES } from '@/constants';

import BarChartWidget from '@/components/widgets/chart/bar-chart-widget.vue';

const stubs = {
  'chart-widget-filters': true,
};

describe('bar-chart-widget', () => {
  const { authModule } = createAuthModule();
  const { activeViewModule } = createActiveViewModule();
  const { alarmModule } = createAlarmModule();
  const { userPreferenceModule } = createUserPreferenceModule();
  const { serviceModule } = createServiceModule();
  const { queryModule } = createQueryModule();
  const { vectorMetricsModule } = createVectorMetricsModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    activeViewModule,
    alarmModule,
    serviceModule,
    queryModule,
    vectorMetricsModule,
  ]);

  const snapshotFactory = generateRenderer(BarChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Renders `bar-chart-widget` with required props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: {
          _id: 'bar-chart-widget-id',
          type: WIDGET_TYPES.barChart,
          title: 'Default map',
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
          title: 'Default map',
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
