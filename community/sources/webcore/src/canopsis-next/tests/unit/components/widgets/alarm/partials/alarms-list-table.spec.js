import { range } from 'lodash';
import Faker from 'faker';

import { flushPromises, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreGetters, createMockedStoreModules } from '@unit/utils/store';
import { fakeAlarm } from '@unit/data/alarm';
import { triggerWindowKeyboardEvent, triggerWindowScrollEvent } from '@unit/utils/events';
import { mockModals } from '@unit/utils/mock-hooks';

import { ALARM_DENSE_TYPES, ALARM_FIELDS } from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';

const stubs = {
  'mass-actions-panel': true,
  'c-empty-data-table-columns': true,
  'alarm-header-cell': true,
  'alarms-list-row': true,
  'alarms-expand-panel': true,
  'c-pagination': true,
  'c-table-pagination': true,
  'c-density-btn-toggle': true,
};

const selectTable = wrapper => wrapper.findComponent({ name: 'VDataTable' });
const selectAlarmsListRow = wrapper => wrapper.findAll('alarms-list-row-stub');
const selectTableHead = wrapper => wrapper.find('thead');
const selectTableBody = wrapper => wrapper.find('tbody');

describe('alarms-list-table', () => {
  const $modals = mockModals();
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
    createMockedStoreGetters({ name: 'info', showHeaderOnKioskMode: false }),
  ]);

  const defaultWidget = generatePreparedDefaultAlarmListWidget();

  const columns = [{
    label: 'Label-1',
    value: 'label',
  }];

  const snapshotFactory = generateRenderer(AlarmsListTable, {
    stubs,
    attachTo: document.body,
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  test('Alarms selected after trigger table', () => {
    const selectedAlarms = alarms.slice(0, -1);
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        alarms,
        columns,
        widget: defaultWidget,
      },
    });

    selectTable(wrapper).triggerCustomEvent('input', selectedAlarms);

    expect(wrapper.vm.selected).toEqual(selectedAlarms);
  });

  test('Pagination update event emitted after trigger update pagination', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        columns,
        widget: defaultWidget,
        alarms: [],
      },
    });

    const options = {
      page: Faker.datatype.number(),
      itemsPerPage: Faker.datatype.number(),
      sortBy: [Faker.datatype.string()],
      sortDesc: [Faker.datatype.boolean()],
      totalItems: Faker.datatype.number(),
    };

    selectTable(wrapper).triggerCustomEvent('update:options', options);

    expect(wrapper).toEmit('update:options', expect.any(Object), options);
  });

  test('Resize listener added after mount and removed after destroy', async () => {
    const addEventListener = jest.spyOn(window, 'addEventListener');
    const removeEventListener = jest.spyOn(window, 'removeEventListener');

    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        columns,
        widget: defaultWidget,
        alarms: [],
        stickyHeader: true,
        selectable: true,
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
      'resize',
      expect.any(Function),
      { passive: true },
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      3,
      'keydown',
      expect.any(Function),
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      4,
      'keyup',
      expect.any(Function),
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      5,
      'resize',
      wrapper.vm.resizeHandler,
      { passive: true },
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      6,
      'scroll',
      wrapper.vm.changeHeaderPosition,
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      7,
      'keydown',
      wrapper.vm.enableSelecting,
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      8,
      'keyup',
      wrapper.vm.disableSelecting,
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      9,
      'mousedown',
      wrapper.vm.mousedownHandler,
    );

    expect(addEventListener).toHaveBeenNthCalledWith(
      10,
      'mouseup',
      wrapper.vm.mouseupHandler,
    );

    await wrapper.setProps({
      stickyHeader: false,
    });

    expect(removeEventListener).toHaveBeenNthCalledWith(1, 'scroll', wrapper.vm.changeHeaderPosition);
    removeEventListener.mockClear();

    wrapper.destroy();

    expect(removeEventListener).toHaveBeenCalledTimes(8);

    expect(removeEventListener).toHaveBeenNthCalledWith(
      1,
      'scroll',
      wrapper.vm.changeHeaderPosition,
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      2,
      'keydown',
      wrapper.vm.enableSelecting,
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      3,
      'keyup',
      wrapper.vm.disableSelecting,
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      4,
      'mousedown',
      wrapper.vm.mousedownHandler,
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      5,
      'mouseup',
      wrapper.vm.mouseupHandler,
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      6,
      'resize',
      wrapper.vm.resizeHandler,
      { passive: true },
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      7,
      'keydown',
      expect.any(Function),
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      8,
      'keyup',
      expect.any(Function),
    );
  });

  test('Timer cleared after disable sticky', async () => {
    const clearTimeout = jest.spyOn(window, 'clearTimeout');

    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        columns,
        widget: defaultWidget,
        alarms: [],
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

    triggerWindowScrollEvent(200);

    await wrapper.setProps({
      stickyHeader: false,
    });

    expect(clearTimeout).toHaveBeenCalled();

    clearTimeout.mockClear();
    headerGetBoundingClientRect.mockClear();
    bodyGetBoundingClientRect.mockClear();
  });

  test('Component adds and removes the same count listeners', async () => {
    const addEventListener = jest.spyOn(window, 'addEventListener');
    const removeEventListener = jest.spyOn(window, 'removeEventListener');

    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        columns,
        widget: defaultWidget,
        alarms: [],
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

    expect(removeEventListener).toHaveBeenCalledTimes(8);
    expect(removeEventListener).toHaveBeenNthCalledWith(
      1,
      'scroll',
      expect.any(Function),
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      2,
      'keydown',
      expect.any(Function),
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      3,
      'keyup',
      expect.any(Function),
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      4,
      'mousedown',
      expect.any(Function),
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      5,
      'mouseup',
      expect.any(Function),
    );
    expect(removeEventListener).toHaveBeenNthCalledWith(
      6,
      'resize',
      expect.any(Function),
      { passive: true },
    );
  });

  test('Header position changed after trigger scroll', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: defaultWidget,
        alarms,
        columns,
        stickyHeader: true,
      },
    });

    const header = selectTableHead(wrapper);
    const body = selectTableBody(wrapper);

    const headerGetBoundingClientRect = jest.spyOn(header.element, 'getBoundingClientRect')
      .mockReturnValue({ top: -200 });

    const bodyGetBoundingClientRect = jest.spyOn(body.element, 'getBoundingClientRect')
      .mockReturnValue({ height: 400 });

    triggerWindowScrollEvent(200);

    await flushPromises();

    expect(header.element.style).toHaveProperty('transform', 'translateY(248px)');

    headerGetBoundingClientRect.mockClear();
    bodyGetBoundingClientRect.mockClear();
  });

  test('Header hidden after trigger start scroll', async () => {
    jest.useFakeTimers();

    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: defaultWidget,
        alarms,
        columns,
        stickyHeader: true,
      },
    });

    const header = selectTableHead(wrapper);
    const body = selectTableBody(wrapper);

    const headerGetBoundingClientRect = jest.spyOn(header.element, 'getBoundingClientRect')
      .mockReturnValue({ top: -200 });

    const bodyGetBoundingClientRect = jest.spyOn(body.element, 'getBoundingClientRect')
      .mockReturnValue({ height: 400 });

    triggerWindowScrollEvent(200);

    expect(+header.element.style.opacity).toBe(0);

    jest.runAllTimers();

    expect(+header.element.style.opacity).toBe(1);

    headerGetBoundingClientRect.mockClear();
    bodyGetBoundingClientRect.mockClear();

    jest.useRealTimers();
  });

  test('Expanded elements works correctly', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: defaultWidget,
        alarms,
        columns,
        stickyHeader: true,
      },
    });

    expect(wrapper.vm.expanded).toEqual({});

    const alarmsListRow = selectAlarmsListRow(wrapper).at(0);

    alarmsListRow.triggerCustomEvent('expand', true);

    const [firstAlarm] = alarms;

    expect(wrapper.vm.expanded).toEqual({
      [firstAlarm._id]: true,
    });

    alarmsListRow.triggerCustomEvent('expand', false);

    expect(wrapper.vm.expanded).toEqual({});
  });

  test('Root cause diagram opened after trigger click state event', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: defaultWidget,
        alarms,
        columns,
        stickyHeader: true,
      },
      mocks: { $modals },
    });

    expect(wrapper.vm.expanded).toEqual({});

    selectAlarmsListRow(wrapper).at(0).triggerCustomEvent('click:state', true);

    /**
     * TODO: Should be tested show modal
     */
  });

  test('Renders `alarms-list-table` with default and required props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: defaultWidget,
        alarms: [],
        totalItems: 0,
        columns: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-table` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarms,
        columns,
        totalItems,
        widget: defaultWidget,
        options: {
          page: 1,
          limit: 10,
        },
        loading: true,
        selectable: true,
        hideChildren: true,
        dense: ALARM_DENSE_TYPES.medium,
        expandable: true,
        stickyHeader: true,
        densable: true,
        hidePagination: true,
        hideActions: true,
        parentAlarm: fakeAlarm(),
        refreshAlarmsList: jest.fn(),
        selectedTag: 'tag',
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-table` with expandable, but without selectable', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarms,
        widget: defaultWidget,
        columns,
        totalItems,
        options: {},
        selectable: false,
        expandable: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-table` with default and required props with compact mode', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: {
          ...defaultWidget,
          parameters: {
            ...defaultWidget.parameters,

            dense: ALARM_DENSE_TYPES.medium,
          },
        },
        alarms: [],
        columns,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
  test('Renders `alarms-list-table` with default and required props with ellipsis headers', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: {
          ...defaultWidget,
          parameters: {
            ...defaultWidget.parameters,

            isEllipsisHeaders: true,
          },
        },
        alarms: [],
        columns,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-table` with default and required props with links column with links in row count', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: {
          ...defaultWidget,
          parameters: {
            ...defaultWidget.parameters,

            dense: ALARM_DENSE_TYPES.medium,
          },
        },
        alarms: [],
        columns: [{
          value: ALARM_FIELDS.links,
          inlineLinksCount: 2,
          linksInRowCount: 3,
        }],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-table` with default and required props with simulate ctrl keydown with selectable = false', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: defaultWidget,
        alarms: [],
        totalItems: 0,
        columns,
        selectable: true,
      },
    });

    triggerWindowKeyboardEvent('keydown', { key: 'Control' });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarms-list-table` with default and required props with simulate ctrl keydown with selectable = true', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        options: {},
        widget: defaultWidget,
        alarms: [],
        totalItems: 0,
        columns,
        selectable: true,
      },
    });

    triggerWindowKeyboardEvent('keydown', { key: 'Control' });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
