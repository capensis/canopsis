<template>
  <div class="c-card-iterator-item">
    <v-layout
      align-center
      justify-space-between
    >
      <v-layout class="c-card-iterator-item__actions pr-2">
        <c-draggable-step-number
          :color="!expanded && hasChildrenError ? 'error' : 'primary'"
          drag-class="item-drag-handler"
        >
          {{ itemNumber }}
        </c-draggable-step-number>
        <c-expand-btn v-model="expanded" />
      </v-layout>
      <slot name="header" />
      <c-action-btn
        type="delete"
        @click="$emit('remove')"
      />
    </v-layout>
    <v-expand-transition mode="out-in">
      <v-layout
        class="c-card-iterator-item__content"
        v-show="expanded"
        :class="{ 'c-card-iterator-item__content--offset': offsetLeft }"
        column
      >
        <slot />
      </v-layout>
    </v-expand-transition>
  </div>
</template>

<script>
import { validationChildrenMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [
    validationChildrenMixin,
  ],
  props: {
    itemNumber: {
      type: [Number, String],
      default: 0,
    },
    offsetLeft: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      expanded: true,
    };
  },
};
</script>

<style lang="scss">
$actionsWidth: 100px;

.c-card-iterator-item {
  &__actions {
    max-width: $actionsWidth;
  }

  &__content {
    &--offset {
      margin-left: $actionsWidth;
    }
  }
}
</style>
