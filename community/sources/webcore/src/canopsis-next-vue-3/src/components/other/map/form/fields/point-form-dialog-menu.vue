<template>
  <v-menu :model-value="value" :target="[positionX, positionY]" :close-on-content-click="false">
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
};
</script>
