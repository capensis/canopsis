<template>
  <v-menu
    :model-value="true"
    :target="[positionX, positionY]"
    :close-on-content-click="false"
    location="top"
    @update:model-value="$emit('close')"
  >
    <point-popup
      :point="point"
      :template="popupTemplate"
      :color-indicator="colorIndicator"
      :actions="popupActions"
      @show:alarms="$emit('show:alarms')"
      @show:map="$emit('show:map')"
      @close="$emit('close')"
    />
  </v-menu>
</template>

<script>
import PointPopup from './point-popup.vue';

export default {
  components: { PointPopup },

  props: {
    point: {
      type: Object,
      required: true,
    },
    positionX: {
      type: Number,
      required: true,
    },
    positionY: {
      type: Number,
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
  },

  emits: ['close', 'show:alarms', 'show:map'],
};
</script>
