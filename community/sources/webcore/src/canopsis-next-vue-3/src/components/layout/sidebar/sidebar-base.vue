<template>
  <sidebar-wrapper :sidebar="sidebar" @close="hideSidebar">
    <component
      v-if="component"
      :is="component"
      ref="sidebarComponent"
      :sidebar="sidebar"
    />
  </sidebar-wrapper>
</template>

<script>
import { mapActions } from 'pinia';

import { useSidebarStore } from '@/stores/modules/sidebar';

import * as Sidebars from '@/components/sidebars';

import SidebarWrapper from './sidebar-wrapper.vue';

export default {
  components: {
    SidebarWrapper,
    ...Sidebars,
  },

  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },

  computed: {
    component() {
      return Object.values(Sidebars).find(component => component.name === this.sidebar.name);
    },
  },

  methods: {
    ...mapActions(useSidebarStore, {
      hide: 'hide',
    }),

    hideSidebar() {
      if (typeof this.$refs.sidebarComponent?.close === 'function') {
        this.$refs.sidebarComponent.close();

        return;
      }

      this.hide({ id: this.sidebar.id });
    },
  },
};
</script>
