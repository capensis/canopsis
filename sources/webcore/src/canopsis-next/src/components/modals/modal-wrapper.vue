<template lang="pug">
  v-card(
    :class="{ 'fill-min-height': fillHeight }",
    :data-test="$attrs['data-test']"
  )
    v-card-title.primary.white--text(v-if="$slots.title", :style="titleStyle")
      div.modal-wrapper__title.headline
        div
          slot(name="title")
        div
          modal-title-buttons(
            :minimize="minimize",
            :close="close"
          )
    template(v-if="!$modal.minimized")
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
  inject: ['$modal'],
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
      default: '',
    },
  },
  computed: {
    titleStyle() {
      const style = {};

      if (this.titleColor) {
        style.backgroundColor = this.titleColor;
      }

      return style;
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
