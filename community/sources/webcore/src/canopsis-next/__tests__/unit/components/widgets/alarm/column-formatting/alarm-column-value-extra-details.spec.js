import { generateShallowRenderer } from '@unit/utils/vue';

import AlarmColumnValueExtraDetails from '@/components/widgets/alarm/columns-formatting/alarm-column-value-extra-details.vue';
const stubs = {
  'extra-details-ack': true,
  'extra-details-last-comment': true,
  'extra-details-ticket': true,
  'extra-details-canceled': true,
  'extra-details-snooze': true,
  'extra-details-pause': {
    name: 'extra-details-pause',
    props: ['pause', 'entity'],
    template: '<div />',
  },
  'extra-details-pbehavior': true,
  'extra-details-parents': true,
  'extra-details-children': true,
  'extra-details-close-delay': true,
  'extra-details-close-delay-value': true,
};

describe('alarm-column-value-extra-details', () => {
  const factory = generateShallowRenderer(AlarmColumnValueExtraDetails, { stubs });

  it('shows new pause details from alarm.v.pause and passes the entity', () => {
    const pause = {
      m: 'Maintenance window',
      pause_type: 'maintenance',
    };
    const entity = { name: 'API', type: 'component' };
    const wrapper = factory({
      propsData: {
        alarm: {
          v: { pause },
          entity,
        },
      },
    });

    expect(wrapper.findComponent({ name: 'extra-details-pause' }).props('pause')).toEqual(pause);
    expect(wrapper.findComponent({ name: 'extra-details-pause' }).props('entity')).toEqual(entity);
  });

  it('shows alarm pause details even when entity.pause is absent', () => {
    const wrapper = factory({
      propsData: {
        alarm: { v: { pause: { m: 'Pause comment' } }, entity: {} },
      },
    });

    expect(wrapper.findComponent({ name: 'extra-details-pause' }).exists()).toBe(true);
  });

  it('keeps a legacy pause pbehavior in the pbehavior details component', () => {
    const wrapper = factory({
      propsData: {
        alarm: {
          v: {
            pbehavior_info: {
              canonical_type: 'pause',
            },
          },
          entity: {},
        },
      },
    });

    expect(wrapper.findComponent({ name: 'extra-details-pause' }).exists()).toBe(false);
    expect(wrapper.findComponent({ name: 'extra-details-pbehavior' }).exists()).toBe(true);
  });
});
