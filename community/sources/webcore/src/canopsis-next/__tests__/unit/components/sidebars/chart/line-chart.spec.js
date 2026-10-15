import { createButtonStub } from '@unit/stubs/button';
import { mockSidebar } from '@unit/utils/mock-hooks';
import { createSettingsMocks, submitWithExpects } from '@unit/utils/settings';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';
import { omit } from 'lodash';

import {
  ALARM_METRIC_PARAMETERS,
  QUICK_RANGES,
  SAMPLINGS,
  SIDE_BARS,
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

const generateDefaultLineChartWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.lineChart })),

  _id: Faker.datatype.string(),
});

const selectWidgetForm = wrapper => wrapper.find('line-chart-widget-form-stub');

describe('line-chart', () => {
  const $sidebar = mockSidebar();

  const {
    createWidget,
    updateWidget,
    fetchActiveView,

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
    id: 'test-sidebar-id',
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

  const timestamp = 1386435600000;

  beforeAll(() => vi.useFakeTimers({ now: timestamp }));

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
});
