<template lang="pug">
  v-menu(
    :value="true",
    :position-x="positionX",
    :position-y="positionY",
    :close-on-content-click="false",
    ignore-click-outside,
    absolute,
    z-index="230"
  )
    point-form-dialog(
      v-on="$listeners",
      :point="point",
      :editing="editing"
    )
</template>

<script>
import PointFormDialog from '@/components/other/map/form/partials/point-form-dialog.vue';

export default {
  inject: ['$validator'],
  components: { PointFormDialog },
  props: {
    point: {
      type: Object,
      required: true,
    },
    editing: {
      type: Boolean,
      default: false,
    },
    positionX: {
      type: Number,
      required: true,
    },
    positionY: {
      type: Number,
      required: true,
    },
  },
};
</script>
