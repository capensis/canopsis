<template lang="pug">
  v-card(:class="{ 'fill-min-height': fillHeight }")
    v-card-title.white--text(v-if="$slots.title", :style="titleStyle")
      div.modal-wrapper__title.headline
        div
          slot(name="title")
        div
          modal-title-buttons(
            :minimize="minimize",
            :close="close"
          )
    template(v-if="!$modal.minimized")
      v-card-text(v-if="$slots.text", :class="textClass", key="text")
        slot(name="text")
      template(v-if="$slots.actions")
        v-divider(key="divider")
        v-card-actions(key="actions")
          v-layout.py-1(justify-end, align-center)
            slot(name="actions")
</template>

<script>
import { CSS_COLORS_VARS } from '@/config';

import ModalTitleButtons from './modal-title-buttons.vue';

export default {
  inject: ['$modal'],
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
      type: [Boolean, Function],
      default: false,
    },
    titleColor: {
      type: String,
      default: CSS_COLORS_VARS.primary,
    },
    textClass: {
      type: String,
      required: false,
    },
  },
  computed: {
    titleStyle() {
      return {
        backgroundColor: this.titleColor,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.modal-wrapper__title {
  display: flex;
  justify-content: space-between;
  width: 100%;
  align-items: center;

  & > div {
    display: flex;
  }
}
</style>
