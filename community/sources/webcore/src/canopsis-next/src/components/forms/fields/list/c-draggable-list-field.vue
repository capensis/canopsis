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
    :class="classes"
    @change="updateOrdering"
    @start="startDragging"
    @end="endDragging"
  >
    <slot />
    <template #footer="">
      <slot name="footer" />
    </template>
  </draggable>
</template>

<script>
import { computed, ref } from 'vue';
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

import { useModelField } from '@/hooks/form';

export default {
  components: { Draggable },
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
      default: '>*',
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    /**
     * Tracks whether a drag operation is currently active.
     */
    const isDragging = ref(false);

    /**
     * Classes applied to the root draggable element depending on state.
     */
    const classes = computed(() => ({
      'c-draggable-list-field--dragging': isDragging.value,
    }));

    /**
     * Handles drag start: set dragging state and re-emit the event.
     *
     * @param {Object} event - Drag start event payload from vuedraggable/Sortable.
     * @returns {void}
     */
    const startDragging = (event) => {
      isDragging.value = true;
      emit('start', event);
    };

    /**
     * Handles drag end: reset dragging state and re-emit the event.
     *
     * @param {Object} event - Drag end event payload from vuedraggable/Sortable.
     * @returns {void}
     */
    const endDragging = (event) => {
      isDragging.value = false;
      emit('end', event);
    };

    /**
     * Processes list reordering and emits updated model value.
     *
     * @param {Object} event - vuedraggable change event (may include moved/added/removed).
     * @returns {void}
     */
    const updateOrdering = event => updateModel(dragDropChangePositionHandler(props.value, event));

    return {
      classes,
      startDragging,
      endDragging,
      updateOrdering,
    };
  },
};
</script>

<style lang="scss">
.c-draggable-list-field--dragging {
  // We need to put it to avoid problem with `dragend` event on text-editon field with content
  .text-editor {
    pointer-events: none;
  }
}
</style>
