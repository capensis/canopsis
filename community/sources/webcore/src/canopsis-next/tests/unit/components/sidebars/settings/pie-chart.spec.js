import { omit } from 'lodash';
import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockDateNow, mockSidebar } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import {
  createSettingsMocks,
  getWidgetRequestWithNewProperty,
  getWidgetRequestWithNewParametersProperty,
  submitWithExpects,
} from '@unit/utils/settings';

import {
  AGGREGATE_FUNCTIONS,
  ALARM_METRIC_PARAMETERS,
  KPI_PIE_CHART_SHOW_MODS,
  QUICK_RANGES,
  SAMPLINGS,
  SIDE_BARS,
  USERS_PERMISSIONS,
  WIDGET_TYPES,
} from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/forms/widgets/common';

import PieChartSettings from '@/components/sidebars/settings/pie-chart.vue';

const stubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-pie-show-mode': true,
  'field-chart-title': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-alarm-metric-aggregate-function': true,
  'field-filters': true,
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-pie-show-mode': true,
  'field-chart-title': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-alarm-metric-aggregate-function': true,
  'field-filters': true,
};

const generateDefaultPieChartWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.pieChart })),

  _id: Faker.datatype.string(),
});
const selectFieldTitle = wrapper => wrapper.find('field-title-stub');
const selectFieldPeriodicRefresh = wrapper => wrapper.find('field-periodic-refresh-stub');
const selectFieldPreset = wrapper => wrapper.find('field-preset-stub');
const selectFieldAlarmMetricPresets = wrapper => wrapper.find('field-alarm-metric-presets-stub');
const selectFieldPieShowModePresets = wrapper => wrapper.find('field-pie-show-mode-stub');
const selectFieldChartTitle = wrapper => wrapper.find('field-chart-title-stub');
const selectFieldQuickDateIntervalType = wrapper => wrapper.find('field-quick-date-interval-type-stub');
const selectFieldSampling = wrapper => wrapper.find('field-sampling-stub');
const selectAlarmMetricAggregateFunction = wrapper => wrapper.find('field-alarm-metric-aggregate-function-stub');
const selectFieldFilters = wrapper => wrapper.find('field-filters-stub');

describe('pie-chart', () => {
  const nowTimestamp = 1386435600000;

  mockDateNow(nowTimestamp);
  const $sidebar = mockSidebar();

  const {
    createWidget,
    updateWidget,
    copyWidget,
    fetchActiveView,
    currentUserPermissionsById,
    activeViewModule,
    widgetModule,
    authModule,
    userPreferenceModule,
    serviceModule,
    widgetTemplateModule,
    infosModule,
  } = createSettingsMocks();

  const widget = {
    ...generateDefaultPieChartWidget(),

    tab: Faker.datatype.string(),
  };

  const sidebar = {
    name: SIDE_BARS.mapSettings,
    config: {
      widget,
    },
    hidden: false,
  };

  const store = createMockedStoreModules([
    userPreferenceModule,
    activeViewModule,
    serviceModule,
    widgetModule,
    authModule,
    widgetTemplateModule,
    infosModule,
  ]);

  const factory = generateShallowRenderer(PieChartSettings, {
    stubs,

    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
    mocks: {
      $sidebar,
    },
  });
  const snapshotFactory = generateRenderer(PieChartSettings, {
    stubs: snapshotStubs,

    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
    mocks: {
      $sidebar,
    },
  });

  test('Create widget with default parameters', async () => {
    const localWidget = getEmptyWidgetByType(WIDGET_TYPES.pieChart);

    localWidget.tab = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            widget: localWidget,
          },
        },
      },
    });

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: createWidget,
      expectData: {
        data: {
          ...formToWidget(widgetToForm(localWidget)),

          tab: localWidget.tab,
        },
      },
    });
  });

  test('Duplicate widget without changes', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            widget,
            duplicate: true,
          },
        },
      },
    });

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: copyWidget,
      expectData: {
        id: widget._id,
        data: omit(widget, ['_id']),
      },
    });
  });

  test('Title changed after trigger field title', async () => {
    const newTitle = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const fieldTitle = selectFieldTitle(wrapper);

    fieldTitle.vm.$emit('input', newTitle);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(widget, 'title', newTitle),
      },
    });
  });

  test('Periodic refresh changed after trigger field periodic refresh', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const fieldPeriodicRefresh = selectFieldPeriodicRefresh(wrapper);

    const periodicRefresh = {
      enabled: Faker.datatype.boolean(),
      value: Faker.datatype.number(),
      unit: Faker.datatype.string(),
    };

    fieldPeriodicRefresh.vm.$emit('input', periodicRefresh);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'periodic_refresh', periodicRefresh),
      },
    });
  });

  test('Preset changed after trigger preset field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const newParameters = {
      ...widget.parameters,
      stacked: !widget.parameters.stacked,
      chart_title: Faker.datatype.string(),
    };

    selectFieldPreset(wrapper).vm.$emit('input', newParameters);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(widget, 'parameters', newParameters),
      },
    });
  });

  test('Metrics changed after trigger field alarm metric presets', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
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

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'metrics', newMetrics),
      },
    });
  });

  test('Show mode changed after trigger field pie show mode', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    selectFieldPieShowModePresets(wrapper).vm.$emit('input', KPI_PIE_CHART_SHOW_MODS.percent);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'show_mode', KPI_PIE_CHART_SHOW_MODS.percent),
      },
    });
  });

  test('Chart title changed after trigger field chart title', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const newChartTitle = Faker.datatype.string();

    selectFieldChartTitle(wrapper).vm.$emit('input', newChartTitle);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'chart_title', newChartTitle),
      },
    });
  });

  test('Filters changed after trigger field filters', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.pieChart.actions.listFilters]: { actions: [] },
    });
    const wrapper = factory({
      propsData: {
        sidebar,
      },
      store: createMockedStoreModules([
        userPreferenceModule,
        activeViewModule,
        serviceModule,
        widgetModule,
        authModule,
        widgetTemplateModule,
        infosModule,
      ]),
    });

    const filters = [Faker.datatype.string()];

    selectFieldFilters(wrapper).vm.$emit('update:filters', filters);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(widget, 'filters', filters),
      },
    });
  });

  test('Quick date interval type changed after trigger field quick date interval', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    selectFieldQuickDateIntervalType(wrapper).vm.$emit('input', QUICK_RANGES.last30Days.value);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'default_time_range', QUICK_RANGES.last30Days.value),
      },
    });
  });

  test('Sampling changed after trigger field sampling', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    selectFieldSampling(wrapper).vm.$emit('input', SAMPLINGS.month);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'default_sampling', SAMPLINGS.month),
      },
    });
  });

  test('Aggregate function changed after trigger field aggregate function', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    selectAlarmMetricAggregateFunction(wrapper).vm.$emit('input', AGGREGATE_FUNCTIONS.sum);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'aggregate_func', AGGREGATE_FUNCTIONS.sum),
      },
    });
  });

  test('Renders `pie-chart` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pie-chart` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            ...sidebar.config,

            widget: {
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
          },
        },
      },
      mocks: {
        $sidebar,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
