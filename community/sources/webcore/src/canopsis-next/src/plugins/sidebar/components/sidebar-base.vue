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
import ClickOutside from '@/services/click-outside';

import SidebarWrapper from './sidebar-wrapper.vue';

export default {
  provide() {
    return {
      $clickOutside: this.$clickOutside,
    };
  },
  components: { SidebarWrapper },
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },
  beforeCreate() {
    this.$clickOutside = new ClickOutside();
  },
};
</script>
