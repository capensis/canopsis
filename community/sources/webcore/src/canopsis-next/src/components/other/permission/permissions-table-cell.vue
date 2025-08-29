<template>
  <td>
    <v-tooltip
      v-for="(checkbox, index) in checkboxes"
      :key="index"
      :disabled="!checkbox.tooltip"
      top
    >
      <template #activator="{ on, attrs }">
        <span v-bind="attrs" v-on="on">
          <v-checkbox
            v-bind="checkbox.bind"
            :disabled="disabled || !role.editable || checkbox.disabled"
            class="ma-0 pa-0"
            color="primary"
            hide-details
            v-on="checkbox.on"
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
  setup(props, { emit }) {
    const { t } = useI18n();

    const checkboxes = computed(() => props.permission.actions.map((action) => {
      const checkboxProps = getPermissionCheckboxProps(props.role, props.permission, action);

      return {
        bind: {
          label: action !== CRUD_ACTIONS.can ? action : undefined,
          inputValue: checkboxProps.inputValue,
          indeterminate: checkboxProps.indeterminate,
        },
        disabled: checkboxProps.disabled,
        tooltip: checkboxProps.tooltipKey ? t(checkboxProps.tooltipKey) : undefined,
        on: {
          change: value => emit('input', value, props.role, props.permission, action),
        },
      };
    }));

    return {
      checkboxes,
    };
  },
};
</script>
