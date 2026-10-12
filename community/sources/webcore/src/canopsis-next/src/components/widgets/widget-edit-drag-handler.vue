<template>
  <div class="widget-edit-drag-handler" @pointerdown.left="$emit('pointerdown', $event)">
    <c-layout class="widget-edit-drag-handler__controls pr-1" justify="end" align="center">
      <span @pointerdown.stop="">
        <c-tooltip location="bottom">
          <template #activator="{ props }">
            <v-btn
              :color="autoHeight ? 'grey-lighten-1' : 'transparent'"
              class="ma-0 mr-1"
              size="small"
              variant="text"
              icon
              v-bind="props"
              @click="$emit('toggle')"
            >
              <v-icon :color="autoHeight ? 'black' : 'grey-darken-1'" size="small"> lock </v-icon>
            </v-btn>
          </template>

          <span>{{ $t('view.autoHeightButton') }}</span>
        </c-tooltip>

        <widget-wrapper-menu :widget="widget" :tab="tab" />
      </span>
    </c-layout>
  </div>
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
  emits: ['pointerdown', 'toggle'],
};
</script>

<style lang="scss" scoped>
.widget-edit-drag-handler {
  position: absolute;
  background-color: rgba(0, 0, 0, 0.12);
  width: 100%;
  height: 36px;
  transition: 0.2s ease-out;
  cursor: move;
  z-index: 3;

  &:hover {
    background-color: rgba(0, 0, 0, 0.15);
  }

  &__controls {
    height: 100%;
  }
}
</style>
