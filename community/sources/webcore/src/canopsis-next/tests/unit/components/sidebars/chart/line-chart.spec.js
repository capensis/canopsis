import { omit } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockDateNow, mockSidebar } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createSettingsMocks, submitWithExpects } from '@unit/utils/settings';

import {
  ALARM_METRIC_PARAMETERS,
  QUICK_RANGES,
  SAMPLINGS,
  SIDE_BARS,
  USERS_PERMISSIONS,
  WIDGET_TYPES,
} from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/entities/widget/form';

import LineChartSettings from '@/components/sidebars/chart/line-chart.vue';

const stubs = {
  'widget-settings': true,
  'line-chart-widget-form': true,
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'widget-settings': true,
  'line-chart-widget-form': true,
};

const generateDefaultLineChartWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.lineChart })),

  _id: Faker.datatype.string(),
});

const selectWidgetForm = wrapper => wrapper.find('line-chart-widget-form-stub');

describe('line-chart', () => {
  const nowTimestamp = 1386435600000;

  mockDateNow(nowTimestamp);
  const $sidebar = mockSidebar();

  const {
    createWidget,
    updateWidget,
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
      widgetMethod: createWidget,
      expectData: {
        data: omit(widget, ['_id']),
      },
    });
  });

  test('All fields changed after input trigger', async () => {
    const newFields = {
      title: Faker.datatype.string(),
      filters: [Faker.datatype.string()],
      parameters: {
        periodic_refresh: {
          enabled: Faker.datatype.boolean(),
          value: Faker.datatype.number(),
          unit: Faker.datatype.string(),
        },
        chart_title: Faker.datatype.string(),
        default_time_range: QUICK_RANGES.last30Days.value,
        default_sampling: SAMPLINGS.month,
        comparison: !widget.parameters.comparison,
        metrics: [
          {
            metric: ALARM_METRIC_PARAMETERS.averageResolve,
          },
          {
            metric: ALARM_METRIC_PARAMETERS.createdAlarms,
          },
        ],
      },
    };

    const updatedWidget = {
      ...omit(widget, ['_id']),
      ...newFields,

      parameters: {
        ...widget.parameters,
        ...newFields.parameters,
      },
    };

    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const form = selectWidgetForm(wrapper);

    form.triggerCustomEvent('input', updatedWidget);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: updatedWidget,
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

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `line-chart` widget settings with custom props and permissions', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.lineChart.actions.listFilters]: { actions: [] },
    });

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

    expect(wrapper).toMatchSnapshot();
  });
});
