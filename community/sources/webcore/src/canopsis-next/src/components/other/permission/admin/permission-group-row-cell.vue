<template>
  <td>
    <v-simple-checkbox
      v-bind="checkbox.bind"
      :disabled="disabled || disabledForRole"
      class="ma-0 pa-0"
      color="primary"
      hide-details="hideDetails"
      v-on="checkbox.on"
    />
  </td>
</template>

<script>
import { CRUD_ACTIONS } from '@/constants';

import { getPermissionActions, getCheckboxValue } from '@/helpers/entities/permissions/list';

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

          if (!acc.indeterminate && checkboxValue !== acc.value && !(index === 0 && maskIndex === 0)) {
            acc.indeterminate = true;
          }

          acc.value = acc.value || checkboxValue;
        });

        return acc;
      }, {
        value: false,
        indeterminate: false,
      });

      const on = {
        input: value => this.group.permissions.forEach((permission) => {
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
