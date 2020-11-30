<template lang="pug">
  v-card(
    :class="{ 'fill-min-height': fillHeight }",
    :data-test="$attrs['data-test']"
  )
    v-card-title.primary.white--text(v-if="$slots.title && !$slots.fullTitle")
      v-layout.headline(justify-space-between, align-center)
        v-flex
          slot(name="title")
        v-flex(v-if="minimize || close")
          modal-title-buttons(
            :modal="modal",
            :minimize="minimize",
            :close="close"
          )
    slot(name="fullTitle")
    template(v-if="!modal.minimized")
      v-card-text(v-if="$slots.text", key="text")
        slot(name="text")
      template(v-if="$slots.actions")
        v-divider(key="divider")
        v-card-actions(key="actions")
          v-layout.py-1(justify-end)
            slot(name="actions")
</template>

<script>
import ModalTitleButtons from './modal-title-buttons.vue';

export default {
  components: { ModalTitleButtons },
  props: {
    fillHeight: {
      type: Boolean,
      default: false,
    },
    minimize: {
      type: Boolean,
      default: false,
    },
    close: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    modal() {
      return this.$parent.modal || this.$parent.$parent.modal;
    },
  },
};
</script>
