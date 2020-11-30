<template lang="pug">
  div
    div(v-if="close")
      v-btn(
        icon,
        small,
        @click="$modals.hide({ id: modal.id })"
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
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...mapGetters(['hasMinimizedModal']),
  },
};
</script>

<style lang="scss" scoped>
.v-btn--minimize {
  pointer-events: auto;
}
</style>
