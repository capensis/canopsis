<template>
  <v-data-table
    :items="groups"
    :headers="headers"
    :expanded.sync="expanded"
    item-key="name"
    loader-height="2"
    hide-default-footer
    disable-pagination
  >
    <template #item="props">
      <permission-group-row
        :expanded="props.isExpanded"
        :group="props.item"
        :roles="roles"
        :changed-roles="changedRoles"
        :disabled="disabled"
        @change="$listeners.change"
        @expand="props.expand"
      />
    </template>
    <template #expanded-item="{ item }">
      <permissions-table
        v-show="expanded.find(({ name }) => item.name === name)"
        :key="`expanded-${item.name}`"
        :permissions="item.permissions"
        :roles="roles"
        :changed-roles="changedRoles"
        :disabled="disabled"
        class="expand-permissions-table"
        @change="$listeners.change"
      />
    </template>
  </v-data-table>
</template>

<script>
import PermissionsTable from './permissions-table.vue';
import PermissionGroupRow from './permission-group-row.vue';

/**
 * @TODO: use group instead of expand
 */
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
  data() {
    return {
      expanded: [],
    };
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

  .expand-permissions-table ::v-deep .v-data-table__wrapper {
    tr td {
      &:first-child {
        padding-left: $titleLeftPadding;
      }
    }

    .v-data-table-header {
      display: none;

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
