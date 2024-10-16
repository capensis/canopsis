<template>
  <div
    :class="{ 'actions-panel--small': small }"
    class="actions-panel"
  >
    <v-layout align-center>
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
      <v-menu
        v-if="preparedActions.dropDown.length"
        v-model="opened"
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
    </v-layout>
  </div>
</template>

<script>
import { DEFAULT_ALARM_ACTIONS_INLINE_COUNT } from '@/constants';

export default {
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    inlineCount: {
      type: Number,
      default: DEFAULT_ALARM_ACTIONS_INLINE_COUNT,
    },
    small: {
      type: Boolean,
      default: false,
    },
    ignoreMediaQuery: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      opened: false,
    };
  },
  computed: {
    preparedActions() {
      if (!this.ignoreMediaQuery && this.$mq !== 'xl') {
        return {
          inline: [],
          dropDown: this.actions,
        };
      }

      if (this.inlineCount < this.actions.length) {
        const inlineCountWithoutMenu = this.inlineCount - 1;

        return {
          inline: this.actions.slice(0, inlineCountWithoutMenu),
          dropDown: this.actions.slice(inlineCountWithoutMenu),
        };
      }

      return {
        inline: this.actions,
        dropDown: [],
      };
    },
  },
  methods: {
    closeMenu() {
      this.opened = false;
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
