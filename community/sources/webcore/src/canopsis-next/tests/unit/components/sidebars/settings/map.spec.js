import { omit } from 'lodash';
import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockDateNow, mockSidebar } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import {
  createSettingsMocks,
  getWidgetRequestWithNewProperty,
  getWidgetRequestWithNewParametersProperty,
  submitWithExpects,
} from '@unit/utils/settings';

import { COLOR_INDICATOR_TYPES, SIDE_BARS, USERS_PERMISSIONS, WIDGET_TYPES } from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/forms/widgets/common';
import { formToWidgetColumns, widgetColumnToForm } from '@/helpers/forms/shared/widget-column';

import MapSettings from '@/components/sidebars/settings/map.vue';

const localVue = createVueInstance();

const stubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-map': true,
  'field-switcher': true,
  'field-color-indicator': true,
  'field-filters': true,
  'field-text-editor': true,
  'field-columns': true,
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-map': true,
  'field-switcher': true,
  'field-color-indicator': true,
  'field-filters': true,
  'field-text-editor': true,
  'field-columns': true,
};

const factory = (options = {}) => shallowMount(MapSettings, {
  localVue,
  stubs,

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(MapSettings, {
  localVue,
  stubs: snapshotStubs,

  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const generateDefaultMapWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.map })),

  _id: Faker.datatype.string(),
});

const selectFieldTitle = wrapper => wrapper.find('field-title-stub');
const selectFieldPeriodicRefresh = wrapper => wrapper.find('field-periodic-refresh-stub');
const selectFieldMap = wrapper => wrapper.find('field-map-stub');
const selectFieldSwitcher = wrapper => wrapper.find('field-switcher-stub');
const selectFieldColorIndicator = wrapper => wrapper.find('field-color-indicator-stub');
const selectFieldFilters = wrapper => wrapper.find('field-filters-stub');
const selectFieldTextEditor = wrapper => wrapper.find('field-text-editor-stub');
const selectFieldColumns = wrapper => wrapper.findAll('field-columns-stub');
const selectAlarmsColumns = wrapper => selectFieldColumns(wrapper).at(0);
const selectEntitiesColumns = wrapper => selectFieldColumns(wrapper).at(1);

describe('map', () => {
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
    dynamicInfoModule,
  } = createSettingsMocks();

  const widget = {
    ...generateDefaultMapWidget(),

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
    dynamicInfoModule,
  ]);

  afterEach(() => {
    createWidget.mockReset();
    updateWidget.mockReset();
    copyWidget.mockReset();
    fetchActiveView.mockReset();
  });

  test('Create widget with default parameters', async () => {
    const localWidget = getEmptyWidgetByType(WIDGET_TYPES.text);

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
      mocks: {
        $sidebar,
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
      mocks: {
        $sidebar,
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
      mocks: {
        $sidebar,
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
      mocks: {
        $sidebar,
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

  test('Map changed after trigger field map', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldMap = selectFieldMap(wrapper);

    const newMap = Faker.datatype.string();

    fieldMap.vm.$emit('input', newMap);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'map', newMap),
      },
    });
  });

  test('Color indicator changed after trigger field color indicator', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldColorIndicator = selectFieldColorIndicator(wrapper);

    const newColorIndicator = COLOR_INDICATOR_TYPES.impactState;

    fieldColorIndicator.vm.$emit('input', newColorIndicator);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'color_indicator', newColorIndicator),
      },
    });
  });

  test('Under pbehavior enabled changed after trigger field switcher', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldSwitcher = selectFieldSwitcher(wrapper);

    fieldSwitcher.vm.$emit('input', true);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'entities_under_pbehavior_enabled', true),
      },
    });
  });

  test('Filters changed after trigger field filters', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.map.actions.listFilters]: { actions: [] },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        userPreferenceModule,
        activeViewModule,
        serviceModule,
        widgetModule,
        authModule,
        widgetTemplateModule,
        dynamicInfoModule,
      ]),
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldFilters = selectFieldFilters(wrapper);

    const filter = Faker.datatype.string();

    fieldFilters.vm.$emit('input', filter);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(widget, 'parameters', {
          ...widget.parameters,

          mainFilter: filter,
        }),
      },
    });
  });

  test('Entity info template changed after trigger field text editor', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldTextEditor = selectFieldTextEditor(wrapper);

    const newTemplate = Faker.datatype.string();

    fieldTextEditor.vm.$emit('input', newTemplate);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'entity_info_template', newTemplate),
      },
    });
  });

  test('Alarms columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldColumns = selectAlarmsColumns(wrapper);

    const newColumns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    fieldColumns.vm.$emit('input', newColumns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'alarmsColumns', formToWidgetColumns(newColumns)),
      },
    });
  });

  test('Entities columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldColumns = selectEntitiesColumns(wrapper);

    const newColumns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    fieldColumns.vm.$emit('input', newColumns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'entitiesColumns', formToWidgetColumns(newColumns)),
      },
    });
  });

  test('Renders `map` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `map` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            ...sidebar.config,

            widget: {
              _id: '_id',
              type: WIDGET_TYPES.map,
              title: 'Map widget',
              parameters: {
                periodic_refresh: {},
                map: 'map_id',
                color_indicator: COLOR_INDICATOR_TYPES.state,
                entities_under_pbehavior_enabled: true,
                filters: [{}],
                entity_info_template: '<div>TEMPLATE</div>',
                alarmsColumns: [{}],
                entitiesColumns: [{}],
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
