import { generateShallowRenderer } from '@unit/utils/vue';

import { EVENT_ENTITY_TYPES } from '@/constants';

import AlarmWebhookExecutionTimelineFlag from '@/components/other/alarm/partials/alarm-webhook-execution-timeline-flag.vue';

const VBadgeStub = {
  name: 'VBadgeStub',
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
  },
  template: '<div><slot name="badge" /><slot /></div>',
};

describe('alarm-webhook-execution-timeline-flag', () => {
  const factory = generateShallowRenderer(AlarmWebhookExecutionTimelineFlag, {
    stubs: {
      'v-badge': VBadgeStub,
      'v-icon': true,
    },
  });

  it('hides the error marker for a successfully completed webhook', () => {
    const wrapper = factory({
      propsData: {
        step: { _t: EVENT_ENTITY_TYPES.webhookComplete },
      },
    });

    expect(wrapper.findComponent(VBadgeStub).props('modelValue')).toBe(false);
  });

  it('shows the error marker for a failed webhook', () => {
    const wrapper = factory({
      propsData: {
        step: { _t: EVENT_ENTITY_TYPES.webhookFail },
      },
    });

    expect(wrapper.findComponent(VBadgeStub).props('modelValue')).toBe(true);
  });
});
