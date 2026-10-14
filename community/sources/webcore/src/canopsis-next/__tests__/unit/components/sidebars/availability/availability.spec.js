import { createButtonStub } from '@unit/stubs/button';
import { createInputStub } from '@unit/stubs/input';
import { mockSidebar } from '@unit/utils/mock-hooks';
import {
  createSettingsMocks,
  getWidgetRequestWithNewProperty,
  submitWithExpects,
} from '@unit/utils/settings';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';
import { omit } from 'lodash';

import {
  AVAILABILITY_DISPLAY_PARAMETERS,
  AVAILABILITY_SHOW_TYPE,
  CUSTOM_WIDGET_TEMPLATE,
  ENTITY_FIELDS,
  CSV_SEPARATORS,
  QUICK_RANGES,
  SIDE_BARS,
  WIDGET_TYPES,
} from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/entities/widget/form';

import AvailabilitySettings from '@/components/sidebars/availability/availability.vue';

const availabilityFormStub = createInputStub('availability-form');

const stubs = {
  'widget-settings': true,
  'availability-form': availabilityFormStub,
  'availability-settings-form': true,
  'v-btn': createButtonStub('v-btn'),
};

const generateDefaultAvailabilityWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.availability })),

  _id: `${WIDGET_TYPES.availability}_id`,
});
const withDynamicExportSettingsColumnKeys = exportSettings => ({
  ...exportSettings,
  widgetExportColumns: exportSettings.widgetExportColumns?.map(column => ({
    ...column,
    key: expect.any(String),
  })),
});
const withDynamicExportColumnKeys = widgetData => ({
  ...widgetData,
  parameters: {
    ...widgetData.parameters,
    export_settings: {
      ...withDynamicExportSettingsColumnKeys(widgetData.parameters.export_settings),
    },
  },
});
const withDynamicParametersExportColumnKeys = parameters => ({
  ...parameters,
  export_settings: withDynamicExportSettingsColumnKeys(parameters.export_settings),
});

describe('availability-settings', () => {
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
    ...generateDefaultAvailabilityWidget(),

    tab: Faker.datatype.string(),
  };

  const sidebar = {
    id: 'test-sidebar-id',
    name: SIDE_BARS.availabilitySettings,
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

  const factory = generateShallowRenderer(AvailabilitySettings, {
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
    const localWidget = getEmptyWidgetByType(WIDGET_TYPES.availability);

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
        data: withDynamicExportColumnKeys({
          ...formToWidget(widgetToForm(localWidget)),

          tab: localWidget.tab,
        }),
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
        data: withDynamicExportColumnKeys(omit(widget, ['_id'])),
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
        widget_columns: [],
        widget_columns_template: CUSTOM_WIDGET_TEMPLATE,
        active_alarms_columns: [],
        active_alarms_columns_template: CUSTOM_WIDGET_TEMPLATE,
        resolved_alarms_columns: [],
        resolved_alarms_columns_template: CUSTOM_WIDGET_TEMPLATE,
        default_time_range: QUICK_RANGES.last30Days.value,
        default_display_parameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        default_show_type: AVAILABILITY_SHOW_TYPE.duration,
        export_settings: {
          exportCsvSeparator: CSV_SEPARATORS.space,
          widgetExportColumns: [],
        },
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

    wrapper.vm.form = updatedWidget;

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(updatedWidget, 'parameters', {
          ...updatedWidget.parameters,
          active_alarms_columns_template: '',
          resolved_alarms_columns_template: '',
          widget_columns_template: '',
          export_settings: {
            exportCsvSeparator: CSV_SEPARATORS.space,
            export_csv_separator: CSV_SEPARATORS.space,
            widgetExportColumns: [],
            widget_export_columns: [],
            widget_export_columns_template: undefined,
          },
        }),
      },
    });
  });

  test('Widget columns changed after trigger update:widget-columns-template', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const templateId = Faker.datatype.string();
    const columns = [
      {
        label: Faker.datatype.string(),
        column: ENTITY_FIELDS.id,
      },
    ];

    wrapper.vm.updateWidgetColumnsTemplate(templateId, columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(
          widget,
          'parameters',
          withDynamicParametersExportColumnKeys({
            ...widget.parameters,
            widget_columns_template: templateId,
            widget_columns: columns.map(({ column, ...rest }) => ({ ...rest, value: column })),
          }),
        ),
      },
    });
  });

  test('Active columns changed after trigger update:active-alarms-columns-template', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const templateId = Faker.datatype.string();
    const columns = [
      {
        label: Faker.datatype.string(),
        column: ENTITY_FIELDS.id,
      },
    ];

    wrapper.vm.updateActiveAlarmsColumnsTemplate(templateId, columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(
          widget,
          'parameters',
          withDynamicParametersExportColumnKeys({
            ...widget.parameters,
            active_alarms_columns_template: templateId,
            active_alarms_columns: columns.map(({ column, ...rest }) => ({
              ...rest,
              value: column,
            })),
          }),
        ),
      },
    });
  });

  test('Resolved columns changed after trigger update:resolved-alarms-columns-template', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const templateId = Faker.datatype.string();
    const columns = [
      {
        label: Faker.datatype.string(),
        column: ENTITY_FIELDS.id,
      },
    ];

    wrapper.vm.updateResolvedAlarmsColumnsTemplate(templateId, columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(
          widget,
          'parameters',
          withDynamicParametersExportColumnKeys({
            ...widget.parameters,
            resolved_alarms_columns_template: templateId,
            resolved_alarms_columns: columns.map(({ column, ...rest }) => ({
              ...rest,
              value: column,
            })),
          }),
        ),
      },
    });
  });
});
