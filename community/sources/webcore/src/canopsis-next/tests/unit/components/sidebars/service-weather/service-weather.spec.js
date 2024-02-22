import { omit } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockDateNow, mockSidebar } from '@unit/utils/mock-hooks';
import {
  createSettingsMocks,
  getWidgetRequestWithNewProperty,
  getWidgetRequestWithNewParametersProperty,
  submitWithExpects,
} from '@unit/utils/settings';

import {
  COLOR_INDICATOR_TYPES,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  SIDE_BARS,
  SORT_ORDERS,
  USERS_PERMISSIONS,
  WIDGET_TYPES,
} from '@/constants';

import ClickOutside from '@/services/click-outside';

import { widgetToForm, formToWidget, getEmptyWidgetByType } from '@/helpers/entities/widget/form';

import ServiceWeatherSettings from '@/components/sidebars/service-weather/service-weather.vue';

const stubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-filters': true,
  'alarms-list-modal-form': true,
  'field-number': true,
  'field-color-indicator': true,
  'field-columns': true,
  'field-sort-column': true,
  'field-default-elements-per-page': true,
  'field-text-editor-with-template': true,
  'field-grid-size': true,
  'margins-form': true,
  'field-slider': true,
  'field-counters-selector': true,
  'field-switcher': true,
  'field-modal-type': true,
  'field-action-required-settings': true,
  'field-tree-of-dependencies-settings': true,
};

const generateDefaultServiceWeatherWidget = () => ({
  ...formToWidget(widgetToForm({ type: WIDGET_TYPES.serviceWeather })),

  _id: Faker.datatype.string(),
});

const selectFieldTitle = wrapper => wrapper.find('field-title-stub');
const selectFieldPeriodicRefresh = wrapper => wrapper.find('field-periodic-refresh-stub');
const selectAlarmsListModalForm = wrapper => wrapper.find('alarms-list-modal-form-stub');
const selectFieldNumber = wrapper => wrapper.find('field-number-stub');
const selectFieldColorIndicator = wrapper => wrapper.find('field-color-indicator-stub');
const selectFieldSortColumn = wrapper => wrapper.find('field-sort-column-stub');
const selectFieldDefaultElementsPerPage = wrapper => wrapper.find('field-default-elements-per-page-stub');
const selectFieldTextEditorWithTemplateItems = wrapper => wrapper.findAll('field-text-editor-with-template-stub');
const selectFieldBlockTemplate = wrapper => selectFieldTextEditorWithTemplateItems(wrapper).at(0);
const selectFieldModalTemplate = wrapper => selectFieldTextEditorWithTemplateItems(wrapper).at(1);
const selectFieldEntityTemplate = wrapper => selectFieldTextEditorWithTemplateItems(wrapper).at(2);
const selectFieldGridSizeItems = wrapper => wrapper.findAll('field-grid-size-stub');
const selectFieldColumnMobile = wrapper => selectFieldGridSizeItems(wrapper).at(0);
const selectFieldColumnTablet = wrapper => selectFieldGridSizeItems(wrapper).at(1);
const selectFieldColumnDesktop = wrapper => selectFieldGridSizeItems(wrapper).at(2);
const selectFieldFilters = wrapper => wrapper.find('field-filters-stub');
const selectMarginsForm = wrapper => wrapper.find('margins-form-stub');
const selectFieldSlider = wrapper => wrapper.find('field-slider-stub');
const selectFieldCounters = wrapper => wrapper.find('field-counters-selector-stub');
const selectFieldSwitcherByIndex = (wrapper, index) => wrapper.findAll('field-switcher-stub').at(index);
const selectIsPriorityField = wrapper => selectFieldSwitcherByIndex(wrapper, 0);
const selectIsHideGrayField = wrapper => selectFieldSwitcherByIndex(wrapper, 1);
const selectEntitiesActionsInQueueField = wrapper => selectFieldSwitcherByIndex(wrapper, 2);
const selectFieldModalType = wrapper => wrapper.find('field-modal-type-stub');
const selectFieldActionRequiredSettingsType = wrapper => wrapper.find('field-action-required-settings-stub');

