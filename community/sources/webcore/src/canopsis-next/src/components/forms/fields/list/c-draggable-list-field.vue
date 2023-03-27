<template lang="pug">
  draggable(
    :value="value",
    :options="options",
    :element="component",
    @change="updateOrdering"
  )
    slot
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { formMixin } from '@/mixins/form';

import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

export default {
  components: { Draggable },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    animation: {
      type: Number,
      default: VUETIFY_ANIMATION_DELAY,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    component: {
      type: String,
      required: false,
    },
    handle: {
      type: String,
      required: false,
    },
    ghostClass: {
      type: String,
      required: false,
    },
    group: {
      type: Object,
      required: false,
    },
  },
  computed: {
    options() {
      return {
        animation: this.animation,
        disabled: this.disabled,
        handle: this.handle,
        ghostClass: 'grey',
        group: this.group,
      };
    },
  },
  methods: {
    updateOrdering(event) {
      this.updateModel(dragDropChangePositionHandler(this.value, event));
    },
  },
};
</script>
