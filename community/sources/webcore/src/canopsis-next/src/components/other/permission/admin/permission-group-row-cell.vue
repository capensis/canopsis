<template lang="pug">
  td
    v-checkbox(
      v-bind="checkbox.bind",
      v-on="checkbox.on",
      :disabled="disabled || disabledForRole",
      color="primary",
      hideDetails
    )
</template>

<script>
import { CRUD_ACTIONS } from '@/constants';

import { getPermissionActions, getCheckboxValue } from '@/helpers/permission';

export default {
  model: {
    prop: 'group',
    event: 'change',
  },
  props: {
    group: {
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
    changedRole: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    disabledForRole() {
      return !this.role.editable;
    },

    checkbox() {
      const bind = this.group.permissions.reduce((acc, permission, index) => {
        const actions = getPermissionActions(permission);

        actions.forEach((action, maskIndex) => {
          const checkboxValue = getCheckboxValue(
            permission._id,
            this.role.permissions,
            this.changedRole,
            action,
          );

          if (!acc.indeterminate && checkboxValue !== acc.inputValue && !(index === 0 && maskIndex === 0)) {
            acc.indeterminate = true;
          }

          acc.inputValue = acc.inputValue || checkboxValue;
        });

        return acc;
      }, {
        inputValue: false,
        indeterminate: false,
      });

      const on = {
        change: value => this.group.permissions.forEach((permission) => {
          const actions = getPermissionActions(permission);

          actions.forEach(
            action => getCheckboxValue(permission._id, this.role.permissions, this.changedRole, action) !== value
            && this.changeCheckboxValue(value, permission, action),
          );
        }),
      };

      return {
        bind,
        on,
      };
    },
  },
  methods: {
    /**
     * Change checkbox value
     *
     * @param {boolean} value
     * @param {Object} permission
     * @param {string} [action = CRUD_ACTIONS.default]
     */
    changeCheckboxValue(value, permission, action = CRUD_ACTIONS.can) {
      this.$emit('change', value, this.role, permission, action);
    },
  },
};
</script>
