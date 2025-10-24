<template>
  <v-layout class="gap-4" column>
    <v-layout class="px-2" justify-space-between align-center>
      <c-search @submit="submitSearch" />
      <v-layout class="gap-2" justify-end>
        <v-btn color="primary" outlined @click="expandAll">
          <v-icon class="mr-2" small>
            $vuetify.icons.expand_all
          </v-icon>
          {{ $t('common.expandAll') }}
        </v-btn>
        <v-btn color="primary" outlined @click="collapseAll">
          <v-icon class="mr-2" small>
            $vuetify.icons.collapse_all
          </v-icon>
          {{ $t('common.collapseAll') }}
        </v-btn>
      </v-layout>
    </v-layout>
    <permissions-table
      :treeview-permissions="filteredTreeviewPermissions"
      :roles="roles"
      :disabled="disabled"
      :search="search"
      @input="updateTreeviewPermissions"
    />
  </v-layout>
</template>

<script>
import {
  computed,
  provide,
  ref,
  watch,
  nextTick,
} from 'vue';

import { filterTreeviewPermissions } from '@/helpers/entities/permissions/list';

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
  },
  setup(props, { emit }) {
    const search = ref('');

    const filteredTreeviewPermissions = computed(() => (
      filterTreeviewPermissions(props.treeviewPermissions, search.value)
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

    const updateTreeviewPermissions = (...args) => emit('input', ...args);

    provide('$allExpandedCounter', allExpandedCounter);

    watch(filteredTreeviewPermissions, () => search.value && nextTick(expandAll));

    return {
      search,

      filteredTreeviewPermissions,

      expandAll,
      collapseAll,
      submitSearch,
      updateTreeviewPermissions,
    };
  },
};
</script>
