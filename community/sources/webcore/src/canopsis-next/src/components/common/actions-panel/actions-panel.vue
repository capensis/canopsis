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
        @click="action.method"
      )
      span.ml-1(v-if="preparedActions.dropDown.length")
        v-menu(
          key="dropdown-menu",
          bottom,
          left,
          @click.native.stop=""
        )
          template(#activator="{ on }")
            v-btn.ma-0(v-on="on", icon)
              v-icon more_vert
          v-list
            v-list-tile(
              v-for="(action, index) in preparedActions.dropDown",
              :key="index",
              :disabled="action.disabled",
              @click.stop="action.method"
            )
              v-list-tile-title
                v-progress-circular.actions-panel__menu-item-loader(
                  v-if="action.loading",
                  :color="action.iconColor",
                  size="16",
                  width="2",
                  left,
                  indeterminate
                )
                v-icon.pr-3(
                  v-else,
                  :color="action.iconColor",
                  :disabled="action.disabled",
                  left,
                  small
                ) {{ action.icon }}
                span.body-1(:class="action.cssClass") {{ action.title }}
</template>

<script>
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
    small: {
      type: Boolean,
      default: false,
    },
    wrap: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    preparedActions() {
      if (this.$mq === 'xl') {
        return {
          inline: this.actions.slice(0, this.inlineCount),
          dropDown: this.actions.slice(this.inlineCount),
        };
      }

      return {
        inline: [],
        dropDown: this.actions,
      };
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
