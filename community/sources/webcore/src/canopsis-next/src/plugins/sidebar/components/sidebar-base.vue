<template>
  <sidebar-wrapper :sidebar="sidebar">
    <component
      v-if="sidebar.name"
      :is="sidebar.name"
      :sidebar="sidebar"
    />
  </sidebar-wrapper>
</template>

<script>
import { provide } from 'vue';

import ClickOutside from '@/services/click-outside';

import SidebarWrapper from './sidebar-wrapper.vue';

/**
 * Root for one sidebar instance (click-outside scope + dynamic panel component).
 */
export default {
  components: { SidebarWrapper },
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const clickOutsideService = new ClickOutside();

    provide('$clickOutside', clickOutsideService);
    provide('$sidebar', props.sidebar);
  },
};
</script>
