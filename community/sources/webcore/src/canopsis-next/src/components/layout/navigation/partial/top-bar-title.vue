<template lang="pug">
  v-toolbar-title.white--text.font-weight-regular.top-bar-title
    span {{ title }}
    div.badge-wrapper(v-if="showBadge")
      v-tooltip(right)
        template(#activator="{ on, attrs }")
          v-btn.badge-button(
            v-on="on",
            v-bind="attrs",
            color="error",
            icon,
            small,
            @click="showInfoModal"
          )
            v-icon(color="white", size="12px") priority_high
        span {{ $t('modals.webSocketError.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';

import Socket from '@/plugins/socket/services/socket';

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
  beforeDestroy() {
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
  position: relative;
  overflow: visible;

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
}
</style>
