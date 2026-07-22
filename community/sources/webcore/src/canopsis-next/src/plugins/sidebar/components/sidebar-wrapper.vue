<template>
  <v-navigation-drawer
    v-model="isOpen"
    v-bind="navigationProps"
    :style="drawerStyle"
    class="sidebar-wrapper"
  >
    <div
      v-if="title"
      class="sidebar-wrapper__sticky-title"
    >
      <v-list :color="sidebar.config.color || 'secondary'">
        <v-list-item>
          <v-list-item-title class="text-subtitle-1">
            <v-layout class="gap-3 align-center">
              <v-icon v-if="sidebar.config.titleIcon" color="white">
                {{ sidebar.config.titleIcon }}
              </v-icon>
              <span class="text-h6 white--text">{{ title }}</span>
            </v-layout>
          </v-list-item-title>
          <v-btn
            icon
            @click.stop="closeHandler"
          >
            <v-icon color="white">
              close
            </v-icon>
          </v-btn>
        </v-list-item>
      </v-list>
      <v-divider />
    </div>
    <slot />
  </v-navigation-drawer>
</template>

<script>
import { computed, inject, ref, onMounted } from 'vue';

import { DEFAULT_SIDEBAR_DRAWER_WIDTH } from '@/config';

import { getMaxZIndex } from '@/helpers/vuetify';

import { useI18n } from '@/hooks/i18n';
import { useStore } from '@/hooks/store';
import { useModals } from '@/hooks/modals';
import { useSidebar } from '@/hooks/sidebar';

/**
 * Wrapper for each sidebar
 */
export default {
  props: {
    sidebar: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const store = useStore();
    const modals = useModals();
    const sidebar = useSidebar();
    const { t } = useI18n();
    const clickOutside = inject('$clickOutside');

    const ready = ref(false);
    const zIndex = getMaxZIndex();

    const hasMaximizedModal = computed(() => store.getters[`${modals.moduleName}/hasMaximizedModal`]);

    const title = computed(() => props.sidebar.config.title || (props.sidebar.name ? t(`settings.titles.${props.sidebar.name}`) : ''));

    const isOpen = computed({
      get: () => props.sidebar.name && !props.sidebar.hidden && ready.value,
      set: (value) => {
        if (!value) {
          sidebar.hide({ id: props.sidebar.id });
        }
      },
    });

    /**
     * Delegates to the injected `$clickOutside` helper so `VNavigationDrawer` can decide
     * whether an outside click (or equivalent) should close the drawer.
     */
    const closeCondition = (...args) => clickOutside.call(...args);

    /**
     * Runs when the user uses in-drawer close controls; hides the sidebar only if
     * `closeCondition` allows closing.
     */
    const closeHandler = () => closeCondition() && sidebar.hide({ id: props.sidebar.id });

    const drawerStyle = computed(() => ({
      zIndex,
    }));

    const navigationProps = computed(() => {
      const { width = DEFAULT_SIDEBAR_DRAWER_WIDTH } = props.sidebar.config;

      return {
        right: true,
        fixed: true,
        width: isOpen.value ? width : 0,
        temporary: true,
        hideOverlay: false,
        ignoreClickOutside: hasMaximizedModal.value,
        customCloseConditional: closeCondition,
      };
    });

    onMounted(() => window.requestAnimationFrame(() => ready.value = true));

    return {
      hasMaximizedModal,
      title,
      isOpen,
      closeHandler,
      drawerStyle,
      navigationProps,
    };
  },
};
</script>

<style lang="scss" scoped>
.sidebar {
  &-wrapper {
    &__sticky-title {
      position: sticky;
      top: 0;
      z-index: 2;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
    }

    &.v-navigation-drawer {
      overflow: visible;

      &__content {
        overflow-x: visible;
      }
    }
  }
}
</style>
