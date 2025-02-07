<template>
  <td>
    <v-checkbox
      v-for="(checkbox, index) in checkboxes"
      v-bind="checkbox.bind"
      :key="index"
      :disabled="disabled || !role.editable"
      class="ma-0 pa-0"
      color="primary"
      hide-details
      v-on="checkbox.on"
    />
  </td>
</template>

<script>
import { computed } from 'vue';

import { CRUD_ACTIONS } from '@/constants';

import { getPermissionCheckboxProps } from '@/helpers/entities/permissions/list';

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
    const checkboxes = computed(() => props.permission.actions.map(action => ({
      bind: {
        label: action !== CRUD_ACTIONS.can ? action : undefined,

        ...getPermissionCheckboxProps(props.role, props.permission, action),
      },
      on: {
        change: value => emit('input', value, props.role, props.permission, action),
      },
    })));

    return {
      checkboxes,
    };
  },
};
</script>
