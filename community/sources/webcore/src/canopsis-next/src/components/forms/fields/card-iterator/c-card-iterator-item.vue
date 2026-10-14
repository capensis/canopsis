<template>
  <div
    :class="{
      'c-card-iterator-item--small': small,
      'c-card-iterator-item--static-delete': staticDelete,
    }"
    class="c-card-iterator-item"
  >
    <v-row
      class="c-card-iterator-item__row"
      align="start"
      justify="space-between"
    >
      <v-col
        class="c-card-iterator-item__actions"
        cols="auto"
      >
        <v-row align="center">
          <v-col cols="auto">
            <c-draggable-step-number
              :color="hasChildrenError ? 'error' : 'primary'"
              :disabled="disabled"
              :hide-number="small"
              :drag-class="dragHandleClass"
              class="c-card-iterator-item__step-number"
            >
              {{ itemNumber }}
            </c-draggable-step-number>
          </v-col>
          <v-col v-if="isExpandable" cols="auto">
            <c-expand-btn v-model="expanded" />
          </v-col>
        </v-row>
      </v-col>
      <v-col class="c-card-iterator-item__header">
        <slot name="header" />
      </v-col>
      <c-action-btn
        v-if="removable"
        :icon="small ? 'close' : 'delete'"
        :small="small"
        :icon-small="small"
        type="delete"
        @click="$emit('remove')"
      />
    </v-row>
    <v-expand-transition v-if="isExpandable" mode="out-in">
      <v-row
        v-show="expanded"
        :class="{ 'c-card-iterator-item__content--offset': offsetLeft }"
        class="c-card-iterator-item__content mt-3"
      >
        <v-col cols="12">
          <slot />
        </v-col>
      </v-row>
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
    disabled: {
      type: Boolean,
      default: false,
    },
    staticDelete: {
      type: Boolean,
      default: false,
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
  --first-field-height: 48px;

  &__row > .c-action-btn {
    position: absolute;
    top: 4px;
    right: 4px;
    z-index: 1;
  }

  &--static-delete {
    .c-card-iterator-item__row > .c-action-btn {
      position: static;
    }
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
    }
  }

  &__actions {
    display: flex;
    align-self: flex-start;
    align-items: center;
    height: var(--first-field-height);
  }

  &__content {
    &--offset {
      margin-left: var(--actions-max-width);
    }
  }
}
</style>
