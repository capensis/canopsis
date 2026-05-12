<template>
  <div :style="{ width: `${maxWidth}px` }" class="mass-actions-panel__wrapper">
    <v-layout ref="actionsElement" class="mass-actions-panel__actions">
      <mass-actions-panel-inline-actions
        v-for="action in preparedActions.inline"
        :key="action.type"
        :action="action"
        :small="small"
      />
      <mass-actions-panel-dropdown-menu-actions
        v-if="preparedActions.dropDown.length"
        key="dropdown-menu-actions"
        :actions="preparedActions.dropDown"
        :small="small"
      />
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

import MassActionsPanelDropdownMenuActions from './mass-actions-panel-dropdown-menu-actions.vue';
import MassActionsPanelInlineActions from './mass-actions-panel-inline-actions.vue';

export default {
  components: {
    MassActionsPanelDropdownMenuActions,
    MassActionsPanelInlineActions,
  },
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
