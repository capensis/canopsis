import { omit } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createButtonStub } from '@unit/stubs/button';
import { createInputStub } from '@unit/stubs/input';
import { mockSidebar } from '@unit/utils/mock-hooks';
import {
  createSettingsMocks,
  getWidgetRequestWithNewProperty,
  getWidgetRequestWithNewParametersProperty,
  submitWithExpects,
} from '@unit/utils/settings';

import {
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  SORT_ORDERS,
  USERS_PERMISSIONS,
  SIDE_BARS,
  COLOR_INDICATOR_TYPES,
  WIDGET_TYPES,
  TREE_OF_DEPENDENCIES_SHOW_TYPES,
  ENTITY_TYPES,
} from '@/constants';

import ClickOutside from '@/services/click-outside';

import { alarmListChartToForm, formToAlarmListChart } from '@/helpers/entities/widget/forms/alarm';
import {
  widgetToForm,
  formToWidget,
  getEmptyWidgetByType,
  widgetParametersToForm,
  formToWidgetParameters,
  generateDefaultContextWidget,
} from '@/helpers/entities/widget/form';
import { formToWidgetColumns, widgetColumnToForm } from '@/helpers/entities/widget/column/form';

import ContextSettings from '@/components/sidebars/context/context.vue';

const stubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': createInputStub('field-title'),
  'field-default-sort-column': createInputStub('field-default-sort-column'),
  'field-columns': createInputStub('field-columns'),
  'field-tree-of-dependencies-settings': createInputStub('field-tree-of-dependencies-settings'),
  'field-root-cause-settings': createInputStub('field-root-cause-settings'),
  'field-filters': createInputStub('field-filters'),
  'field-context-entities-types-filter': createInputStub('field-context-entities-types-filter'),
  'export-csv-form': createInputStub('export-csv-form'),
  'charts-form': createInputStub('charts-form'),
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-default-sort-column': true,
  'field-columns': true,
  'field-tree-of-dependencies-settings': true,
  'field-root-cause-settings': true,
  'field-filters': true,
  'field-context-entities-types-filter': true,
  'export-csv-form': true,
  'charts-form': true,
};

const selectFieldTitle = wrapper => wrapper.find('input.field-title');
const selectFieldDefaultSortColumn = wrapper => wrapper.find('input.field-default-sort-column');
const selectFieldWidgetColumns = wrapper => wrapper.findAll('input.field-columns').at(0);
const selectFieldServiceDependenciesColumns = wrapper => wrapper.findAll('input.field-columns').at(1);
const selectFieldActiveAlarmsColumns = wrapper => wrapper.findAll('input.field-columns').at(2);
const selectFieldResolvedAlarmsColumns = wrapper => wrapper.findAll('input.field-columns').at(3);
const selectFieldTreeOfDependenciesSettings = wrapper => wrapper.find('input.field-tree-of-dependencies-settings');
const selectFieldContextEntitiesTypesFilter = wrapper => wrapper.find('input.field-context-entities-types-filter');
const selectFieldExportCsvForm = wrapper => wrapper.find('input.export-csv-form');
const selectChartsForm = wrapper => wrapper.find('input.charts-form');
const selectFieldFilters = wrapper => wrapper.find('input.field-filters');
const selectFieldRootCauseSettings = wrapper => wrapper.find('.field-root-cause-settings');

