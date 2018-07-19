<template lang="pug">
  div(
    v-if="expanded || (!expanded && text.length<=maxLetters)",
    @click.stop="expanded = !expanded"
  )
    span {{ text }}
  div(v-else)
    span {{ text.substr(0, maxLetters) }}
    span(@click.stop="expanded = !expanded") ...
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
  computed: {},
};
</script>
