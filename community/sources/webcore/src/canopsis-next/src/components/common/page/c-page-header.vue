<template lang="pug">
  div
    h2.text-xs-center.display-1.font-weight-medium.mt-3.mb-2
      slot {{ $t(`pageHeaders.${name}.title`) }}
      v-btn.mr-0(v-if="hasMessage", icon, @click="toggleMessageVisibility")
        v-icon(color="info") help_outline
    v-expand-transition
      div(v-if="hasMessage && shownMessage")
        v-layout.pb-2(justify-center)
          div.subheading.page-header__message.pre-wrap(v-html="message")
        v-layout.pb-2(v-show="!messageWasHidden", justify-center)
          v-btn(color="primary", @click="hideMessage") {{ $t('pageHeaders.hideMessage') }}
</template>

<script>
import { get, isFunction } from 'lodash';

import { DOCUMENTATION_BASE_URL } from '@/config';
import { DOCUMENTATION_LINKS } from '@/constants';

import { removeTrailingSlashes } from '@/helpers/url';

import { tourBaseMixin } from '@/mixins/tour/base';

export default {
  mixins: [tourBaseMixin],
  props: {
    name: {
      type: String,
      default() {
        const name = get(this.$route, 'meta.requiresPermission.id');

        return isFunction(name)
          ? name(this.$route)
          : name;
      },
    },
  },
  data() {
    return {
      shownMessage: false,
    };
  },
  computed: {
    hasMessage() {
      return this.$te(`pageHeaders.${this.name}.message`);
    },

    messageWasHidden() {
      return !!get(this.currentUser, ['ui_tours', this.name]);
    },

    learMoreMessage() {
      if (!DOCUMENTATION_LINKS[this.name]) {
        return '';
      }

      const link = removeTrailingSlashes(`${DOCUMENTATION_BASE_URL}${DOCUMENTATION_LINKS[this.name]}`);
      const linkMessage = `<a href="${link}" target="_blank"><strong>${link}</strong></a>`;

      return this.$t('pageHeaders.learnMore', { link: linkMessage });
    },

    message() {
      const message = this.$t(`pageHeaders.${this.name}.message`);

      return this.learMoreMessage ? `${message}\n${this.learMoreMessage}` : message;
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
        await this.finishTourByName(this.name);
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
