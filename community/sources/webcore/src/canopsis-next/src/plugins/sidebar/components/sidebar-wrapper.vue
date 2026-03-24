<template>
  <div class="sidebar-wrapper">
    <v-navigation-drawer
      v-model="isOpen"
      v-bind="navigationProps"
      :style="drawerStyle"
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
                <v-btn class="white--text" color="white" outlined>
                  <v-icon class="mr-2" color="white" small>
                    refresh
                  </v-icon>
                  {{ $t('common.restart') }}
                </v-btn>
              </v-layout>
            </v-list-item-title>
            <v-btn
              v-if="sidebar.config?.minimizable"
              icon
              @click.stop="minimize"
            >
              <v-icon color="white">
                $vuetify.icons.hide_sidebar
              </v-icon>
            </v-btn>
            <v-btn
              v-else
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
      <v-slide-x-transition v-if="sidebar.config?.minimizable">
        <v-layout
          v-if="sidebar.minimized"
          :style="minimizedHeaderStyle"
          class="sidebar--minimized__header gap-5"
          column
          justify-center
          align-center
        >
          <v-btn
            icon
            @click.stop="maximize"
          >
            <v-icon color="white">
              $vuetify.icons.show_sidebar
            </v-icon>
          </v-btn>
          <span class="text-h6 white--text">{{ sidebar.config.titleMinimized }}</span>
        </v-layout>
      </v-slide-x-transition>
      <slot />
    </v-navigation-drawer>
  </div>
</template>

<script>
import { computed, inject, onMounted, ref } from 'vue';

import { DEFAULT_SIDEBAR_DRAWER_WIDTH, CSS_COLORS_VARS } from '@/config';
import { SIDE_BARS_WITH_OVERFLOW_Y_HIDDEN } from '@/constants';

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

    const minimizedHeaderStyle = computed(() => ({
      backgroundColor: CSS_COLORS_VARS[props.sidebar.config.color] || CSS_COLORS_VARS.secondary,
    }));

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

    const minimize = () => sidebar.minimize({ id: props.sidebar.id });

    const maximize = () => sidebar.maximize({ id: props.sidebar.id });

    const drawerStyle = computed(() => ({
      zIndex,
    }));

    const navigationProps = computed(() => {
      const { minimizable, width = DEFAULT_SIDEBAR_DRAWER_WIDTH } = props.sidebar.config;

      return {
        right: true,
        fixed: true,
        class: SIDE_BARS_WITH_OVERFLOW_Y_HIDDEN.includes(props.sidebar.name) ? 'sidebar--overflow-y-hidden' : '',
        width: isOpen.value && !props.sidebar.minimized ? width : 0,
        temporary: !minimizable,
        hideOverlay: minimizable,
        ignoreClickOutside: hasMaximizedModal.value,
        customCloseConditional: closeCondition,
      };
    });

    onMounted(() => window.requestAnimationFrame(() => ready.value = true));

    return {
      hasMaximizedModal,
      title,
      minimizedHeaderStyle,
      isOpen,
      closeHandler,
      minimize,
      maximize,
      drawerStyle,
      navigationProps,
    };
  },
};
</script>

<style lang="scss" scoped>
.sidebar {
  &-wrapper {
    position: relative;

    &__sticky-title {
      position: sticky;
      top: 0;
      z-index: 2;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
    }

    ::v-deep .v-navigation-drawer {
      overflow: visible;

      &__content {
        overflow-x: visible;
      }

      &.sidebar--overflow-y-hidden .v-navigation-drawer__content  {
        overflow-y: hidden;
      }
    }
  }

  &--minimized {
    &__header {
      position: absolute;
      width: 58px;
      height: 200px;
      border-top-left-radius: 10px;
      border-bottom-left-radius: 10px;
      right: 100%;
      top: 50%;
      transform: translateY(-50%);
    }
  }
}
</style>
