<template lang="pug">
  v-card.card-with-see-alarms-btn(
    v-on="$listeners",
    :class="itemClasses",
    tile,
    dark
  )
    slot
    v-btn.card-with-see-alarms-btn__btn(
      v-if="showButton",
      flat,
      @click.stop="$emit('show:alarms')"
    ) {{ $t('common.seeAlarms') }}
</template>

<script>
export default {
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
    font-size: .6em;
    height: var(--see-alarms-btn-height);
    margin: 0;
    background-color: rgba(0, 0, 0, .2);

    &.v-btn--active:before, &.v-btn:focus:before, &.v-btn:hover:before {
      background-color: rgba(0, 0, 0, .5);
    }
  }
}
</style>
