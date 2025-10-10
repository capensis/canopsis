<template>
  <div>
    <span @click.stop="textClicked">{{ shortenedText }}</span>
    <v-menu
      v-if="!isShort"
      :close-on-content-click="false"
    >
      <template #activator="{ on }">
        <span
          class="ml-1"
          v-on="on"
        >...</span>
      </template>
      <v-card>
        <v-card-title class="pre-wrap">
          {{ text }}
        </v-card-title>
      </v-card>
    </v-menu>
  </div>
</template>

<script>
import { computed } from 'vue';

import { EXPAND_DEFAULT_MAX_LETTERS } from '@/config';

export default {
  props: {
    maxLetters: {
      type: Number,
      default: EXPAND_DEFAULT_MAX_LETTERS,
    },
    text: {
      type: [String, Number, Array],
      default: '',
    },
  },
  setup(props, { emit, listeners }) {
    const preparedText = computed(() => String(props.text));

    const isShort = computed(() => preparedText.value.length <= props.maxLetters);

    const shortenedText = computed(() => {
      if (isShort.value) {
        return preparedText.value;
      }

      return preparedText.value.substring(0, props.maxLetters);
    });

    const textClicked = () => {
      emit('textClicked');

      listeners.click?.();
    };

    return {
      preparedText,
      isShort,
      shortenedText,
      textClicked,
    };
  },
};
</script>
