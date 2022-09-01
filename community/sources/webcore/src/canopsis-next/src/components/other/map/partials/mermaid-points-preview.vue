<template lang="pug">
  div
    mermaid-point-marker(
      ref="points",
      v-for="point in points",
      :key="point._id",
      :x="point.x",
      :y="point.y",
      :entity="point.entity",
      :size="iconSize",
      :color-indicator="colorIndicator",
      :pbehavior-enabled="pbehaviorEnabled",
      @click="openPopup(point, $event)"
    )
    point-popup-dialog(
      v-if="activePoint",
      :point="activePoint",
      :position-x="positionX",
      :position-y="positionY",
      :popup-template="popupTemplate",
      :color-indicator="colorIndicator",
      :popup-actions="popupActions",
      @show:alarms="showAlarms",
      @show:map="showLinkedMap",
      @close="closePopup"
    )
</template>

<script>
import { mapInformationPopup } from '@/mixins/map/map-information-popup';

import MermaidPointMarker from './mermaid-point-marker.vue';
import PointPopupDialog from './point-popup-dialog.vue';

export default {
  components: { MermaidPointMarker, PointPopupDialog },
  mixins: [mapInformationPopup],
  props: {
    points: {
      type: Array,
      required: true,
    },
  },
};
</script>
