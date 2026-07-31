<template>
  <draggable
    :model-value="currentValue"
    :group="group"
    :tag="component"
    :animation="animation"
    :disabled="disabled"
    :handle="handle"
    :ghost-class="ghostClass"
    :drag-class="dragClass"
    :chosen-class="chosenClass || dragClass"
    :component-data="componentData"
    :move="itemMove"
    :item-key="itemKey"
    :class="classes"
    @update:model-value="updateModel"
    @start="startDragging"
    @end="endDragging"
  >
    <template #item="{ element, index }">
      <div class="c-draggable-list-field__item">
        <slot
          :element="element"
          :index="index"
          name="item"
        />
      </div>
    </template>

    <template #footer>
      <slot name="footer" />
    </template>
  </draggable>
</template>

<script>
import { computed, ref } from 'vue';
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

export default {
  components: { Draggable },

  props: {
    modelValue: {
      type: Array,
      default: undefined,
    },
    value: {
      type: Array,
      default: undefined,
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
      default: 'div',
    },
    componentData: {
      type: Object,
      default: undefined,
    },
    handle: {
      type: String,
      default: undefined,
    },
    ghostClass: {
      type: String,
      default: 'grey',
    },
    dragClass: {
      type: String,
      default: undefined,
    },
    chosenClass: {
      type: String,
      default: undefined,
    },
    group: {
      type: [Object, String],
      default: undefined,
    },
    itemMove: {
      type: Function,
      default: undefined,
    },
    draggable: {
      type: String,
      default: '>*',
    },
    itemKey: {
      type: [String, Function],
      default: '_id',
    },
  },

  emits: ['update:modelValue', 'update:value', 'change', 'start', 'end'],

  setup(props, { emit }) {
    /**
     * Tracks whether a drag operation is currently active.
     */
    const isDragging = ref(false);
    const currentValue = computed(() => props.modelValue ?? props.value ?? []);

    /**
     * Classes applied to the root draggable element depending on state.
     */
    const classes = computed(() => ({
      'c-draggable-list-field': true,
      'c-draggable-list-field--dragging': isDragging.value,
    }));

    /**
     * vuedraggable@4 already scopes draggable items with [data-draggable].
     * The old Vuetify 2 wrapper default used `>*`, which would become
     * `[data-draggable]>*` and make Sortable drag child DOM nodes instead of items.
     */
    const sortableDraggable = computed(() => (props.draggable === '>*' ? undefined : props.draggable));

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
     * Emits the updated list produced by vuedraggable.
     *
     * @param {Array} value - Updated list value.
     * @returns {void}
     */
    const updateModel = (value) => {
      emit('update:modelValue', value);
      emit('update:value', value);
      emit('change', value);
    };

    return {
      classes,
      currentValue,
      sortableDraggable,
      startDragging,
      endDragging,
      updateModel,
    };
  },
};
</script>

<style lang="scss">
.c-draggable-list-field {
  &--dragging {
    // We need to put it to avoid problem with `dragend` event on text-editon field with content
    .text-editor {
      pointer-events: none;
    }
  }

  &__item {
    width: 100%;
  }
}
</style>
