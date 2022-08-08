import { range } from 'lodash';
import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { fakeAlarm } from '@unit/data/alarm';

import { generateDefaultAlarmListWidget } from '@/helpers/entities';

import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';

const localVue = createVueInstance();

const stubs = {
  'mass-actions-panel': true,
  'c-empty-data-table-columns': true,
  'alarm-header-cell': true,
  'alarms-list-row': true,
  'alarms-expand-panel': true,
};

const factory = (options = {}) => shallowMount(AlarmsListTable, {
  localVue,
  stubs,
  attachTo: document.body,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmsListTable, {
  localVue,
  stubs,
  attachTo: document.body,

  ...options,
});

const selectTable = wrapper => wrapper.find('v-data-table-stub');
const selectMassActionsPanel = wrapper => wrapper.find('mass-actions-panel-stub');
const selectAlarmsListRow = wrapper => wrapper.findAll('alarms-list-row-stub');
const selectTableHead = wrapper => wrapper.find('thead');
const selectTableBody = wrapper => wrapper.find('tbody');

const dispatchScrollEvent = (detail) => {
  window.dispatchEvent(
    new CustomEvent('scroll', { detail }),
  );
};

describe('alarms-list-table', () => {
  const timestamp = 1386435600;
  const totalItems = 5;
  const alarms = range(totalItems).map(value => ({
    _id: `alarm-${value}`,
    t: timestamp,
    entity: {
      _id: `entity-${value}`,
      name: `entity-name-${value}`,
      impact: [],
      depends: [],
      enable_history: [],
      measurements: null,
      enabled: true,
      type: 'resource',
      component: `component-${value}`,
    },
    v: {
      state: {
        _t: 'stateinc',
        t: timestamp,
        a: `author-${value}`,
        m: `message-${value}`,
        val: 3,
      },
      status: {
        _t: 'statusinc',
        t: timestamp,
        a: `author-${value}`,
        m: `message-${value}`,
        val: 1,
      },
      component: `component-${value}`,
      connector: `connector-${value}`,
      connector_name: `connector_name-${value}`,
      creation_date: timestamp,
      activation_date: timestamp,
      display_name: `display_name-${value}`,
      initial_output: `initial_output-${value}`,
      output: `output-${value}`,
      initial_long_output: `initial_long_output-${value}`,
      long_output: `long_output-${value}`,
      long_output_history: [],
      last_update_date: timestamp,
      last_event_date: timestamp,
      resource: `resource-${value}`,
      tags: [],
      parents: [],
      children: [],
      total_state_changes: 1,
      extra: {},
      infos_rule_version: {},
      duration: 270,
      current_state_duration: 270,
      infos: {},
    },
    assigned_instructions: [],
    infos: {},
    links: {},
  }));

  const associativeTableModule = {
    name: 'associativeTable',
    actions: {
      fetch: jest.fn(() => ({})),
    },
  };

  const store = createMockedStoreModules([
    associativeTableModule,
  ]);

  const defaultWidget = generateDefaultAlarmListWidget();

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('Alarms selected after trigger table', () => {
    const selectedAlarms = alarms.slice(0, -1);
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        widget: defaultWidget,
        columns: [],
        hasColumns: true,
      },
    });

    const table = selectTable(wrapper);

    table.vm.$emit('input', selectedAlarms);

    expect(wrapper.vm.selected).toEqual(selectedAlarms);
  });

  it('Selected alarms cleared after trigger mass actions', () => {
    const selectedAlarms = alarms.slice(0, -1);
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        widget: defaultWidget,
        columns: [],
        hasColumns: true,
      },
    });

    const table = selectTable(wrapper);
    table.vm.$emit('input', selectedAlarms);

    const massActionsPanel = selectMassActionsPanel(wrapper);
    massActionsPanel.vm.$emit('clear:items');

    expect(wrapper.vm.selected).toEqual([]);
  });

  it('Pagination update event emitted after trigger update pagination', () => {
    const wrapper = factory({
      store,
      propsData: {
        widget: defaultWidget,
        alarms: [],
        columns: [],
        hasColumns: true,
      },
    });

    const pagination = {
      descending: Faker.datatype.boolean(),
      multiSortBy: [],
      page: Faker.datatype.number(),
      rowsPerPage: Faker.datatype.number(),
      sortBy: Faker.datatype.string(),
      totalItems: Faker.datatype.number(),
    };

    const table = selectTable(wrapper);
    table.vm.$emit('update:pagination', pagination);

    const updatePaginationEvents = wrapper.emitted('update:pagination');
    expect(updatePaginationEvents).toHaveLength(1);

    const [eventData] = updatePaginationEvents[0];
    expect(eventData).toEqual(pagination);
  });

  it('Resize listener added after mount and removed after destroy', async () => {
    const addEventListener = jest.spyOn(window, 'addEventListener');
    const removeEventListener = jest.spyOn(window, 'removeEventListener');

    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: defaultWidget,
        alarms: [],
        columns: [],
        hasColumns: true,
        stickyHeader: true,
      },
    });

    expect(addEventListener).toHaveBeenNthCalledWith(
      1,
      'resize',
      expect.any(Function),
      { passive: true },
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      2,
      'scroll',
      expect.any(Function),
    );

    await wrapper.setProps({
      stickyHeader: false,
    });

    expect(removeEventListener).toHaveBeenCalledTimes(1);
    removeEventListener.mockClear();

    wrapper.destroy();

    expect(removeEventListener).toHaveBeenCalledTimes(2);
    expect(removeEventListener).toHaveBeenNthCalledWith(
      1,
      'scroll',
      expect.any(Function),
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      2,
      'resize',
      expect.any(Function),
      { passive: true },
    );
  });

  it('Timer cleared after disable sticky', async () => {
    const clearTimeout = jest.spyOn(window, 'clearTimeout');

    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: defaultWidget,
        alarms: [],
        columns: [],
        hasColumns: true,
      },
    });

    const header = selectTableHead(wrapper);
    const body = selectTableBody(wrapper);

    const headerGetBoundingClientRect = jest.spyOn(header.element, 'getBoundingClientRect')
      .mockReturnValue({ top: -200 });

    const bodyGetBoundingClientRect = jest.spyOn(body.element, 'getBoundingClientRect')
      .mockReturnValue({ height: 400 });

    await wrapper.setProps({
      stickyHeader: true,
    });

    dispatchScrollEvent(200);

    await wrapper.setProps({
      stickyHeader: false,
    });

    expect(clearTimeout).toHaveBeenCalled();

    clearTimeout.mockClear();
    headerGetBoundingClientRect.mockClear();
    bodyGetBoundingClientRect.mockClear();
  });

  it('Component adds and removes the same count listeners', async () => {
    const addEventListener = jest.spyOn(window, 'addEventListener');
    const removeEventListener = jest.spyOn(window, 'removeEventListener');

    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: defaultWidget,
        alarms: [],
        columns: [],
        hasColumns: true,
      },
    });

    expect(addEventListener).toHaveBeenCalledWith(
      'resize',
      expect.any(Function),
      { passive: true },
    );

    addEventListener.mockClear();

    await wrapper.setProps({
      stickyHeader: true,
    });

    expect(addEventListener).toHaveBeenCalledTimes(1);

    await wrapper.setProps({
      stickyHeader: false,
    });

    expect(removeEventListener).toHaveBeenCalledTimes(1);
    removeEventListener.mockClear();

    wrapper.destroy();

    expect(removeEventListener).toHaveBeenCalledTimes(2);
    expect(removeEventListener).toHaveBeenNthCalledWith(
      1,
      'scroll',
      expect.any(Function),
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      2,
      'resize',
      expect.any(Function),
      { passive: true },
    );
  });

  it('Header position changed after trigger scroll', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: defaultWidget,
        alarms,
        columns: [{
          label: 'Label-1',
          value: 'label',
        }],
        hasColumns: true,
        stickyHeader: true,
      },
    });

    const header = selectTableHead(wrapper);
    const body = selectTableBody(wrapper);

    const headerGetBoundingClientRect = jest.spyOn(header.element, 'getBoundingClientRect')
      .mockReturnValue({ top: -200 });

    const bodyGetBoundingClientRect = jest.spyOn(body.element, 'getBoundingClientRect')
      .mockReturnValue({ height: 400 });

    dispatchScrollEvent(200);

    await flushPromises();

    expect(header.element.style).toHaveProperty('transform', 'translateY(248px)');

    headerGetBoundingClientRect.mockClear();
    bodyGetBoundingClientRect.mockClear();
  });

  it('Header hidden after trigger start scroll', async () => {
    jest.useFakeTimers();

    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: defaultWidget,
        alarms,
        columns: [{
          label: 'Label-1',
          value: 'label',
        }],
        hasColumns: true,
        stickyHeader: true,
      },
    });

    const header = selectTableHead(wrapper);
    const body = selectTableBody(wrapper);

    const headerGetBoundingClientRect = jest.spyOn(header.element, 'getBoundingClientRect')
      .mockReturnValue({ top: -200 });

    const bodyGetBoundingClientRect = jest.spyOn(body.element, 'getBoundingClientRect')
      .mockReturnValue({ height: 400 });

    dispatchScrollEvent(200);

    expect(+header.element.style.opacity).toBe(0);

    jest.runAllTimers();

    expect(+header.element.style.opacity).toBe(1);

    headerGetBoundingClientRect.mockClear();
    bodyGetBoundingClientRect.mockClear();

    jest.useRealTimers();
  });

  it('Expanded elements works correctly', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: defaultWidget,
        alarms,
        columns: [{
          label: 'Label-1',
          value: 'label',
        }],
        hasColumns: true,
        stickyHeader: true,
      },
    });

    expect(wrapper.vm.expanded).toEqual({});

    const alarmsListRow = selectAlarmsListRow(wrapper).at(0);

    alarmsListRow.vm.row.expanded = true;

    const [firstAlarm] = alarms;

    expect(wrapper.vm.expanded).toEqual({
      [firstAlarm._id]: true,
    });

    alarmsListRow.vm.row.expanded = false;

    expect(wrapper.vm.expanded).toEqual({
      [firstAlarm._id]: false,
    });
  });

  it('Renders `alarms-list-table` with default and required props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: defaultWidget,
        alarms: [],
        columns: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-table` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarms,
        widget: defaultWidget,
        columns: [{
          label: 'Label-1',
          value: 'label',
        }],
        totalItems,
        pagination: {},
        isTourEnabled: true,
        loading: true,
        hasColumns: true,
        selectable: true,
        hideGroups: true,
        expandable: true,
        stickyHeader: true,
        parentAlarm: fakeAlarm(),
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-table` with expandable, but without selectable', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarms,
        widget: defaultWidget,
        columns: [],
        totalItems,
        pagination: {},
        hasColumns: true,
        selectable: false,
        expandable: true,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-table` with default and required props with compact mode', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: {
          ...defaultWidget,
          parameters: {
            ...defaultWidget.parameters,

            dense: true,
          },
        },
        alarms: [],
        columns: [],
        hasColumns: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
