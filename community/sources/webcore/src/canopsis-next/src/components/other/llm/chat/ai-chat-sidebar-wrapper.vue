<template>
  <v-navigation-drawer
    v-model="opened"
    v-bind="navigationProps"
    :style="drawerStyle"
    class="ai-chat-sidebar-wrapper"
  >
    <div class="ai-chat-sidebar-wrapper__sticky-title">
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
        class="ai-chat-sidebar-wrapper--minimized__header gap-5"
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
      class: { 'ai-chat-sidebar-wrapper--overflow-y-hidden': true, 'ai-chat-sidebar-wrapper--minimized': minimized.value },
      width: !minimized.value ? LLM_AI_CHAT_WIDTH : 0,
      hideOverlay: true,
    }));

    /**
     * Updates minimized state and notifies the parent via `update:minimized`.

     * @param {boolean} [newMinimized] - When omitted, keeps the current minimized value (re-sync only).
     */
    const updateMinimized = (newMinimized = minimized.value) => {
      minimized.value = newMinimized;

      emit('update:minimized', newMinimized);
    };

    /**
     * Collapses the sidebar to the minimized strip.
     */
    const minimize = () => updateMinimized(true);

    /**
     * Expands the sidebar from the minimized strip.
     */
    const maximize = () => updateMinimized(false);

    /**
     * Emits `restart` so the parent can reset the chat session.
     */
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
.ai-chat-sidebar-wrapper {
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

    &--close .ai-chat-sidebar-wrapper--minimized__header {
      transform: translate(60px, -50%);
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
