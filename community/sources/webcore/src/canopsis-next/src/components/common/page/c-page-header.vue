<template lang="pug">
  div
    h2.text-xs-center.display-1.font-weight-medium.mt-3.mb-2
      slot
      v-btn(icon, @click="toggleMessageVisibility")
        v-icon(color="info") help_outline
    v-expand-transition
      div(v-if="shownMessage")
        v-layout.pb-2(justify-center)
          div.subheading.page-header__message.pre-wrap(v-html="$t('pageHeaders.helpMessages.idleRules')")
        v-layout.pb-2(v-show="!messageWasHidden", justify-center)
          v-btn(color="primary", @click="hideMessage") {{ $t('pageHeaders.hideHelpMessage') }}
</template>

<script>
import { get } from 'lodash';

import { tourBaseMixin } from '@/mixins/tour/base';

export default {
  mixins: [tourBaseMixin],
  props: {
    messageName: {
      type: String,
      default: 'fsad',
    },
  },
  data() {
    return {
      shownMessage: false,
    };
  },
  computed: {
    messageWasHidden() {
      return !!get(this.currentUser, ['ui_tours', this.messageName]);
    },
  },
  created() {
    if (!this.messageWasHidden) {
      this.shownMessage = true;
    }
  },
  methods: {
    toggleMessageVisibility() {
      this.shownMessage = !this.shownMessage;
    },

    async hideMessage() {
      if (!this.messageWasHidden) {
        await this.finishTourByName(this.messageName);
      }

      this.shownMessage = false;
    },
  },
};
</script>

<style lang="scss" scoped>
$messageMaxWidth: 1050px;

.page-header__message {
  max-width: $messageMaxWidth;
  text-align: center;
}
</style>
