<template>
  <v-card flat>
    <v-card-text>
      <v-layout align-center>
        <span>{{ $t('eventFilter.event') }}:</span>
        <c-copy-btn
          :value="eventString"
          :tooltip="$t('eventFilter.copyEventToClipboard')"
          left
          small
          icon-small
          @success="onSuccessCopied"
          @error="onErrorCopied"
        />
      </v-layout>
      <c-json-treeview :json="eventString" />
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  props: {
    failure: {
      type: Object,
      required: true,
    },
  },
  computed: {
    eventString() {
      return JSON.stringify(this.failure.event);
    },
  },
  methods: {
    onSuccessCopied() {
      this.$popups.success({ text: this.$t('eventFilter.eventCopied') });
    },

    onErrorCopied() {
      this.$popups.success({ text: this.$t('errors.default') });
    },
  },
};
</script>
