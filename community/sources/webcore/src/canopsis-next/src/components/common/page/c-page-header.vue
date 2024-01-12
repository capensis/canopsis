<template>
  <div>
    <h2 class="text-center text-h4 font-weight-medium mt-4 mb-2">
      <slot>{{ $t(`pageHeaders.${name}.title`) }}</slot>
      <v-btn
        class="ml-2 my-2"
        v-if="hasMessage"
        icon
        @click="toggleMessageVisibility"
      >
        <v-icon color="info">
          help_outline
        </v-icon>
      </v-btn>
    </h2>
    <v-expand-transition>
      <div v-show="hasMessage && shownMessage">
        <v-layout
          class="pb-2"
          justify-center
        >
          <c-compiled-template
            class="text-subtitle-1 page-header__message pre-wrap"
            :template="message"
          />
        </v-layout>
        <v-layout
          class="pb-2"
          v-if="!messageWasHidden"
          justify-center
        >
          <v-btn
            class="my-2"
            :loading="isHidePending"
            color="primary"
            @click="hideMessage"
          >
            {{ $t('pageHeaders.hideMessage') }}
          </v-btn>
        </v-layout>
      </div>
    </v-expand-transition>
  </div>
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
      isHidePending: false,
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
      this.isHidePending = true;

      if (!this.messageWasHidden) {
        await this.finishTourByName(this.name);
      }

      this.isHidePending = false;
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
