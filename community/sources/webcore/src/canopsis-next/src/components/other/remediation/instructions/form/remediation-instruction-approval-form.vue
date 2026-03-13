<template>
  <v-layout column>
    <v-checkbox
      v-model="needApprove"
      :label="$t('remediation.instruction.requestApproval')"
      :disabled="disabled || required"
      color="primary"
      hide-details
    />
    <template v-if="needApprove">
      <v-layout v-if="disabled">
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
import { ref, computed } from 'vue';

import { REMEDIATION_INSTRUCTION_APPROVAL_TYPES, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';

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
  setup(props) {
    const { tc } = useI18n();
    const validator = useValidator();

    const needApprove = ref(!!(props.approval && props.approval.comment) || props.required);

    const approvePermission = computed(() => USER_PERMISSIONS.technical.exploitation.remediationInstructionApprove);

    const isRoleType = computed(() => (
      props.approval
      && props.approval.type === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role
    ));

    const roleFieldName = computed(() => `${props.name}.role`);
    const userFieldName = computed(() => `${props.name}.user`);

    const assignLabel = computed(() => (isRoleType.value ? tc('common.role') : tc('common.user')));

    const assignValue = computed(() => {
      if (isRoleType.value) {
        return props.approval && props.approval.role ? props.approval.role.name : '';
      }
      return props.approval && props.approval.user ? props.approval.user.display_name : '';
    });

    /**
     * Clear validate errors for the irrelevant approval assignment field
     * when the approval type changes.
     *
     * If the type is 'role', it clears errors for the user field; otherwise it
     * clears errors for the role field.
     *
     * @param {'role'|'user'} type - Newly selected approval type.
     * @returns {void}
     */
    const resetErrors = (type) => {
      const removingField = type === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role
        ? userFieldName.value
        : roleFieldName.value;

      if (validator.errors && typeof validator.errors.remove === 'function') {
        validator.errors.remove(removingField);
      }
    };

    return {
      needApprove,
      approvePermission,
      isRoleType,
      roleFieldName,
      userFieldName,
      assignLabel,
      assignValue,
      resetErrors,
    };
  },
};
</script>
