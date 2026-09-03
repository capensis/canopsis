<template>
  <td class="permissions-table-cell">
    <v-tooltip
      v-for="(checkbox, index) in checkboxes"
      :key="index"
      :disabled="!checkbox.tooltip"
      location="top"
    >
      <template #activator="{ props: activatorProps }">
        <span v-bind="activatorProps" class="permissions-table-cell__checkbox">
          <c-light-checkbox
            :disabled="disabled || !role.editable || checkbox.disabled"
            :indeterminate="checkbox.indeterminate"
            :label="checkbox.label"
            :model-value="checkbox.modelValue"
            @update:model-value="checkbox.onChange"
          />
        </span>
      </template>

      <span>{{ checkbox.tooltip }}</span>
    </v-tooltip>
  </td>
</template>

<script>
import { computed } from 'vue';

import { CRUD_ACTIONS } from '@/constants';

import { getPermissionCheckboxProps } from '@/helpers/entities/permissions/list';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    permission: {
      type: Object,
      required: true,
    },
    role: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:permission'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const checkboxes = computed(() => props.permission.actions.map((action) => {
        const checkboxProps = getPermissionCheckboxProps(props.role, props.permission, action);

        return {
          label: action !== CRUD_ACTIONS.can ? action : undefined,
          modelValue: checkboxProps.inputValue,
          indeterminate: checkboxProps.indeterminate,
          disabled: checkboxProps.disabled,
          tooltip: checkboxProps.tooltipKey ? t(checkboxProps.tooltipKey) : undefined,
          onChange: value => emit('update:permission', value, props.role, props.permission, action),
        };
      }),
    );

    return {
      checkboxes,
    };
  },
};
</script>

<style lang="scss" scoped>
.permissions-table-cell {
  white-space: nowrap;

  &__checkbox {
    display: block;

    & + & {
      margin-top: 4px;
    }
  }
}
</style>
