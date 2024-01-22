import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { getWidgetRequestWithNewProperty, getWidgetRequestWithNewParametersProperty } from '@unit/utils/settings';

import {
  AGGREGATE_FUNCTIONS,
  ALARM_METRIC_PARAMETERS,
  KPI_PIE_CHART_SHOW_MODS,
  QUICK_RANGES,
  SAMPLINGS,
  WIDGET_TYPES,
} from '@/constants';

import { widgetToForm, formToWidget } from '@/helpers/entities/widget/form';

import PieChartWidgetForm from '@/components/sidebars/chart/form/pie-chart-widget-form.vue';

const stubs = {
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-pie-show-mode': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-alarm-metric-aggregate-function': true,
  'field-filters': true,
};

const snapshotStubs = {
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-pie-show-mode': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-alarm-metric-aggregate-function': true,
  'field-filters': true,
};

const selectFieldTitle = wrapper => wrapper.findAll('field-title-stub').at(0);
const selectFieldPeriodicRefresh = wrapper => wrapper.find('field-periodic-refresh-stub');
const selectFieldPreset = wrapper => wrapper.find('field-preset-stub');
const selectFieldAlarmMetricPresets = wrapper => wrapper.find('field-alarm-metric-presets-stub');
const selectFieldPieShowModePresets = wrapper => wrapper.find('field-pie-show-mode-stub');
const selectFieldChartTitle = wrapper => wrapper.findAll('field-title-stub').at(1);
const selectFieldQuickDateIntervalType = wrapper => wrapper.find('field-quick-date-interval-type-stub');
const selectFieldSampling = wrapper => wrapper.find('field-sampling-stub');
const selectAlarmMetricAggregateFunction = wrapper => wrapper.find('field-alarm-metric-aggregate-function-stub');
const selectFieldFilters = wrapper => wrapper.find('field-filters-stub');

describe('pie-chart-widget-form', () => {
  const form = formToWidget(widgetToForm({ type: WIDGET_TYPES.pieChart }));
  const factory = generateShallowRenderer(PieChartWidgetForm, { stubs });
  const snapshotFactory = generateRenderer(PieChartWidgetForm, { stubs: snapshotStubs });

  test('Title changed after trigger field title', async () => {
    const newTitle = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const fieldTitle = selectFieldTitle(wrapper);

    fieldTitle.triggerCustomEvent('input', newTitle);

    expect(wrapper).toEmit('input', getWidgetRequestWithNewProperty(form, 'title', newTitle));
  });

  test('Periodic refresh changed after trigger field periodic refresh', async () => {
    const wrapper = factory({
      propsData: {
        form,
        withPeriodicRefresh: true,
      },
    });

    const fieldPeriodicRefresh = selectFieldPeriodicRefresh(wrapper);

    const periodicRefresh = {
      enabled: Faker.datatype.boolean(),
      value: Faker.datatype.number(),
      unit: Faker.datatype.string(),
    };

    fieldPeriodicRefresh.triggerCustomEvent('input', {
      ...wrapper.vm.form.parameters,
      periodic_refresh: periodicRefresh,
    });

    expect(wrapper).toEmit('input', getWidgetRequestWithNewParametersProperty(form, 'periodic_refresh', periodicRefresh));
  });

  test('Preset changed after trigger preset field', async () => {
    const wrapper = factory({
      propsData: {
        form,
        withPreset: true,
      },
    });

    const newParameters = {
      ...form.parameters,
      stacked: !form.parameters.stacked,
      chart_title: Faker.datatype.string(),
    };

    selectFieldPreset(wrapper).triggerCustomEvent('input', newParameters);
    expect(wrapper).toEmit('input', getWidgetRequestWithNewProperty(form, 'parameters', newParameters));
  });

  test('Metrics changed after trigger field alarm metric presets', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newMetrics = [
      {
        metric: ALARM_METRIC_PARAMETERS.averageResolve,
      },
      {
        metric: ALARM_METRIC_PARAMETERS.createdAlarms,
      },
    ];

    selectFieldAlarmMetricPresets(wrapper).triggerCustomEvent('input', newMetrics);
    expect(wrapper).toEmit('input', getWidgetRequestWithNewParametersProperty(form, 'metrics', newMetrics));
  });

  test('Show mode changed after trigger field pie show mode', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectFieldPieShowModePresets(wrapper).triggerCustomEvent('input', KPI_PIE_CHART_SHOW_MODS.percent);
    expect(wrapper).toEmit(
      'input',
      getWidgetRequestWithNewParametersProperty(form, 'show_mode', KPI_PIE_CHART_SHOW_MODS.percent),
    );
  });

  test('Chart title changed after trigger field chart title', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newChartTitle = Faker.datatype.string();

    selectFieldChartTitle(wrapper).triggerCustomEvent('input', newChartTitle);
    expect(wrapper).toEmit('input', getWidgetRequestWithNewParametersProperty(form, 'chart_title', newChartTitle));
  });

  test('Filters changed after trigger field filters', async () => {
    const wrapper = factory({
      propsData: {
        form,
        withFilters: true,
      },
    });

    const filters = [Faker.datatype.string()];

    selectFieldFilters(wrapper).triggerCustomEvent('update:filters', filters);
    expect(wrapper).toEmit('input', getWidgetRequestWithNewProperty(form, 'filters', filters));
  });

  test('Quick date interval type changed after trigger field quick date interval', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectFieldQuickDateIntervalType(wrapper).triggerCustomEvent('input', QUICK_RANGES.last30Days.value);
    expect(wrapper).toEmit(
      'input',
      getWidgetRequestWithNewParametersProperty(form, 'default_time_range', QUICK_RANGES.last30Days.value),
    );
  });

  test('Sampling changed after trigger field sampling', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectFieldSampling(wrapper).triggerCustomEvent('input', SAMPLINGS.month);
    expect(wrapper).toEmit('input', getWidgetRequestWithNewParametersProperty(form, 'default_sampling', SAMPLINGS.month));
  });

  test('Aggregate function changed after trigger field aggregate function', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectAlarmMetricAggregateFunction(wrapper).triggerCustomEvent('input', AGGREGATE_FUNCTIONS.sum);
    expect(wrapper).toEmit(
      'input',
      getWidgetRequestWithNewParametersProperty(form, 'aggregate_func', AGGREGATE_FUNCTIONS.sum),
    );
  });

  test('Renders `pie-chart-widget-form` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pie-chart-widget-form` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          _id: '_id',
          type: WIDGET_TYPES.pieChart,
          title: 'Map widget',
          filters: [{ title: 'Filter' }],
          parameters: {
            periodic_refresh: {},
            metrics: [{ metric: ALARM_METRIC_PARAMETERS.averageResolve }],
            show_mode: KPI_PIE_CHART_SHOW_MODS.numbers,
            chart_title: 'Chart title',
            default_time_range: QUICK_RANGES.last30Days.value,
            default_sampling: SAMPLINGS.day,
            aggregate_func: AGGREGATE_FUNCTIONS.avg,
          },
        },
        withPeriodicRefresh: true,
        withFilters: true,
        withPreset: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