describe('service-weather', () => {
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
    ...generateDefaultServiceWeatherWidget(),

    tab: Faker.datatype.string(),
  };

  const sidebar = {
    name: SIDE_BARS.serviceWeatherSettings,
    config: {
      widget,
    },
    hidden: false,
  };

  const store = createMockedStoreModules([
    userPreferenceModule,
    activeViewModule,
    serviceModule,
    widgetTemplateModule,
    infosModule,
    widgetModule,
    authModule,
  ]);

  const factory = generateShallowRenderer(ServiceWeatherSettings, {

    stubs,
    store,
    propsData: {
      sidebar,
    },
    mocks: {
      $sidebar,
    },

    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  const snapshotFactory = generateRenderer(ServiceWeatherSettings, {

    stubs,
    store,
    propsData: {
      sidebar,
    },
    mocks: {
      $sidebar,
    },

    parentComponent: {
      provide: {
        $clickOutside: new ClickOutside(),
      },
    },
  });

  afterEach(() => {
    createWidget.mockReset();
    updateWidget.mockReset();
    fetchActiveView.mockReset();
  });

  test('Create widget with default parameters', async () => {
    const localWidget = getEmptyWidgetByType(WIDGET_TYPES.serviceWeather);

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

    const wrapper = factory();

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
    const wrapper = factory();

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

  test('Filters changed after trigger field filters', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.listFilters]: { actions: [] },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        userPreferenceModule,
        activeViewModule,
        serviceModule,
        widgetTemplateModule,
        infosModule,
        widgetModule,
        authModule,
      ]),
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const filter = Faker.datatype.string();

    selectFieldFilters(wrapper).triggerCustomEvent('input', filter);

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

  test('Alarms list modal changed after trigger alarms list modal form', async () => {
    const wrapper = factory();

    const newAlarmsList = {
      itemsPerPage: Faker.datatype.number(),
      moreInfoTemplate: Faker.datatype.string(),
      infoPopups: [],
      widgetColumns: [],
    };

    selectAlarmsListModalForm(wrapper).triggerCustomEvent('input', newAlarmsList);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'alarmsList', newAlarmsList),
      },
    });
  });

  test('Limit changed after trigger number field', async () => {
    const wrapper = factory();

    const newLimit = Faker.datatype.number();

    selectFieldNumber(wrapper).triggerCustomEvent('input', newLimit);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'limit', newLimit),
      },
    });
  });

  test('Color indicator changed after trigger field color indicator', async () => {
    const wrapper = factory();

    const fieldColorIndicator = selectFieldColorIndicator(wrapper);

    const newColorIndicator = COLOR_INDICATOR_TYPES.impactState;

    fieldColorIndicator.triggerCustomEvent('input', newColorIndicator);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'colorIndicator', newColorIndicator),
      },
    });
  });

  test('Sort changed after trigger sort column field', async () => {
    const wrapper = factory();

    const newSort = {
      column: Faker.datatype.string(),
      order: SORT_ORDERS.asc,
    };

    selectFieldSortColumn(wrapper).triggerCustomEvent('input', newSort);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'sort', newSort),
      },
    });
  });

  test('Modal items per page changed after trigger default per page field', async () => {
    const wrapper = factory();

    const newItemsPerPage = Faker.datatype.number();

    selectFieldDefaultElementsPerPage(wrapper).triggerCustomEvent('input', newItemsPerPage);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'modalItemsPerPage', newItemsPerPage),
      },
    });
  });

  test('Block template changed after trigger field template', async () => {
    const wrapper = factory();

    const blockTemplate = Faker.datatype.string();
    const blockTemplateTemplate = Faker.datatype.string();

    selectFieldBlockTemplate(wrapper).triggerCustomEvent('input', blockTemplate, blockTemplateTemplate);

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
            blockTemplate,
            blockTemplateTemplate,
          },
        ),
      },
    });
  });

  test('Modal template changed after trigger field template', async () => {
    const wrapper = factory();

    const modalTemplate = Faker.datatype.string();
    const modalTemplateTemplate = Faker.datatype.string();

    selectFieldModalTemplate(wrapper).triggerCustomEvent('input', modalTemplate, modalTemplateTemplate);

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
            modalTemplate,
            modalTemplateTemplate,
          },
        ),
      },
    });
  });

  test('Entity template changed after trigger field template', async () => {
    const wrapper = factory();

    const entityTemplate = Faker.datatype.string();
    const entityTemplateTemplate = Faker.datatype.string();

    selectFieldEntityTemplate(wrapper).triggerCustomEvent('input', entityTemplate, entityTemplateTemplate);

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
            entityTemplate,
            entityTemplateTemplate,
          },
        ),
      },
    });
  });

  test('Column mobile changed after trigger field grid size', async () => {
    const wrapper = factory();

    const newSize = Faker.datatype.number();

    selectFieldColumnMobile(wrapper).triggerCustomEvent('input', newSize);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'columnMobile', newSize),
      },
    });
  });

  test('Column tablet changed after trigger field grid size', async () => {
    const wrapper = factory();

    const newSize = Faker.datatype.number();

    selectFieldColumnTablet(wrapper).triggerCustomEvent('input', newSize);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'columnTablet', newSize),
      },
    });
  });

  test('Column desktop changed after trigger field grid size', async () => {
    const wrapper = factory();

    const newSize = Faker.datatype.number();

    selectFieldColumnDesktop(wrapper).triggerCustomEvent('input', newSize);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'columnDesktop', newSize),
      },
    });
  });

  test('Margins changed after trigger margins form', async () => {
    const wrapper = factory();

    const newMargins = {
      top: Faker.datatype.number(),
      right: Faker.datatype.number(),
      bottom: Faker.datatype.number(),
      left: Faker.datatype.number(),
    };

    selectMarginsForm(wrapper).triggerCustomEvent('input', newMargins);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'margin', newMargins),
      },
    });
  });

  test('Height factor changed after trigger slider field', async () => {
    const wrapper = factory();

    const newHeightFactor = Faker.datatype.number();

    selectFieldSlider(wrapper).triggerCustomEvent('input', newHeightFactor);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'heightFactor', newHeightFactor),
      },
    });
  });

  test('Counters changed after trigger counters field', async () => {
    const wrapper = factory();

    const newCounters = {
      pbehavior_enabled: Faker.datatype.boolean(),
      pbehavior_types: [Faker.datatype.string()],
      state_enabled: Faker.datatype.boolean(),
      state_types: [Faker.datatype.string()],
    };

    selectFieldCounters(wrapper).triggerCustomEvent('input', newCounters);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'counters', newCounters),
      },
    });
  });

  test('Is priority enabled changed after trigger switcher field', async () => {
    const wrapper = factory();

    selectIsPriorityField(wrapper).triggerCustomEvent('input', true);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'isPriorityEnabled', true),
      },
    });
  });

  test('Is priority enabled changed after trigger switcher field', async () => {
    const wrapper = factory();

    selectIsHideGrayField(wrapper).triggerCustomEvent('input', true);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'isHideGrayEnabled', true),
      },
    });
  });

  test('Is entities actions in queue changed after trigger switcher field', async () => {
    const wrapper = factory();

    selectEntitiesActionsInQueueField(wrapper).triggerCustomEvent('input', true);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'entitiesActionsInQueue', true),
      },
    });
  });

  test('Modal type changed after trigger modal type field', async () => {
    const wrapper = factory();

    selectFieldModalType(wrapper).triggerCustomEvent('input', SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'modalType', SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo),
      },
    });
  });

  test('Action required settings changed after trigger action required settings field', async () => {
    const wrapper = factory();

    const newValue = {
      is_blinking: Faker.datatype.boolean(),
      icon_name: Faker.datatype.string(),
      color: Faker.internet.color(),
    };

    selectFieldActionRequiredSettingsType(wrapper).triggerCustomEvent('input', newValue);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'actionRequiredSettings', newValue),
      },
    });
  });

  test('Renders `service-weather` widget settings with default props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-weather` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            ...sidebar.config,

            widget: {
              _id: '_id',
              type: WIDGET_TYPES.serviceWeather,
              title: 'Map widget',
              filters: [{}, {}],
              parameters: {
                periodic_refresh: {},
                mainFilter: 'main_filter',
                alarmsList: {},
                limit: 12,
                colorIndicator: COLOR_INDICATOR_TYPES.state,
                serviceDependenciesColumns: [{}, {}],
                sort: {},
                modalItemsPerPage: 11,
                blockTemplate: '<div>block-template</div>',
                modalTemplate: '<div>modal-template</div>',
                entityTemplate: '<div>entity-template</div>',
                columnMobile: 1,
                columnTablet: 2,
                columnDesktop: 3,
                margin: {},
                heightFactor: 22,
                counters: {},
                isPriorityEnabled: true,
                isHideGrayEnabled: true,
                modalType: SERVICE_WEATHER_WIDGET_MODAL_TYPES.both,
                actionRequiredSettings: {
                  is_blinking: true,
                  icon_name: 'menu',
                  color: '#123',
                },
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
