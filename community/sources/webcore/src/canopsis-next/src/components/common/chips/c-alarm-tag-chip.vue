<template lang="pug">
  v-chip.c-alarm-tag-chip(
    :class="{ 'c-alarm-tag-chip--closable': close, 'c-alarm-tag-chip--small': small }",
    :color="color",
    small,
    @click="$emit('click')"
  )
    span.white--text
      slot
    v-icon.cursor-pointer.ml-2(
      v-if="close",
      color="white",
      @click.stop="$emit('close')"
    ) cancel
</template>

<script>
export default {
  props: {
    color: {
      type: String,
      required: false,
    },
    close: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.c-alarm-tag-chip {
  border-radius: 5px;
  font-size: 12px;

  &--closable {
    .v-chip__content {
      padding-right: 4px;
    }
  }

  &--small {
    height: 20px !important;
    margin: 2px !important;
  }

  .v-datatable thead th.column.sortable & .v-icon {
    opacity: .6;
  }

  .v-icon {
    transform: unset !important;
  }
}
</style>
