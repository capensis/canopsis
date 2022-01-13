import { omit } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';

import { ENTITIES_STATES, ENTITIES_STATUSES } from '@/constants';
import AlarmColumnValueStatus from '@/components/widgets/alarm/columns-formatting/alarm-column-value-status.vue';

const localVue = createVueInstance();

const stubs = {
  'c-no-events-icon': true,
};

const snapshotFactory = (options = {}) => mount(AlarmColumnValueStatus, {
  localVue,
  stubs,

  ...options,
});

describe('alarm-column-value-status', () => {
  it.each(Object.entries(ENTITIES_STATES))('Renders `alarm-column-value-status` with ongoing status and state: %s', (_, state) => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
          v: {
            state: {
              val: state,
            },
            status: {
              val: ENTITIES_STATUSES.ongoing,
            },
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it.each(
    Object.entries(omit(ENTITIES_STATUSES, ['ongoing', 'noEvents'])),
  )('Renders `alarm-column-value-status` with status: %s', (_, status) => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {},
          v: {
            state: {
              val: ENTITIES_STATES.ok,
            },
            status: {
              val: status,
            },
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-status` with status: noEvents', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          entity: {
            idle_since: 1386435600000,
          },
          v: {
            status: {
              val: ENTITIES_STATUSES.noEvents,
            },
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
