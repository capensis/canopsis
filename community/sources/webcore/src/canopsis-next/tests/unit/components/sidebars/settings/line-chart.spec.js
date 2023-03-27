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
  ALARM_METRIC_PARAMETERS,
  QUICK_RANGES,
  SAMPLINGS,
  SIDE_BARS,
  USERS_PERMISSIONS,
  WIDGET_TYPES,
} from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/forms/widgets/common';

import LineChartSettings from '@/components/sidebars/settings/line-chart.vue';

const stubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-preset': true,
  'field-alarm-metric-presets': true,
  'field-chart-title': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-filters': true,
  'field-switcher': true,
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
  'field-chart-title': true,
  'field-quick-date-interval-type': true,
  'field-sampling': true,
  'field-filters': true,
  'field-switcher': true,
};

const generateDefaultLineChartWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.lineChart })),

  _id: Faker.datatype.string(),
});
const selectFieldTitle = wrapper => wrapper.find('field-title-stub');
const selectFieldPeriodicRefresh = wrapper => wrapper.find('field-periodic-refresh-stub');
const selectFieldAlarmMetricPresets = wrapper => wrapper.find('field-alarm-metric-presets-stub');
const selectFieldChartTitle = wrapper => wrapper.find('field-chart-title-stub');
const selectFieldQuickDateIntervalType = wrapper => wrapper.find('field-quick-date-interval-type-stub');
const selectFieldSampling = wrapper => wrapper.find('field-sampling-stub');
const selectFieldFilters = wrapper => wrapper.find('field-filters-stub');
const selectFieldSwitcher = wrapper => wrapper.find('field-switcher-stub');

describe('line-chart', () => {
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
    ...generateDefaultLineChartWidget(),

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

  const factory = generateShallowRenderer(LineChartSettings, {
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
  const snapshotFactory = generateRenderer(LineChartSettings, {
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
    const localWidget = getEmptyWidgetByType(WIDGET_TYPES.lineChart);

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
      [USERS_PERMISSIONS.business.lineChart.actions.listFilters]: { actions: [] },
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

  test('Comparison changed after trigger field comparison', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const newComparison = !widget.parameters.comparison;

    selectFieldSwitcher(wrapper).vm.$emit('input', newComparison);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'comparison', newComparison),
      },
    });
  });

  test('Renders `line-chart` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `line-chart` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            ...sidebar.config,

            widget: {
              _id: '_id',
              type: WIDGET_TYPES.lineChart,
              title: 'Map widget',
              filters: [{ title: 'Filter' }],
              parameters: {
                periodic_refresh: {},
                metrics: [{ metric: ALARM_METRIC_PARAMETERS.averageResolve }],
                chart_title: 'Chart title',
                default_time_range: QUICK_RANGES.last30Days.value,
                default_sampling: SAMPLINGS.day,
                comparison: false,
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
