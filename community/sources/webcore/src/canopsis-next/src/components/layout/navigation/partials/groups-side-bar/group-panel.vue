<template>
  <v-expansion-panel
    class="secondary group-panel"
    :hide-actions="hideActions"
    :class="{ 'group-panel--editing': isEditing }"
    :rounded="false"
  >
    <v-expansion-panel-header
      class="pa-0 pr-6"
      :hide-actions="hideActions"
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
          depressed
          small
          icon
          @click.stop="handleChange"
        >
          <v-icon small>
            edit
          </v-icon>
        </v-btn>
      </div>
    </v-expansion-panel-header>
    <v-expansion-panel-content class="group-item__content">
      <slot />
    </v-expansion-panel-content>
  </v-expansion-panel>
</template>

<script>
export default {
  props: {
    isEditing: {
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
  & ::v-deep .v-expansion-panel-header {
    height: 48px;
    min-height: 48px;
  }

  & ::v-deep .v-expansion-panel-content .v-card {
    border-radius: 0;
    box-shadow: 0 0 0 0 rgba(0,0,0,.2), 0 0 0 0 rgba(0,0,0,.14), 0 0 0 0 rgba(0,0,0,.12) !important;
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
    padding: 5px 0;
  }

  &--editing {
    & ::v-deep .v-expansion-panel-header {
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
