<template lang="pug">
  v-data-table(
    :items="sortedGroups",
    :headers="headers",
    item-key="name",
    expand,
    hide-actions
  )
    template(#items="props")
      permission-group-row(
        :expanded="props.expanded",
        :group="props.item",
        :roles="roles",
        :changed-roles="changedRoles",
        :disabled="disabled",
        @change="$listeners.change",
        @expand="props.expanded = !props.expanded"
      )
    template(#expand="{ item }")
      permissions-table.expand-permissions-table(
        :permissions="item.permissions",
        :roles="roles",
        :changed-roles="changedRoles",
        :disabled="disabled",
        :sort-by="sortBy",
        @change="$listeners.change"
      )
</template>

<script>
import { sortBy } from 'lodash';

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
    sortBy: {
      type: [Function, Array, String],
      default: () => ['name'],
    },
  },
  computed: {
    headers() {
      return [
        { text: '', sortable: false },

        ...this.roles.map(role => ({ text: role._id, sortable: false })),
      ];
    },

    groupsWithName() {
      return this.groups.map(({ key, name, permissions }) => ({ permissions, name: name ?? this.$tc(key) }));
    },

    sortedGroups() {
      return sortBy(this.groupsWithName, this.sortBy);
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
