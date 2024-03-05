<template>
  <v-data-table
    :items="permissions"
    :headers="headers"
    :items-per-page="-1"
    item-key="_id"
    hide-default-footer
  >
    <template #item="{ item }">
      <permission-row
        :permission="item"
        :roles="roles"
        :changed-roles="changedRoles"
        :disabled="disabled"
        @change="$listeners.change"
      />
    </template>
  </v-data-table>
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
