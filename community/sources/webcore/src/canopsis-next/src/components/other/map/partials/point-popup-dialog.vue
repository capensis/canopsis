<template>
  <v-menu
    :value="true"
    :position-x="positionX"
    :position-y="positionY"
    :close-on-content-click="false"
    ignore-click-outside
    absolute
    top
  >
    <point-popup
      v-click-outside="clickOutsideDirective"
      :point="point"
      :template="popupTemplate"
      :color-indicator="colorIndicator"
      :actions="popupActions"
      v-on="$listeners"
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
  computed: {
    clickOutsideDirective() {
      return {
        handler: () => this.$emit('close'),
        include: () => Array.from(document.querySelectorAll('.point-icon')),
        closeConditional: () => true,
      };
    },
  },
};
</script>
