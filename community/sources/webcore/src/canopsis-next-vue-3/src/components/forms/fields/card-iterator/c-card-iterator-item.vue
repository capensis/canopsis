<template>
  <div :class="{ 'c-card-iterator-item--small': small }" class="c-card-iterator-item">
    <c-layout
      class="c-card-iterator-item__row"
      gap="2"
      align="start"
      justify="space-between"
    >
      <c-layout
        :class="{ 'c-card-iterator-item__actions--draggable-only': !isExpandable }"
        class="c-card-iterator-item__actions"
        align="center"
      >
        <c-draggable-step-number
          :color="!expanded && hasChildrenError ? 'error' : 'primary'"
          :hide-number="small"
          :drag-class="dragHandleClass"
          class="c-card-iterator-item__step-number"
        >
          {{ itemNumber }}
        </c-draggable-step-number>
        <c-expand-btn v-if="isExpandable" v-model="expanded" />
      </c-layout>
      <div class="c-card-iterator-item__header">
        <slot name="header" />
      </div>
      <c-action-btn
        v-if="removable"
        :icon="small ? 'close' : 'delete'"
        :small="small"
        :icon-small="small"
        type="delete"
        @click="$emit('remove')"
      />
    </c-layout>
    <v-expand-transition v-if="isExpandable" mode="out-in">
      <c-layout
        v-show="expanded"
        :class="{ 'c-card-iterator-item__content--offset': offsetLeft }"
        class="c-card-iterator-item__content"
        column
      >
        <slot />
      </c-layout>
    </v-expand-transition>
  </div>
</template>

<script>
import { computed, ref, useSlots } from 'vue';

import { useValidationChildren } from '@/hooks/validator/validation-children';

export default {
  props: {
    defaultExpanded: {
      type: Boolean,
      default: false,
    },
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
    removable: {
      type: Boolean,
      default: true,
    },
    expandable: {
      type: Boolean,
      default: true,
    },
  },
  setup(props) {
    const slots = useSlots();
    const expanded = ref(props.defaultExpanded);
    const { hasChildrenError } = useValidationChildren();
    const isExpandable = computed(() => props.expandable && !!slots.default);

    return {
      expanded,
      hasChildrenError,
      isExpandable,
    };
  },
};
</script>

<style lang="scss">
.c-card-iterator-item {
  --actions-max-width: 100px;
  --actions-min-width: 60px;
  --one-action-width: 26px;
  --first-field-height: 48px;

  &__row > .c-action-btn {
    position: absolute;
    top: 4px;
    right: 4px;
    z-index: 1;
  }

  &--small {
    .c-card-iterator-item__row > .c-action-btn {
      width: 24px;
      min-width: 24px;
      height: 24px;

      .v-btn {
        width: 24px;
        height: 24px;
        min-width: 24px;
        padding: 0;
      }

      .v-icon {
        font-size: 16px;
      }
    }
  }

  &__row {
    align-items: flex-start;
  }

  &__header {
    display: flex;
    flex: 1 1 auto;
    min-width: 0;

    > * {
      flex: 1 1 auto;
      min-width: 0;
      width: 100%;
    }
  }

  &__actions {
    align-self: flex-start;
    align-items: center;
    gap: 8px;
    height: var(--first-field-height);
    max-width: var(--actions-max-width);
    min-width: var(--actions-min-width);

    &--draggable-only {
      min-width: var(--one-action-width);
      max-width: var(--one-action-width);
    }
  }

  &__step-number {
    gap: 8px;
  }

  &__content {
    &--offset {
      margin-left: var(--actions-max-width);
    }
  }
}
</style>
