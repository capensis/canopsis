<template lang="pug">
  v-layout(column)
    v-checkbox(
      v-field="approval.need_approve",
      :label="$t('remediationInstructions.requestApproval')",
      color="primary",
      hide-details
    )
    template(v-if="approval.need_approve")
      v-layout(row, align-center)
        v-flex(xs6)
          remediation-instruction-approval-type-field(v-field="approval.type", @input="resetErrors")
        v-flex(xs5)
          c-role-picker-field(
            v-show="isRoleType",
            v-field="approval.role",
            :required="isRoleType",
            :name="roleFieldName",
            :label="$tc('common.role')"
          )
          c-user-picker-field(
            v-show="!isRoleType",
            v-field="approval.user",
            :required="!isRoleType",
            :name="userFieldName",
            :label="$tc('common.user')"
          )
      v-textarea(
        v-field="approval.comment",
        v-validate="'required'",
        :label="$tc('common.comment')",
        :error-messages="errors.collect('comment')",
        name="comment"
      )
</template>

<script>
import { REMEDIATION_INSTRUCTION_APPROVAL_TYPES } from '@/constants';

import RemediationInstructionApprovalTypeField from './fields/remediation-instruction-approval-type-field.vue';

export default {
  inject: ['$validator'],
  components: { RemediationInstructionApprovalTypeField },
  model: {
    prop: 'approval',
    event: 'input',
  },
  props: {
    approval: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'approval',
    },
  },
  computed: {
    isRoleType() {
      return this.approval.type === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role;
    },

    roleFieldName() {
      return `${this.name}.role`;
    },

    userFieldName() {
      return `${this.name}.user`;
    },
  },
  methods: {
    resetErrors(type) {
      const removingField = type === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role
        ? this.userFieldName
        : this.roleFieldName;

      this.errors.remove(removingField);
    },
  },
};
</script>
