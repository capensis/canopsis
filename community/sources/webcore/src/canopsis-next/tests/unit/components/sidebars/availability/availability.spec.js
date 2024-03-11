import { omit } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockSidebar } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import { createSettingsMocks, getWidgetRequestWithNewProperty, submitWithExpects } from '@unit/utils/settings';

import {
  AVAILABILITY_DISPLAY_PARAMETERS,
  AVAILABILITY_SHOW_TYPE,
  CUSTOM_WIDGET_TEMPLATE,
  ENTITY_FIELDS,
  QUICK_RANGES,
  SIDE_BARS,
  USERS_PERMISSIONS,
  WIDGET_TYPES,
} from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/entities/widget/form';

import AvailabilitySettings from '@/components/sidebars/availability/availability.vue';

const stubs = {
  'widget-settings': true,
  'availability-settings-form': true,
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'widget-settings': true,
  'availability-form': true,
};

const generateDefaultAvailabilityWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.availability })),

  _id: `${WIDGET_TYPES.availability}_id`,
});
const selectWidgetForm = wrapper => wrapper.find('availability-form-stub');

describe('availability-settings', () => {
  const nowTimestamp = 1386435600000;
  jest.useFakeTimers({ now: nowTimestamp });

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
    ...generateDefaultAvailabilityWidget(),

    tab: Faker.datatype.string(),
  };

  const sidebar = {
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
  const snapshotFactory = generateRenderer(AvailabilitySettings, {
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
        widget_columns: [],
        widget_columns_template: CUSTOM_WIDGET_TEMPLATE,
        active_alarms_columns: [],
        active_alarms_columns_template: CUSTOM_WIDGET_TEMPLATE,
        resolved_alarms_columns: [],
        resolved_alarms_columns_template: CUSTOM_WIDGET_TEMPLATE,
        default_time_range: QUICK_RANGES.last30Days.value,
        default_display_parameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
        default_show_type: AVAILABILITY_SHOW_TYPE.duration,
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

    selectWidgetForm(wrapper).triggerCustomEvent('input', updatedWidget);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(
          updatedWidget,
          'parameters',
          {
            ...updatedWidget.parameters,
            active_alarms_columns_template: '',
            resolved_alarms_columns_template: '',
            widget_columns_template: '',
          },
        ),
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

    const form = selectWidgetForm(wrapper);

    const templateId = Faker.datatype.string();
    const columns = [{
      label: Faker.datatype.string(),
      column: ENTITY_FIELDS.id,
    }];

    await form.triggerCustomEvent('update:widget-columns-template', templateId, columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(
          widget,
          'parameters',
          {
            ...widget.parameters,
            widget_columns_template: templateId,
            widget_columns: columns
              .map(({ column, ...rest }) => ({ ...rest, value: column })),
          },
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

    const form = selectWidgetForm(wrapper);

    const templateId = Faker.datatype.string();
    const columns = [{
      label: Faker.datatype.string(),
      column: ENTITY_FIELDS.id,
    }];

    await form.triggerCustomEvent('update:active-alarms-columns-template', templateId, columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(
          widget,
          'parameters',
          {
            ...widget.parameters,
            active_alarms_columns_template: templateId,
            active_alarms_columns: columns
              .map(({ column, ...rest }) => ({ ...rest, value: column })),
          },
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

    const form = selectWidgetForm(wrapper);

    const templateId = Faker.datatype.string();
    const columns = [{
      label: Faker.datatype.string(),
      column: ENTITY_FIELDS.id,
    }];

    await form.triggerCustomEvent('update:resolved-alarms-columns-template', templateId, columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(
          widget,
          'parameters',
          {
            ...widget.parameters,
            resolved_alarms_columns_template: templateId,
            resolved_alarms_columns: columns
              .map(({ column, ...rest }) => ({ ...rest, value: column })),
          },
        ),
      },
    });
  });

  test('Renders `availability-settings` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-settings` widget settings with custom props and permissions', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.availability.actions.listFilters]: { actions: [] },
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
              type: WIDGET_TYPES.availability,
              title: 'Abailability widget',
              filters: [{ title: 'Filter' }],
              parameters: {
                periodic_refresh: {},
                widget_columns: [{}],
                widget_columns_template: CUSTOM_WIDGET_TEMPLATE,
                active_alarms_columns: [{}, {}],
                active_alarms_columns_template: CUSTOM_WIDGET_TEMPLATE,
                default_time_range: QUICK_RANGES.last6Months.value,
                default_display_parameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
                default_show_type: AVAILABILITY_SHOW_TYPE.duration,
              },
            },
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
