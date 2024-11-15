<template>
  <v-data-table
    :items="items"
    :headers="headers"
    :hide-default-header="indent !== 0"
    class="permissions-table"
    hide-default-footer
  >
    <template #item="{ item, index, isExpanded, expand }">
      <tr>
        <td :class="{ [`pl-${indent * 2}`]: true, 'cursor-pointer': item.children?.length }">
          <c-expand-btn
            v-if="item.children?.length"
            :expanded="isExpanded"
            class="mr-2"
            @expand="expand"
          />
          <span>{{ item.text }}</span>
        </td>
        <td v-for="role in roles" :key="role.value">
          {{ items[index].value[role.value] }}
          <v-simple-checkbox
            v-field="items[index].value[role.value]"
            color="primary"
          />
        </td>
      </tr>
    </template>
    <template #expanded-item="{ item, index }">
      <permissions-table
        v-if="item.children?.length"
        v-field="items[index].children"
        :items="item.children"
        :headers="headers"
        :indent="indent + 1"
      />
    </template>
  </v-data-table>
</template>

<script>
import { computed } from 'vue';

export default {
  name: 'permissions-table',
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    headers: {
      type: Array,
      default: () => [],
    },
    indent: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const roles = computed(() => props.headers.slice(1));

    return {
      roles,
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
