import { generateRenderer } from '@unit/utils/vue';

import { EVENT_ENTITY_TYPES } from '@/constants';

import AlarmColumnValueState from '@/components/widgets/alarm/columns-formatting/alarm-column-value-state.vue';

const stubs = {
  'c-alarm-chip': true,
};

describe('alarm-column-value-state', () => {
  const snapshotFactory = generateRenderer(AlarmColumnValueState, { stubs });

  it('Renders `alarm-column-value-state` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-state` with alarm state', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          v: {
            events_count: 'Events count',
            state: {
              val: 'custom-state-val',
              _t: EVENT_ENTITY_TYPES.changeState,
            },
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-state` with custom propertyKey', () => {
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
