import { omit } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createButtonStub } from '@unit/stubs/button';
import { createInputStub } from '@unit/stubs/input';
import { mockDateNow, mockSidebar } from '@unit/utils/mock-hooks';
import {
  createSettingsMocks,
  getWidgetRequestWithNewProperty,
  getWidgetRequestWithNewParametersProperty,
  submitWithExpects,
} from '@unit/utils/settings';

import {
  ALARMS_OPENED_VALUES,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  SORT_ORDERS,
  TIME_UNITS,
  USERS_PERMISSIONS,
  SIDE_BARS,
  COLOR_INDICATOR_TYPES,
  WIDGET_TYPES,
} from '@/constants';

import ClickOutside from '@/services/click-outside';

import { alarmListChartToForm, formToAlarmListChart } from '@/helpers/entities/widget/forms/alarm';
import {
  generateDefaultAlarmListWidget,
  widgetToForm,
  formToWidget,
  getEmptyWidgetByType,
  widgetParametersToForm,
  formToWidgetParameters,
} from '@/helpers/entities/widget/form';
import { formToWidgetColumns, widgetColumnToForm } from '@/helpers/entities/widget/column/form';

import AlarmSettings from '@/components/sidebars/alarm/alarm.vue';

const stubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': createInputStub('field-title'),
  'field-periodic-refresh': createInputStub('field-periodic-refresh'),
  'field-default-sort-column': createInputStub('field-default-sort-column'),
  'field-columns': createInputStub('field-columns'),
  'field-default-elements-per-page': createInputStub('field-default-elements-per-page'),
  'field-opened-resolved-filter': createInputStub('field-opened-resolved-filter'),
  'field-filters': createInputStub('field-filters'),
  'field-remediation-instructions-filters': createInputStub('field-remediation-instructions-filters'),
  'field-live-reporting': createInputStub('field-live-reporting'),
  'field-info-popup': createInputStub('field-info-popup'),
  'field-text-editor-with-template': createInputStub('field-text-editor-with-template'),
  'field-grid-range-size': createInputStub('field-grid-range-size'),
  'field-switcher': createInputStub('field-switcher'),
  'field-fast-action-output': createInputStub('field-fast-action-output'),
  'field-number': createInputStub('field-number'),
  'field-density': createInputStub('field-density'),
  'export-csv-form': createInputStub('export-csv-form'),
  'charts-form': createInputStub('charts-form'),
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
  'widget-settings': true,
  'widget-settings-item': true,
  'widget-settings-group': true,
  'field-title': true,
  'field-periodic-refresh': true,
  'field-default-sort-column': true,
  'field-columns': true,
  'field-default-elements-per-page': true,
  'field-opened-resolved-filter': true,
  'field-filters': true,
  'field-remediation-instructions-filters': true,
  'field-live-reporting': true,
  'field-info-popup': true,
  'field-text-editor-with-template': true,
  'field-grid-range-size': true,
  'field-switcher': true,
  'field-fast-action-output': true,
  'field-number': true,
  'field-density': true,
  'export-csv-form': true,
  'charts-form': true,
  'field-resize-column-behavior': true,
};

