<template>
  <c-layout
    gap="4"
    column
  >
    <c-layout
      class="permissions-table-wrapper__toolbar px-2"
      justify="space-between"
      wrap
    >
      <v-row>
        <c-search
          class="permissions-table-wrapper__search-field"
          @submit="submitSearch"
        />
      </v-row>

      <v-row justify="end">
        <v-btn
          color="primary"
          variant="outlined"
          @click="expandAll"
        >
          <c-icon
            class="mr-2"
            small
          >
            expand_all
          </c-icon>
          {{ $t('common.expandAll') }}
        </v-btn>

        <v-btn
          color="primary"
          variant="outlined"
          @click="collapseAll"
        >
          <c-icon
            class="mr-2"
            small
          >
            collapse_all
          </c-icon>
          {{ $t('common.collapseAll') }}
        </v-btn>
      </v-row>
    </c-layout>

    <permissions-table
      :disabled="disabled"
      :items="filteredItems"
      :roles="roles"
      :search="search"
      :search-depth="searchDepth"
      @update:permission="updateTreeviewPermissions"
    />
  </c-layout>
</template>

<script>
import { sortBy } from 'lodash';
import {
  computed,
  nextTick,
  provide,
  ref,
  watch,
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

  emits: ['update:permission'],

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

    const expandAll = () => (
      allExpandedCounter.value = allExpandedCounter.value > 0 ? allExpandedCounter.value + 1 : 1
    );
    const collapseAll = () => (
      allExpandedCounter.value = allExpandedCounter.value > 0 ? 0 : allExpandedCounter.value - 1
    );
    const submitSearch = (value = '') => search.value = value;
    const updateTreeviewPermissions = (...args) => emit('update:permission', ...args);

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

<style lang="scss" scoped>
.permissions-table-wrapper {
  &__toolbar {
    width: 100%;
  }

  &__search {
    min-width: 260px;
  }

  &__search-field {
    width: 100%;

    :deep(.v-input) {
      width: 100%;
    }
  }

  &__actions {
    flex: 0 0 auto;
  }
}
</style>
