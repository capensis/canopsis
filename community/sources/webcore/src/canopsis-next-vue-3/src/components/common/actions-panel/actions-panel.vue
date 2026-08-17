<template>
  <div :class="{ 'actions-panel--small': small }" class="actions-panel">
    <c-layout align="center">
      <actions-panel-btn
        v-for="action in preparedActions.inline"
        :key="action.title"
        :action="action"
        :small="small"
      />
      <actions-panel-menu ref="menu" :actions="preparedActions.dropDown" :small="small" />
    </c-layout>
  </div>
</template>

<script>
import { computed, ref } from 'vue';
import { useDisplay } from 'vuetify';

import { DEFAULT_ALARM_ACTIONS_INLINE_COUNT } from '@/constants';

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
    const { xl } = useDisplay();
    const menu = ref(null);

    const preparedActions = computed(() => {
      if (!props.ignoreMediaQuery && !xl.value) {
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

    const closeMenu = () => menu.value?.closeMenu?.();

    return {
      closeMenu,
      menu,
      preparedActions,
    };
  },
};
</script>

<style lang="scss">
.actions-panel {
  &__menu-item {
    display: inline-flex;
    align-items: center;
    vertical-align: middle;
  }

  &__menu-item-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    flex-shrink: 0;
  }

  &--small {
    .v-btn--icon.v-btn--density-default,
    .v-btn--icon.v-btn--density-comfortable,
    .v-btn--icon.v-btn--density-compact {
      width: 24px;
      height: 24px;

      .v-icon {
        font-size: 20px;
      }
    }
  }
}
</style>
