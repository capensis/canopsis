<template>
  <draggable
    :value="value"
    :group="group"
    :tag="component"
    :animation="animation"
    :disabled="disabled"
    :handle="handle"
    :ghost-class="ghostClass"
    :drag-class="dragClass"
    :chosen-class="dragClass"
    :component-data="componentData"
    :move="itemMove"
    :draggable="draggable"
    @change="updateOrdering"
    @start="$emit('start', $event)"
    @end="$emit('end', $event)"
  >
    <slot />
    <template #footer="">
      <slot name="footer" />
    </template>
  </draggable>
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

import { formMixin } from '@/mixins/form';

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
    componentData: {
      type: Object,
      required: false,
    },
    handle: {
      type: String,
      required: false,
    },
    ghostClass: {
      type: String,
      default: 'grey',
    },
    dragClass: {
      type: String,
      required: false,
    },
    chosenClass: {
      type: String,
      required: false,
    },
    group: {
      type: [Object, String],
      required: false,
    },
    itemMove: {
      type: Function,
      required: false,
    },
    draggable: {
      type: String,
      required: false,
    },
  },
  methods: {
    updateOrdering(event) {
      this.updateModel(dragDropChangePositionHandler(this.value, event));
    },
  },
};
</script>
