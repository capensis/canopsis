import { generateShallowRenderer } from '@unit/utils/vue';

import { TIME_UNITS } from '@/constants';

import RemediationInstructionStepField from '@/components/other/remediation/instructions/form/fields/remediation-instruction-step-field.vue';

const stubs = {
  'c-card-iterator-item': {
    name: 'CCardIteratorItem',
    props: ['disabled', 'removable'],
    template: '<div><slot name="header" /><slot /></div>',
  },
  'c-name-field': true,
  'c-text-field': true,
  'c-form-block': true,
  'c-form-block-row': true,
  'c-workflow-field': true,
  'remediation-instruction-step-endpoint-field': true,
  'remediation-instruction-operations-form': true,
};

const selectCardItem = wrapper => wrapper.findComponent({ name: 'CCardIteratorItem' });

describe('remediation-instruction-step-field', () => {
  const factory = generateShallowRenderer(RemediationInstructionStepField, { stubs });
  const step = {
    name: 'Step',
    operations: [{
      time_to_complete: {
        value: 1,
        unit: TIME_UNITS.minute,
      },
    }],
  };

  test.each([
    { disabled: false, removable: true },
    { disabled: true, removable: false },
  ])('Removable is $removable when disabled is $disabled', ({ disabled, removable }) => {
    const wrapper = factory({
      props: {
        disabled,
        name: 'steps.step',
        step,
        stepNumber: 1,
      },
    });

    expect(selectCardItem(wrapper).props('removable')).toBe(removable);
    expect(selectCardItem(wrapper).props('disabled')).toBe(disabled);
  });
});
