<template lang="pug">
  div.permissions-table-wrapper
    permissions-groups-table(
      v-if="isGroup",
      :groups="permissions",
      :roles="roles",
      :changed-roles="changedRoles",
      :disabled="disabled",
      :sort-by="sortBy",
      @change="$listeners.change"
    )
    permissions-table(
      v-else,
      :permissions="permissions",
      :roles="roles",
      :changed-roles="changedRoles",
      :disabled="disabled",
      :sort-by="sortBy",
      @change="$listeners.change"
    )
</template>

<script>
import PermissionsTable from './permissions-table.vue';
import PermissionsGroupsTable from './permissions-groups-table.vue';

export default {
  components: {
    PermissionsTable,
    PermissionsGroupsTable,
  },
  props: {
    permissions: {
      type: Array,
      required: true,
    },
    roles: {
      type: Array,
      required: true,
    },
    changedRoles: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    sortBy: {
      type: [Function, Array, String],
      required: false,
    },
  },
  computed: {
    isGroup() {
      return this.permissions.length && this.permissions[0].permissions;
    },
  },
};
</script>

<style lang="scss" scoped>
  $checkboxCellWidth: 112px;
  $cellPadding: 8px 8px;

  .permissions-table-wrapper ::v-deep {
    .v-table__overflow {
      overflow: visible;

      td, th {
        padding: $cellPadding;

        &:not(:first-child) {
          width: $checkboxCellWidth;
        }
      }

      th {
        transition: none;
        position: sticky;
        top: 48px;
        z-index: 1;

        background: white;

        .theme--dark & {
          background: #424242;
        }
      }
    }

    .v-expansion-panel__body {
      overflow: auto;
    }
  }
</style>
