<template>
  <v-card :class="{ 'fill-min-height': fillHeight }">
    <v-card-title
      class="white--text"
      v-if="$slots.title"
      :style="titleStyle"
    >
      <div class="modal-wrapper__title text-h5">
        <div>
          <slot name="title" />
        </div>
        <div>
          <modal-title-buttons
            :minimize="minimize"
            :close="close"
          />
        </div>
      </div>
    </v-card-title>
    <template v-if="!$modal.minimized">
      <v-card-text
        v-if="$slots.text"
        :class="textClass"
        key="text"
      >
        <slot name="text" />
      </v-card-text>
      <template v-if="$slots.actions">
        <v-divider key="divider" />
        <v-card-actions
          key="actions"
          class="justify-end align-center px-2 py-3"
        >
          <slot name="actions" />
        </v-card-actions>
      </template>
    </template>
  </v-card>
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
