import { mount, createVueInstance } from '@unit/utils/vue';

import TimeLineFlag from '@/components/widgets/alarm/time-line/time-line-flag.vue';
import { ENTITIES_STATES, ENTITIES_STATUSES, EVENT_ENTITY_TYPES } from '@/constants';

const localVue = createVueInstance();

const stubs = {
  'c-alarm-chip': true,
};

const snapshotFactory = (options = {}) => mount(TimeLineFlag, {
  localVue,
  stubs,

  ...options,
});

describe('time-line-flag', () => {
  it.each(
    Object.entries(ENTITIES_STATUSES),
  )('Renders `time-line-flag` with status: %s', (_, status) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: 'status',
          val: status,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it.each(
    Object.entries(ENTITIES_STATES),
  )('Renders `time-line-flag` with state: %s', (_, status) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: 'state',
          val: status,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it.each(
    Object.entries(EVENT_ENTITY_TYPES),
  )('Renders `time-line-flag` with type: %s', (_, type) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: type,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
