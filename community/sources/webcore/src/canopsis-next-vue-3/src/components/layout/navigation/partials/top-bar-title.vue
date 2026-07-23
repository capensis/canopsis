<template>
  <v-toolbar-title class="text-white font-weight-regular top-bar-title">
    <c-compiled-template :template="title" parent-element="span" />

    <div v-if="showBadge" class="badge-wrapper">
      <v-tooltip location="right">
        <template #activator="{ props }">
          <v-btn
            class="badge-button"
            color="error"
            size="x-small"
            icon
            v-bind="props"
            @click="showInfoModal"
          >
            <v-icon color="white" size="12px"> priority_high </v-icon>
          </v-btn>
        </template>

        <span>{{ $t('modals.webSocketError.title') }}</span>
      </v-tooltip>
    </div>
  </v-toolbar-title>
</template>

<script>
import { MODALS } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import { entitiesInfoMixin } from '@/mixins/entities/info';

export default {
  mixins: [entitiesInfoMixin],

  props: {
    title: {
      type: String,
      default: '',
    },
  },

  data() {
    return {
      showBadge: false,
    };
  },

  created() {
    this.$socket.on(Socket.EVENTS_TYPES.networkError, this.socketNetworkErrorHandler);
  },

  beforeUnmount() {
    this.$socket.off(Socket.EVENTS_TYPES.networkError, this.socketNetworkErrorHandler);
  },

  methods: {
    socketNetworkErrorHandler() {
      this.showBadge = true;
    },

    showInfoModal() {
      this.$modals.show({
        name: MODALS.info,
        config: {
          title: this.$t('modals.webSocketError.title'),
          text: this.isProVersion
            ? this.$t('modals.webSocketError.text')
            : this.$t('modals.webSocketError.shortText'),
        },
      });
    },
  },
};
</script>

<style lang="scss">
.top-bar-title {
  flex: 0 0 auto;
  overflow: visible;
  position: relative;
  white-space: nowrap;

  .badge-wrapper {
    position: absolute;
    top: -7px;
    right: -17px;

    .badge-button {
      margin: 0;
      width: 16px;
      height: 16px;
    }
  }

  p {
    margin: 0;
    display: inline-flex;
  }
}
</style>
