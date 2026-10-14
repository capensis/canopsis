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

import { FEATURES_SERVICE_KEYS } from '@/constants';

import { featuresService } from '@/services/features';

import { useSidebarStore } from '@/stores/modules/sidebar';

import { normalizeFeatureComponents } from '@/helpers/component';

import * as Sidebars from '@/components/sidebars';

import SidebarWrapper from './sidebar-wrapper.vue';

const featureSidebarComponents = normalizeFeatureComponents(
  featuresService.get(FEATURES_SERVICE_KEYS.components.sidebars.components, {}),
);

export default {
  components: {
    SidebarWrapper,
    ...Sidebars,
    ...featureSidebarComponents,
  },

  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },

  computed: {
    component() {
      return (
        featureSidebarComponents[this.sidebar.name] ??
        Object.values(Sidebars).find(component => component.name === this.sidebar.name)
      );
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
