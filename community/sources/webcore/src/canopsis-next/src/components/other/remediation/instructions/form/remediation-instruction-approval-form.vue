<template>
  <v-layout column>
    <v-checkbox
      v-field="needApprove"
      :label="$t('remediation.instruction.requestApproval')"
      :disabled="disabled || required"
      color="primary"
      hide-details
    />
    <template v-if="needApprove">
      <v-layout
        v-if="disabled"
      >
        <span class="text-subtitle-1 grey--text my-4">{{ assignLabel }}: {{ assignValue }}</span>
      </v-layout>
      <v-layout
        v-else
        align-center
      >
        <v-flex xs6>
          <remediation-instruction-approval-type-field
            v-field="approval.type"
            @input="resetErrors"
          />
        </v-flex>
        <v-flex xs5>
          <c-role-field
            v-show="isRoleType"
            v-field="approval.role"
            :required="isRoleType"
            :name="roleFieldName"
            :permission="approvePermission"
            autocomplete
          />
          <c-user-picker-field
            v-show="!isRoleType"
            v-field="approval.user"
            :required="!isRoleType"
            :name="userFieldName"
            :label="$tc('common.user')"
            :permission="approvePermission"
            return-object
          />
        </v-flex>
      </v-layout>
      <v-textarea
        v-field="approval.comment"
        v-validate="'required'"
        :label="$tc('common.comment')"
        :error-messages="errors.collect('comment')"
        :disabled="disabled"
        name="comment"
      />
    </template>
  </v-layout>
</template>

<script>
import { REMEDIATION_INSTRUCTION_APPROVAL_TYPES, USERS_PERMISSIONS } from '@/constants';

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
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      needApprove: !!this.approval?.comment || this.required,
    };
  },
  computed: {
    approvePermission() {
      return USERS_PERMISSIONS.api.remediation.instructionApprove;
    },

    isRoleType() {
      return this.approval.type === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role;
    },

    roleFieldName() {
      return `${this.name}.role`;
    },

    userFieldName() {
      return `${this.name}.user`;
    },

    assignLabel() {
      return this.isRoleType ? this.$tc('common.role') : this.$tc('common.user');
    },

    assignValue() {
      return this.isRoleType ? this.approval.role.name : this.approval.user.display_name;
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
