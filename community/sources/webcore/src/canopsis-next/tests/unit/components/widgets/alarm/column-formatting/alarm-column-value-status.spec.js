import { omit } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

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

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it.each(
    Object.entries(omit(ENTITIES_STATUSES, ['ongoing', 'noEvents'])),
  )('Renders `alarm-column-value-status` with status: %s', async (_, status) => {
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

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `alarm-column-value-status` with status: noEvents', async () => {
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

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
