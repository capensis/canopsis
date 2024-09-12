import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { getWidgetRequestWithNewProperty, getWidgetRequestWithNewParametersProperty } from '@unit/utils/settings';
import { expectsOneInput } from '@unit/utils/form';

import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS, WIDGET_TYPES } from '@/constants';

import { widgetToForm, formToWidget } from '@/helpers/entities/widget/form';

import BarChartWidgetForm from '@/components/sidebars/chart/form/bar-chart-widget-form.vue';

const stubs = {
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-bar-graph-type': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-filters': true,
  'field-switcher': true,
};

const snapshotStubs = {
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-bar-graph-type': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-filters': true,
  'field-switcher': true,
};

const selectFieldTitle = wrapper => wrapper.findAll('field-title-stub').at(0);
const selectFieldPeriodicRefresh = wrapper => wrapper.find('field-periodic-refresh-stub');
const selectFieldPreset = wrapper => wrapper.find('field-preset-stub');
const selectFieldAlarmMetricPresets = wrapper => wrapper.find('field-alarm-metric-presets-stub');
const selectFieldBarGraphType = wrapper => wrapper.find('field-bar-graph-type-stub');
const selectFieldChartTitle = wrapper => wrapper.findAll('field-title-stub').at(1);
const selectFieldQuickDateIntervalType = wrapper => wrapper.find('field-quick-date-interval-type-stub');
const selectFieldSampling = wrapper => wrapper.find('field-sampling-stub');
const selectFieldFilters = wrapper => wrapper.find('field-filters-stub');
const selectFieldSwitcher = wrapper => wrapper.find('field-switcher-stub');

describe('bar-chart-widget-form', () => {
  const form = formToWidget(widgetToForm({ type: WIDGET_TYPES.barChart }));
  const factory = generateShallowRenderer(BarChartWidgetForm, { stubs });
  const snapshotFactory = generateRenderer(BarChartWidgetForm, { stubs: snapshotStubs });

  test('Title changed after trigger field title', async () => {
    const newTitle = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const fieldTitle = selectFieldTitle(wrapper);

    fieldTitle.vm.$emit('input', newTitle);

    expectsOneInput(wrapper, getWidgetRequestWithNewProperty(form, 'title', newTitle));
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

    fieldPeriodicRefresh.vm.$emit('input', {
      ...wrapper.vm.form.parameters,
      periodic_refresh: periodicRefresh,
    });

    expectsOneInput(wrapper, getWidgetRequestWithNewParametersProperty(form, 'periodic_refresh', periodicRefresh));
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

    selectFieldPreset(wrapper).vm.$emit('input', newParameters);
    expectsOneInput(wrapper, getWidgetRequestWithNewProperty(form, 'parameters', newParameters));
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

    selectFieldAlarmMetricPresets(wrapper).vm.$emit('input', newMetrics);
    expectsOneInput(wrapper, getWidgetRequestWithNewParametersProperty(form, 'metrics', newMetrics));
  });

  test('Graph type changed after trigger field bar graph type', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });
    const newStacked = true;

    selectFieldBarGraphType(wrapper).vm.$emit('input', newStacked);
    expectsOneInput(wrapper, getWidgetRequestWithNewParametersProperty(form, 'stacked', newStacked));
  });

  test('Chart title changed after trigger field chart title', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newChartTitle = Faker.datatype.string();

    selectFieldChartTitle(wrapper).vm.$emit('input', newChartTitle);
    expectsOneInput(wrapper, getWidgetRequestWithNewParametersProperty(form, 'chart_title', newChartTitle));
  });

  test('Filters changed after trigger field filters', async () => {
    const wrapper = factory({
      propsData: {
        form,
        withFilters: true,
      },
    });

    const filters = [Faker.datatype.string()];

    selectFieldFilters(wrapper).vm.$emit('update:filters', filters);
    expectsOneInput(wrapper, getWidgetRequestWithNewProperty(form, 'filters', filters));
  });

  test('Quick date interval type changed after trigger field quick date interval', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectFieldQuickDateIntervalType(wrapper).vm.$emit('input', QUICK_RANGES.last30Days.value);
    expectsOneInput(
      wrapper,
      getWidgetRequestWithNewParametersProperty(form, 'default_time_range', QUICK_RANGES.last30Days.value),
    );
  });

  test('Sampling changed after trigger field sampling', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectFieldSampling(wrapper).vm.$emit('input', SAMPLINGS.month);
    expectsOneInput(wrapper, getWidgetRequestWithNewParametersProperty(form, 'default_sampling', SAMPLINGS.month));
  });

  test('Comparison changed after trigger field comparison', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newComparison = !form.parameters.comparison;

    selectFieldSwitcher(wrapper).vm.$emit('input', newComparison);
    expectsOneInput(wrapper, getWidgetRequestWithNewParametersProperty(form, 'comparison', newComparison));
  });

  test('Renders `bar-chart-widget-form` widget settings with default props', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `bar-chart-widget-form` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          _id: '_id',
          type: WIDGET_TYPES.barChart,
          title: 'Bar chart',
          filters: [{ title: 'Filter' }],
          parameters: {
            periodic_refresh: {},
            metrics: [{ metric: ALARM_METRIC_PARAMETERS.averageResolve }],
            stacked: true,
            chart_title: 'Chart title',
            default_time_range: QUICK_RANGES.last30Days.value,
            default_sampling: SAMPLINGS.day,
            comparison: false,
          },
        },
        withFilters: true,
        withPeriodicRefresh: true,
        withPreset: true,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
