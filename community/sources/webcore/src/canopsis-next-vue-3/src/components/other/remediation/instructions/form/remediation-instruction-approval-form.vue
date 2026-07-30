<template>
  <v-layout column>
    <c-enabled-field
      v-model="needApprove"
      :label="$t('remediation.instruction.requestApproval')"
      :disabled="disabled || required"
      hide-details
    />
    <template v-if="needApprove">
      <v-layout v-if="disabled">
        <span class="text-subtitle-1 grey--text my-4"> {{ assignLabel }}: {{ assignValue }} </span>
      </v-layout>
      <v-layout v-else align-center>
        <v-flex xs6>
          <remediation-instruction-approval-type-field
            v-model="typeField"
            @update:model-value="resetErrors"
          />
        </v-flex>
        <v-flex xs5>
          <c-role-field
            v-show="isRoleType"
            v-model="roleField"
            :required="isRoleType"
            :name="roleFieldName"
            :permission="approvePermission"
            autocomplete
          />
          <c-user-picker-field
            v-show="!isRoleType"
            v-model="userField"
            :required="!isRoleType"
            :name="userFieldName"
            :label="$t('common.user')"
            :permission="approvePermission"
            return-object
          />
        </v-flex>
      </v-layout>
      <v-textarea
        v-model="commentField"
        v-validate="'required'"
        :label="$t('common.comment')"
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
  emits: ['update:approval'],
  setup(props, { emit }) {
    const { t } = useI18n();
    const validator = useValidator();
    // TODO: fix that
    const createFieldModel = field =>
      computed({
        get: () => props.approval[field],
        set: value =>
          emit('update:approval', {
            ...props.approval,
            [field]: value,
          }),
      });

    const needApprove = ref(!!(props.approval && props.approval.comment) || props.required);
    const typeField = createFieldModel('type');
    const roleField = createFieldModel('role');
    const userField = createFieldModel('user');
    const commentField = createFieldModel('comment');

    const approvePermission = computed(
      () => USER_PERMISSIONS.technical.exploitation.remediationInstructionApprove,
    );

    const isRoleType = computed(
      () => props.approval && props.approval.type === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role,
    );

    const roleFieldName = computed(() => `${props.name}.role`);
    const userFieldName = computed(() => `${props.name}.user`);

    const assignLabel = computed(() => (isRoleType.value ? t('common.role') : t('common.user')));

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
      const removingField =
        type === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role
          ? userFieldName.value
          : roleFieldName.value;

      if (validator.errors && typeof validator.errors.remove === 'function') {
        validator.errors.remove(removingField);
      }
    };

    return {
      needApprove,
      typeField,
      roleField,
      userField,
      commentField,
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
