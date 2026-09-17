<template>
  <v-row align="center">
    <v-col cols="12">
      <c-enabled-field
        v-model="needApprove"
        :label="t('remediation.instruction.requestApproval')"
        :disabled="disabled || required"
        hide-details
      />
    </v-col>
    <template v-if="needApprove">
      <v-col v-if="disabled" cols="12">
        <span class="text-subtitle-1 text-grey my-4"> {{ assignLabel }}: {{ assignValue }} </span>
      </v-col>
      <template v-else>
        <v-col>
          <remediation-instruction-approval-type-field
            v-model="typeField"
            @update:model-value="resetErrors"
          />
        </v-col>
        <v-col>
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
            :label="t('common.user')"
            :permission="approvePermission"
            return-object
          />
        </v-col>
      </template>
      <v-col cols="12">
        <c-textarea
          v-model="commentField"
          :label="t('common.comment')"
          :disabled="disabled"
          name="comment"
          required
        />
      </v-col>
    </template>
  </v-row>
</template>

<script>
import { ref, computed } from 'vue';

import { REMEDIATION_INSTRUCTION_APPROVAL_TYPES, USER_PERMISSIONS } from '@/constants';

import { useFormContext, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import RemediationInstructionApprovalTypeField from './fields/remediation-instruction-approval-type-field.vue';

export default {
  components: { RemediationInstructionApprovalTypeField },
  props: {
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
    const { t } = useI18n();
    const form = useFormContext();
    const { typeField, roleField, userField, commentField } = useFormModels({
      typeField: 'approval.type',
      roleField: 'approval.role',
      userField: 'approval.user',
      commentField: 'approval.comment',
    });

    const needApprove = ref(!!commentField.value || props.required);

    const approvePermission = computed(
      () => USER_PERMISSIONS.technical.exploitation.remediationInstructionApprove,
    );

    const isRoleType = computed(
      () => typeField.value === REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role,
    );

    const roleFieldName = computed(() => `${props.name}.role`);
    const userFieldName = computed(() => `${props.name}.user`);

    const assignLabel = computed(() => (isRoleType.value ? t('common.role') : t('common.user')));

    const assignValue = computed(() => {
      if (isRoleType.value) {
        return roleField.value?.name ?? '';
      }

      return userField.value?.display_name ?? '';
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

      form.setFieldError(removingField, undefined);
    };

    return {
      t,
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
