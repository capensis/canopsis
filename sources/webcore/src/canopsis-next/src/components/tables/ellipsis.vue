<template lang="pug">
  div(
    v-if="expanded || (!expanded && text.length<=maxLetters)"
  )
    span(@click.stop="textClicked") {{ text }}
  div(v-else)
    span(@click.stop="textClicked") {{ text.substr(0, maxLetters) }}
    span ...
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
      expanded: false,
    };
  },
  methods: {
    textClicked() {
      this.$emit('textClicked');
      this.expanded = !this.expanded;
    },
  },
};
</script>
