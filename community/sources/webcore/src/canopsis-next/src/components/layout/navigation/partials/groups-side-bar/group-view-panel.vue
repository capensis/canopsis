<template>
  <v-card :color="cardColor">
    <v-card-text class="panel-item-content">
      <v-layout
        align-center
        justify-space-between
      >
        <v-flex :class="{ 'panel-view-title--editing': isEditing }">
          <v-layout align-center>
            <span
              class="pl-2"
              :class="{ ellipsis: ellipsis }"
            >
              <slot name="title">{{ view.title }}</slot></span>
          </v-layout>
        </v-flex>
        <v-flex>
          <v-layout
            v-if="allowEditing"
            justify-end
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
          </v-layout>
        </v-flex>
      </v-layout>
    </v-card-text>
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
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    cursor: pointer;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    position: relative;
    padding: 12px 24px;
    height: 48px;

    & > div {
      max-width: 100%;
    }

    & ::v-deep .v-btn:not(:last-child) {
      margin-right: 0;
    }
  }

  .panel-view-title {
    &--editing {
      max-width: 73%;
    }
  }
</style>
