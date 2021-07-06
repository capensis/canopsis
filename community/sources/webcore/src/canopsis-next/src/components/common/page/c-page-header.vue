<template lang="pug">
  div
    h2.text-xs-center.display-1.font-weight-medium.mt-3.mb-2
      slot
      v-btn(icon, @click="toggleHelpMessage")
        v-icon(color="info") help_outline
    v-expand-transition
      div(v-if="shownHelpMessage")
        v-layout(justify-center)
          div.subheading.help-message.pre-wrap(v-html="$t('pageHeaders.helpMessages.idleRules')")
        v-layout.pb-2(justify-center)
          v-btn(color="primary", @click="hideHelpMessage") {{ $t('pageHeaders.hideHelpMessage') }}
</template>

<script>
import { get } from 'lodash';

import { tourBaseMixin } from '@/mixins/tour/base';

export default {
  mixins: [tourBaseMixin],
  props: {
    helpMessageName: {
      type: String,
      default: 'dsa',
    },
  },
  data() {
    return {
      shownHelpMessage: false,
    };
  },
  computed: {
    helpMessageWasHidden() {
      return !!get(this.currentUser, ['ui_tours', this.helpMessageName]);
    },
  },
  created() {
    if (!this.helpMessageWasHidden) {
      this.shownHelpMessage = true;
    }
  },
  methods: {
    toggleHelpMessage() {
      this.shownHelpMessage = !this.shownHelpMessage;
    },

    async hideHelpMessage() {
      if (!this.helpMessageWasHidden) {
        await this.finishTourByName(this.helpMessageName);
      }

      this.shownHelpMessage = false;
    },
  },
};
</script>

<style lang="scss" scoped>
$helpMessageMaxWidth: 1050px;

.help-message {
  max-width: $helpMessageMaxWidth;
  text-align: center;
}
</style>
