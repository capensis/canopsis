<template lang="pug">
  div
    span(@click.stop="textClicked") {{ shortenedText }}
    v-menu(
    v-if="!isShort",
    v-model="isFullTextMenuOpen",
    :close-on-content-click="false",
    :open-on-click="false"
    )
      span.ml-1(slot="activator", small, depressed, @click.stop="openFullTextMenu") ...
      v-card(dark)
        v-card-title {{ text }}
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
      type: String,
      default: '',
    },
  },
  data() {
    return {
      isFullTextMenuOpen: false,
    };
  },
  computed: {
    isShort() {
      return this.text.length <= this.maxLetters;
    },
    shortenedText() {
      if (this.isShort) {
        return this.text;
      }

      return this.text.substr(0, this.maxLetters);
    },
  },
  methods: {
    textClicked() {
      this.$emit('textClicked');
    },
    openFullTextMenu() {
      this.isFullTextMenuOpen = true;
    },
  },
};
</script>
