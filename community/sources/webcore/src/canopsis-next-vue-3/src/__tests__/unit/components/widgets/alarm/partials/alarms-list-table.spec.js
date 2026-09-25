import { triggerWindowScrollEvent } from '@unit/utils/events';
import { mockModals } from '@unit/utils/mock-hooks';
import {
  createMockedStoreGetters,
  createMockedStoreModules,
  createEntityInfoPropertyModule,
  createActiveViewModule,
} from '@unit/utils/store';
import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { range } from 'lodash';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';

const VDataTableServerStub = {
  name: 'VDataTableServer',
  props: {
    headers: {
      type: Array,
      default: () => [],
    },
    items: {
      type: Array,
      default: () => [],
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:options'],
  template: `
    <div class="v-data-table-server-stub">
      <div class="v-table__wrapper">
        <table>
          <thead>
            <tr>
              <th
                v-for="header in headers"
                :key="header.key"
                :data-value="header.value"
              >
                <slot
                  :name="'header.' + header.key"
                  :column="{ ...header, sortable: header.sortable !== false }"
                  :is-sorted="() => false"
                  :sort-by="[]"
                  :get-sort-icon="() => ''"
                />
              </th>
            </tr>
          </thead>
          <tbody>
            <slot name="body" :items="items" />
          </tbody>
        </table>
      </div>
    </div>
  `,
};

const AlarmsListRowStub = {
  name: 'AlarmsListRow',
  props: ['alarm'],
  emits: ['update:selected', 'update:expanded', 'select:tag', 'remove:tag'],
  template: '<tr class="alarms-list-row-stub"><td>{{ alarm._id }}</td></tr>',
};

const TableExpandRowStub = {
  name: 'CTableExpandRow',
  props: ['expanded', 'colspan'],
  template:
    '<tr v-if="expanded" class="c-table-expand-row-stub"><td :colspan="colspan"><slot /></td></tr>',
};

const MassActionsPanelStub = {
  name: 'CMassActionsPanel',
  emits: ['clear:selected'],
  template: '<div class="c-mass-actions-panel-stub"><slot name="actions" /></div>',
};

const PaginationStub = {
  name: 'CPagination',
  emits: ['update:page'],
  template: '<div class="c-pagination-stub" />',
};

const stubs = {
  'v-data-table-server': VDataTableServerStub,
  'mass-actions-panel': true,
  'c-mass-actions-panel': MassActionsPanelStub,
  'c-empty-data-table-columns': true,
  'alarm-header-cell': true,
  'alarms-list-row': AlarmsListRowStub,
  'alarms-expand-panel': true,
  'c-pagination': PaginationStub,
  'c-table-pagination': true,
  'c-density-btn-toggle': true,
  'c-grid-btns': true,
  'c-light-checkbox': true,
  'c-table-expand-row': TableExpandRowStub,
  'c-horizontal-scrollbar': true,
  'c-layout': {
    template: '<div><slot /></div>',
  },
};

const selectTable = wrapper => wrapper.findComponent({ name: 'VDataTableServer' });
const selectAlarmsListRow = wrapper => wrapper.findAllComponents({ name: 'AlarmsListRow' });
const selectTableHead = wrapper => wrapper.find('thead');
const selectTableBody = wrapper => wrapper.find('tbody');
const selectPagination = wrapper => wrapper.findComponent({ name: 'CPagination' });
const selectMassActionsPanel = wrapper => wrapper.findComponent({ name: 'CMassActionsPanel' });

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
      fetch: vi.fn(() => ({})),
    },
  };

  const { entityInfoPropertyModule } = createEntityInfoPropertyModule();
  const { activeViewModule } = createActiveViewModule();

  const store = createMockedStoreModules([
    associativeTableModule,
    entityInfoPropertyModule,
    activeViewModule,
    createMockedStoreGetters({ name: 'info', showHeaderOnKioskMode: false }),
  ]);

  const defaultWidget = generatePreparedDefaultAlarmListWidget();

  const columns = [
    {
      label: 'Label-1',
      value: 'label',
    },
  ];

  const defaultQuery = {
    page: 1,
    itemsPerPage: 10,
    sortBy: [],
    sortDesc: [],
    search: '',
  };

  const defaultMeta = {
    total_count: totalItems,
  };

  const renderer = generateRenderer(AlarmsListTable, {
    stubs,
    attachTo: document.body,
    propsData: {
      alarms: [],
      columns,
      widget: defaultWidget,
      meta: defaultMeta,
      query: defaultQuery,
    },
  });

  const factory = options => renderer(options);

  afterEach(() => {
    vi.clearAllMocks();
  });

  test('Alarms selected after trigger table', () => {
    const selectedAlarms = alarms.slice(0, -1);
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        selectable: true,
      },
    });

    selectedAlarms.forEach(alarm => wrapper.vm.toggleAlarmSelection(alarm, true));

    expect(wrapper.vm.selectedAlarmIds).toEqual(selectedAlarms.map(({ _id }) => _id));
    expect(wrapper.vm.selectedAlarms).toEqual(selectedAlarms);
  });

  test('Pagination update event emitted after trigger update pagination', async () => {
    const wrapper = factory({
      store,
    });

    const options = {
      page: Faker.datatype.number(),
      itemsPerPage: Faker.datatype.number(),
      sortBy: [
        {
          key: Faker.datatype.string(),
          order: Faker.datatype.boolean() ? 'desc' : 'asc',
        },
      ],
    };

    selectTable(wrapper).vm.$emit('update:options', options);

    expect(wrapper).toEmit('update:query', {
      ...defaultQuery,
      page: options.itemsPerPage <= defaultQuery.itemsPerPage ? options.page : 1,
      itemsPerPage: options.itemsPerPage,
      sortBy: [options.sortBy[0].key],
      sortDesc: [options.sortBy[0].order === 'desc'],
    });
  });

  test('Resize listener added after mount and removed after destroy', async () => {
    const addEventListener = vi.spyOn(window, 'addEventListener');
    const removeEventListener = vi.spyOn(window, 'removeEventListener');

    const wrapper = factory({
      store,
      propsData: {
        stickyHeader: true,
        selectable: true,
      },
    });

    expect(addEventListener).toHaveBeenCalledWith('scroll', wrapper.vm.stickyScrollHandler);
    expect(addEventListener).toHaveBeenCalledWith('resize', wrapper.vm.stickyScrollHandler);
    expect(addEventListener).toHaveBeenCalledWith('keydown', wrapper.vm.handleKeyDown);
    expect(addEventListener).toHaveBeenCalledWith('keyup', wrapper.vm.handleKeyUp);
    expect(addEventListener).toHaveBeenCalledWith('mousedown', wrapper.vm.mousedownHandler);
    expect(addEventListener).toHaveBeenCalledWith('mouseup', wrapper.vm.mouseupHandler);
    expect(addEventListener).toHaveBeenCalledWith('blur', wrapper.vm.disableSelecting);

    await wrapper.setProps({
      stickyHeader: false,
    });

    expect(removeEventListener).toHaveBeenNthCalledWith(
      1,
      'scroll',
      wrapper.vm.stickyScrollHandler,
    );
    removeEventListener.mockClear();

    wrapper.destroy();

    expect(removeEventListener).toHaveBeenCalledWith('scroll', wrapper.vm.stickyScrollHandler);
    expect(removeEventListener).toHaveBeenCalledWith('resize', wrapper.vm.stickyScrollHandler);
    expect(removeEventListener).toHaveBeenCalledWith('keydown', wrapper.vm.handleKeyDown);
    expect(removeEventListener).toHaveBeenCalledWith('keyup', wrapper.vm.handleKeyUp);
    expect(removeEventListener).toHaveBeenCalledWith('mousedown', wrapper.vm.mousedownHandler);
    expect(removeEventListener).toHaveBeenCalledWith('mouseup', wrapper.vm.mouseupHandler);
    expect(removeEventListener).toHaveBeenCalledWith('blur', wrapper.vm.disableSelecting);
  });

  test('Timer cleared after disable sticky', async () => {
    const clearTimeout = vi.spyOn(window, 'clearTimeout');

    const wrapper = factory({
      store,
    });

    const header = selectTableHead(wrapper);
    const body = selectTableBody(wrapper);

    const headerGetBoundingClientRect = vi
      .spyOn(header.element, 'getBoundingClientRect')
      .mockReturnValue({ top: -200 });

    const bodyGetBoundingClientRect = vi
      .spyOn(body.element, 'getBoundingClientRect')
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
    const addEventListener = vi.spyOn(window, 'addEventListener');
    const removeEventListener = vi.spyOn(window, 'removeEventListener');

    const wrapper = factory({
      store,
      propsData: {
        stickyHeader: true,
        selectable: true,
      },
    });

    expect(addEventListener).toHaveBeenCalledWith('scroll', wrapper.vm.stickyScrollHandler);
    expect(addEventListener).toHaveBeenCalledWith('keydown', wrapper.vm.handleKeyDown);

    await wrapper.setProps({
      stickyHeader: false,
    });

    expect(removeEventListener).toHaveBeenCalledWith('scroll', wrapper.vm.stickyScrollHandler);
    removeEventListener.mockClear();

    wrapper.destroy();

    expect(removeEventListener).toHaveBeenCalledWith('scroll', expect.any(Function));
    expect(removeEventListener).toHaveBeenCalledWith('resize', expect.any(Function));
    expect(removeEventListener).toHaveBeenCalledWith('keydown', expect.any(Function));
    expect(removeEventListener).toHaveBeenCalledWith('keyup', expect.any(Function));
    expect(removeEventListener).toHaveBeenCalledWith('mousedown', expect.any(Function));
    expect(removeEventListener).toHaveBeenCalledWith('mouseup', expect.any(Function));
    expect(removeEventListener).toHaveBeenCalledWith('blur', expect.any(Function));
  });

  test('Header position changed after trigger scroll', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        stickyHeader: true,
      },
    });

    const header = selectTableHead(wrapper);
    const body = selectTableBody(wrapper);

    const headerGetBoundingClientRect = vi
      .spyOn(header.element, 'getBoundingClientRect')
      .mockReturnValue({ top: -200 });

    const bodyGetBoundingClientRect = vi
      .spyOn(body.element, 'getBoundingClientRect')
      .mockReturnValue({ height: 400 });

    triggerWindowScrollEvent(200);

    await flushPromises();

    expect(header.element.style).toHaveProperty('transform', 'translateY(248px)');

    headerGetBoundingClientRect.mockClear();
    bodyGetBoundingClientRect.mockClear();
  });

  test('Header hidden after trigger start scroll', async () => {
    vi.useFakeTimers();

    const wrapper = factory({
      store,
      propsData: {
        alarms,
        stickyHeader: true,
      },
    });

    const header = selectTableHead(wrapper);
    const body = selectTableBody(wrapper);

    const headerGetBoundingClientRect = vi
      .spyOn(header.element, 'getBoundingClientRect')
      .mockReturnValue({ top: -200 });

    const bodyGetBoundingClientRect = vi
      .spyOn(body.element, 'getBoundingClientRect')
      .mockReturnValue({ height: 400 });

    triggerWindowScrollEvent(200);

    expect(+header.element.style.opacity).toBe(0);

    vi.runAllTimers();

    expect(+header.element.style.opacity).toBe(1);

    headerGetBoundingClientRect.mockClear();
    bodyGetBoundingClientRect.mockClear();

    vi.useRealTimers();
  });

  test('Expanded elements works correctly', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        stickyHeader: true,
        expandable: true,
      },
    });

    expect(wrapper.vm.expandedAlarmIds).toEqual([]);

    const alarmsListRow = selectAlarmsListRow(wrapper).at(0);

    alarmsListRow.vm.$emit('update:expanded', true);

    await flushPromises();

    const [firstAlarm] = alarms;
    expect(wrapper.vm.expandedAlarmIds).toEqual([firstAlarm._id]);
    expect(wrapper.vm.isAlarmExpanded(firstAlarm)).toBe(true);

    alarmsListRow.vm.$emit('update:expanded', false);
    await flushPromises();

    expect(wrapper.vm.expandedAlarmIds).toEqual([]);
  });

  test('Root cause diagram opened after trigger click state event', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        stickyHeader: true,
      },
      mocks: { $modals },
    });

    expect(wrapper.vm.expandedAlarmIds).toEqual([]);
    expect(selectAlarmsListRow(wrapper)).toHaveLength(alarms.length);
  });

  test('Limits the initial rows rendering and defers the remaining rows', async () => {
    const animationFrames = [];
    const requestAnimationFrameSpy = vi
      .spyOn(window, 'requestAnimationFrame')
      .mockImplementation((callback) => {
        animationFrames.push(callback);

        return animationFrames.length;
      });
    const manyAlarms = range(25).map(value => ({
      ...alarms[value % alarms.length],
      _id: `alarm-${value}`,
    }));
    const wrapper = factory({
      store,
      propsData: {
        alarms: manyAlarms,
      },
    });

    animationFrames.length = 0;
    requestAnimationFrameSpy.mockClear();
    wrapper.vm.setBootedRows(manyAlarms, 0);

    expect(Object.keys(wrapper.vm.bootedRows)).toEqual(
      manyAlarms.slice(0, 20).map(({ _id }) => _id),
    );
    expect(animationFrames).toHaveLength(1);

    await animationFrames.shift()();

    expect(wrapper.vm.bootedRows[manyAlarms[20]._id]).toBe(true);
    expect(animationFrames.length).toBeGreaterThanOrEqual(1);

    requestAnimationFrameSpy.mockRestore();
  });

  test('Top pagination emits query update after page changes', () => {
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        stickyHeader: true,
      },
      mocks: { $modals },
    });

    selectPagination(wrapper).vm.$emit('update:page', 3);

    expect(wrapper).toEmit('update:query', {
      ...defaultQuery,
      page: 3,
    });
  });

  test('Fetches entity info properties list on mount', async () => {
    const customFetchList = vi.fn(() => ({}));
    const customEntityInfoPropertyModule = {
      name: 'entityInfoProperty',
      getters: {
        items: [],
        itemsWithAlias: [],
        itemsWithoutAlias: [],
        meta: {},
        pending: false,
      },
      actions: {
        fetchList: customFetchList,
      },
    };

    const customStore = createMockedStoreModules([
      associativeTableModule,
      customEntityInfoPropertyModule,
      activeViewModule,
      createMockedStoreGetters({ name: 'info', showHeaderOnKioskMode: false }),
    ]);

    factory({
      store: customStore,
      propsData: {
        alarms,
      },
    });

    await flushPromises();

    expect(customFetchList).toHaveBeenCalled();
    expect(
      customFetchList.mock.calls.some(
        call => call[1]?.params?.paginate === false || call[0]?.params?.paginate === false,
      ),
    ).toBe(true);
  });

  test('Shows loading state when entityInfoPropertyPending is true', () => {
    const customEntityInfoPropertyModule = {
      name: 'entityInfoProperty',
      getters: {
        items: [],
        itemsWithAlias: [],
        itemsWithoutAlias: [],
        meta: {},
        pending: true,
      },
      actions: {
        fetchList: vi.fn(() => ({})),
      },
    };

    const customStore = createMockedStoreModules([
      associativeTableModule,
      customEntityInfoPropertyModule,
      activeViewModule,
      createMockedStoreGetters({ name: 'info', showHeaderOnKioskMode: false }),
    ]);

    const wrapper = factory({
      store: customStore,
      propsData: {
        alarms,
      },
    });

    const table = selectTable(wrapper);

    expect(table.props('loading')).toBe(true);
  });

  test('clearSelected with force=true clears selection regardless of keepSelectedAfterAction', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        selectable: true,
      },
    });

    alarms.forEach(alarm => wrapper.vm.toggleAlarmSelection(alarm, true));
    await flushPromises();

    wrapper.vm.keepSelectedAfterAction = true;
    wrapper.vm.clearSelected();

    expect(wrapper.vm.selectedAlarmIds).toEqual([]);
  });

  test('clearSelectedAfterAction clears selection when keepSelectedAfterAction is false', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        selectable: true,
      },
    });

    alarms.forEach(alarm => wrapper.vm.toggleAlarmSelection(alarm, true));
    await flushPromises();

    wrapper.vm.keepSelectedAfterAction = false;
    wrapper.vm.clearSelectedAfterAction();

    expect(wrapper.vm.selectedAlarmIds).toEqual([]);
  });

  test('clearSelectedAfterAction keeps selection when keepSelectedAfterAction is true', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        selectable: true,
      },
    });

    alarms.forEach(alarm => wrapper.vm.toggleAlarmSelection(alarm, true));
    await flushPromises();

    wrapper.vm.keepSelectedAfterAction = true;
    wrapper.vm.clearSelectedAfterAction();

    expect(wrapper.vm.selectedAlarmIds).toEqual(alarms.map(({ _id }) => _id));
  });

  test('keepSelectedAfterAction initialized from defaultKeepSelectedAfterAction prop', () => {
    const wrapperWithKeepSelected = factory({
      store,
      propsData: {
        alarms: [],
        defaultKeepSelectedAfterAction: true,
      },
    });

    expect(wrapperWithKeepSelected.vm.keepSelectedAfterAction).toBe(true);

    const wrapperWithoutKeepSelected = factory({
      store,
      propsData: {
        alarms: [],
        widget: defaultWidget,
      },
    });

    expect(wrapperWithoutKeepSelected.vm.keepSelectedAfterAction).toBe(false);
  });

  test('defaultKeepSelectedAfterAction prop updates keepSelectedAfterAction value', async () => {
    const wrapper = factory({
      store,
    });

    expect(wrapper.vm.keepSelectedAfterAction).toBe(false);

    await wrapper.setProps({ defaultKeepSelectedAfterAction: true });

    expect(wrapper.vm.keepSelectedAfterAction).toBe(true);

    await wrapper.setProps({ defaultKeepSelectedAfterAction: false });

    expect(wrapper.vm.keepSelectedAfterAction).toBe(false);
  });

  test('c-mass-actions-panel clear:selected clears selection', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarms,
        selectable: true,
      },
    });

    alarms.forEach(alarm => wrapper.vm.toggleAlarmSelection(alarm, true));
    await flushPromises();

    const massActionsPanel = selectMassActionsPanel(wrapper);

    massActionsPanel.vm.$emit('clear:selected');

    expect(wrapper.vm.selectedAlarmIds).toEqual([]);
  });
});
