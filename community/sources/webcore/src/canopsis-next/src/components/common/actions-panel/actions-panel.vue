<template>
  <div
    :class="{ 'actions-panel--small': small }"
    class="actions-panel"
  >
    <v-layout align-center>
      <actions-panel-btn
        v-for="action in preparedActions.inline"
        :key="action.title"
        :action="action"
      />
      <actions-panel-menu
        key="dropdown-menu"
        :actions="preparedActions.dropDown"
      />
    </v-layout>
  </div>
</template>

<script>
import { computed } from 'vue';

import { DEFAULT_ALARM_ACTIONS_INLINE_COUNT } from '@/constants';

import { useComponentInstance } from '@/hooks/vue';

import ActionsPanelBtn from './actions-panel-btn.vue';
import ActionsPanelMenu from './actions-panel-menu.vue';

export default {
  components: {
    ActionsPanelBtn,
    ActionsPanelMenu,
  },
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    inlineCount: {
      type: Number,
      default: DEFAULT_ALARM_ACTIONS_INLINE_COUNT,
    },
    small: {
      type: Boolean,
      default: false,
    },
    ignoreMediaQuery: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const instance = useComponentInstance();

    const preparedActions = computed(() => {
      if (!props.ignoreMediaQuery && instance.$mq !== 'xl') {
        return {
          inline: [],
          dropDown: props.actions,
        };
      }

      if (props.inlineCount < props.actions.length) {
        const inlineCountWithoutMenu = props.inlineCount - 1;

        return {
          inline: props.actions.slice(0, inlineCountWithoutMenu),
          dropDown: props.actions.slice(inlineCountWithoutMenu),
        };
      }

      return {
        inline: props.actions,
        dropDown: [],
      };
    });

    return {
      preparedActions,
    };
  },
};
</script>

<style lang="scss">
.actions-panel {
  &__menu-item-loader {
    margin-right: 32px;
  }

  &--small {
    .v-btn--icon {
      width: 24px;
      height: 24px;

      .v-icon {
        font-size: 20px;
      }
    }
  }
}
</style>
