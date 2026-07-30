<template>
  <v-expansion-panel
    :class="{ 'group-panel--editing': isEditing }"
    class="group-panel"
    color="secondary"
    elevation="0"
    rounded="0"
  >
    <v-expansion-panel-title
      :class="{ 'group-panel__active': active }"
      :hide-actions="hideActions"
      class="pa-0 pr-6"
    >
      <div class="px-6 py-4 group-panel__title">
        <slot name="title">
          {{ group.title }}
        </slot>
      </div>
      <div class="group-panel__actions">
        <v-btn
          v-show="editable"
          :disabled="orderChanged"
          size="small"
          variant="text"
          icon
          @click.stop="handleChange"
        >
          <c-icon size="18">
edit
</c-icon>
        </v-btn>
      </div>
    </v-expansion-panel-title>
    <v-expansion-panel-text class="group-item__content" eager>
      <slot />
    </v-expansion-panel-text>
  </v-expansion-panel>
</template>

<script>
export default {
  props: {
    isEditing: {
      type: Boolean,
      default: false,
    },
    active: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    group: {
      type: Object,
      required: true,
    },
    orderChanged: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['change'],

  setup(props, { emit }) {
    const handleChange = () => {
      emit('change');
    };

    return {
      handleChange,
    };
  },
};
</script>

<style lang="scss" scoped>
.group-panel {
  :deep(.v-expansion-panel-title) {
    height: 48px;
    min-height: 48px;
  }

  :deep(.v-expansion-panel-text__wrapper) {
    padding: 0;
  }

  &__actions {
    flex-shrink: 0;
  }

  &__title {
    width: 100%;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    display: inline-block;
    padding: 16px 24px;
  }

  &--editing {
    :deep(.v-expansion-panel-title) {
      cursor: move;
    }

    .views-panel--empty {
      &:after {
        content: '';
        display: block;
        height: 48px;
        border: 4px dashed #4f6479;
        border-radius: 5px;
        position: relative;
      }
    }
  }
}
</style>
