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
        :roles="sortedRoles",
        :changedRoles="changedRoles",
        :disabled="disabled",
        @change="$listeners.change"
      )
</template>

<script>
import { sortBy } from 'lodash';

import sortRightsMixin from '@/mixins/rights/entities/sort-headers';

import RightRow from './right-row.vue';

export default {
  components: {
    RightRow,
  },
  mixins: [sortRightsMixin],
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
    sortedRights() {
      return sortBy(this.rights, ['name']);
    },

    sortedRoles() {
      return sortBy(this.roles, [({ _id: name }) => name.toLowerCase()]);
    },
  },
};
</script>
