<template>
  <v-menu
    :model-value="value"
    :close-on-content-click="false"
    :content-props="contentProps"
    location-strategy="static"
    scroll-strategy="none"
  >
    <point-form-dialog
      v-if="value && point"
      :point="point"
      :editing="editing"
      :exists-entities="existsEntities"
      v-bind="$attrs"
    />
  </v-menu>
</template>

<script>
import PointFormDialog from './point-form-dialog.vue';

export default {
  components: { PointFormDialog },
  inheritAttrs: false,
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    point: {
      type: Object,
      required: false,
    },
    positionX: {
      type: Number,
      required: true,
    },
    positionY: {
      type: Number,
      required: true,
    },
    editing: {
      type: Boolean,
      default: false,
    },
    existsEntities: {
      type: Array,
      required: false,
    },
  },
  computed: {
    contentProps() {
      return {
        style: {
          position: 'fixed',
          left: `${this.positionX}px`,
          top: `${this.positionY}px`,
        },
      };
    },
  },
};
</script>
