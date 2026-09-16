import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { featuresService } from '@/services/features';

import AlarmsListRow from '@/components/widgets/alarm/partials/alarms-list-row.vue';

const stubs = {
  'alarms-list-row-instructions-icon': true,
  'alarms-list-row-bookmark-icon': true,
  'alarms-expand-panel-btn': true,
  'alarm-column-value': true,
  'actions-panel': true,
  'c-booted-placeholder-loader': true,
};

const selectExpandButton = wrapper => wrapper.find('alarms-expand-panel-btn-stub');
const selectTableRow = wrapper => wrapper.find('tr');
const selectCheckbox = wrapper => wrapper.find('.v-simple-checkbox');
const selectAlarmColumnValue = wrapper => wrapper.find('alarm-column-value-stub');

describe('alarms-list-row', () => {
  const fetchItem = vi.fn();
  const updateQuery = vi.fn();
  const alarmModule = {
    name: 'alarm',
    getters: {},
    actions: {
      fetchItem,
    },
  };
  const queryModule = {
    name: 'query',
    getters: { getQueryById: () => () => ({}) },
    actions: {
      update: updateQuery,
    },
  };

  const featureHasSpy = vi.spyOn(featuresService, 'has').mockReturnValueOnce(true);
  const featureGetSpy = vi.spyOn(featuresService, 'get').mockReturnValueOnce(undefined);
  const featureCallSpy = vi.spyOn(featuresService, 'call').mockReturnValueOnce(() => {});

  const factory = generateShallowRenderer(AlarmsListRow, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
        $intersectionObserver: {
          observe: vi.fn(),
          unobserve: vi.fn(),
        },
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
    const wrapper = factory({
      propsData: {
        alarm: {
          v: {
            status: {},
          },
        },
        booted: true,
        expanded: false,
        widget: {},
        headers: [],
        selectable: true,
      },
    });

    selectCheckbox(wrapper).trigger('click');

    expect(wrapper).toEmitInput(true);
  });

  it('Listeners from feature called after trigger', () => {
    const clickListener = vi.fn();
    const hasFeatureSpy = vi
      .spyOn(featuresService, 'has')
      .mockReturnValueOnce(false)
      .mockReturnValueOnce(true);
    const callFeatureSpy = vi
      .spyOn(featuresService, 'call')
      .mockReturnValueOnce({ click: clickListener });

    const wrapper = factory({
      propsData: {
        alarm: {
          v: {
            status: {},
          },
        },
        booted: true,
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
      store: createMockedStoreModules([alarmModule, queryModule]),
      propsData: {
        alarm,
        booted: true,
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
        booted: true,
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
        booted: true,
        expand: true,
        widget: {},
        headers: [{ value: 'first' }, { value: 'second' }],
        expandable: true,
      },
    });

    await flushPromises();

    selectAlarmColumnValue(wrapper).triggerCustomEvent('click:state');

    expect(wrapper).toHaveBeenEmit('click:state');
  });

  it('Expand panel opened after double click on row when openExpandPanelByDoubleClick enabled', async () => {
    const showExpandPanel = vi.fn();
    const hasFeatureSpy = vi.spyOn(featuresService, 'has').mockReturnValue(false);

    const alarm = {
      _id: 'alarm-id',
      v: {
        status: {},
      },
    };

    const wrapper = factory({
      propsData: {
        alarm,
        booted: true,
        expanded: false,
        widget: {
          parameters: {
            openExpandPanelByDoubleClick: true,
          },
        },
        headers: [{ value: 'first' }],
        expandable: true,
      },
    });

    await flushPromises();

    wrapper.vm.$refs.expandPanelBtn = { showExpandPanel };

    selectTableRow(wrapper).trigger('dblclick');

    expect(showExpandPanel).toHaveBeenCalledTimes(1);

    hasFeatureSpy.mockRestore();
  });
});
