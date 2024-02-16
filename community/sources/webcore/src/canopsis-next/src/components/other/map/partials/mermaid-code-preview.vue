<template>
  <div
    class="mermaid-code-preview"
    :class="{ 'mermaid-preview--error': !parsed }"
    v-html="svg"
  />
</template>

<script>
import { merge } from 'lodash';

import { MERMAID_THEME_PROPERTIES_BY_NAME } from '@/constants';

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
    theme: {
      type: String,
      required: false,
    },
  },
  data() {
    return {
      svg: '',
      parsed: false,
    };
  },
  computed: {
    resultConfig() {
      const themeProperties = MERMAID_THEME_PROPERTIES_BY_NAME[this.theme] ?? {};

      return merge({
        theme: this.theme,
        ...themeProperties,

        er: {
          useMaxWidth: false,
        },
        pie: {
          useMaxWidth: false,
        },
        sequence: {
          useMaxWidth: false,
        },
        requirement: {
          useMaxWidth: false,
        },
      }, this.config);
    },
  },
  watch: {
    value: {
      immediate: true,
      handler: 'renderMermaid',
    },
    resultConfig: 'renderMermaid',
  },
  methods: {
    renderMermaid() {
      try {
        this.svg = renderMermaid(this.value, this.resultConfig);

        this.parsed = true;
      } catch (err) {
        this.parsed = false;
      }
    },
  },
};
</script>

<style lang="scss">
.mermaid-code-preview {
  height: 100%;

  svg {
    width: 800px;
    max-width: 800px !important;
  }

  &--error svg {
    opacity: 0.5;
  }
}
</style>
