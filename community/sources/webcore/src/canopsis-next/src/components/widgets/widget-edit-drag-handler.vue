<template lang="pug">
  div.widget-edit-drag-handler(@pointerdown.left="$emit('pointerdown', $event)")
    v-layout.widget-edit-drag-handler__controls.pr-1(justify-end, align-center)
      span(@pointerdown.stop="")
        v-tooltip(bottom)
          template(#activator="{ on }")
            v-btn.ma-0.mr-1(
              v-on="on",
              :color="autoHeight ? 'grey lighten-1' : 'transparent'",
              icon,
              small,
              @click="$emit('toggle')"
            )
              v-icon(
                :color="autoHeight ? 'black' : 'grey darken-1'",
                small
              ) lock
          span {{ $t('view.autoHeightButton') }}
        widget-wrapper-menu(
          :widget="widget",
          :tab="tab"
        )
</template>

<script>
import WidgetWrapperMenu from '@/components/widgets/partials/widget-wrapper-menu.vue';

export default {
  components: { WidgetWrapperMenu },
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
    autoHeight: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.widget-edit-drag-handler {
  position: absolute;
  background-color: rgba(0, 0, 0, .12);
  width: 100%;
  height: 36px;
  transition: .2s ease-out;
  cursor: move;
  z-index: 3;

  &:hover {
    background-color: rgba(0, 0, 0, .15);
  }

  &__controls {
    height: 100%;
  }
}
</style>
