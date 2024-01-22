import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import featuresService from '@/services/features';

import AlarmsListRow from '@/components/widgets/alarm/partials/alarms-list-row.vue';

const stubs = {
  'alarms-list-row-instructions-icon': true,
  'alarms-list-row-bookmark-icon': true,
  'alarms-expand-panel-btn': true,
  'alarm-column-value': true,
  'actions-panel': true,
};

const selectExpandButton = wrapper => wrapper.find('alarms-expand-panel-btn-stub');
const selectTableRow = wrapper => wrapper.find('tr');
const selectCheckbox = wrapper => wrapper.find('.v-simple-checkbox');
const selectAlarmColumnValue = wrapper => wrapper.find('alarm-column-value-stub');

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

  const factory = generateShallowRenderer(AlarmsListRow, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(AlarmsListRow, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

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
        alarm: {
          v: {
            status: {},
          },
        },
        expanded: false,
        widget: {},
        headers: [],
        selectable: true,
      },
    });

    selectCheckbox(wrapper).trigger('click');

    expect(wrapper).toEmit('input', true);
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
        alarm: {
          v: {
            status: {},
          },
        },
        expanded: false,
        widget: {},
        headers: [],
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

    const newExpanded = false;
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        queryModule,
      ]),
      propsData: {
        alarm,
        expanded: false,
        widget: {},
        headers: [{}, {}],
        expandable: true,
        hideGroups: true,
      },
    });

    const expandButton = selectExpandButton(wrapper);

    expandButton.triggerCustomEvent('input', newExpanded);

    await flushPromises();

    expect(wrapper).toEmit('expand', newExpanded);
  });

  it('Row closed after trigger expand button with expanded: true', async () => {
    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
    };

    const newExpanded = false;
    const wrapper = factory({
      propsData: {
        alarm,
        expand: true,
        widget: {},
        headers: [{}, {}],
        expandable: true,
      },
    });

    const expandButton = selectExpandButton(wrapper);

    expandButton.triggerCustomEvent('input', newExpanded);

    await flushPromises();

    expect(wrapper).toEmit('expand', newExpanded);
  });

  it('Click state emitted after trigger click state event', async () => {
    const entity = {
      _id: 'alarm-entity',
    };
    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
      entity,
    };

    const wrapper = factory({
      propsData: {
        alarm,
        expand: true,
        widget: {},
        headers: [{ value: 'first' }, { value: 'second' }],
        expandable: true,
      },
    });

    await flushPromises();

    selectAlarmColumnValue(wrapper).triggerCustomEvent('click:state');

    expect(wrapper).toEmit('click:state');
  });

  it('Renders `alarms-list-row` with default and required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {},
        alarm: {
          v: {
            status: {},
          },
        },
        expanded: false,
        headers: [{ value: 'value1' }, { value: 'value2' }, { value: 'actions' }],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selected: true,
        selectable: true,
        expandable: true,
        alarm: {
          v: {
            status: {},
          },
        },
        expanded: false,
        widget: {},
        headers: [{ value: 'value1' }, { value: 'value2' }, { value: 'actions' }],
        columnsFilters: [{}, {}],
        selectedTag: 'tag',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with resolved alarm', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selected: true,
        selectable: true,
        alarm: {
          v: {
            status: {
              val: 0,
            },
          },
        },
        expanded: false,
        widget: {},
        headers: [{ value: 'actions' }],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with expand button', () => {
    const wrapper = snapshotFactory({
      propsData: {
        selected: true,
        expandable: true,
        alarm: {
          v: {
            status: {},
          },
        },
        expanded: false,
        widget: {},
        headers: [{ value: 'actions' }],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with instructions', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          assigned_instructions: [{}],
          v: {
            status: {},
          },
        },
        expanded: false,
        widget: {},
        headers: [{ value: 'actions' }],
        parentAlarm: {
          children_instructions: true,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
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
        alarm,
        expanded: false,
        showInstructionIcon: true,
        widget: {},
        headers: [{ value: 'actions' }],
        parentAlarm: {
          children_instructions: true,
          filtered_children: [alarm._id],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with bookmarked alarm', () => {
    const alarm = {
      _id: 'alarm-id',
      bookmark: true,
      v: {
        status: {},
      },
    };
    const wrapper = snapshotFactory({
      propsData: {
        alarm,
        widget: {},
        headers: [{ value: 'actions' }],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list-row` with feature classes', () => {
    const hasFeatureSpy = jest.spyOn(featuresService, 'has')
      .mockReturnValueOnce(true)
      .mockReturnValueOnce(false);
    const callFeatureSpy = jest.spyOn(featuresService, 'call')
      .mockReturnValueOnce({ 'class-2': true });

    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          _id: 'alarm-id',
          v: {
            status: {},
          },
        },
        expanded: false,
        widget: {},
        showInstructionIcon: true,
        headers: [{ value: 'actions' }],
      },
    });

    expect(wrapper).toMatchSnapshot();

    hasFeatureSpy.mockClear();
    callFeatureSpy.mockClear();
  });
});
