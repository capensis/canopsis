<template>
  <v-data-table
    :items="groups"
    :headers="headers"
    item-key="name"
    show-expand
    hide-default-footer
  >
    <template #items="props">
      <permission-group-row
        :expanded="props.expanded"
        :group="props.item"
        :roles="roles"
        :changed-roles="changedRoles"
        :disabled="disabled"
        @change="$listeners.change"
        @expand="props.expanded = !props.expanded"
      />
    </template>
    <template #expand="{ item }">
      <permissions-table
        class="expand-permissions-table"
        :permissions="item.permissions"
        :roles="roles"
        :changed-roles="changedRoles"
        :disabled="disabled"
        @change="$listeners.change"
      />
    </template>
  </v-data-table>
</template>

<script>
import PermissionsTable from './permissions-table.vue';
import PermissionGroupRow from './permission-group-row.vue';

export default {
  components: {
    PermissionsTable,
    PermissionGroupRow,
  },
  props: {
    groups: {
      type: Array,
      default: () => [],
    },
    roles: {
      type: Array,
      default: () => [],
    },
    changedRoles: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    headers() {
      return [
        { text: '', sortable: false },

        ...this.roles.map(role => ({ text: role.name, sortable: false })),
      ];
    },
  },
};
</script>

<style lang="scss" scoped>
  $titleLeftPadding: 36px;

  .expand-permissions-table ::v-deep .v-table__overflow {
    tr td {
      &:first-child {
        padding-left: $titleLeftPadding;
      }
    }

    thead tr {
      height: 0;
      visibility: hidden;

      th {
        position: relative;
        height: 0;
        line-height: 0;
        padding-top: 0;
        padding-bottom: 0;
      }
    }
  }
</style>
