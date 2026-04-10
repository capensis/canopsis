<template>
  <v-layout>
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
</template>

<script>
import { computed } from 'vue';

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
      required: true,
    },
    inlineCount: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const instance = useComponentInstance();

    const preparedActions = computed(() => {
      if (['m', 't'].includes(instance.$mq)) {
        return {
          inline: [],
          dropDown: props.actions,
        };
      }

      if (!props.inlineCount || props.inlineCount >= props.actions.length) {
        return {
          inline: props.actions,
          dropDown: [],
        };
      }

      const inlineCountWithoutMenu = props.inlineCount - 1;

      return {
        inline: props.actions.slice(0, inlineCountWithoutMenu),
        dropDown: props.actions.slice(inlineCountWithoutMenu),
      };
    });

    return {
      preparedActions,
    };
  },
};
</script>
