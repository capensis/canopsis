<template>
  <v-data-table
    :items="items"
    :headers="headers"
    :hide-default-header="indent !== 0"
    :items-per-page="items.length"
    class="permissions-table"
    item-key="_id"
    hide-default-footer
  >
    <template #item="{ item, isExpanded, expand }">
      <tr>
        <td :class="{ [`pl-${indent * 3 + 2}`]: true, 'cursor-pointer': item.children }">
          <c-expand-btn
            v-if="item.children"
            :expanded="isExpanded"
            class="mr-2"
            @expand="expand"
          />
          <span :class="{ 'font-weight-medium': item.children }">
            {{ item.title ? item.title : item._id }}
          </span>
        </td>
        <td v-for="role in roles" :key="role.value">
          <permissions-table-cell
            :role="role"
            :permission="item"
            :disabled="disabled"
            @input="$listeners.input"
          />
        </td>
      </tr>
    </template>
    <template #expanded-item="{ item }">
      <permissions-table
        v-if="item.children"
        :treeview-permissions="item.children"
        :roles="roles"
        :indent="indent + 1"
        @input="$listeners.input"
      />
    </template>
  </v-data-table>
</template>

<script>
import { computed } from 'vue';

import PermissionsTableCell from './permissions-table-cell.vue';

export default {
  name: 'permissions-table',
  components: { PermissionsTableCell },
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    treeviewPermissions: {
      type: Object,
      default: () => ({}),
    },
    roles: {
      type: Array,
      default: () => [],
    },
    indent: {
      type: Number,
      default: 0,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const items = computed(() => Object.values(props.treeviewPermissions));
    const headers = computed(() => [
      { text: '', sortable: false },

      ...props.roles.map(role => ({ text: role.name, value: role._id, sortable: false })),
    ]);

    return {
      items,
      headers,
    };
  },
};
</script>

<style lang="scss" scoped>
$checkboxCellWidth: 112px; // TODO: move to css vars
$cellPadding: 8px 8px; // TODO: move to css vars

.permissions-table ::v-deep {
  .v-data-table__wrapper {
    overflow: unset !important;
    padding-top: 0;

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
      z-index: 1;

      background: var(--v-table-background-base);

      .theme--dark & {
        background: var(--v-table-background-base);
      }
    }
  }

  .v-expansion-panel__body {
    overflow: auto;
  }
}
</style>
