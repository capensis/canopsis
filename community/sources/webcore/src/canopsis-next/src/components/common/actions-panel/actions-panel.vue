<template lang="pug">
  div.actions-panel(:class="{ 'actions-panel--small': small }")
    v-layout(:wrap="wrap", row, align-center)
      c-action-btn(
        v-for="(action, index) in preparedActions.inline",
        :key="index",
        :tooltip="action.title",
        :disabled="action.disabled",
        :loading="action.loading",
        :icon="action.icon",
        :color="action.iconColor",
        :badge-value="action.badgeValue",
        :badge-tooltip="action.badgeTooltip",
        lazy,
        @click="action.method"
      )
      span(v-if="preparedActions.dropDown.length")
        v-menu(
          key="dropdown-menu",
          bottom,
          left,
          lazy,
          @click.native.stop=""
        )
          template(#activator="{ on }")
            v-btn.ma-0(v-on="on", icon)
              v-icon more_vert
          v-list
            v-list-tile(
              v-for="(action, index) in preparedActions.dropDown",
              :key="index",
              :disabled="action.disabled || action.loading",
              @click.stop="action.method"
            )
              v-list-tile-title
                span.mr-4
                  v-progress-circular(
                    v-if="action.loading",
                    :color="action.iconColor",
                    :size="16",
                    :width="2",
                    indeterminate
                  )
                  v-icon.ma-0.pa-0(
                    v-else,
                    :color="action.iconColor",
                    :disabled="action.disabled",
                    left,
                    small
                  ) {{ action.icon }}
                span.body-1(:class="action.cssClass") {{ action.title }}
</template>

<script>
import { ALARM_ACTION_BUTTON_MARGINS, ALARM_ACTION_BUTTON_WIDTHS, ALARM_DENSE_TYPES } from '@/constants';

export default {
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    inlineCount: {
      type: Number,
      default: 3,
    },
    medium: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    wrap: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      width: 0,
    };
  },
  computed: {
    denseType() {
      if (this.small) {
        return ALARM_DENSE_TYPES.small;
      }

      return this.medium
        ? ALARM_DENSE_TYPES.medium
        : ALARM_DENSE_TYPES.large;
    },

    preparedActions() {
      const actionWidth = ALARM_ACTION_BUTTON_WIDTHS[this.denseType] + ALARM_ACTION_BUTTON_MARGINS[this.denseType];
      const maxActionsCount = Math.floor(this.width / actionWidth) - 1;
      const countForSlice = Math.min(maxActionsCount, this.inlineCount);
      const countForSliceWithMenu = countForSlice + 1;

      if (countForSliceWithMenu < this.actions.length) {
        return {
          inline: this.actions.slice(0, countForSlice),
          dropDown: this.actions.slice(countForSlice),
        };
      }

      return {
        inline: this.actions,
        dropDown: [],
      };
    },
  },
  mounted() {
    this.$resizeObserver = new ResizeObserver(this.resizeObserverHandler);
    this.$resizeObserver.observe(this.$el);
  },
  beforeDestroy() {
    this.$resizeObserver.disconnect();
  },
  methods: {
    resizeObserverHandler([entry]) {
      const { target: { offsetWidth } } = entry ?? {};

      if (this.width !== offsetWidth) {
        this.width = offsetWidth;
      }
    },
  },
};
</script>

<style lang="scss">
.actions-panel {
  &__menu-item-loader {
    margin-right: 32px;
  }

  &--small {
    .v-btn--icon {
      width: 24px;
      height: 24px;

      .v-icon {
        font-size: 20px;
      }
    }
  }
}
</style>
