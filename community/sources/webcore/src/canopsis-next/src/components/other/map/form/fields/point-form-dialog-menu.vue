<template>
  <v-menu
    :value="value"
    :position-x="positionX"
    :position-y="positionY"
    :close-on-content-click="false"
    ignore-click-outside
    absolute
  >
    <point-form-dialog
      v-on="$listeners"
      v-if="value && point"
      :point="point"
      :editing="editing"
      :exists-entities="existsEntities"
    />
  </v-menu>
</template>

<script>
import PointFormDialog from './point-form-dialog.vue';

export default {
  components: { PointFormDialog },
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