const selectSwitcherFieldByTitle = (wrapper, title) => wrapper.find(`input.field-switcher[title="${title}"]`);
const selectFieldTitle = wrapper => wrapper.find('input.field-title');
const selectFieldPeriodicRefresh = wrapper => wrapper.find('input.field-periodic-refresh');
const selectFieldDefaultSortColumn = wrapper => wrapper.find('input.field-default-sort-column');
const selectFieldWidgetColumns = wrapper => wrapper.findAll('input.field-columns').at(0);
const selectFieldWidgetGroupColumns = wrapper => wrapper.findAll('input.field-columns').at(1);
const selectFieldServiceDependenciesColumns = wrapper => wrapper.findAll('input.field-columns').at(2);
const selectFieldDefaultElementsPerPage = wrapper => wrapper.find('input.field-default-elements-per-page');
const selectFieldOpenedResolvedFilter = wrapper => wrapper.find('input.field-opened-resolved-filter');
const selectFieldFilters = wrapper => wrapper.find('input.field-filters');
const selectFieldRemediationInstructionsFilters = wrapper => wrapper.find('input.field-remediation-instructions-filters');
const selectFieldLiveReporting = wrapper => wrapper.find('input.field-live-reporting');
const selectFieldInfoPopups = wrapper => wrapper.find('input.field-info-popup');
const selectFieldTextEditorWithTemplate = wrapper => wrapper.find('input.field-text-editor-with-template');
const selectFieldGridRangeSize = wrapper => wrapper.find('input.field-grid-range-size');
const selectFieldClearFilterDisabled = wrapper => selectSwitcherFieldByTitle(
  wrapper,
  'Disable possibility to clear selected filter',
);
const selectFieldHtmlEnabledSwitcher = wrapper => selectSwitcherFieldByTitle(wrapper, 'HTML enabled on timeline?');
const selectFieldAckNoteRequired = wrapper => selectSwitcherFieldByTitle(wrapper, 'Note field required when ack?');
const selectFieldMultiAckEnabled = wrapper => selectSwitcherFieldByTitle(wrapper, 'Multiple ack');
const selectFieldFastAckOutput = wrapper => wrapper.findAll('input.field-fast-action-output').at(0);
const selectFieldFastCancelOutput = wrapper => wrapper.findAll('input.field-fast-action-output').at(1);
const selectFieldSnoozeNoteRequired = wrapper => selectSwitcherFieldByTitle(wrapper, 'Note field required when snooze?');
const selectFieldRemoveAlarmsFromMetaAlarmCommentRequired = wrapper => selectSwitcherFieldByTitle(
  wrapper,
  'Comment field required when remove alarms from manual meta alarm?',
);
const selectFieldExportCsvForm = wrapper => wrapper.find('input.export-csv-form');
const selectFieldStickyHeader = wrapper => selectSwitcherFieldByTitle(wrapper, 'Sticky header');
const selectFieldKioskHideActions = wrapper => selectSwitcherFieldByTitle(wrapper, 'Hide actions');
const selectFieldKioskHideMassSelection = wrapper => selectSwitcherFieldByTitle(wrapper, 'Hide mass selection');
const selectFieldKioskHideToolbar = wrapper => selectSwitcherFieldByTitle(wrapper, 'Hide toolbar');
const selectFieldActionsAllowWithOkState = wrapper => selectSwitcherFieldByTitle(
  wrapper,
  'Actions allowed when state OK?',
);
const selectFieldShowRootCauseByStateClick = wrapper => selectSwitcherFieldByTitle(
  wrapper,
  'Show root cause diagram called from Severity column',
);
const selectChartsForm = wrapper => wrapper.findAll('input.charts-form').at(0);

