<template>
  <div
    :class="{ 'actions-panel--small': small }"
    class="actions-panel"
  >
    <v-layout
      :wrap="wrap"
      align-center
    >
      <c-action-btn
        v-for="(action, index) in preparedActions.inline"
        :key="index"
        :tooltip="action.title"
        :disabled="action.disabled"
        :loading="action.loading"
        :icon="action.icon"
        :color="action.iconColor"
        :badge-value="action.badgeValue"
        :badge-tooltip="action.badgeTooltip"
        @click="action.method"
      />
      <span v-if="preparedActions.dropDown.length">
        <v-menu
          key="dropdown-menu"
          bottom
          left
          @click.native.stop=""
        >
          <template #activator="{ on }">
            <v-btn
              icon
              v-on="on"
            >
              <v-icon>more_vert</v-icon>
            </v-btn>
          </template>
          <v-list>
            <v-list-item
              v-for="(action, index) in preparedActions.dropDown"
              :key="index"
              :disabled="action.disabled || action.loading"
              @click.stop="action.method"
            >
              <v-list-item-title>
                <span class="mr-4">
                  <v-progress-circular
                    v-if="action.loading"
                    :color="action.iconColor"
                    :size="16"
                    :width="2"
                    indeterminate
                  />
                  <v-icon
                    v-else
                    :color="action.iconColor"
                    :disabled="action.disabled"
                    class="ma-0 pa-0"
                    left
                    small
                  >
                    {{ action.icon }}
                  </v-icon>
                </span>
                <span
                  :class="action.cssClass"
                  class="text-body-1"
                >
                  {{ action.title }}
                </span>
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </span>
    </v-layout>
  </div>
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
