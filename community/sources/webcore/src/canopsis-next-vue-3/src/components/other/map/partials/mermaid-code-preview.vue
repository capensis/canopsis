<template>
  <div v-html="svg" :class="{ 'mermaid-preview--error': !parsed }" class="mermaid-code-preview" />
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

      return merge(
        {
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
        },
        this.config,
      );
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
    normalizeSvgHeight(svg) {
      const wrapper = document.createElement('div');

      wrapper.innerHTML = svg;

      const svgElement = wrapper.querySelector('svg');
      const [, , , height] = svgElement?.getAttribute('viewBox')?.split(/\s+/).map(Number) ?? [];

      if (height) {
        svgElement.setAttribute('height', height);
        svgElement.style.height = `${height}px`;
      }

      return wrapper.innerHTML;
    },

    renderMermaid() {
      try {
        this.svg = this.normalizeSvgHeight(renderMermaid(this.value, this.resultConfig));

        this.parsed = true;
      } catch (err) {
        console.error(err);

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
