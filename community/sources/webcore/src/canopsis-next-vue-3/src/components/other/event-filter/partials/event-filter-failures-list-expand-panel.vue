<template>
  <v-card flat>
    <v-card-text>
      <c-layout align="center">
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
      </c-layout>
      <c-object-treeview
        :root-name="$t('eventFilter.event')"
        :value="failure.event"
        value-max-width="520px"
      />
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
