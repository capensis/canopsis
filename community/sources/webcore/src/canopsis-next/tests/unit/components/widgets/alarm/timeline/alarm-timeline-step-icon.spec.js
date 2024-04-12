import { generateRenderer } from '@unit/utils/vue';

import { ALARM_LIST_STEPS, ALARM_STATUSES } from '@/constants';

import AlarmTimelineStepIcon from '@/components/widgets/alarm/timeline/alarm-timeline-step-icon.vue';

const stubs = {
  'c-alarm-chip': true,
  'c-alarm-pbehavior-chip': true,
  'c-alarm-extra-details-chip': true,
  'v-icon': true,
};

describe('alarm-timeline-step-icon', () => {
  const snapshotFactory = generateRenderer(AlarmTimelineStepIcon, { stubs });

  test.each(Object.values(ALARM_LIST_STEPS))('Renders `alarm-timeline-step-icon` for step `%s` correctly', (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: type,
          val: 1,
          color: 'color',
          icon_name: 'icon_name',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test.each(Object.values(ALARM_STATUSES))('Renders `alarm-timeline-step-icon` for step `statusinc` and status `%s` correctly', (status) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: ALARM_LIST_STEPS.statusinc,
          val: status,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
