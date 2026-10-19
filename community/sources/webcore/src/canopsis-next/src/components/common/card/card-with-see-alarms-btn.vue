<template>
  <v-card
    :class="itemClasses"
    class="card-with-see-alarms-btn"
    theme="dark"
    tile
    v-bind="$attrs"
  >
    <slot />

    <v-btn
      v-if="showButton"
      class="card-with-see-alarms-btn__btn"
      variant="text"
      @click.stop="$emit('show:alarms')"
    >
      {{ $t('common.seeAlarms') }}
    </v-btn>
  </v-card>
</template>

<script>
export default {
  inheritAttrs: false,

  emits: ['show:alarms'],

  props: {
    showButton: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    itemClasses() {
      return {
        'card-with-see-alarms-btn--with-btn': this.showButton,
      };
    },
  },
};
</script>

<style lang="scss">
.card-with-see-alarms-btn {
  --see-alarms-btn-height: 18px;

  min-height: unset !important;

  a {
    color: white;
  }

  &--with-btn {
    padding-bottom: var(--see-alarms-btn-height);
  }

  &__btn {
    position: absolute;
    bottom: 0;
    width: 100%;
    min-width: 0;
    height: var(--see-alarms-btn-height) !important;
    margin: 0;
    color: white;
    background-color: rgba(0, 0, 0, 0.2);

    .v-btn__content {
      color: white;
    }

    &.v-btn--active,
    &:focus,
    &:hover {
      .v-btn__overlay {
        background-color: rgba(0, 0, 0, 0.5);
      }
    }
  }
}
</style>
