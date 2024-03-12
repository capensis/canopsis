<template>
  <mermaid-preview
    :map="map"
    class="mermaid-expand-panel"
  />
</template>

<script>
import MermaidPreview from './mermaid-preview.vue';

export default {
  components: { MermaidPreview },
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
};
</script>
