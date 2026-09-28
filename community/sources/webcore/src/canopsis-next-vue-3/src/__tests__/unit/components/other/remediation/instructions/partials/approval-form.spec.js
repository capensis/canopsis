import { generateShallowRenderer } from '@unit/utils/vue';

import RemediationInstructionApprovalForm from '@/components/other/remediation/instructions/partials/approval-form.vue';

const stubs = {
  'v-form': {
    template: '<form><slot /></form>',
  },
  'v-btn': {
    props: ['color', 'type', 'variant'],
    template: `
      <button :data-color="color" :data-type="type" :data-variant="variant">
        <slot />
      </button>
    `,
  },
  'v-icon': true,
  'v-row': {
    inheritAttrs: false,
    template: '<div v-bind="$attrs"><slot /></div>',
  },
  'v-col': {
    inheritAttrs: false,
    template: '<div v-bind="$attrs"><slot /></div>',
  },
  'c-description-field': true,
};

const selectButtons = wrapper => wrapper.findAll('button');

describe('remediation-instruction-approval-form', () => {
  const factory = generateShallowRenderer(RemediationInstructionApprovalForm, { stubs });

  test('Approve and dismiss actions use the same variant and are aligned to the end', () => {
    const wrapper = factory();
    const [dismissButton, approveButton] = selectButtons(wrapper);

    expect(wrapper.find('.justify-end').exists()).toBe(true);
    expect(dismissButton.attributes('data-variant')).toBeUndefined();
    expect(approveButton.attributes('data-variant')).toBeUndefined();
  });

  test('Dismiss confirmation uses the approve button variant', async () => {
    const wrapper = factory();

    await selectButtons(wrapper)[0].trigger('click');

    const [, dismissButton] = selectButtons(wrapper);

    expect(dismissButton.attributes('data-variant')).toBeUndefined();
    expect(dismissButton.attributes('data-color')).toBe('warning');
    expect(dismissButton.attributes('data-type')).toBe('submit');
  });
});
