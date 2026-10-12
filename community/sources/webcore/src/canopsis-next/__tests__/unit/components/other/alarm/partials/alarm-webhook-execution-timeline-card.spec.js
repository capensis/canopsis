import { generateShallowRenderer } from '@unit/utils/vue';

import AlarmWebhookExecutionTimelineCard from '@/components/other/alarm/partials/alarm-webhook-execution-timeline-card.vue';

const CExpandBtnStub = {
  name: 'CExpandBtn',
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:modelValue'],
  template: '<button class="c-expand-btn" />',
};

const selectExpandButton = wrapper => wrapper.findComponent(CExpandBtnStub);

describe('alarm-webhook-execution-timeline-card', () => {
  const factory = generateShallowRenderer(AlarmWebhookExecutionTimelineCard, {
    stubs: {
      'c-expand-btn': CExpandBtnStub,
      'c-request-text-information': true,
      'v-card': true,
      'v-card-text': true,
      'v-expand-transition': true,
    },
  });

  it('passes the expanded state to the expand button', () => {
    const wrapper = factory({
      propsData: {
        step: { _t: 'webhookfail' },
        expandable: true,
        expanded: true,
      },
    });

    expect(selectExpandButton(wrapper).props('modelValue')).toBe(true);
  });

  it('emits the requested expanded state when the button is clicked', async () => {
    const wrapper = factory({
      propsData: {
        step: { _t: 'webhookfail' },
        expandable: true,
      },
    });

    await selectExpandButton(wrapper).triggerCustomEvent('update:modelValue', true);

    expect(wrapper).toEmit('expand', true);
  });
});
