<template>
  <td>
    <v-simple-checkbox
      v-for="(checkbox, index) in checkboxes"
      v-bind="checkbox.bind"
      :key="index"
      :disabled="disabled || disabledForRole"
      class="ma-0 pa-0"
      color="primary"
      hide-details
      v-on="checkbox.on"
    />
  </td>
</template>

<script>
import { CRUD_ACTIONS } from '@/constants';

import { getCheckboxValue, getPermissionActions } from '@/helpers/entities/permissions/list';

export default {
  model: {
    prop: 'permission',
    event: 'change',
  },
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
    changedRole: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    disabledForRole() {
      return !this.role.editable;
    },

    checkboxes() {
      const actions = getPermissionActions(this.permission);

      if (actions) {
        return actions.map(action => ({
          bind: {
            label: action !== CRUD_ACTIONS.can ? action : undefined,

            value: getCheckboxValue(
              this.permission._id,
              this.role.permissions,
              this.changedRole,
              action,
            ),
          },
          on: {
            input: value => this.changeCheckboxValue(value, action),
          },
        }));
      }

      return [
        {
          bind: {
            value: getCheckboxValue(
              this.permission._id,
              this.role.permissions,
              this.changedRole,
            ),
          },
          on: {
            input: value => this.changeCheckboxValue(value),
          },
        },
      ];
    },
  },
  methods: {
    /**
     * Change checkbox value
     *
     * @param {boolean} value
     * @param {string} [action = CRUD_ACTIONS.can]
     */
    changeCheckboxValue(value, action = CRUD_ACTIONS.can) {
      this.$emit('change', value, this.role, this.permission, action);
    },
  },
};
</script>
