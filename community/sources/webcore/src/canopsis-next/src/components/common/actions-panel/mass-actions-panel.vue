<template>
  <div :style="{ width: `${maxWidth}px` }" class="mass-actions-panel__wrapper">
    <v-layout ref="actionsElement" class="mass-actions-panel__actions">
      <c-action-btn
        v-for="action in preparedActions.inline"
        :key="action.type"
        :tooltip="action.title"
        :disabled="action.disabled"
        :loading="action.loading"
        :icon="action.icon"
        :color="action.iconColor"
        :badge-value="action.badgeValue"
        :badge-tooltip="action.badgeTooltip"
        :small="small"
        @click="action.method"
      />
      <v-menu
        v-if="preparedActions.dropDown.length"
        bottom
        left
        @click.native.stop=""
      >
        <template #activator="{ on }">
          <v-btn
            :small="small"
            icon
            v-on="on"
          >
            <v-icon>more_vert</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item
            v-for="action in preparedActions.dropDown"
            :key="action.type"
            :disabled="action.disabled || action.loading"
            @click.stop="action.method"
          >
            <v-list-item-title>
              <span class="mr-4">
                <v-icon
                  :color="action.iconColor"
                  :disabled="action.disabled"
                  class="ma-0 pa-0"
                  left
                  small
                >
                  {{ action.icon }}
                </v-icon>
              </span>
              <span
                :class="action.cssClass"
                class="text-body-1"
              >
                {{ action.title }}
              </span>
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-layout>
  </div>
</template>

<script>
import {
  computed,
  ref,
  watch,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { MASS_ACTIONS_BUTTON_WIDTH } from '@/constants';

export default {
  props: {
    actions: {
      type: Array,
      required: true,
    },
    inlineCount: {
      type: Number,
      default: 0,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const actionsElement = ref(null);
    const localInlineCount = ref(props.inlineCount);

    const maxInlineCount = computed(() => Math.min(props.inlineCount || props.actions.length, props.actions.length));
    const minInlineCount = computed(() => (props.actions.length ? 1 : 0));
    const maxWidth = computed(() => MASS_ACTIONS_BUTTON_WIDTH * maxInlineCount.value);
    const preparedActions = computed(() => {
      if (!localInlineCount.value && !minInlineCount.value) {
        return {
          inline: [],
          dropDown: [],
        };
      }

      if (props.actions.length === localInlineCount.value) {
        return {
          inline: props.actions,
          dropDown: [],
        };
      }

      const inlineCountWithoutMenu = localInlineCount.value - 1;

      return {
        inline: props.actions.slice(0, inlineCountWithoutMenu),
        dropDown: props.actions.slice(inlineCountWithoutMenu),
      };
    });

    const observer = new ResizeObserver(([entry]) => {
      const availableInlineItems = Math.max(
        minInlineCount.value,
        Math.min(
          maxInlineCount.value,
          Math.floor(entry.contentRect.width / MASS_ACTIONS_BUTTON_WIDTH),
        ),
      );

      if (localInlineCount.value === availableInlineItems) {
        return;
      }

      localInlineCount.value = availableInlineItems;
    });

    watch(maxInlineCount, (newMaxInlineCount) => {
      if (localInlineCount.value > newMaxInlineCount) {
        localInlineCount.value = newMaxInlineCount;
      }
    });

    onMounted(() => actionsElement.value && observer.observe(actionsElement.value));

    onBeforeUnmount(() => observer.disconnect());

    return {
      actionsElement,
      localInlineCount,
      maxWidth,
      preparedActions,
    };
  },
};
</script>

<style lang="scss">
.mass-actions-panel {
  &__wrapper {
    position: relative;
    height: 28px;
  }

  &__actions {
    position: absolute;
    top: 0;
    left: 0;
    height: 28px;
    width: 100%;
  }
}
</style>
