<template lang="pug">
  c-zoom-overlay
    network-graph.fill-height(
      ref="networkGraph",
      :options="options",
      :node-html-label-options="nodeHtmlLabelsOptions"
    )
</template>

<script>
import { COLORS } from '@/config';

export default {
  props: {
    entities: {
      type: Array,
      default: () => [],
    },
    impact: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    styleOption() {
      return [
        {
          selector: 'node',
          style: {
            'font-size': '10px',
            'background-color': COLORS.secondary,
            'border-color': COLORS.secondary,
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'curve-style': 'bezier',
            'line-color': 'silver',
          },
        },
      ];
    },
  },
};
</script>
