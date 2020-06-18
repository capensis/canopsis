<template lang="pug">
  v-data-table(
    :items="groups",
    :headers="headers",
    item-key="key",
    expand,
    hide-actions
  )
    template(slot="items", slot-scope="props")
      right-group-row(
        :group="props.item",
        :roles="roles",
        :changedRoles="changedRoles",
        @change="change",
        @click="props.expanded = !props.expanded"
      )
    template(slot="expand", slot-scope="{ item }")
      rights-table.expand-rights-table(
        :rights="item.rights",
        :roles="roles",
        :changedRoles="changedRoles",
        @change="change"
      )
</template>

<script>
import RightsTable from './rights-table.vue';
import RightGroupRow from './right-group-row.vue';

export default {
  components: {
    RightsTable,
    RightGroupRow,
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
  },
  computed: {
    headers() {
      return [
        { text: '', sortable: false },

        ...this.roles.map(role => ({ text: role._id, sortable: false })),
      ];
    },
  },
  methods: {
    change(...args) {
      this.$emit('change', ...args);
    },
  },
};
</script>
