<template lang="pug">
  v-data-table(
    :items="rights",
    :headers="headers",
    item-key="_id",
    expand,
    hide-actions
  )
    template(slot="items", slot-scope="{ item }")
      right-row(
        :right="item",
        :roles="roles",
        :changedRoles="changedRoles",
        :disabled="disabled",
        @change="$listeners.change"
      )
</template>

<script>
import RightRow from './right-row.vue';

export default {
  components: {
    RightRow,
  },
  props: {
    rights: {
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

        ...this.roles.map(role => ({ text: role._id, sortable: false })),
      ];
    },
  },
};
</script>
