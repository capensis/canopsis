<template>
  <teleport to="body">
    <aside
      :style="wrapperStyle"
      :class="wrapperClasses"
      class="ai-chat-sidebar-wrapper"
    >
      <div class="ai-chat-sidebar-wrapper__sticky-title bg-primary">
        <c-layout
          class="ai-chat-sidebar-wrapper__header px-4"
          gap="3"
          justify="space-between"
          align="center"
        >
          <c-layout class="min-width-0" gap="3" align="center">
            <v-icon class="ai-chat-sidebar-wrapper__ai-icon" color="white">
              $vuetify.icons.ai
            </v-icon>
            <span class="ai-chat-sidebar-wrapper__title text-white text-truncate">
              {{ $t('llm.chat.title') }}
            </span>
            <v-btn
              v-if="restartable"
              class="text-white"
              color="white"
              variant="outlined"
              @click="restart"
            >
              <v-icon class="mr-2" color="white" size="small">
                refresh
              </v-icon>
              {{ $t('common.restart') }}
            </v-btn>
          </c-layout>
          <v-btn variant="text" icon @click.stop="minimize">
            <v-icon color="white">
              $vuetify.icons.hide_sidebar
            </v-icon>
          </v-btn>
        </c-layout>
      </div>
      <v-slide-x-transition>
        <div
          v-if="minimized"
          :style="minimizedHeaderStyle"
          class="ai-chat-sidebar-wrapper--minimized__header"
        >
          <v-btn
            class="ai-chat-sidebar-wrapper--minimized__button"
            variant="text"
            icon
            @click.stop="maximize"
          >
            <v-icon color="white">
              $vuetify.icons.show_sidebar
            </v-icon>
          </v-btn>
          <span class="ai-chat-sidebar-wrapper--minimized__title text-white">
            AI
          </span>
        </div>
      </v-slide-x-transition>
      <div class="ai-chat-sidebar-wrapper__content">
        <slot />
      </div>
    </aside>
  </teleport>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { CSS_COLORS_VARS } from '@/config';
import { LLM_AI_CHAT_WIDTH } from '@/constants';

import { getMaxZIndex } from '@/helpers/vuetify';

export default {
  emits: ['update:minimized', 'restart'],
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
    const minimized = ref(props.initialMinimized);
    const zIndex = getMaxZIndex();

    const minimizedHeaderStyle = computed(() => ({
      backgroundColor: CSS_COLORS_VARS.primary,
    }));

    const wrapperStyle = computed(() => ({
      zIndex,
      width: `${LLM_AI_CHAT_WIDTH}px`,
    }));

    const wrapperClasses = computed(() => ({
      'ai-chat-sidebar-wrapper--minimized': minimized.value,
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

    onMounted(updateMinimized);

    return {
      minimized,
      minimizedHeaderStyle,
      wrapperStyle,
      wrapperClasses,

      minimize,
      maximize,
      restart,
    };
  },
};
</script>

<style lang="scss" scoped>
.ai-chat-sidebar-wrapper {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: visible;
  background: rgb(var(--v-theme-surface));
  box-shadow: 0 0 8px rgba(0, 0, 0, 0.24);
  transform: translateX(0);
  transition: transform 0.2s ease-in-out;

  &__title {
    font-size: 22px;
    font-weight: 400;
    line-height: 1.5;
  }

  &__ai-icon {
    display: inline-flex;
    width: 24px;
    min-width: 24px;
    align-items: center;
    justify-content: center;

    :deep(svg) {
      display: block;
      margin-right: auto;
      margin-left: auto;
    }
  }

  &__sticky-title {
    position: sticky;
    top: 0;
    flex: 0 0 auto;
    z-index: 2;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
  }

  &__header {
    min-height: 64px;
  }

  &__content {
    flex: 1 1 auto;
    min-height: 0;
    overflow-x: hidden;
    overflow-y: auto;
  }

  &--minimized {
    transform: translateX(100%);

    &__header {
      position: absolute;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 12px;
      width: 58px;
      height: 220px;
      border-top-left-radius: 10px;
      border-bottom-left-radius: 10px;
      left: -58px;
      top: 50%;
      transform: translate(0, -50%);
      transition: transform 0.1s ease-in-out;
    }

    &__button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      margin: 0;
    }

    &__title {
      width: 100%;
      font-size: 22px;
      font-weight: 400;
      line-height: 1.5;
      text-align: center;
    }
  }
}
</style>
