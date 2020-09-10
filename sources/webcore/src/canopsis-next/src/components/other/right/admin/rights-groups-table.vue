<template lang="pug">
  v-data-table(
    :items="sortedGroups",
    :headers="headers",
    item-key="name",
    expand,
    hide-actions
  )
    template(slot="items", slot-scope="props")
      right-group-row(
        :expanded="props.expanded",
        :group="props.item",
        :roles="roles",
        :changedRoles="changedRoles",
        :disabled="disabled",
        @change="$listeners.change",
        @expand="props.expanded = !props.expanded"
      )
    template(slot="expand", slot-scope="{ item }")
      rights-table.expand-rights-table(
        :rights="item.rights",
        :roles="roles",
        :changedRoles="changedRoles",
        :disabled="disabled",
        @change="$listeners.change"
      )
</template>

<script>
import { sortBy } from 'lodash';

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

    groupsWithName() {
      return this.groups.map(({ key, rights }) => ({ rights, name: this.$t(key) }));
    },

    sortedGroups() {
      return sortBy(this.groupsWithName, ['name']);
    },
  },
};
</script>

<style lang="scss" scoped>
  $titleLeftPadding: 36px;

  .expand-rights-table /deep/ .v-table__overflow {
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
