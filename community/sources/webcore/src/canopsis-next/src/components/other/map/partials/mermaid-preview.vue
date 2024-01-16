<template>
  <c-zoom-overlay
    skip-alt
    skip-shift
  >
    <panzoom
      class="mermaid-preview"
      ref="panzoom"
      :style="containerStyles"
      :help-text="$t('mermaid.panzoom.helpText')"
    >
      <mermaid-code-preview
        class="mermaid-preview__preview"
        :value="map.parameters.code"
      />
      <mermaid-points-preview
        class="mermaid-preview__points"
        :points="map.parameters.points"
        :popup-template="popupTemplate"
        :popup-actions="popupActions"
        :color-indicator="colorIndicator"
        :pbehavior-enabled="pbehaviorEnabled"
        @show:map="$emit('show:map', $event)"
        @show:alarms="$emit('show:alarms', $event)"
      />
    </panzoom>
  </c-zoom-overlay>
</template>

<script>
import Panzoom from '@/components/common/panzoom/panzoom.vue';

import MermaidPointsPreview from './mermaid-points-preview.vue';
import MermaidCodePreview from './mermaid-code-preview.vue';

export default {
  components: { Panzoom, MermaidPointsPreview, MermaidCodePreview },
  props: {
    map: {
      type: Object,
      required: true,
    },
    popupTemplate: {
      type: String,
      required: false,
    },
    popupActions: {
      type: Boolean,
      default: false,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    pbehaviorEnabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    containerStyles() {
      const points = this.map.parameters.points ?? [];
      const minHeight = Math.max.apply(null, points.map(({ y }) => y));

      return {
        minHeight: `${minHeight}px`,
      };
    },
  },
  watch: {
    map() {
      this.$refs.panzoom.reset();
    },
  },
};
</script>

<style lang="scss">
.mermaid-preview {
  background: white;
  border-radius: 5px;

  &__preview {
    width: 800px;
  }

  &__points {
    pointer-events: none;
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
  }
}
</style>