describe('alarm', () => {
  const parentComponent = {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  };

  const factory = generateShallowRenderer(AlarmSettings, {
    stubs,
    parentComponent,
  });

  const snapshotFactory = generateRenderer(AlarmSettings, {
    stubs: snapshotStubs,
    parentComponent,
  });

  const nowTimestamp = 1386435600000;

  mockDateNow(nowTimestamp);

  const $sidebar = mockSidebar();

  const {
    createWidget,
    updateWidget,
    fetchActiveView,
    fetchUserPreference,
    activeViewModule,
    widgetModule,
    authModule,
    userPreferenceModule,
    widgetTemplateModule,
    serviceModule,
    infosModule,
  } = createSettingsMocks();

  const widget = {
    ...generateDefaultAlarmListWidget(),

    _id: '3f8dba7c-f39e-42ae-912c-e78cb39669c5',
    tab: Faker.datatype.string(),
  };

  widget.parameters.usedAlarmProperties = [
    'v.connector',
    'v.connector_name',
    'v.component',
    'v.resource',
    'v.output',
    'extra_details',
    'v.state.val',
    'v.status.val',
  ];

  const sidebar = {
    name: SIDE_BARS.alarmSettings,
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
    const localWidget = getEmptyWidgetByType(WIDGET_TYPES.alarmList);

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

    fieldTitle.setValue(newTitle);

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

    const periodicRefresh = {
      enabled: Faker.datatype.boolean(),
      value: Faker.datatype.number(),
      unit: Faker.datatype.string(),
    };

    selectFieldPeriodicRefresh(wrapper).triggerCustomEvent('input', {
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

  test('Live watching changed after trigger field periodic refresh', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const liveWatching = Faker.datatype.boolean();

    selectFieldPeriodicRefresh(wrapper).triggerCustomEvent('input', {
      ...wrapper.vm.form.parameters,
      liveWatching,
    });

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'liveWatching', liveWatching),
      },
    });
  });

  test('Default sort column changed after trigger field default sort column', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldDefaultSortColumn = selectFieldDefaultSortColumn(wrapper);

    const sort = {
      order: SORT_ORDERS.desc,
      column: Faker.datatype.string(),
    };

    fieldDefaultSortColumn.triggerCustomEvent('input', sort);

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
      mocks: {
        $sidebar,
      },
    });

    const fieldColumns = selectFieldWidgetColumns(wrapper);

    const columns = [{
      ...widgetColumnToForm(),

      value: Faker.datatype.string(),
    }];

    fieldColumns.triggerCustomEvent('input', columns);

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
            usedAlarmProperties: [''],
          },
        ),
      },
    });
  });

  test('Widget columns with `alarm.` prefix changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldColumns = selectFieldWidgetColumns(wrapper);

    const columns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: `alarm.${Faker.datatype.string()}`,
      isHtml: false,
      colorIndicator: COLOR_INDICATOR_TYPES.state,
    }];

    fieldColumns.triggerCustomEvent('input', columns);

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
            usedAlarmProperties: [''],
          },
        ),
      },
    });
  });

  test('Widget group columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldColumns = selectFieldWidgetGroupColumns(wrapper);

    const columns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    fieldColumns.triggerCustomEvent('input', columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'widgetGroupColumns', formToWidgetColumns(columns)),
      },
    });
  });

  test('Service dependencies columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldColumns = selectFieldServiceDependenciesColumns(wrapper);

    const columns = [{
      ...widgetColumnToForm(),

      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    fieldColumns.triggerCustomEvent('input', columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data:
          getWidgetRequestWithNewParametersProperty(widget, 'serviceDependenciesColumns', formToWidgetColumns(columns)),
      },
    });
  });

  test('Default elements per page changed after trigger field default elements per page', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldDefaultElementsPerPage = selectFieldDefaultElementsPerPage(wrapper);

    const itemsPerPage = Faker.datatype.number();

    fieldDefaultElementsPerPage.triggerCustomEvent('input', itemsPerPage);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'itemsPerPage', itemsPerPage),
      },
    });
  });

  test('Opened changed after trigger opened and resolved field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldOpenedResolvedFilter = selectFieldOpenedResolvedFilter(wrapper);

    fieldOpenedResolvedFilter.triggerCustomEvent('input', ALARMS_OPENED_VALUES.all);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'opened', ALARMS_OPENED_VALUES.all),
      },
    });
  });

  test('Filters changed after trigger update:filters on filters field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        infosModule,
        {
          ...authModule,
          getters: {
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.listFilters]: {
                actions: [],
              },
            },
          },
        },
      ]),
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldFilters = selectFieldFilters(wrapper);

    const filters = [{
      title: Faker.datatype.string(),
      filter: Faker.helpers.createTransaction(),
    }];

    fieldFilters.triggerCustomEvent('update:filters', filters);

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
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        infosModule,
        {
          ...authModule,
          getters: {
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.listFilters]: {
                actions: [],
              },
            },
          },
        },
      ]),
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldFilters = selectFieldFilters(wrapper);

    const filter = {
      title: Faker.datatype.string(),
      filter: Faker.helpers.createTransaction(),
    };

    fieldFilters.triggerCustomEvent('input', filter);

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

  test('Instruction filters changed after trigger remediation instructions field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        infosModule,
        {
          ...authModule,
          getters: {
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.listRemediationInstructionsFilters]: {
                actions: [],
              },
            },
          },
        },
      ]),
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldRemediationInstructionsFilters = selectFieldRemediationInstructionsFilters(wrapper);

    const remediationInstruction = {
      _id: 'instruction_1',
      name: 'instruction-1',
      type: {
        _id: 'instruction-type-1',
      },
    };
    const remediationInstructionsFilters = [{
      with: true,
      all: false,
      auto: true,
      manual: false,
      locked: true,
      disabled: false,
      instructions: [remediationInstruction],
      _id: 'id1',
    }];

    fieldRemediationInstructionsFilters.triggerCustomEvent('input', remediationInstructionsFilters);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'remediationInstructionsFilters', remediationInstructionsFilters),
      },
    });
  });

  test('Live reporting changed after trigger live reporting field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldLiveReporting = selectFieldLiveReporting(wrapper);

    const liveReporting = {
      tstart: 1,
      tstop: 2,
    };

    fieldLiveReporting.triggerCustomEvent('input', liveReporting);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'liveReporting', liveReporting),
      },
    });
  });

  test('Info popup changed after trigger info popup field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldInfoPopups = selectFieldInfoPopups(wrapper);

    const infoPopups = [{
      column: 'alarm.v.connector',
      template: 'Info popup',
    }];

    fieldInfoPopups.triggerCustomEvent('input', infoPopups);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'infoPopups', infoPopups),
      },
    });
  });

  test('More info template changed after trigger text field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldTextEditor = selectFieldTextEditorWithTemplate(wrapper);

    const moreInfoTemplate = 'More info template';
    const moreInfoTemplateTemplate = 'template-id';

    fieldTextEditor.triggerCustomEvent('input', moreInfoTemplate, moreInfoTemplateTemplate);

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
            moreInfoTemplate,
            moreInfoTemplateTemplate,
          },
        ),
      },
    });
  });

  test('Grid range changed after trigger grid range field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldGridRangeSize = selectFieldGridRangeSize(wrapper);

    const expandGridRangeSize = [1, 11];

    fieldGridRangeSize.triggerCustomEvent('input', expandGridRangeSize);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'expandGridRangeSize', expandGridRangeSize),
      },
    });
  });

  test('Clear filter disabled changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldHtmlEnabledSwitcher = selectFieldClearFilterDisabled(wrapper);

    const clearFilterDisabled = Faker.datatype.boolean();

    fieldHtmlEnabledSwitcher.triggerCustomEvent('input', clearFilterDisabled);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'clearFilterDisabled', clearFilterDisabled),
      },
    });
  });

  test('Clear filter disabled changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const showRootCauseByStateClick = Faker.datatype.boolean();

    selectFieldShowRootCauseByStateClick(wrapper).triggerCustomEvent('input', showRootCauseByStateClick);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'showRootCauseByStateClick', showRootCauseByStateClick),
      },
    });
  });

  test('Html enabled on timeline changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldHtmlEnabledSwitcher = selectFieldHtmlEnabledSwitcher(wrapper);

    const isHtmlEnabledOnTimeLine = Faker.datatype.boolean();

    fieldHtmlEnabledSwitcher.triggerCustomEvent('input', isHtmlEnabledOnTimeLine);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'isHtmlEnabledOnTimeLine', isHtmlEnabledOnTimeLine),
      },
    });
  });

  test('Ack note required changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldAckNoteRequired = selectFieldAckNoteRequired(wrapper);

    const isAckNoteRequired = Faker.datatype.boolean();

    fieldAckNoteRequired.triggerCustomEvent('input', isAckNoteRequired);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'isAckNoteRequired', isAckNoteRequired),
      },
    });
  });

  test('Multi ack enabled changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldAckNoteRequired = selectFieldMultiAckEnabled(wrapper);

    const isMultiAckEnabled = Faker.datatype.boolean();

    fieldAckNoteRequired.triggerCustomEvent('input', isMultiAckEnabled);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'isMultiAckEnabled', isMultiAckEnabled),
      },
    });
  });

  test('Fast ack output changed after trigger fast ack output field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldFastAckOutput = selectFieldFastAckOutput(wrapper);

    const fastAckOutput = {
      enabled: true,
      output: Faker.datatype.string(),
    };

    fieldFastAckOutput.triggerCustomEvent('input', fastAckOutput);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'fastAckOutput', fastAckOutput),
      },
    });
  });

  test('Fast cancel output changed after trigger fast cancel output field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldFastCancelOutput = selectFieldFastCancelOutput(wrapper);

    const fastAckOutput = {
      enabled: true,
      output: Faker.datatype.string(),
    };

    fieldFastCancelOutput.triggerCustomEvent('input', fastAckOutput);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'fastCancelOutput', fastAckOutput),
      },
    });
  });

  test('Snooze note required changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldSnoozeNoteRequired = selectFieldSnoozeNoteRequired(wrapper);

    const isSnoozeNoteRequired = Faker.datatype.boolean();

    fieldSnoozeNoteRequired.triggerCustomEvent('input', isSnoozeNoteRequired);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'isSnoozeNoteRequired', isSnoozeNoteRequired),
      },
    });
  });

  test('Remove alarms from meta alarms comment required changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const isRemoveAlarmsFromMetaAlarmCommentRequired = Faker.datatype.boolean();

    selectFieldRemoveAlarmsFromMetaAlarmCommentRequired(wrapper).triggerCustomEvent('input', isRemoveAlarmsFromMetaAlarmCommentRequired);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(
          widget,
          'isRemoveAlarmsFromMetaAlarmCommentRequired',
          isRemoveAlarmsFromMetaAlarmCommentRequired,
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
      mocks: {
        $sidebar,
      },
    });

    const fieldExportCsvForm = selectFieldExportCsvForm(wrapper);

    const exportProperties = {
      ...widgetParametersToForm(widget),
      exportCsvSeparator: EXPORT_CSV_SEPARATORS.semicolon,
      exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
      widgetExportColumns: [],
    };

    fieldExportCsvForm.triggerCustomEvent('input', exportProperties);

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
              type: WIDGET_TYPES.alarmList,
              parameters: exportProperties,
            }),
          ),
      },
    });
  });

  test('Sticky header changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldStickyHeader = selectFieldStickyHeader(wrapper);

    const stickyHeader = Faker.datatype.boolean();

    fieldStickyHeader.triggerCustomEvent('input', stickyHeader);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'sticky_header', stickyHeader),
      },
    });
  });

  /**
   * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/4390
   */
  test('Invalid periodic refresh converted to valid object', async () => {
    const periodicRefresh = {
      value: 1,
      unit: {},
      enabled: false,
    };
    const wrapper = factory({
      store,
      propsData: {
        sidebar: {
          ...sidebar,

          config: {
            widget: {
              ...widget,
              parameters: {
                ...widget.parameters,
                periodic_refresh: periodicRefresh,
              },
            },
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
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'periodic_refresh', {
          ...periodicRefresh,
          unit: TIME_UNITS.second,
        }),
      },
    });
  });

  test('Kiosk fields changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldKioskHideActions = selectFieldKioskHideActions(wrapper);
    const fieldKioskHideMassSelection = selectFieldKioskHideMassSelection(wrapper);
    const fieldKioskHideToolbar = selectFieldKioskHideToolbar(wrapper);

    const hideActions = Faker.datatype.boolean();
    const hideMassSelection = Faker.datatype.boolean();
    const hideToolbar = Faker.datatype.boolean();

    fieldKioskHideActions.triggerCustomEvent('input', hideActions);
    fieldKioskHideMassSelection.triggerCustomEvent('input', hideMassSelection);
    fieldKioskHideToolbar.triggerCustomEvent('input', hideToolbar);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'kiosk', {
          hideActions,
          hideMassSelection,
          hideToolbar,
        }),
      },
    });
  });

  test('Actions allowed with state ok changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const isActionsAllowWithOkState = Faker.datatype.boolean();

    selectFieldActionsAllowWithOkState(wrapper).triggerCustomEvent('input', isActionsAllowWithOkState);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'isActionsAllowWithOkState', isActionsAllowWithOkState),
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

    const chartsForm = selectChartsForm(wrapper);
    const newCharts = [formToAlarmListChart(alarmListChartToForm())];

    chartsForm.triggerCustomEvent('input', newCharts);

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

  test('Renders `alarm` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm` widget settings with all rights', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        infosModule,
        serviceModule,
        {
          ...authModule,
          getters: {
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.listFilters]: { actions: [] },
              [USERS_PERMISSIONS.business.alarmsList.actions.editFilter]: { actions: [] },
              [USERS_PERMISSIONS.business.alarmsList.actions.addFilter]: { actions: [] },
              [USERS_PERMISSIONS.business.alarmsList.actions.userFilter]: { actions: [] },
              [USERS_PERMISSIONS.business.alarmsList.actions.listRemediationInstructionsFilters]: {
                actions: [],
              },
              [USERS_PERMISSIONS.business.alarmsList.actions.addRemediationInstructionsFilter]: {
                actions: [],
              },
              [USERS_PERMISSIONS.business.alarmsList.actions.editRemediationInstructionsFilter]: {
                actions: [],
              },
              [USERS_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter]: {
                actions: [],
              },
            },
          },
        },
      ]),
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar: {
          config: {
            widget: {
              _id: 'widget_AlarmsList_2',
              type: 'AlarmsList',
              title: 'Alarms list',
              parameters: {
                itemsPerPage: 12,
                infoPopups: [],
                moreInfoTemplate: 'More info template',
                isAckNoteRequired: true,
                isSnoozeNoteRequired: true,
                isMultiAckEnabled: true,
                isHtmlEnabledOnTimeLine: true,
                sticky_header: true,
                fastAckOutput: { enabled: true, value: 'output' },
                widgetColumns: [{ label: 'connector', value: 'v.connector' }],
                widgetGroupColumns: [{ label: 'connector', value: 'v.connector' }],
                serviceDependenciesColumns: [{ label: 'connector', value: 'v.connector' }],
                linksCategoriesAsList: { enabled: true, limit: 3 },
                periodic_refresh: { value: 30, unit: 's', enabled: true },
                viewFilters: [],
                mainFilter: null,
                liveReporting: {},
                sort: { order: SORT_ORDERS.desc, column: 'connector' },
                opened: true,
                expandGridRangeSize: [1, 11],
                exportCsvSeparator: 'comma',
                exportCsvDatetimeFormat: 'YYYY-DDThh:mm:ssZ',
                widgetExportColumns: [{ label: 'connector', value: 'v.connector' }],
              },
              grid_parameters: {
                mobile: { x: 1, y: 1, h: 2, w: 11, autoHeight: true },
                tablet: { x: 1, y: 1, h: 2, w: 11, autoHeight: true },
                desktop: { x: 1, y: 1, h: 2, w: 11, autoHeight: true },
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
