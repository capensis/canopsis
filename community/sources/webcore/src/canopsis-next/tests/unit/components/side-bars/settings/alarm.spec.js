import { omit } from 'lodash';
import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { createButtonStub } from '@unit/stubs/button';
import { createInputStub } from '@unit/stubs/input';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { alarmListWidgetToForm } from '@/helpers/forms/widgets/alarm';
import {
  ALARMS_OPENED_VALUES,
  CANOPSIS_EDITION, EXPORT_CSV_DATETIME_FORMATS, EXPORT_CSV_SEPARATORS,
  REMEDIATION_INSTRUCTION_TYPES,
  SORT_ORDERS,
  USERS_PERMISSIONS,
} from '@/constants';
import ClickOutside from '@/services/click-outside';

import AlarmSettings from '@/components/side-bars/settings/alarm.vue';

const localVue = createVueInstance();

const stubs = {
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
  'export-csv-form': createInputStub('export-csv-form'),
  'v-btn': createButtonStub('v-btn'),
};

const snapshotStubs = {
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
  'export-csv-form': true,
};

const $clickOutside = new ClickOutside();

const factory = (options = {}) => shallowMount(AlarmSettings, {
  localVue,
  stubs,
  parentComponent: {
    provide: {
      $clickOutside,
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmSettings, {
  localVue,
  stubs: snapshotStubs,
  parentComponent: {
    provide: {
      $clickOutside,
    },
  },

  ...options,
});

const selectSubmitButton = wrapper => wrapper.find('button.v-btn');
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
const selectFieldHtmlEnabledSwitcher = wrapper => wrapper.findAll('input.field-switcher').at(0);
const selectFieldAckNoteRequired = wrapper => wrapper.findAll('input.field-switcher').at(1);
const selectFieldMultiAckEnabled = wrapper => wrapper.findAll('input.field-switcher').at(2);
const selectFieldFastAckOutput = wrapper => wrapper.find('input.field-fast-ack-output');
const selectFieldSnoozeNoteRequired = wrapper => wrapper.findAll('input.field-switcher').at(3);
const selectFieldLinksCategoriesAsList = wrapper => wrapper.find('input.field-enabled-limit');
const selectFieldExportCsvForm = wrapper => wrapper.find('input.export-csv-form');
const selectFieldStickyHeader = wrapper => wrapper.findAll('input.field-switcher').at(4);

describe('alarm', () => {
  const nowTimestamp = 1386435600000;
  mockDateNow(nowTimestamp);

  const userPreferences = {
    content: {
      itemsPerPage: 13,
      category: '_id',
    },
  };
  const widget = alarmListWidgetToForm();
  const view = {
    enabled: true,
    title: 'Alarm widgets',
    description: 'Alarm widgets',
    tabs: [
      {
        widgets: [widget],
      },
    ],
    tags: ['alarm'],
    periodic_refresh: {
      value: 1,
      unit: 's',
      enabled: false,
    },
    author: 'root',
    group: {
      _id: 'text-widget-group',
    },
  };
  const updateUserPreference = jest.fn();
  const updateView = jest.fn();
  const updateQuery = jest.fn();
  const hideSideBar = jest.fn();
  const sideBarModule = {
    name: 'sideBar',
    actions: {
      hide: hideSideBar,
    },
  };
  const infoModule = {
    name: 'info',
    getters: { edition: CANOPSIS_EDITION.cat },
  };
  const queryModule = {
    name: 'query',
    getters: { getQueryById: () => () => ({}) },
    actions: {
      update: updateQuery,
    },
  };
  const viewModule = {
    name: 'view',
    getters: {
      item: view,
    },
    actions: {
      update: updateView,
    },
  };
  const userPreferenceModule = {
    name: 'userPreference',
    getters: {
      getItemByWidget: () => () => userPreferences,
    },
    actions: {
      update: updateUserPreference,
    },
  };
  const authModule = {
    name: 'auth',
    getters: {
      currentUserPermissionsById: {},
    },
  };
  const defaultQuery = {
    active_columns: widget.parameters.widgetColumns.map(v => v.value),
    category: userPreferences.content.category,
    limit: userPreferences.content.itemsPerPage,
    opened: widget.parameters.opened,
    correlation: widget.parameters.isCorrelationEnabled ?? false,
    multiSortBy: [],
    page: 1,
    with_instructions: true,
  };

  const store = createMockedStoreModules([
    sideBarModule,
    infoModule,
    queryModule,
    viewModule,
    userPreferenceModule,
    authModule,
  ]);

  const getViewRequest = () => ({
    ...omit(view, ['_id']),
    group: view.group._id,
  });

  const getViewRequestWithNewWidgetProperty = (parameter, value) => ({
    ...getViewRequest(),
    tabs: [{
      widgets: [{
        ...widget,
        [parameter]: value,
      }],
    }],
  });

  const getViewRequestWithNewWidgetParameter = (parameter, value) => getViewRequestWithNewWidgetProperty('parameters', {
    ...widget.parameters,
    [parameter]: value,
  });

  const submitWithExpects = async (wrapper, { viewData, userPreferencesData, queryData }) => {
    const submitButton = selectSubmitButton(wrapper);

    submitButton.trigger('click');

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledTimes(1);
    expect(updateUserPreference).toHaveBeenLastCalledWith(
      expect.any(Object),
      { data: userPreferencesData },
      undefined,
    );

    expect(updateView).toHaveBeenCalledTimes(1);
    expect(updateView).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: view._id,
        data: viewData,
      },
      undefined,
    );

    expect(updateQuery).toHaveBeenCalledTimes(1);
    expect(updateQuery).toHaveBeenLastCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: queryData,
      },
      undefined,
    );

    expect(hideSideBar).toHaveBeenCalledTimes(1);
  };

  afterEach(() => {
    updateUserPreference.mockReset();
    updateView.mockReset();
    updateQuery.mockReset();
    hideSideBar.mockReset();
  });

  it('Items per page changed after mount', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    expect(wrapper.vm.settings.userPreferenceContent.itemsPerPage).toEqual(
      userPreferences.content.itemsPerPage,
    );
  });

  it('Title changed after trigger field title', async () => {
    const newTitle = Faker.datatype.string();

    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldTitle = selectFieldTitle(wrapper);

    fieldTitle.setValue(newTitle);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetProperty('title', newTitle),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Periodic refresh changed after trigger field periodic refresh', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
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
      viewData: getViewRequestWithNewWidgetParameter('periodic_refresh', periodicRefresh),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Default sort column changed after trigger field default sort column', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldDefaultSortColumn = selectFieldDefaultSortColumn(wrapper);

    const sort = {
      order: SORT_ORDERS.desc,
      column: Faker.datatype.string(),
    };

    fieldDefaultSortColumn.vm.$emit('input', sort);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('sort', sort),
      userPreferencesData: userPreferences,
      queryData: {
        ...defaultQuery,
        multiSortBy: [{
          descending: true,
          sortBy: sort.column,
        }],
      },
    });
  });

  it('Widget columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldColumns = selectFieldWidgetColumns(wrapper);

    const columns = [{
      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    fieldColumns.vm.$emit('input', columns);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('widgetColumns', columns),
      userPreferencesData: userPreferences,
      queryData: {
        ...defaultQuery,
        active_columns: columns.map(v => v.value),
      },
    });
  });

  it('Widget group columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldColumns = selectFieldWidgetGroupColumns(wrapper);

    const columns = [{
      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    fieldColumns.vm.$emit('input', columns);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('widgetGroupColumns', columns),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Service dependencies columns changed after trigger field columns', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldColumns = selectFieldServiceDependenciesColumns(wrapper);

    const columns = [{
      label: Faker.datatype.string(),
      value: Faker.datatype.string(),
    }];

    fieldColumns.vm.$emit('input', columns);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('serviceDependenciesColumns', columns),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Default elements per page changed after trigger field default elements per page', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldDefaultElementsPerPage = selectFieldDefaultElementsPerPage(wrapper);

    const itemsPerPage = Faker.datatype.number();

    fieldDefaultElementsPerPage.vm.$emit('input', itemsPerPage);

    await submitWithExpects(wrapper, {
      viewData: getViewRequest(),
      userPreferencesData: {
        ...userPreferences,
        content: {
          ...userPreferences.content,
          itemsPerPage,
        },
      },
      queryData: {
        ...defaultQuery,
        limit: itemsPerPage,
      },
    });
  });

  it('Opened changed after trigger opened and resolved field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldOpenedResolvedFilter = selectFieldOpenedResolvedFilter(wrapper);

    fieldOpenedResolvedFilter.vm.$emit('input', ALARMS_OPENED_VALUES.all);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('opened', ALARMS_OPENED_VALUES.all),
      userPreferencesData: userPreferences,
      queryData: {
        ...defaultQuery,
        opened: ALARMS_OPENED_VALUES.all,
      },
    });
  });

  it('Filters changed after trigger filters field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        {
          ...authModule,
          getters: {
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.listFilters]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldFilters = selectFieldFilters(wrapper);

    const filter = {
      title: Faker.datatype.string(),
      filter: Faker.helpers.createTransaction(),
    };

    fieldFilters.vm.$emit('input', filter);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetProperty('parameters', {
        ...widget.parameters,
        mainFilterUpdatedAt: nowTimestamp,
        mainFilter: filter,
      }),
      userPreferencesData: userPreferences,
      queryData: {
        ...defaultQuery,
        filter: filter.filter,
      },
    });
  });

  it('Instruction filters changed after trigger remediation instructions field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
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
        config: {
          widget,
        },
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
      viewData: getViewRequestWithNewWidgetParameter(
        'remediationInstructionsFilters',
        remediationInstructionsFilters,
      ),
      userPreferencesData: userPreferences,
      queryData: {
        ...defaultQuery,
        include_instructions: [remediationInstruction._id],
        include_instruction_types: [REMEDIATION_INSTRUCTION_TYPES.auto],
      },
    });
  });

  it('Live reporting changed after trigger live reporting field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldLiveReporting = selectFieldLiveReporting(wrapper);

    const liveReporting = {
      tstart: 1,
      tstop: 2,
    };

    fieldLiveReporting.vm.$emit('input', liveReporting);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('liveReporting', liveReporting),
      userPreferencesData: userPreferences,
      queryData: {
        ...defaultQuery,
        ...liveReporting,
      },
    });
  });

  it('Info popup changed after trigger info popup field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldInfoPopups = selectFieldInfoPopups(wrapper);

    const infoPopups = [{
      column: 'alarm.v.connector',
      template: 'Info popup',
    }];

    fieldInfoPopups.vm.$emit('input', infoPopups);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('infoPopups', infoPopups),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('More info template popup changed after trigger text field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldTextEditor = selectFieldTextEditor(wrapper);

    const moreInfoTemplate = 'More info template';

    fieldTextEditor.vm.$emit('input', moreInfoTemplate);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('moreInfoTemplate', moreInfoTemplate),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Grid range changed after trigger grid range field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldGridRangeSize = selectFieldGridRangeSize(wrapper);

    const expandGridRangeSize = [1, 11];

    fieldGridRangeSize.vm.$emit('input', expandGridRangeSize);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter(
        'expandGridRangeSize',
        expandGridRangeSize,
      ),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Html enabled on timeline changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldHtmlEnabledSwitcher = selectFieldHtmlEnabledSwitcher(wrapper);

    const isHtmlEnabledOnTimeLine = Faker.datatype.boolean();

    fieldHtmlEnabledSwitcher.vm.$emit('input', isHtmlEnabledOnTimeLine);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter(
        'isHtmlEnabledOnTimeLine',
        isHtmlEnabledOnTimeLine,
      ),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Ack note required changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldAckNoteRequired = selectFieldAckNoteRequired(wrapper);

    const isAckNoteRequired = Faker.datatype.boolean();

    fieldAckNoteRequired.vm.$emit('input', isAckNoteRequired);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('isAckNoteRequired', isAckNoteRequired),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Multi ack enabled changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldAckNoteRequired = selectFieldMultiAckEnabled(wrapper);

    const isMultiAckEnabled = Faker.datatype.boolean();

    fieldAckNoteRequired.vm.$emit('input', isMultiAckEnabled);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('isMultiAckEnabled', isMultiAckEnabled),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Fast ack output changed after trigger fast ack output field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldAckNoteRequired = selectFieldFastAckOutput(wrapper);

    const fastAckOutput = {
      enabled: true,
      output: Faker.datatype.string(),
    };

    fieldAckNoteRequired.vm.$emit('input', fastAckOutput);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('fastAckOutput', fastAckOutput),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Snooze note required changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldSnoozeNoteRequired = selectFieldSnoozeNoteRequired(wrapper);

    const isSnoozeNoteRequired = Faker.datatype.boolean();

    fieldSnoozeNoteRequired.vm.$emit('input', isSnoozeNoteRequired);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('isSnoozeNoteRequired', isSnoozeNoteRequired),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Snooze note required changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldLinksCategoriesAsList = selectFieldLinksCategoriesAsList(wrapper);

    const linksCategoriesAsList = {
      enabled: Faker.datatype.boolean(),
      limit: Faker.datatype.number(),
    };

    fieldLinksCategoriesAsList.vm.$emit('input', linksCategoriesAsList);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('linksCategoriesAsList', linksCategoriesAsList),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Export parameters changed after trigger export csv form', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget: { ...widget },
        },
      },
    });

    const fieldExportCsvForm = selectFieldExportCsvForm(wrapper);

    const exportProperties = {
      ...widget.parameters,
      exportCsvSeparator: EXPORT_CSV_SEPARATORS.semicolon,
      exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
      widgetExportColumns: [],
    };

    fieldExportCsvForm.vm.$emit('input', exportProperties);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetProperty('parameters', {
        ...widget.parameters,
        ...exportProperties,
      }),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Sticky header changed after trigger switcher field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        config: {
          widget,
        },
      },
    });

    const fieldStickyHeader = selectFieldStickyHeader(wrapper);

    const stickyHeader = Faker.datatype.boolean();

    fieldStickyHeader.vm.$emit('input', stickyHeader);

    await submitWithExpects(wrapper, {
      viewData: getViewRequestWithNewWidgetParameter('sticky_header', stickyHeader),
      userPreferencesData: userPreferences,
      queryData: defaultQuery,
    });
  });

  it('Renders `alarm` widget settings with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        config: {
          widget: alarmListWidgetToForm(),
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm` widget settings with all rights', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
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
        config: {
          widget: alarmListWidgetToForm(),
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm` widget settings with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
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
              mainFilterUpdatedAt: 0,
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
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
