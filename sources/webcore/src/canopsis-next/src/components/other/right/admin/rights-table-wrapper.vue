<template lang="pug">
  div.pa-3.rights-table-wrapper
    rights-groups-table(
      v-if="isGroup",
      :groups="rights",
      :roles="roles",
      :changedRoles="changedRoles",
      :disabled="disabled",
      @change="changeCheckboxValue"
    )
    rights-table(
      v-else,
      :rights="rights",
      :roles="roles",
      :changedRoles="changedRoles",
      :disabled="disabled",
      @change="changeCheckboxValue"
    )
</template>

<script>
import RightsTable from './rights-table.vue';
import RightsGroupsTable from './rights-groups-table.vue';

export default {
  components: {
    RightsTable,
    RightsGroupsTable,
  },
  props: {
    rights: {
      type: Array,
      required: true,
    },
    roles: {
      type: Array,
      required: true,
    },
    changedRoles: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isGroup() {
      return this.rights.length && this.rights[0].rights;
    },
  },
  methods: {
    changeCheckboxValue(...args) {
      this.$emit('change', ...args);
    },
  },
};
</script>

<style lang="scss" scoped>
  $checkboxCellWidth: 112px;
  $cellPadding: 8px 8px;

  .rights-table-wrapper /deep/ {
    .v-table__overflow {
      background: black;
      overflow: visible;

      td, th {
        padding: $cellPadding;

        &:not(:first-child) {
          width: $checkboxCellWidth;
        }
      }

      th {
        transition: none;
        position: sticky;
        top: 48px;
        background: white;
        z-index: 1;
      }

      .v-datatable__expand-content .v-table {
        background: #f3f3f3;
      }
    }

    .v-expansion-panel__body {
      overflow: auto;
    }
  }
</style>
