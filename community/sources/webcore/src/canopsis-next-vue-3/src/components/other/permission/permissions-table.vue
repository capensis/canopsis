<template>
  <table :class="{ 'permissions-table--collapsed-header': indent !== 0 }" class="permissions-table">
    <thead>
      <tr>
        <th />
        <th v-for="role in roles" :key="role._id">
          {{ role.name }}
        </th>
      </tr>
    </thead>

    <tbody>
      <template v-for="item in preparedItems" :key="item._id">
        <tr class="permissions-table__item-row">
          <td :class="`pl-${indent * 3 + 2}`">
            <c-expand-btn
              v-if="item.children"
              :model-value="isExpanded(item)"
              class="mr-2"
              @update:model-value="toggleExpanded(item)"
            />

            <span
              :class="{ 'font-weight-medium cursor-pointer': item.children }"
              @click="item.children && toggleExpanded(item)"
            >
              <mark v-if="item.hasMask && search">{{ item.title }}</mark>
              <span v-else>{{ item.title }}</span>
            </span>
          </td>

          <permissions-table-cell
            v-for="role in roles"
            :key="role._id"
            :disabled="disabled"
            :permission="item"
            :role="role"
            @update:permission="forwardInput"
          />
        </tr>

        <tr v-if="item.children && isExpanded(item)" class="permissions-table__nested-row">
          <td :colspan="roles.length + 1" class="permissions-table__nested-cell">
            <div class="permissions-table__nested-wrapper">
              <permissions-table
                :disabled="disabled"
                :indent="indent + 1"
                :items="item.children"
                :roles="roles"
                :search="search"
                :search-depth="searchDepth"
                @update:permission="forwardInput"
              />
            </div>
          </td>
        </tr>
      </template>
    </tbody>
  </table>
</template>

<script>
import { computed, inject, ref, watch } from 'vue';

import PermissionsTableCell from './permissions-table-cell.vue';

export default {
  name: 'permissions-table',

  components: { PermissionsTableCell },

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

  emits: ['update:permission'],

  setup(props, { emit }) {
    const preparedItems = computed(() => props.items.map(item => ({
        ...item,

        hasMask:
          props.search &&
          (props.searchDepth === props.indent || (!props.searchDepth && !item.children)),
      })),
    );

    const allExpandedCounter = inject('$allExpandedCounter', ref(0));
    const expandedIds = ref([]);

    const isExpanded = item => expandedIds.value.includes(item._id);
    const toggleExpanded = (item) => {
      expandedIds.value = isExpanded(item)
        ? expandedIds.value.filter(id => id !== item._id)
        : [...expandedIds.value, item._id];
    };
    const checkExpanded = () => {
      expandedIds.value =
        allExpandedCounter.value > 0 ? preparedItems.value.map(({ _id: id }) => id) : [];
    };
    const forwardInput = (...args) => emit('update:permission', ...args);

    watch(allExpandedCounter, () => window.requestAnimationFrame(checkExpanded), {
      immediate: true,
    });

    return {
      preparedItems,
      isExpanded,
      toggleExpanded,
      forwardInput,
    };
  },
};
</script>

<style lang="scss" scoped>
.permissions-table {
  --checkboxCellWidth: 112px;
  --cellPadding: 8px;
  --rowHeight: 48px;
  --topBarHeight: 48px;

  border-collapse: separate;
  border-spacing: 0;
  table-layout: fixed;
  width: 100%;

  &--collapsed-header {
    > thead {
      display: none;
    }
  }

  th,
  td {
    height: var(--rowHeight);
    padding: var(--cellPadding);
    text-align: left;
    vertical-align: middle;

    &:not(:first-child) {
      width: var(--checkboxCellWidth);
    }
  }

  th {
    background: rgb(var(--v-theme-surface));
    font-weight: 500;
    position: sticky;
    top: calc(var(--topBarHeight) - 1px);
    transition: none;
    z-index: 2;
  }

  thead {
    position: sticky;
    top: calc(var(--topBarHeight) - 1px);
    z-index: 2;
  }

  &__item-row {
    background-color: transparent;
  }

  mark {
    background: rgba(var(--v-theme-warning), 0.24);
    color: inherit;
  }

  &__nested-cell {
    padding: 4px 0 4px 24px !important;
  }

  &__nested-wrapper {
    border-left: 1px solid rgba(var(--v-theme-on-surface), 0.12);
    padding-left: 8px;
  }

  &__nested-row tr:last-of-type td {
    border-bottom: none;
  }
}
</style>
