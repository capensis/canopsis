<template lang="pug">
  div.c-card-iterator-item
    v-layout(row, align-center, justify-space-between)
      v-layout.c-card-iterator-item__actions.pr-2
        c-draggable-step-number(
          :color="!expanded && hasChildrenError ? 'error' : 'primary'",
          drag-class="item-drag-handler"
        ) {{ itemNumber }}
        c-expand-btn(v-model="expanded")
      slot(name="header")
      c-action-btn(type="delete", @click="$emit('remove')")
    v-expand-transition(mode="out-in")
      v-layout.c-card-iterator-item__content(
        v-show="expanded",
        :class="{ 'c-card-iterator-item__content--offset': offsetLeft }",
        column
      )
        slot
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
