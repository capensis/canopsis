<template>
  <v-navigation-drawer
    v-model="opened"
    v-bind="navigationProps"
    :style="drawerStyle"
    class="sidebar-wrapper"
  >
    <div class="sidebar-wrapper__sticky-title">
      <v-list color="primary">
        <v-list-item>
          <v-list-item-title class="text-subtitle-1">
            <v-layout class="gap-3 align-center">
              <v-icon color="white">
                $vuetify.icons.ai
              </v-icon>
              <span class="text-h6 white--text">{{ $t('llm.chat.title') }}</span>
              <v-btn
                v-if="restartable"
                class="white--text"
                color="white"
                outlined
                @click="restart"
              >
                <v-icon class="mr-2" color="white" small>
                  refresh
                </v-icon>
                {{ $t('common.restart') }}
              </v-btn>
            </v-layout>
          </v-list-item-title>
          <v-btn icon @click.stop="minimize">
            <v-icon color="white">
              $vuetify.icons.hide_sidebar
            </v-icon>
          </v-btn>
        </v-list-item>
      </v-list>
      <v-divider />
    </div>
    <v-slide-x-transition>
      <v-layout
        v-if="minimized"
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
        <span class="text-h6 white--text">AI</span>
      </v-layout>
    </v-slide-x-transition>
    <slot />
  </v-navigation-drawer>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { CSS_COLORS_VARS } from '@/config';
import { LLM_AI_CHAT_WIDTH } from '@/constants';

import { getMaxZIndex } from '@/helpers/vuetify';

/**
 * Wrapper for each sidebar
 */
export default {
  props: {
    initialMinimized: {
      type: Boolean,
      default: false,
    },
    restartable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const opened = ref(false);
    const minimized = ref(props.initialMinimized);
    const zIndex = getMaxZIndex();

    const minimizedHeaderStyle = computed(() => ({
      backgroundColor: CSS_COLORS_VARS.primary,
    }));

    const drawerStyle = computed(() => ({ zIndex }));

    const navigationProps = computed(() => ({
      right: true,
      fixed: true,
      class: 'sidebar--overflow-y-hidden',
      width: !minimized.value ? LLM_AI_CHAT_WIDTH : 0,
      hideOverlay: true,
    }));

    const updateMinimized = (newMinimized = minimized.value) => {
      minimized.value = newMinimized;

      emit('update:minimized', newMinimized);
    };

    const minimize = () => updateMinimized(true);

    const maximize = () => updateMinimized(false);

    const restart = () => emit('restart');

    onMounted(() => window.requestAnimationFrame(() => {
      opened.value = true;
      updateMinimized();
    }));

    return {
      opened,
      minimized,
      minimizedHeaderStyle,
      drawerStyle,
      navigationProps,

      minimize,
      maximize,
      restart,
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

      &--close .sidebar--minimized__header {
        transform: translate(60px, -50%);
      }

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
      transform: translate(0, -50%);
      transition: transform 0.1s ease-in-out;
    }
  }
}
</style>
