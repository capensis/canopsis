import { omit } from 'lodash';
import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { ENTITIES_STATES, ENTITIES_STATUSES } from '@/constants';

import AlarmColumnValueStatus from '@/components/widgets/alarm/columns-formatting/alarm-column-value-status.vue';

const stubs = {
  'c-no-events-icon': true,
};

describe('alarm-column-value-status', () => {
  const snapshotFactory = generateRenderer(AlarmColumnValueStatus, {
    stubs,
    attachTo: document.body,
  });

  it.each(Object.entries(ENTITIES_STATES))('Renders `alarm-column-value-status` with ongoing status and state: %s', async (_, state) => {
    snapshotFactory({
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

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it.each(
    Object.entries(omit(ENTITIES_STATUSES, ['ongoing', 'noEvents'])),
  )('Renders `alarm-column-value-status` with status: %s', async (_, status) => {
    snapshotFactory({
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

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-status` with status: noEvents', async () => {
    snapshotFactory({
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

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});
