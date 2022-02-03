import { mount, createVueInstance } from '@unit/utils/vue';

import AlarmColumnValueExtraDetails from '@/components/widgets/alarm/columns-formatting/alarm-column-value-extra-details.vue';

const localVue = createVueInstance();

const stubs = {
  'extra-details-ack': true,
  'extra-details-last-comment': true,
  'extra-details-ticket': true,
  'extra-details-canceled': true,
  'extra-details-snooze': true,
  'extra-details-pbehavior': true,
  'extra-details-causes': true,
  'extra-details-consequences': true,
};

const snapshotFactory = (options = {}) => mount(AlarmColumnValueExtraDetails, {
  localVue,
  stubs,

  ...options,
});

describe('alarm-column-value-extra-details', () => {
  it('Renders `alarm-column-value-extra-details` with empty alarm', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          v: {},
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarm-column-value-extra-details` with full alarm', () => {
    const wrapper = snapshotFactory({
      propsData: {
        alarm: {
          rule: {},
          pbehavior: {
            name: 'pbehavior-name',
            author: 'pbehavior-author',
            tstart: '',
            tstop: '',
            rrule: 'rrule',
            type: {
              name: 'pbehavior-type-name',
              icon_name: 'pbehavior-type-icon',
            },
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
          causes: {},
          consequences: {},
          v: {
            ack: {},
            lastComment: {
              m: ' ',
            },
            ticket: {},
            canceled: {},
            snooze: {},
          },
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
