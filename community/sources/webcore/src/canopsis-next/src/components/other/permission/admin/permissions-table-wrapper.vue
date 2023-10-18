<template>
  <div class="permissions-table-wrapper">
    <permissions-groups-table
      :groups="permissions"
      :roles="roles"
      :changed-roles="changedRoles"
      :disabled="disabled"
      @change="$listeners.change"
    />
  </div>
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
    .v-data-table__wrapper {
      overflow: unset !important;

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

        background: var(--v-table-background-base);

        .theme--dark & {
          background: var(--v-table-background-base);
        }
      }
    }

    .v-expansion-panel__body {
      overflow: auto;
    }
  }
</style>
