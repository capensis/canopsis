<template lang="pug">
  div.permissions-table-wrapper
    permissions-groups-table(
      :groups="permissions",
      :roles="roles",
      :changed-roles="changedRoles",
      :disabled="disabled",
      @change="$listeners.change"
    )
</template>

<script>
import PermissionsGroupsTable from './permissions-groups-table.vue';

export default {
  components: {
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
