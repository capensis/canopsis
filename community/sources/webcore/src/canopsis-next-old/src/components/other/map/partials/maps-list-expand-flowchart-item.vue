<template>
  <flowchart-preview :map="map" />
</template>

<script>
import FlowchartPreview from './flowchart-preview.vue';

export default {
  components: { FlowchartPreview },
  props: {
    map: {
      type: Object,
      required: true,
    },
  },
};
</script>
