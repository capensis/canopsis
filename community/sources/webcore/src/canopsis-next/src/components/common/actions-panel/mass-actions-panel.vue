<template>
  <v-layout>
    <c-action-btn
      v-for="(action, index) in preparedActions.inline"
      :key="index"
      :tooltip="action.title"
      :disabled="action.disabled"
      :loading="action.loading"
      :icon="action.icon"
      :color="action.iconColor"
      :badge-value="action.badgeValue"
      :badge-tooltip="action.badgeTooltip"
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
          icon
          v-on="on"
        >
          <v-icon>more_vert</v-icon>
        </v-btn>
      </template>
      <v-list>
        <v-list-item
          v-for="(action, index) in preparedActions.dropDown"
          :key="index"
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
</template>

<script>
import { computed } from 'vue';

import { useComponentInstance } from '@/hooks/vue';

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
