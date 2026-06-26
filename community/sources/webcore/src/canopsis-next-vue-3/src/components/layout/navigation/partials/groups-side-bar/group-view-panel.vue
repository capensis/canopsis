<template>
  <v-card
    :class="[
      'group-view-panel',
      'secondary',
      isViewActive ? 'lighten-3' : 'lighten-1',
    ]"
    elevation="0"
  >
    <v-card-text class="panel-item-content">
      <div :class="['panel-item-content__title pl-2', { 'text-truncate': ellipsis }]">
        <slot name="title">
          {{ view.title }}
        </slot>
      </div>

      <div class="panel-item-content__actions">
        <c-icon
          v-if="view.is_private"
          size="18"
        >
          lock
        </c-icon>
        <template v-if="editable || duplicable">
          <v-btn
            v-if="editable"
            :disabled="isOrderChanged"
            size="small"
            variant="text"
            icon
            @click.prevent="$emit('change')"
          >
            <c-icon size="18">edit</c-icon>
          </v-btn>
          <v-btn
            v-if="duplicable"
            :disabled="isOrderChanged"
            size="small"
            variant="text"
            icon
            @click.prevent="$emit('duplicate')"
          >
            <c-icon size="18">file_copy</c-icon>
          </v-btn>
        </template>
      </div>
    </v-card-text>
    <v-divider class="group-view-panel__divider" />
  </v-card>
</template>

<script>
export default {
  props: {
    view: {
      type: Object,
      required: true,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
    isOrderChanged: {
      type: Boolean,
      default: false,
    },
    isViewActive: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['change', 'duplicate'],

  computed: {
    ellipsis() {
      return !this.$slots.title;
    },
  },
};
</script>

<style lang="scss" scoped>
.group-view-panel {
  border-radius: 0;

  &__divider {
    border-color: rgba(0, 0, 0, .12) !important;
    opacity: 1;
  }
}

.panel-item-content {
  display: flex;
  cursor: pointer;
  align-items: center;
  justify-content: space-between;
  position: relative;
  padding: 12px 24px;
  height: 48px;

  &__title {
    width: 100%;
  }

  &__actions {
    flex-shrink: 0;

    display: flex;
    align-items: center;
  }

  :deep(.v-btn:not(:last-child)) {
    margin-right: 0;
  }
}
</style>
