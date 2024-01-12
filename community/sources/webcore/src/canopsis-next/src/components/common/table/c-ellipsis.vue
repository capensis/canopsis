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
import { EXPAND_DEFAULT_MAX_LETTERS } from '@/config';

export default {
  props: {
    maxLetters: {
      type: Number,
      default: EXPAND_DEFAULT_MAX_LETTERS,
    },
    text: {
      type: [String, Number],
      default: '',
    },
  },
  computed: {
    preparedText() {
      return String(this.text);
    },

    isShort() {
      return this.preparedText.length <= this.maxLetters;
    },

    shortenedText() {
      if (this.isShort) {
        return this.preparedText;
      }

      return this.preparedText.substring(0, this.maxLetters);
    },
  },
  methods: {
    textClicked() {
      this.$emit('textClicked');
    },
  },
};
</script>
