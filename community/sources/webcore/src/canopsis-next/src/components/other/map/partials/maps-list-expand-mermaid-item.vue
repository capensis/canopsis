<template lang="pug">
  panzoom.mermaid-expand-panel(:style="containerStyles", :help-text="$t('mermaid.panzoom.helpText')")
    mermaid-preview.mermaid-expand-panel__preview(:value="map.parameters.code")
    mermaid-points-preview.mermaid-expand-panel__points(:points="map.parameters.points")
</template>

<script>
import Panzoom from '@/components/common/panzoom/panzoom.vue';
import MermaidPreview from '@/components/other/map/partials/mermaid-preview.vue';
import MermaidPointsPreview from '@/components/other/map/partials/mermaid-points-preview.vue';

export default {
  components: { Panzoom, MermaidPointsPreview, MermaidPreview },
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
  computed: {
    containerStyles() {
      const minHeight = Math.max.apply(null, this.map.parameters.points.map(({ y }) => y));

      return {
        minHeight: `${minHeight}px`,
      };
    },
  },
};
</script>

<style lang="scss">
.mermaid-expand-panel {
  background: #F9F9F9;

  &__preview {
    width: 800px;
  }

  &__points {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
  }
}
</style>
