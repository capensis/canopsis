<template>
  <v-data-table
    :items="preparedItems"
    :headers="headers"
    :hide-default-header="indent !== 0"
    :items-per-page="preparedItems.length"
    :expanded.sync="expanded"
    class="permissions-table"
    item-key="_id"
    hide-default-footer
  >
    <template #item="{ item, isExpanded, expand }">
      <tr>
        <td :class="{ [`pl-${indent * 3 + 2}`]: true }">
          <c-expand-btn
            v-if="item.children"
            :expanded="isExpanded"
            class="mr-2"
            @expand="expand"
          />
          <span
            :class="{ 'font-weight-medium': item.children, 'cursor-pointer': item.children }"
            @click="item.children && expand(!isExpanded)"
          >
            <v-list-item-mask v-if="item.hasMask" :text="item.title" :mask="search" />
            <span v-else>{{ item.title }}</span>
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
        :items="item.children"
        :roles="roles"
        :indent="indent + 1"
        :search="search"
        :search-depth="searchDepth"
        @input="$listeners.input"
      />
    </template>
  </v-data-table>
</template>

<script>
import { computed, ref, inject, watch } from 'vue';

import PermissionsTableCell from './permissions-table-cell.vue';

export default {
  name: 'permissions-table',
  components: { PermissionsTableCell },
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    items: {
      type: Array,
      default: () => [],
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
    search: {
      type: String,
      default: '',
    },
    searchDepth: {
      type: Number,
      required: false,
    },
  },
  setup(props) {
    const preparedItems = computed(() => props.items.map(item => ({
      ...item,

      hasMask: props.search
      && (props.searchDepth === props.indent || (!props.searchDepth && !item.children)),
    })));

    const headers = computed(() => [
      { text: '', sortable: false },

      ...props.roles.map(role => ({ text: role.name, value: role._id, sortable: false })),
    ]);

    /**
     * Expand/collapse functionality for permissions table
     * Listens to global allExpandedCounter from parent and updates expanded state accordingly
     * Uses requestAnimationFrame for smooth UI updates based on indent level
     */
    const allExpandedCounter = inject('$allExpandedCounter', 0);

    const expanded = ref([]);

    /**
     * Updates the expanded state based on allExpandedCounter value
     * Expands all items if counter is positive, collapses all otherwise
     */
    const checkExpanded = () => expanded.value = allExpandedCounter.value > 0 ? [...preparedItems.value] : [];

    watch(allExpandedCounter, () => window.requestAnimationFrame(checkExpanded, props.indent), { immediate: true });

    return {
      expanded,

      preparedItems,
      headers,
    };
  },
};
</script>

<style lang="scss" scoped>
.permissions-table ::v-deep {
  --topBarHeight: 48px;
  --checkboxCellWidth: 112px;
  --cellPadding: 8px 8px;

  .v-data-table__wrapper {
    overflow: unset !important;
    padding-top: 0;

    td, th {
      padding: var(--cellPadding);

      &:not(:first-child) {
        width: var(--checkboxCellWidth);
      }
    }

    th {
      transition: none;
      z-index: 1;

      background: var(--v-table-background-base);

      .theme--dark & {
        background: var(--v-table-background-base);
      }

      .v-window__container:not(.v-window__container--is-active) & {
        position: sticky;
        top: calc((var(--topBarHeight) * 2) - 1px);

        .v-app--side-bar-groups & {
          top: calc(var(--topBarHeight) - 1px);
        }
      }
    }
  }

  .v-expansion-panel__body {
    overflow: auto;
  }

  .v-input--selection-controls__input {
    margin: 0;
  }

  &.v-data-table:not(.v-data-table--expand) tbody tr {
    &:nth-of-type(2n + 1) {
    background-color: transparent !important;

    &:hover {
      background: var(--v-table-hover-row-color-base, #eeeeee) !important;
    }
  }
  }
}
</style>
