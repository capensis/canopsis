<template lang="pug">
  v-data-table(
    :items="permissions",
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
