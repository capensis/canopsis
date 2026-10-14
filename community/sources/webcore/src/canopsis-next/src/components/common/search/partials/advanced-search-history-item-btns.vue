<template>
  <c-layout>
    <c-action-btn type="delete" class="advanced-search__history__item__btn" @click="remove" />
    <c-action-btn
      :class="{
        'advanced-search__history__item__btn--pinned': pinned,
      }"
      :color="pinned ? 'primary' : 'grey'"
      :tooltip="$t('common.pin')"
      class="advanced-search__history__item__btn"
      icon="push_pin"
      @click="$emit('toggle-pin', id)"
    />
  </c-layout>
</template>

<script>
export default {
  props: {
    id: {
      type: String,
      required: true,
    },
    pinned: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['remove', 'toggle-pin'],

  methods: {
    remove() {
      this.$emit('remove', this.id);
    },
  },
};
</script>

<style lang="scss">
.v-list-item {
  .advanced-search__history__item__btn.c-action-btn:not(.advanced-search__history__item__btn--pinned) {
    opacity: 0;
  }

  &:hover .advanced-search__history__item__btn.c-action-btn {
    opacity: 1;
  }
}
</style>