describe('context', () => {
  const $sidebar = mockSidebar();
  const parentComponent = {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  };

  const factory = generateShallowRenderer(ContextSettings, {
    stubs,
    parentComponent,
    mocks: {
      $sidebar,
    },
  });

  const snapshotFactory = generateRenderer(ContextSettings, {
    stubs: snapshotStubs,
    parentComponent,
  });

  const nowTimestamp = 1386435600000;
  jest.useFakeTimers({ now: nowTimestamp });

  const {
    createWidget,
    updateWidget,
    fetchActiveView,
    fetchUserPreference,
    activeViewModule,
    widgetModule,
    authModule,
    currentUserPermissionsById,
    userPreferenceModule,
    widgetTemplateModule,
    serviceModule,
    infosModule,
  } = createSettingsMocks();

  const widget = {
    ...generateDefaultContextWidget(),

    _id: '3f8dba7c-f39e-42ae-912c-e78cb39669c5',
    tab: Faker.datatype.string(),
  };

  const sidebar = {
    name: SIDE_BARS.contextSettings,
    config: {
      widget,
    },
    hidden: false,
  };

  const store = createMockedStoreModules([
    activeViewModule,
    widgetModule,
    userPreferenceModule,
    authModule,
    userPreferenceModule,
    widgetTemplateModule,
    serviceModule,
    infosModule,
  ]);

  afterEach(() => {
    createWidget.mockReset();
    updateWidget.mockReset();
    fetchActiveView.mockReset();
    fetchUserPreference.mockReset();
  });

  test('Create widget with default parameters', async () => {
    const localWidget = getEmptyWidgetByType(WIDGET_TYPES.context);

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

  test('Title changed after trigger field title', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const newTitle = Faker.datatype.string();
    selectFieldTitle(wrapper).setValue(newTitle);

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

  test('Default sort column changed after trigger field default sort column', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const sort = {
      order: SORT_ORDERS.desc,
      column: Faker.datatype.string(),
    };

    selectFieldDefaultSortColumn(wrapper).triggerCustomEvent('input', sort);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'sort', sort),
      },
    });
  });

  test('Widget columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const columns = [{
      ...widgetColumnToForm(),

      value: Faker.datatype.string(),
    }];

    selectFieldWidgetColumns(wrapper).triggerCustomEvent('input', columns);

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
            widgetColumns: formToWidgetColumns(columns),
          },
        ),
      },
    });
  });

  test('Widget columns with `entity.` prefix changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const columns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: `entity.${Faker.datatype.string()}`,
      isHtml: false,
      colorIndicator: COLOR_INDICATOR_TYPES.state,
    }];

    selectFieldWidgetColumns(wrapper).triggerCustomEvent('input', columns);

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
            widgetColumns: formToWidgetColumns(columns),
          },
        ),
      },
    });
  });

  test('Service dependencies columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const columns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    selectFieldServiceDependenciesColumns(wrapper).triggerCustomEvent('input', columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data:
          getWidgetRequestWithNewParametersProperty(
            widget,
            'serviceDependenciesColumns',
            formToWidgetColumns(columns),
          ),
      },
    });
  });

  test('Widget active alarms columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const columns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    selectFieldActiveAlarmsColumns(wrapper).triggerCustomEvent('input', columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'activeAlarmsColumns', formToWidgetColumns(columns)),
      },
    });
  });

  test('Widget resolved alarms columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const columns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    selectFieldResolvedAlarmsColumns(wrapper).triggerCustomEvent('input', columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(
          widget,
          'resolvedAlarmsColumns',
          formToWidgetColumns(columns),
        ),
      },
    });
  });

  test('Widget tree of dependencies settings changed after trigger field tree of deps', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    selectFieldTreeOfDependenciesSettings(wrapper).triggerCustomEvent(
      'input',
      TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
    );

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(
          widget,
          'treeOfDependenciesShowType',
          TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
        ),
      },
    });
  });

  test('Widget tree of dependencies settings changed after trigger field tree of deps', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const newTypes = [ENTITY_TYPES.component, ENTITY_TYPES.connector];
    selectFieldContextEntitiesTypesFilter(wrapper).triggerCustomEvent('input', newTypes);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(
          widget,
          'selectedTypes',
          newTypes,
        ),
      },
    });
  });

  test('Filters changed after trigger update:filters on filters field', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.context.actions.listFilters]: {
        actions: [],
      },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        infosModule,
        authModule,
      ]),
      propsData: {
        sidebar,
      },
    });

    const filters = [{
      title: Faker.datatype.string(),
      filter: Faker.helpers.createTransaction(),
    }];

    selectFieldFilters(wrapper).triggerCustomEvent('update:filters', filters);

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

  test('Filter changed after trigger input on filters field', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.context.actions.listFilters]: {
        actions: [],
      },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        infosModule,
        authModule,
      ]),
      propsData: {
        sidebar,
      },
    });

    const filter = {
      title: Faker.datatype.string(),
      filter: Faker.helpers.createTransaction(),
    };

    selectFieldFilters(wrapper).triggerCustomEvent('input', filter);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'mainFilter', filter),
      },
    });
  });

  test('Root cause settings changed after trigger root cause settings field', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.context.actions.listFilters]: {
        actions: [],
      },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        infosModule,
        authModule,
      ]),
      propsData: {
        sidebar,
      },
    });

    const newParameters = {
      ...widget.parameters,
      showRootCauseByStateClick: false,
      rootCauseColorIndicator: COLOR_INDICATOR_TYPES.impactState,
    };
    selectFieldRootCauseSettings(wrapper).triggerCustomEvent('input', newParameters);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewProperty(
          widget,
          'parameters',
          formToWidgetParameters({
            type: WIDGET_TYPES.context,
            parameters: newParameters,
          }),
        ),
      },
    });
  });

  test('Export parameters changed after trigger export csv form', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
    });

    const exportProperties = {
      ...widgetParametersToForm(widget),
      exportCsvSeparator: EXPORT_CSV_SEPARATORS.semicolon,
      exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
      widgetExportColumns: [],
    };

    selectFieldExportCsvForm(wrapper).triggerCustomEvent('input', exportProperties);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data:
          getWidgetRequestWithNewProperty(
            widget,
            'parameters',
            formToWidgetParameters({
              type: WIDGET_TYPES.context,
              parameters: exportProperties,
            }),
          ),
      },
    });
  });

  test('Charts fields changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const newCharts = [formToAlarmListChart(alarmListChartToForm())];

    selectChartsForm(wrapper).triggerCustomEvent('input', newCharts);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'charts', newCharts),
      },
    });
  });

  test('Renders `context` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `context` widget settings with all rights', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.context.actions.listFilters]: { actions: [] },
      [USERS_PERMISSIONS.business.context.actions.editFilter]: { actions: [] },
      [USERS_PERMISSIONS.business.context.actions.addFilter]: { actions: [] },
      [USERS_PERMISSIONS.business.context.actions.userFilter]: { actions: [] },
    });
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        infosModule,
        serviceModule,
        authModule,
      ]),
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `context` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar: {
          config: {
            widget: {
              ...formToWidget(
                widgetToForm({
                  type: WIDGET_TYPES.context,
                  title: 'Context widget title',
                  parameters: {
                    sort: {
                      order: SORT_ORDERS.desc,
                      column: 'connector',
                    },
                    widgetColumns: [],
                    serviceDependenciesColumns: [{ label: 'connector', value: 'v.connector' }],
                    treeOfDependenciesShowType: TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState,
                    selectedTypes: [ENTITY_TYPES.service],
                  },
                }),
              ),

              _id: 'context-widget-id',
            },
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
