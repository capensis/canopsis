<template lang="pug">
  div.modal-title-buttons(
    :class="{ 'close': close, 'minimize': minimize }"
  )
    div.modal-title-button__wrapper(v-if="minimize")
      v-tooltip(
        v-if="!modal.minimized",
        :disabled="!hasMinimizedModal",
        left
      )
        v-btn.v-btn--minimize(
          slot="activator",
          :disabled="hasMinimizedModal",
          icon,
          @click="$modals.minimize({ id: modal.id })"
        )
          v-icon(color="white") minimize
        span {{ $t('modals.common.titleButtons.minimizeTooltip') }}
      v-btn(
        v-else,
        icon,
        @click="$modals.maximize({ id: modal.id })"
      )
        v-icon(color="white") maximize
    div.modal-title-button__wrapper(v-if="close")
      v-btn(
        slot="activator",
        icon,
        @click="closeHandler"
      )
        v-icon(color="white") close
</template>

<script>
import { isFunction } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('modals');

export default {
  props: {
    modal: {
      type: Object,
      required: true,
    },
    minimize: {
      type: Boolean,
      default: false,
    },
    close: {
      type: [Boolean, Function],
      default: false,
    },
  },
  computed: {
    ...mapGetters(['hasMinimizedModal']),

    closeHandler() {
      if (!this.close) {
        return false;
      }

      return isFunction(this.close)
        ? this.close
        : () => this.$modals.hide({ id: this.modal.id });
    },
  },
};
</script>

<style lang="scss" scoped>
.modal-title-buttons {
  display: flex;

  .v-btn {
    margin: 0 8px;
  }

  &.close, &.minimize {
    width: 48px;
  }

  &.close.minimize {
    width: 96px;
  }

  .v-dialog:not(.v-dialog--fullscreen) & {
    .modal-title-button__wrapper .v-btn {
      font-size: 13px;
      height: 28px;
      width: 28px;

      .v-icon {
        font-size: 16px;
      }
    }

    &.close, &.minimize {
      width: 40px;
    }

    &.close.minimize {
      width: 80px;
    }
  }

  .v-btn--minimize {
    pointer-events: auto;
  }
}
</style>
