import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_LIST_ACTIONS_TYPES } from '@/constants';

import AlarmColumnValueState from '@/components/widgets/alarm/columns-formatting/alarm-column-value-state.vue';

const stubs = {
  'c-alarm-chip': true,
};

const selectAlarmChip = wrapper => wrapper.find('c-alarm-chip-stub');

describe('alarm-column-value-state', () => {
  const snapshotFactory = generateRenderer(AlarmColumnValueState, { stubs });
  const factory = generateShallowRenderer(AlarmColumnValueState, { stubs });

  test('Click emitted after trigger click on chip', () => {
    const wrapper = factory({
      propsData: {
        alarm: {},
      },
    });

    selectAlarmChip(wrapper).triggerCustomEvent('click');

    expect(wrapper).toHaveBeenEmit('click');
  });

  test('Renders `alarm-column-value-state` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-column-value-state` with alarm state', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          v: {
            events_count: 'Events count',
            state: {
              val: 'custom-state-val',
              _t: ALARM_LIST_ACTIONS_TYPES.changeState,
            },
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-column-value-state` with custom propertyKey', () => {
    const wrapper = snapshotFactory({
      propsData: {
        propertyKey: 'customPropertyKey',
        alarm: {
          customPropertyKey: 'custom-property-value',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
