<template lang="pug">
  div
    template(v-if="text.length <= maxLetters")
      span(@click.stop="textClicked") {{ text }}
    template(v-else)
      span(@click.stop="textClicked") {{ text.substr(0, maxLetters) }}
      v-menu(v-model="fullTextMenu", :close-on-content-click="false")
        span.ml-1(slot="activator", small, depressed, @click.stop="fullTextMenu = true") ...
        v-card
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
      required: true,
    },
  },
  data() {
    return {
      fullTextMenu: false,
    };
  },
  methods: {
    textClicked() {
      this.$emit('textClicked');
    },
  },
};
</script>
