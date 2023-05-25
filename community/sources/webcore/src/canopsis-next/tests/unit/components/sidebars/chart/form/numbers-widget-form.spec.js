import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  getWidgetRequestWithNewProperty,
  getWidgetRequestWithNewParametersProperty,
} from '@unit/utils/settings';
import { expectsOneInput } from '@unit/utils/form';

import {
  ALARM_METRIC_PARAMETERS,
  QUICK_RANGES,
  SAMPLINGS,
  WIDGET_TYPES,
} from '@/constants';

import { widgetToForm, formToWidget } from '@/helpers/forms/widgets/common';

import NumbersWidgetForm from '@/components/sidebars/chart/form/numbers-widget-form.vue';

const stubs = {
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-filters': true,
  'field-switcher': true,
  'field-font-size': true,
};

const snapshotStubs = {
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-filters': true,
  'field-switcher': true,
  'field-font-size': true,
};

const selectFieldTitle = wrapper => wrapper.findAll('field-title-stub').at(0);
const selectFieldPeriodicRefresh = wrapper => wrapper.find('field-periodic-refresh-stub');
const selectFieldPreset = wrapper => wrapper.find('field-preset-stub');
const selectFieldAlarmMetricPresets = wrapper => wrapper.find('field-alarm-metric-presets-stub');
const selectFieldChartTitle = wrapper => wrapper.findAll('field-title-stub').at(1);
const selectFieldQuickDateIntervalType = wrapper => wrapper.find('field-quick-date-interval-type-stub');
const selectFieldSampling = wrapper => wrapper.find('field-sampling-stub');
const selectFieldFilters = wrapper => wrapper.find('field-filters-stub');
const selectFieldSwitcher = wrapper => wrapper.find('field-switcher-stub');
const selectFieldFontSize = wrapper => wrapper.find('field-font-size-stub');

describe('numbers-widget-form', () => {
  const form = formToWidget(widgetToForm({ type: WIDGET_TYPES.numbers }));
  const factory = generateShallowRenderer(NumbersWidgetForm, { stubs });
  const snapshotFactory = generateRenderer(NumbersWidgetForm, { stubs: snapshotStubs });

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

    fieldPeriodicRefresh.vm.$emit('input', periodicRefresh);

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

  test('Trend changed after trigger field switcher', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newShowTrend = !form.parameters.show_trend;

    selectFieldSwitcher(wrapper).vm.$emit('input', newShowTrend);
    expectsOneInput(wrapper, getWidgetRequestWithNewParametersProperty(form, 'show_trend', newShowTrend));
  });

  test('Font size changed after trigger field font size', async () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newFontSize = Faker.datatype.number();

    selectFieldFontSize(wrapper).vm.$emit('input', newFontSize);
    expectsOneInput(wrapper, getWidgetRequestWithNewParametersProperty(form, 'font_size', newFontSize));
  });

  test('Renders `numbers-widget-form` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `numbers-widget-form` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          _id: '_id',
          type: WIDGET_TYPES.numbers,
          title: 'Map widget',
          filters: [{ title: 'Filter' }],
          parameters: {
            periodic_refresh: {},
            metrics: [{ metric: ALARM_METRIC_PARAMETERS.averageResolve }],
            chart_title: 'Chart title',
            default_time_range: QUICK_RANGES.last30Days.value,
            default_sampling: SAMPLINGS.day,
            show_trend: true,
          },
        },
        withPeriodicRefresh: true,
        withFilters: true,
        withPreset: true,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
