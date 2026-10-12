<template>
  <v-layout class="gap-4" column>
    <v-layout class="px-2" justify-space-between align-center>
      <c-search @submit="submitSearch" />
      <v-layout class="gap-2" justify-end>
        <v-btn color="primary" outlined @click="expandAll">
          <v-icon class="mr-2" small>
            expand_all
          </v-icon>
          {{ $t('common.expandAll') }}
        </v-btn>
        <v-btn color="primary" outlined @click="collapseAll">
          <v-icon class="mr-2" small>
            collapse_all
          </v-icon>
          {{ $t('common.collapseAll') }}
        </v-btn>
      </v-layout>
    </v-layout>
    <permissions-table
      :items="filteredItems"
      :roles="roles"
      :disabled="disabled"
      :search="search"
      :search-depth="searchDepth"
      @input="updateTreeviewPermissions"
    />
  </v-layout>
</template>

<script>
import { sortBy } from 'lodash';
import {
  computed,
  provide,
  ref,
  watch,
  nextTick,
} from 'vue';

import { filterTreeviewPermissions } from '@/helpers/entities/permissions/list';

import { useI18n } from '@/hooks/i18n';

import PermissionsTable from './permissions-table.vue';

export default {
  components: { PermissionsTable },
  props: {
    treeviewPermissions: {
      type: Object,
      default: () => ({}),
    },
    roles: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    searchDepth: {
      type: Number,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const search = ref('');

    const getTranslatedAndSortedItems = (items = {}) => sortBy(Object.values(items).map((item) => {
      const newItem = { ...item };
      const { title, name, children, _id: id } = item;

      if (!title) {
        newItem.title = name ? t(`permission.title.${name}`) : id;
      }

      if (children) {
        newItem.children = getTranslatedAndSortedItems(children);
      }

      return newItem;
    }), ['position', 'title']);

    const itemsWithTranslations = computed(() => getTranslatedAndSortedItems(props.treeviewPermissions));

    const filteredItems = computed(() => (
      filterTreeviewPermissions(itemsWithTranslations.value, search.value, props.searchDepth)
    ));

    const allExpandedCounter = ref(0);

    /**
     * Expands all permissions in the treeview
     */
    const expandAll = () => (
      allExpandedCounter.value = allExpandedCounter.value > 0 ? allExpandedCounter.value + 1 : 1
    );

    /**
     * Collapses all permissions in the treeview
     */
    const collapseAll = () => (
      allExpandedCounter.value = allExpandedCounter.value > 0 ? 0 : allExpandedCounter.value - 1
    );

    /**
     * Sets the search value for filtering permissions
     *
     * @param {string} [value=''] - Search query string
     */
    const submitSearch = (value = '') => search.value = value;

    /**
     * Emits input event to update treeview permissions in parent component
     *
     * @param {...*} args - Updated permissions data to be passed to parent
     */
    const updateTreeviewPermissions = (...args) => emit('input', ...args);

    provide('$allExpandedCounter', allExpandedCounter);

    watch(filteredItems, () => search.value && nextTick(expandAll));

    return {
      search,

      filteredItems,

      expandAll,
      collapseAll,
      submitSearch,
      updateTreeviewPermissions,
    };
  },
};
</script>
