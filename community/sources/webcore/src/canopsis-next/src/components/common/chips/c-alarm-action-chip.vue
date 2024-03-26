<template>
  <v-chip
    :class="chipClass"
    :color="color"
    :text-color="textColor"
    :outlined="outlined"
    class="c-alarm-action-chip"
    small
    @click="$emit('click')"
  >
    <span class="c-alarm-action-chip__text">
      <slot />
    </span>
    <v-icon
      v-if="closable"
      class="cursor-pointer ml-2"
      @click.stop="$emit('close')"
    >
      cancel
    </v-icon>
  </v-chip>
</template>

<script>
export default {
  props: {
    color: {
      type: String,
      required: false,
    },
    closable: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    outlined: {
      type: Boolean,
      default: false,
    },
    textColor: {
      type: String,
      default: 'white',
    },
  },
  computed: {
    chipClass() {
      return {
        'c-alarm-action-chip--closable': this.closable,
        'c-alarm-action-chip--small': this.small,
      };
    },
  },
};
</script>

<style lang="scss">
.c-alarm-action-chip.v-chip {
  border-radius: 5px;
  font-size: 12px;
  min-height: 24px;
  height: unset !important;
  padding: 0;

  &__text {
    white-space: initial;
    word-wrap: break-word;
    max-width: 100%;
    overflow: hidden;
  }

  .v-chip__content {
    height: unset !important;
    cursor: pointer;
    max-width: 100%;
  }

  &--closable {
    .v-chip__content {
      padding-right: 4px;
    }
  }

  &--small {
    min-height: 20px !important;
    margin: 2px;
  }

  .v-data-table thead th.column.sortable & .v-icon {
    opacity: .6;
  }

  .v-icon {
    transform: unset !important;
  }
}
</style>
