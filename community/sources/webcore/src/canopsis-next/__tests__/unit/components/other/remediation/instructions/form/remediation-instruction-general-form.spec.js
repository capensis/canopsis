import { generateShallowRenderer } from '@unit/utils/vue';

import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import RemediationInstructionGeneralForm from '@/components/other/remediation/instructions/form/remediation-instruction-general-form.vue';

const stubs = {
  'c-name-field': true,
  'c-form-block': true,
  'c-form-block-row': true,
  'c-instruction-type-field': true,
  'c-description-field': true,
  'c-duration-field': true,
  'c-priority-field': true,
  'c-enabled-field': true,
  'c-number-field': true,
  'c-triggers-field': true,
  'c-help-icon': true,
  'v-expand-transition': true,
  'remediation-instruction-approval-form': true,
  'remediation-instruction-steps-form': true,
  'remediation-instruction-jobs-form': true,
};

const selectStepsForm = wrapper => wrapper.find('remediation-instruction-steps-form-stub');
const selectJobsForm = wrapper => wrapper.find('remediation-instruction-jobs-form-stub');

describe('remediation-instruction-general-form', () => {
  const factory = generateShallowRenderer(RemediationInstructionGeneralForm, { stubs });

  test('Additional form is not rendered by default', () => {
    const wrapper = factory({
      props: {
        form: { type: REMEDIATION_INSTRUCTION_TYPES.manual },
      },
    });

    expect(selectStepsForm(wrapper).exists()).toBe(false);
    expect(selectJobsForm(wrapper).exists()).toBe(false);
  });

  test('Steps are rendered for a manual instruction when additional content is enabled', () => {
    const wrapper = factory({
      props: {
        form: { type: REMEDIATION_INSTRUCTION_TYPES.manual },
        withAdditional: true,
      },
    });

    expect(selectStepsForm(wrapper).exists()).toBe(true);
    expect(selectJobsForm(wrapper).exists()).toBe(false);
  });

  test.each([REMEDIATION_INSTRUCTION_TYPES.auto, REMEDIATION_INSTRUCTION_TYPES.simpleManual])(
    'Jobs are rendered for instruction type %s when additional content is enabled',
    (type) => {
      const wrapper = factory({
        props: {
          form: { type },
          disabled: true,
          withAdditional: true,
        },
      });

      expect(selectStepsForm(wrapper).exists()).toBe(false);
      expect(selectJobsForm(wrapper).exists()).toBe(true);
    },
  );
});
