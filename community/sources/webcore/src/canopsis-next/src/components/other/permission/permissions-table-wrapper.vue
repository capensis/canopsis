<template>
  <v-layout class="gap-4" column>
    <v-layout class="px-2" justify-space-between align-center>
      <c-search-field v-model="search" />
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
      @input="emit('input', $event)"
    />
  </v-layout>
</template>

<script>
import { computed, provide, ref } from 'vue';

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
  setup(props) {
    const search = ref('');

    const filteredTreeviewPermissions = computed(() => ({ ...props.treeviewPermissions }));

    const allExpanded = ref(false);

    const expandAll = () => allExpanded.value = true;
    const collapseAll = () => allExpanded.value = false;

    provide('$allExpanded', allExpanded);

    return {
      search,

      filteredTreeviewPermissions,

      expandAll,
      collapseAll,
    };
  },
};
</script>
