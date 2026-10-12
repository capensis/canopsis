import { generateShallowRenderer } from '@unit/utils/vue';

import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import RemediationInstructionApprovalTabs from '@/components/other/remediation/instructions/partials/approval-tabs.vue';

const stubs = {
  'v-tabs': true,
  'v-tab': true,
  'v-window': true,
  'v-window-item': true,
  'remediation-instruction-general-form': true,
};

const selectInstructionForms = wrapper =>
  wrapper.findAllComponents({
    name: 'remediation-instruction-general-form',
  });

describe('remediation-instruction-approval-tabs', () => {
  const factory = generateShallowRenderer(RemediationInstructionApprovalTabs, { stubs });

  test('Original instruction form includes its steps or jobs', () => {
    const wrapper = factory({
      props: {
        original: { type: REMEDIATION_INSTRUCTION_TYPES.manual },
      },
    });

    expect(selectInstructionForms(wrapper)).toHaveLength(1);
    expect(selectInstructionForms(wrapper)[0].attributes()).toMatchObject({
      disabled: '',
      'disabled-common': '',
      'with-additional': '',
    });
  });

  test('Updated and original instruction forms include their steps or jobs', () => {
    const wrapper = factory({
      props: {
        original: { type: REMEDIATION_INSTRUCTION_TYPES.manual },
        updated: { type: REMEDIATION_INSTRUCTION_TYPES.auto },
      },
    });

    expect(selectInstructionForms(wrapper)).toHaveLength(2);
    expect(
      selectInstructionForms(wrapper).every(form => 'with-additional' in form.attributes()),
    ).toBe(true);
  });
});
