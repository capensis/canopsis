import { generateRenderer } from '@unit/utils/vue';

import AlarmColumnValueExtraDetails from '@/components/widgets/alarm/columns-formatting/alarm-column-value-extra-details.vue';

const stubs = {
  'extra-details-ack': true,
  'extra-details-last-comment': true,
  'extra-details-ticket': true,
  'extra-details-canceled': true,
  'extra-details-snooze': true,
  'extra-details-pbehavior': true,
  'extra-details-parents': true,
  'extra-details-children': true,
};

describe('alarm-column-value-extra-details', () => {
  const fullAlarm = {
    rule: {},
    pbehavior: {
      name: 'pbehavior-name',
      author: 'pbehavior-author',
      tstart: '',
      tstop: '',
      rrule: 'rrule',
      reason: {
        name: 'pbehavior-reason-name',
      },
      comments: [
        {
          _id: 'pbehavior-comment-1-id',
          author: 'pbehavior-comment-1-author',
          message: 'pbehavior-comment-1-message',
        },
        {
          _id: 'pbehavior-comment-2-id',
          author: 'pbehavior-comment-2-author',
          message: 'pbehavior-comment-2-message',
        },
      ],
    },
    v: {
      ack: {},
      last_comment: {
        m: ' ',
      },
      tickets: [{}],
      canceled: {},
      snooze: {},
      pbehavior_info: {
        icon_name: 'icon',
        type_name: 'type',
      },
    },
  };

  const snapshotFactory = generateRenderer(AlarmColumnValueExtraDetails, { stubs });

  it('Renders `alarm-column-value-extra-details` with empty alarm', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          v: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-extra-details` with full alarm (only parent)', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          ...fullAlarm,

          parent: 3,
          meta_alarm_rules: [
            {
              id: 'parent-rule-1-id',
              name: 'parent-rule-1-name',
            },
            {
              id: 'parent-rule-2-id',
              name: 'parent-rule-2-name',
            },
            {
              id: 'parent-rule-3-id',
              name: 'parent-rule-3-name',
            },
          ],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-extra-details` with full alarm (only children)', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          ...fullAlarm,

          children: 3,
          opened_children: 2,
          closed_children: 1,
          meta_alarm_rule: {
            id: 'child-rule-id',
            name: 'child-rule-name',
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
