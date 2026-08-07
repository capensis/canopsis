<template>
  <div class="sidebars-wrapper position-relative">
    <sidebar-base
      v-for="sidebar in sidebars"
      :key="sidebar.id"
      :sidebar="sidebar"
    />
  </div>
</template>

<script>
import { mapState } from 'pinia';

import { useSidebarStore } from '@/stores/modules/sidebar';

import SidebarBase from './sidebar-base.vue';

export default {
  name: 'TheSidebar',

  components: {
    SidebarBase,
  },

  computed: {
    ...mapState(useSidebarStore, ['sidebars']),
  },
};
</script>
