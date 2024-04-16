import { generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityListExpandPanel from '@/components/other/availability/partials/availability-list-expand-panel.vue';

const stubs = {
  'availability-history': true,
  'entity-alarms-list-table': true,
};

describe('availability-list-expand-panel', () => {
  const snapshotFactory = generateRenderer(AvailabilityListExpandPanel, { stubs });

  test('Renders `availability-list-expand-panel` with required props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          entity: {},
        },
        activeAlarmsColumns: [{}, {}],
        resolvedAlarmsColumns: [{}],
        interval: {},
      },
    });

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-list-expand-panel` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          entity: {},
        },
        activeAlarmsColumns: [{}],
        resolvedAlarmsColumns: [{}, {}],
        interval: {},
        defaultShowType: AVAILABILITY_SHOW_TYPE.duration,
        displayParameter: AVAILABILITY_DISPLAY_PARAMETERS.downtime,
      },
    });

    await wrapper.activateAllTabs();

    expect(wrapper).toMatchSnapshot();
  });
});
