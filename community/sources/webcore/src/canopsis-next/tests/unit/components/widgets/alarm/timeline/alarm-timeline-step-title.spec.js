import { omit } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';

import { ALARM_LIST_STEPS, ALARM_STATUSES } from '@/constants';

import AlarmTimelineStepTitle from '@/components/widgets/alarm/timeline/alarm-timeline-step-title.vue';

describe('alarm-timeline-step-title', () => {
  const snapshotFactory = generateRenderer(AlarmTimelineStepTitle);

  const statusTestCases = Object.values(ALARM_STATUSES).reduce((acc, status) => {
    acc.push([ALARM_LIST_STEPS.statusinc, status], [ALARM_LIST_STEPS.statusdec, status]);

    return acc;
  }, []);

  const testCases = Object.values(omit(ALARM_LIST_STEPS, ['statusinc', 'statusdec', 'snooze']));

  test('Renders `alarm-timeline-step-title` for step `snooze` correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: ALARM_LIST_STEPS.snooze,
          val: 120,
          a: 'root',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-timeline-step-title` for step `webhookstart` with deep flag correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: ALARM_LIST_STEPS.webhookStart,
          a: 'root',
        },
        deep: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test.each(statusTestCases)('Renders `alarm-timeline-step-title` for step `%s` and status `%s` correctly', (type, status) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: type,
          val: status,
          a: 'root',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test.each(testCases)('Renders `alarm-timeline-step-title` for step `%s` correctly', (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        step: {
          _t: type,
          a: 'root',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
