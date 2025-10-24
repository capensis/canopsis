<template>
  <div
    :class="{ 'close': close, 'minimize': minimize }"
    class="modal-title-buttons"
  >
    <div
      v-if="minimize"
      class="modal-title-button__wrapper"
    >
      <v-tooltip
        v-if="!$modal.minimized"
        left
      >
        <template #activator="{ on }">
          <v-btn
            class="v-btn--minimize my-0"
            icon
            v-on="on"
            @click="$modals.minimize({ id: $modal.id })"
          >
            <v-icon
              color="white"
              large
            >
              minimize
            </v-icon>
          </v-btn>
        </template>
        <span>{{ $t('modals.common.titleButtons.minimizeTooltip') }}</span>
      </v-tooltip>
      <v-btn
        v-else
        class="v-btn-legacy-m--x"
        icon
        small
        @click="$modals.maximize({ id: $modal.id })"
      >
        <v-icon color="white">
          maximize
        </v-icon>
      </v-btn>
    </div>
    <div
      v-if="close"
      class="modal-title-button__wrapper"
    >
      <v-btn
        :small="$modal.minimized"
        icon
        @click="closeHandler"
      >
        <v-icon
          :large="!$modal.minimized"
          color="white"
        >
          close
        </v-icon>
      </v-btn>
    </div>
  </div>
</template>

<script>
import { isFunction } from 'lodash';

export default {
  inject: {
    $modal: {
      default() {
        return {};
      },
    },
    $closeModal: {
      default() {
        return () => this.$modals.hide({ id: this.$modal.id });
      },
    },
  },
  props: {
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
    closeHandler() {
      if (isFunction(this.close)) {
        return this.close;
      }

      return this.$closeModal;
    },
  },
};
</script>

<style lang="scss" scoped>
.modal-title-buttons {
  display: flex;

  .v-btn--minimize {
    pointer-events: auto;
  }
}
</style>
