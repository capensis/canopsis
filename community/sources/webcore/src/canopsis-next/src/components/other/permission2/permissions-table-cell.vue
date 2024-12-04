<template>
  <td>
    <v-simple-checkbox
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
import { difference } from 'lodash';
import { computed } from 'vue';

import { CRUD_ACTIONS } from '@/constants';

const getPermissionCheckboxProps = (role, permission, action) => {
  if (permission.allChildren) {
    const childrenPermissionsDiffs = permission.allChildren
      .map(({ _id: id, actions }) => difference(actions ?? [], role.permissions[id] ?? []).length);
    const value = childrenPermissionsDiffs.every(v => !v);

    return {
      value,
      indeterminate: !value && childrenPermissionsDiffs.some(v => !v),
    };
  }

  return {
    value: role.permissions[permission._id]?.includes(action) ?? false,
  };
};

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
        input: value => emit('input', value, props.role, props.permission, action),
      },
    })));

    return {
      checkboxes,
    };
  },
};
</script>
