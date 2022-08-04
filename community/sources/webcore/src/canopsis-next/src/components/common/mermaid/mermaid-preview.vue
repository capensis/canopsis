<template lang="pug">
  div.mermaid-preview(:class="{ 'mermaid-preview--error': error }", v-html="svg")
</template>

<script>
import { renderMermaid } from '@/helpers/mermaid';

export default {
  props: {
    value: {
      type: String,
      required: false,
    },
    config: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'mermaid',
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      svg: '',
    };
  },
  watch: {
    value: {
      immediate: true,
      handler: 'renderMermaid',
    },
    config: 'renderMermaid',
  },
  methods: {
    renderMermaid() {
      try {
        this.svg = renderMermaid(this.value, this.config);
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
</script>

<style lang="scss">
.mermaid-preview {
  height: 100%;
  background-color: #F9F9F9;

  svg {
    width: 800px;
    max-width: 800px !important;
  }

  &--error svg {
    opacity: 0.5;
  }
}
</style>
