import flushPromises from 'flush-promises';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import featuresService from '@/services/features';

import AlarmsListRow from '@/components/widgets/alarm/partials/alarms-list-row.vue';

const localVue = createVueInstance();

const stubs = {
  'v-checkbox-functional': true,
  'alarms-list-row-icon': true,
  'alarms-expand-panel-btn': true,
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

const selectExpandButton = wrapper => wrapper.find('alarms-expand-panel-btn-stub');
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

  const featureHasSpy = jest.spyOn(featuresService, 'has')
    .mockReturnValueOnce(true);
  const featureGetSpy = jest.spyOn(featuresService, 'get')
    .mockReturnValueOnce(undefined);
  const featureCallSpy = jest.spyOn(featuresService, 'call')
    .mockReturnValueOnce(() => {});

  afterEach(() => {
    updateQuery.mockClear();
    fetchItem.mockClear();
    featureHasSpy.mockClear();
    featureGetSpy.mockClear();
    featureCallSpy.mockClear();
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

    expandButton.vm.$emit('input', true);

    await flushPromises();

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

    expandButton.vm.$emit('input', false);

    await flushPromises();

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
        selectedTag: 'tag',
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
