import { generateRenderer } from '@unit/utils/vue';
import { ENTITIES_STATES, ENTITIES_STATUSES, EVENT_ENTITY_TYPES } from '@/constants';

import TimeLineFlag from '@/components/widgets/alarm/time-line/alarms-time-line-flag.vue';

const stubs = {
  'c-alarm-chip': true,
};

describe('alarms-time-line-flag', () => {
  const snapshotFactory = generateRenderer(TimeLineFlag, { stubs });

  it.each(
    Object.entries(ENTITIES_STATUSES),
  )('Renders `alarms-time-line-flag` with status: %s', (_, status) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: 'status',
          val: status,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it.each(
    Object.entries(ENTITIES_STATES),
  )('Renders `alarms-time-line-flag` with state: %s', (_, status) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: 'state',
          val: status,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it.each(
    Object.entries(EVENT_ENTITY_TYPES),
  )('Renders `alarms-time-line-flag` with type: %s', (_, type) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: type,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
