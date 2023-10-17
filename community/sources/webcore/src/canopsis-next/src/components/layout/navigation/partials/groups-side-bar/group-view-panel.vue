<template>
  <v-card :color="cardColor">
    <div class="panel-item-content">
      <div
        :class="['panel-item-content__title pl-2', {
          'panel-view__title--editing': isEditing,
          'text-truncate': ellipsis
        }]"
      >
        <slot name="title">
          {{ view.title }}
        </slot>
      </div>

      <div
        v-if="allowEditing"
        class="panel-item-content__actions"
      >
        <v-btn
          class="ma-0"
          v-show="hasEditAccess"
          :disabled="isOrderChanged"
          depressed
          small
          icon
          @click.prevent="$emit('change')"
        >
          <v-icon small>
            edit
          </v-icon>
        </v-btn>
        <v-btn
          class="ma-0"
          v-show="isEditing"
          :disabled="isOrderChanged"
          depressed
          small
          icon
          @click.prevent="$emit('duplicate')"
        >
          <v-icon small>
            file_copy
          </v-icon>
        </v-btn>
      </div>
    </div>
    <v-divider />
  </v-card>
</template>

<script>
export default {
  props: {
    view: {
      type: Object,
      required: true,
    },
    allowEditing: {
      type: Boolean,
      default: false,
    },
    hasEditAccess: {
      type: Boolean,
      default: false,
    },
    isEditing: {
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
  computed: {
    cardColor() {
      return `secondary ${this.isViewActive ? 'lighten-3' : 'lighten-1'}`;
    },

    ellipsis() {
      return !this.$slots.title && !this.$scopedSlots.title;
    },
  },
};
</script>

<style lang="scss" scoped>
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

    & ::v-deep .v-btn:not(:last-child) {
      margin-right: 0;
    }
  }
</style>
