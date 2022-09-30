<template lang="pug">
  v-data-table(
    :items="sortedPermissions",
    :headers="headers",
    item-key="_id",
    expand,
    hide-actions
  )
    template(#items="{ item }")
      permission-row(
        :permission="item",
        :roles="roles",
        :changed-roles="changedRoles",
        :disabled="disabled",
        @change="$listeners.change"
      )
</template>

<script>
import { sortBy } from 'lodash';

import PermissionRow from './permission-row.vue';

export default {
  components: {
    PermissionRow,
  },
  props: {
    permissions: {
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
      default: () => ['description'],
    },
  },
  computed: {
    headers() {
      return [
        { text: '', sortable: false },

        ...this.roles.map(role => ({ text: role._id, sortable: false })),
      ];
    },

    sortedPermissions() {
      return sortBy(this.permissions, this.sortBy);
    },
  },
};
</script>
