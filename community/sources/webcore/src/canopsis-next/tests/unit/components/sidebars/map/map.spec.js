import { omit } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockSidebar } from '@unit/utils/mock-hooks';
import { createButtonStub } from '@unit/stubs/button';
import {
  createSettingsMocks,
  getWidgetRequestWithNewProperty,
  getWidgetRequestWithNewParametersProperty,
  submitWithExpects,
} from '@unit/utils/settings';

import { COLOR_INDICATOR_TYPES, SIDE_BARS, USERS_PERMISSIONS, WIDGET_TYPES } from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/entities/widget/form';
import { formToWidgetColumns, widgetColumnToForm } from '@/helpers/entities/widget/column/form';

import MapSettings from '@/components/sidebars/map/map.vue';

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
    infosModule,
  ]);

  const factory = generateShallowRenderer(MapSettings, { stubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });
  const snapshotFactory = generateRenderer(MapSettings, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  const timestamp = 1386435600000;

  beforeAll(() => jest.useFakeTimers({ now: timestamp }));

  afterEach(() => {
    createWidget.mockReset();
    updateWidget.mockReset();
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
      widgetMethod: createWidget,
      expectData: {
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

    fieldTitle.triggerCustomEvent('input', newTitle);

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

    fieldPeriodicRefresh.triggerCustomEvent('input', {
      ...wrapper.vm.form.parameters,
      periodic_refresh: periodicRefresh,
    });

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

    fieldMap.triggerCustomEvent('input', newMap);

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

    fieldColorIndicator.triggerCustomEvent('input', newColorIndicator);

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

    fieldSwitcher.triggerCustomEvent('input', true);

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
        infosModule,
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

    fieldFilters.triggerCustomEvent('input', filter);

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

    fieldTextEditor.triggerCustomEvent('input', newTemplate);

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

    fieldColumns.triggerCustomEvent('input', newColumns);

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

    fieldColumns.triggerCustomEvent('input', newColumns);

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

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
