<template>
  <v-layout>
    <v-btn
      class="advanced-search__history__item__btn"
      small
      icon
      @click.prevent.stop="remove"
    >
      <v-icon color="grey" small>
        delete
      </v-icon>
    </v-btn>
    <v-btn
      :class="{ 'advanced-search__history__item__btn--pinned': pinned }"
      small
      icon
      @click.prevent.stop="togglePin"
    >
      <v-icon :color="pinned ? 'inherit' : 'grey'" small>
        $vuetify.icons.push_pin
      </v-icon>
    </v-btn>
  </v-layout>
</template>

<script>
export default {
  props: {
    id: {
      type: [Number, String],
      default: '',
    },
    pinned: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const remove = () => emit('remove', props.id);
    const togglePin = () => emit('toggle-pin', props.id);

    return {
      remove,
      togglePin,
    };
  },
};
</script>

<style lang="scss">
.v-list-item {
  .advanced-search__history__item__btn.v-btn:not(.advanced-search__history__item__btn--pinned) {
    opacity: 0;
  }

  &:hover .advanced-search__history__item__btn.v-btn {
    opacity: 1;
  }
}
</style>
