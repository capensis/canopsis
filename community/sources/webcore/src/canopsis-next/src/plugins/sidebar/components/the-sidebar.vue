<template>
  <div>
    <sidebar-base
      v-for="sidebarItem in sidebars"
      :key="sidebarItem.id"
      :sidebar="sidebarItem"
    />
  </div>
</template>

<script>
import { computed } from 'vue';

import { useSidebar } from '@/hooks/sidebar';
import { useStore } from '@/hooks/store';

/**
 * Wrapper for all sidebars
 */
export default {
  setup() {
    const store = useStore();
    const sidebar = useSidebar();

    const sidebars = computed(() => store.getters[`${sidebar.moduleName}/sidebars`]);

    return {
      sidebars,
    };
  },
};
</script>
