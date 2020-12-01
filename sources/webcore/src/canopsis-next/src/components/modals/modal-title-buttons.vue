<template lang="pug">
  div.modal-title-buttons(:class="{ 'close': close, 'minimize': minimize }")
    div(v-if="close")
      v-btn(
        icon,
        small,
        @click="closeHandler"
      )
        v-icon(color="white", small) close
    div(v-if="minimize")
      v-tooltip(
        v-if="!modal.minimized",
        :disabled="!hasMinimizedModal",
        left
      )
        v-btn.v-btn--minimize(
          slot="activator",
          :disabled="hasMinimizedModal",
          icon,
          small,
          @click="$modals.minimize({ id: modal.id })"
        )
          v-icon(color="white", small) minimize
        span {{ $t('modals.common.titleButtons.minimizeTooltip') }}
      v-btn(
        v-else,
        icon,
        small,
        @click="$modals.maximize({ id: modal.id })"
      )
        v-icon(color="white", small) maximize
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
  float: right;

  &.close, &.minimize {
    width: 45px;
  }

  &.close.minimize {
    width: 90px;
  }

  .v-btn--minimize {
    pointer-events: auto;
  }
}
</style>
