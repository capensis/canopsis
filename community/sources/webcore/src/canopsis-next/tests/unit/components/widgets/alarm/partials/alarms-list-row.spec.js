import flushPromises from 'flush-promises';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import featuresService from '@/services/features';

import AlarmsListRow from '@/components/widgets/alarm/partials/alarms-list-row.vue';

const localVue = createVueInstance();

const stubs = {
  'v-checkbox-functional': true,
  'alarm-list-row-icon': true,
  'c-expand-btn': true,
  'alarm-column-value': true,
  'actions-panel': true,
};

const factory = (options = {}) => shallowMount(AlarmsListRow, {
  localVue,
  stubs,
  parentComponent: {
    provide: {
      $system: {},
    },
  },

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmsListRow, {
  localVue,
  stubs,
  parentComponent: {
    provide: {
      $system: {},
    },
  },

  ...options,
});

const selectExpandButton = wrapper => wrapper.find('c-expand-btn-stub');
const selectTableRow = wrapper => wrapper.find('tr');
const selectCheckbox = wrapper => wrapper.find('v-checkbox-functional-stub');

describe('alarms-list-row', () => {
  const fetchItem = jest.fn();
  const updateQuery = jest.fn();
  const alarmModule = {
    name: 'alarm',
    getters: {},
    actions: {
      fetchItem,
    },
  };
  const queryModule = {
    name: 'query',
    getters: { getQueryById: () => () => ({ }) },
    actions: {
      update: updateQuery,
    },
  };

  afterEach(() => {
    updateQuery.mockClear();
    fetchItem.mockClear();
  });

  it('Alarm selected after trigger checkbox', () => {
    const wrapper = snapshotFactory({
      propsData: {
        row: {
          item: {
            v: {
              status: {},
            },
          },
          expanded: false,
        },
        widget: {},
        columns: [],
        selectable: true,
      },
    });

    const checkbox = selectCheckbox(wrapper);

    checkbox.vm.$emit('change', true);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(true);
  });

  it('Listeners from feature called after trigger', () => {
    const clickListener = jest.fn();
    const hasFeatureSpy = jest.spyOn(featuresService, 'has')
      .mockReturnValueOnce(false)
      .mockReturnValueOnce(true);
    const callFeatureSpy = jest.spyOn(featuresService, 'call')
      .mockReturnValueOnce({ click: clickListener });

    const wrapper = snapshotFactory({
      propsData: {
        row: {
          item: {
            v: {
              status: {},
            },
          },
          expanded: false,
        },
        widget: {},
        columns: [],
      },
    });

    expect(hasFeatureSpy).toHaveBeenCalledTimes(2);
    expect(callFeatureSpy).toHaveBeenCalledTimes(1);
    expect(hasFeatureSpy).toHaveBeenNthCalledWith(1, 'components.alarmListRow.computed.classes');
    expect(hasFeatureSpy).toHaveBeenNthCalledWith(2, 'components.alarmListRow.computed.listeners');

    const row = selectTableRow(wrapper);

    row.trigger('click');

    expect(clickListener).toHaveBeenCalled();

    hasFeatureSpy.mockClear();
    callFeatureSpy.mockClear();
  });

  it('Row expanded after trigger expand button with correlation, resolved and filtered alarms', async () => {
    const alarm = {
      _id: 'alarm-id',
      causes: {},
      consequences: {},
      filtered_children: ['test'],
      v: {
        resolved: {},
        status: {},
      },
    };
    const row = {
      item: alarm,
      expanded: false,
    };
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        {
          ...queryModule,
          getters: { getQueryById: () => () => ({ correlation: true }) },
        },
      ]),
      propsData: {
        row,
        widget: {},
        columns: [{}, {}],
        expandable: true,
      },
    });

    const expandButton = selectExpandButton(wrapper);

    expandButton.vm.$emit('expand');

    await flushPromises();

    expect(fetchItem).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: alarm._id,
        params: {
          correlation: true,
          sort_dir: 'desc',
          sort_key: 't',
          limit: 1,
          with_instructions: true,
          opened: false,
          with_causes: true,
          with_consequences: true,
          with_steps: true,
        },
        dataPreparer: expect.any(Function),
      },
      undefined,
    );

    const [, options] = fetchItem.mock.calls[0];

    const fetchedAlarm = {
      _id: 'fetched-alarm',
    };

    expect(options.dataPreparer({ data: [fetchedAlarm] })).toEqual([{
      ...fetchedAlarm,
      filtered_children: alarm.filtered_children,
    }]);

    expect(options.dataPreparer({ data: [] })).toEqual([]);

    expect(row.expanded).toBe(true);
  });

  it('Row expanded after trigger expand button with hidden groups and without filtered alarms', async () => {
    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
    };
    const row = {
      item: alarm,
      expanded: false,
    };
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        queryModule,
      ]),
      propsData: {
        row,
        widget: {},
        columns: [{}, {}],
        expandable: true,
        hideGroups: true,
      },
    });

    const expandButton = selectExpandButton(wrapper);

    expandButton.vm.$emit('expand');

    await flushPromises();

    expect(fetchItem).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: alarm._id,
        params: {
          sort_dir: 'desc',
          sort_key: 't',
          limit: 1,
          with_instructions: true,
          correlation: false,
          with_steps: true,
        },
        dataPreparer: expect.any(Function),
      },
      undefined,
    );

    const [, options] = fetchItem.mock.calls[0];

    const fetchedAlarms = [{
      _id: 'fetched-alarm-1',
    }, {
      _id: 'fetched-alarm-2',
    }];

    expect(options.dataPreparer({ data: fetchedAlarms })).toEqual(fetchedAlarms);

    expect(options.dataPreparer({})).toEqual([]);

    expect(row.expanded).toBe(true);
  });

  it('Row expanded after trigger expand button without causes and consequences', async () => {
    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
    };
    const row = {
      item: alarm,
      expanded: false,
    };
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        {
          ...queryModule,
          getters: { getQueryById: () => () => ({ correlation: true }) },
        },
      ]),
      propsData: {
        row,
        widget: {},
        columns: [{}, {}],
        expandable: true,
      },
    });

    const expandButton = selectExpandButton(wrapper);

    expandButton.vm.$emit('expand');

    await flushPromises();

    expect(fetchItem).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: alarm._id,
        params: {
          sort_dir: 'desc',
          sort_key: 't',
          limit: 1,
          with_instructions: true,
          correlation: true,
          with_steps: true,
        },
        dataPreparer: expect.any(Function),
      },
      undefined,
    );

    expect(row.expanded).toBe(true);
  });

  it('Row closed after trigger expand button with expanded: true', async () => {
    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
    };
    const row = {
      item: alarm,
      expanded: true,
    };
    const wrapper = factory({
      propsData: {
        row,
        widget: {},
        columns: [{}, {}],
        expandable: true,
      },
    });

    const expandButton = selectExpandButton(wrapper);

    expandButton.vm.$emit('expand');

    await flushPromises();

    expect(fetchItem).not.toHaveBeenCalled();

    expect(row.expanded).toBe(false);
  });

  it('Renders `alarms-list-row` with default and required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {},
        row: {
          item: {
            v: {
              status: {},
            },
          },
          expanded: false,
        },
        columns: [{}, {}],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selected: true,
        selectable: true,
        expandable: true,
        row: {
          item: {
            v: {
              status: {},
            },
          },
          expanded: false,
        },
        widget: {},
        columns: [{}, {}],
        columnsFilters: [{}, {}],
        isTourEnabled: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with resolved alarm', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selected: true,
        selectable: true,
        row: {
          item: {
            v: {
              status: {
                val: 0,
              },
            },
          },
          expanded: false,
        },
        widget: {},
        columns: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with expand button', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selected: true,
        expandable: true,
        row: {
          item: {
            v: {
              status: {},
            },
          },
          expanded: false,
        },
        widget: {},
        columns: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with instructions', () => {
    const wrapper = snapshotFactory({
      propsData: {
        row: {
          item: {
            assigned_instructions: [{}],
            v: {
              status: {},
            },
          },
          expanded: false,
        },
        widget: {},
        columns: [],
        parentAlarm: {
          children_instructions: true,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with filtered children in parent alarm', () => {
    const alarm = {
      _id: 'alarm-id',
      assigned_instructions: [{}],
      v: {
        status: {},
      },
    };
    const wrapper = snapshotFactory({
      propsData: {
        row: {
          item: alarm,
          expanded: false,
        },
        widget: {},
        columns: [],
        parentAlarm: {
          children_instructions: true,
          filtered_children: [alarm._id],
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with feature classes', () => {
    const hasFeatureSpy = jest.spyOn(featuresService, 'has')
      .mockReturnValueOnce(true)
      .mockReturnValueOnce(false);
    const callFeatureSpy = jest.spyOn(featuresService, 'call')
      .mockReturnValueOnce({ 'class-2': true });

    const wrapper = snapshotFactory({
      propsData: {
        row: {
          item: {
            _id: 'alarm-id',
            v: {
              status: {},
            },
          },
          expanded: false,
        },
        widget: {},
        columns: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();

    hasFeatureSpy.mockClear();
    callFeatureSpy.mockClear();
  });
});
