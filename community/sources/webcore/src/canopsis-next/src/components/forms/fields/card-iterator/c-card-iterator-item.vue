<template>
  <div
    :class="{ 'c-card-iterator-item--small': small }"
    class="c-card-iterator-item"
  >
    <v-layout
      class="gap-2"
      align-center
      justify-space-between
    >
      <v-layout
        :class="{ 'c-card-iterator-item__actions--draggable-only': !$slots.default }"
        class="c-card-iterator-item__actions"
      >
        <c-draggable-step-number
          :color="!expanded && hasChildrenError ? 'error' : 'primary'"
          :hide-number="small"
          :drag-class="dragHandleClass"
        >
          {{ itemNumber }}
        </c-draggable-step-number>
        <c-expand-btn v-if="$slots.default" v-model="expanded" />
      </v-layout>
      <slot name="header" />
      <c-action-btn
        :icon="small ? 'close' : 'delete'"
        :small="small"
        :icon-small="small"
        class="c-card-iterator-item__remove-btn"
        type="delete"
        @click="$emit('remove')"
      />
    </v-layout>
    <v-expand-transition v-if="$slots.default" mode="out-in">
      <v-layout
        v-show="expanded"
        :class="{ 'c-card-iterator-item__content--offset': offsetLeft }"
        class="c-card-iterator-item__content"
        column
      >
        <slot />
      </v-layout>
    </v-expand-transition>
  </div>
</template>

<script>
import { ref } from 'vue';

import { useValidationChildren } from '@/hooks/validator/validation-children';

export default {
  inject: ['$validator'],
  props: {
    itemNumber: {
      type: [Number, String],
      default: 0,
    },
    offsetLeft: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    dragHandleClass: {
      type: String,
      default: 'item-drag-handler',
    },
  },
  setup() {
    const expanded = ref(false);

    const { hasChildrenError } = useValidationChildren();

    return {
      expanded,
      hasChildrenError,
    };
  },
};
</script>

<style lang="scss">
.c-card-iterator-item {
  --actions-max-width: 100px;
  --actions-min-width: 60px;
  --one-action-width: 26px;

  &__actions {
    max-width: var(--actions-max-width);
    min-width: var(--actions-min-width);

    &--draggable-only {
      min-width: var(--one-action-width);
      max-width: var(--one-action-width);
    }
  }

  &__content {
    &--offset {
      margin-left: var(--actions-max-width);
    }
  }

  &--small .c-card-iterator-item__remove-btn {
    position: absolute;
    top: 2px;
    right: 2px;
    margin: 0 !important;
  }
}
</style>
