import { omit } from 'lodash';
import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
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
import { generateDefaultAlarmListWidget } from '@/helpers/entities';
import {
  widgetToForm,
  formToWidget,
  getEmptyWidgetByType,
  widgetParametersToForm, formToWidgetParameters,
} from '@/helpers/forms/widgets/common';
import { formToWidgetColumns, widgetColumnToForm } from '@/helpers/forms/shared/widget-column';

import AlarmSettings from '@/components/sidebars/settings/alarm.vue';

const localVue = createVueInstance();

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
  'field-text-editor': createInputStub('field-text-editor'),
  'field-grid-range-size': createInputStub('field-grid-range-size'),
  'field-switcher': createInputStub('field-switcher'),
  'field-fast-ack-output': createInputStub('field-fast-ack-output'),
  'field-enabled-limit': createInputStub('field-enabled-limit'),
  'field-density': createInputStub('field-density'),
  'export-csv-form': createInputStub('export-csv-form'),
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
  'field-text-editor': true,
  'field-grid-range-size': true,
  'field-switcher': true,
  'field-fast-ack-output': true,
  'field-enabled-limit': true,
  'field-density': true,
  'export-csv-form': true,
};

const factory = (options = {}) => shallowMount(AlarmSettings, {
  localVue,
  stubs,
  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmSettings, {
  localVue,
  stubs: snapshotStubs,
  parentComponent: {
    provide: {
      $clickOutside: new ClickOutside(),
    },
  },

  ...options,
});

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
const selectFieldTextEditor = wrapper => wrapper.find('input.field-text-editor');
const selectFieldGridRangeSize = wrapper => wrapper.find('input.field-grid-range-size');
const selectFieldClearFilterDisabled = wrapper => wrapper.findAll('input.field-switcher').at(0);
const selectFieldHtmlEnabledSwitcher = wrapper => wrapper.findAll('input.field-switcher').at(1);
const selectFieldAckNoteRequired = wrapper => wrapper.findAll('input.field-switcher').at(2);
const selectFieldMultiAckEnabled = wrapper => wrapper.findAll('input.field-switcher').at(3);
const selectFieldFastAckOutput = wrapper => wrapper.find('input.field-fast-ack-output');
const selectFieldSnoozeNoteRequired = wrapper => wrapper.findAll('input.field-switcher').at(4);
const selectFieldLinksCategoriesAsList = wrapper => wrapper.find('input.field-enabled-limit');
const selectFieldExportCsvForm = wrapper => wrapper.find('input.export-csv-form');
const selectFieldStickyHeader = wrapper => wrapper.findAll('input.field-switcher').at(6);

describe('alarm', () => {
  const nowTimestamp = 1386435600000;

  mockDateNow(nowTimestamp);

  const $sidebar = mockSidebar();

  const {
    createWidget,
    updateWidget,
    copyWidget,
    fetchActiveView,
    fetchUserPreference,
    activeViewModule,
    widgetModule,
    authModule,
    userPreferenceModule,
    widgetTemplateModule,
    serviceModule,
    dynamicInfoModule,
  } = createSettingsMocks();

  const widget = {
    ...generateDefaultAlarmListWidget(),

    _id: '3f8dba7c-f39e-42ae-912c-e78cb39669c5',
    tab: Faker.datatype.string(),
  };

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
    dynamicInfoModule,
  ]);

  afterEach(() => {
    createWidget.mockReset();
    updateWidget.mockReset();
    copyWidget.mockReset();
    fetchActiveView.mockReset();
    fetchUserPreference.mockReset();
  });

  it('Create widget with default parameters', async () => {
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

  it('Duplicate widget without changes', async () => {
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

  it('Title changed after trigger field title', async () => {
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

  it('Periodic refresh changed after trigger field periodic refresh', async () => {
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

  it('Default sort column changed after trigger field default sort column', async () => {
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

    fieldDefaultSortColumn.vm.$emit('input', sort);

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

  it('Widget columns changed after trigger field columns', async () => {
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

    fieldColumns.vm.$emit('input', columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'widgetColumns', formToWidgetColumns(columns)),
      },
    });
  });

  it('Widget columns with `alarm.` prefix changed after trigger field columns', async () => {
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

    fieldColumns.vm.$emit('input', columns);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'widgetColumns', formToWidgetColumns(columns)),
      },
    });
  });

  it('Widget group columns changed after trigger field columns', async () => {
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

    fieldColumns.vm.$emit('input', columns);

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

  it('Service dependencies columns changed after trigger field columns', async () => {
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

    fieldColumns.vm.$emit('input', columns);

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

  it('Default elements per page changed after trigger field default elements per page', async () => {
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

    fieldDefaultElementsPerPage.vm.$emit('input', itemsPerPage);

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

  it('Opened changed after trigger opened and resolved field', async () => {
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

    fieldOpenedResolvedFilter.vm.$emit('input', ALARMS_OPENED_VALUES.all);

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

  it('Filters changed after trigger update:filters on filters field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        dynamicInfoModule,
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

    fieldFilters.vm.$emit('update:filters', filters);

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

  it('Filter changed after trigger input on filters field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        dynamicInfoModule,
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

  it('Instruction filters changed after trigger remediation instructions field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        serviceModule,
        dynamicInfoModule,
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

    fieldRemediationInstructionsFilters.vm.$emit('input', remediationInstructionsFilters);

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

  it('Live reporting changed after trigger live reporting field', async () => {
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

    fieldLiveReporting.vm.$emit('input', liveReporting);

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

  it('Info popup changed after trigger info popup field', async () => {
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

    fieldInfoPopups.vm.$emit('input', infoPopups);

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

  it('More info template popup changed after trigger text field', async () => {
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

    const moreInfoTemplate = 'More info template';

    fieldTextEditor.vm.$emit('input', moreInfoTemplate);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'moreInfoTemplate', moreInfoTemplate),
      },
    });
  });

  it('Grid range changed after trigger grid range field', async () => {
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

    fieldGridRangeSize.vm.$emit('input', expandGridRangeSize);

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

  it('Clear filter disabled changed after trigger switcher field', async () => {
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

    fieldHtmlEnabledSwitcher.vm.$emit('input', clearFilterDisabled);

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

  it('Html enabled on timeline changed after trigger switcher field', async () => {
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

    fieldHtmlEnabledSwitcher.vm.$emit('input', isHtmlEnabledOnTimeLine);

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

  it('Ack note required changed after trigger switcher field', async () => {
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

    fieldAckNoteRequired.vm.$emit('input', isAckNoteRequired);

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

  it('Multi ack enabled changed after trigger switcher field', async () => {
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

    fieldAckNoteRequired.vm.$emit('input', isMultiAckEnabled);

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

  it('Fast ack output changed after trigger fast ack output field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldAckNoteRequired = selectFieldFastAckOutput(wrapper);

    const fastAckOutput = {
      enabled: true,
      output: Faker.datatype.string(),
    };

    fieldAckNoteRequired.vm.$emit('input', fastAckOutput);

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

  it('Snooze note required changed after trigger switcher field', async () => {
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

    fieldSnoozeNoteRequired.vm.$emit('input', isSnoozeNoteRequired);

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

  it('Link categories as list required changed after trigger links categories as list field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        sidebar,
      },
      mocks: {
        $sidebar,
      },
    });

    const fieldLinksCategoriesAsList = selectFieldLinksCategoriesAsList(wrapper);

    const linksCategoriesAsList = {
      enabled: Faker.datatype.boolean(),
      limit: Faker.datatype.number(),
    };

    fieldLinksCategoriesAsList.vm.$emit('input', linksCategoriesAsList);

    await submitWithExpects(wrapper, {
      fetchActiveView,
      hideSidebar: $sidebar.hide,
      widgetMethod: updateWidget,
      expectData: {
        id: widget._id,
        data: getWidgetRequestWithNewParametersProperty(widget, 'linksCategoriesAsList', linksCategoriesAsList),
      },
    });
  });

  it('Export parameters changed after trigger export csv form', async () => {
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

    fieldExportCsvForm.vm.$emit('input', exportProperties);

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

  it('Sticky header changed after trigger switcher field', async () => {
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

    fieldStickyHeader.vm.$emit('input', stickyHeader);

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
  it('Invalid periodic refresh converted to valid object', async () => {
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

  it('Renders `alarm` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        sidebar,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm` widget settings with all rights', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        activeViewModule,
        widgetModule,
        userPreferenceModule,
        widgetTemplateModule,
        dynamicInfoModule,
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

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm` widget settings with custom props', async () => {
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
