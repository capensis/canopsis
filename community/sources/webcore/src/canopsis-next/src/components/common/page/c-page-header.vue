<template>
  <div>
    <h2 class="text-center text-h4 font-weight-medium mt-4 mb-2">
      <slot>{{ $t(`pageHeaders.${pageName}.title`) }}</slot>
      <v-btn
        v-if="hasMessage"
        class="ml-2 my-2"
        icon
        @click="toggleMessageVisibility"
      >
        <v-icon color="info">
          help_outline
        </v-icon>
      </v-btn>
    </h2>
    <v-expand-transition>
      <div
        v-if="hasMessage"
        v-show="shownMessage"
      >
        <v-layout
          class="pb-2"
          justify-center
        >
          <c-compiled-template
            :template="message"
            class="text-subtitle-1 page-header__message pre-wrap"
          />
        </v-layout>
        <v-layout
          v-if="!messageWasHidden"
          class="pb-2"
          justify-center
        >
          <v-btn
            :loading="isHidePending"
            class="my-2"
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
import { isFunction } from 'lodash';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router/composables';

import { DOCUMENTATION_BASE_URL } from '@/config';
import { DOCUMENTATION_LINKS } from '@/constants';

import { removeTrailingSlashes } from '@/helpers/url';

import { useI18n } from '@/hooks/i18n';
import { useTourBase } from '@/hooks/tour/tour-base';

export default {
  props: {
    name: {
      type: String,
      default: undefined,
    },
  },
  setup(props) {
    const route = useRoute();
    const { t, te } = useI18n();
    const { currentUser, finishTourByName } = useTourBase();

    const shownMessage = ref(false);
    const isHidePending = ref(false);

    const pageName = computed(() => {
      if (props.name) {
        return props.name;
      }

      const id = route.meta?.requiresPermission?.id;

      return isFunction(id) ? id(route) : id;
    });

    const hasMessage = computed(() => {
      const key = pageName.value;

      return !key || te(`pageHeaders.${key}.message`);
    });

    const messageWasHidden = computed(() => !!currentUser.value?.ui_tours?.[pageName.value]);

    const learMoreMessage = computed(() => {
      const key = pageName.value;

      if (!key || !DOCUMENTATION_LINKS[key]) {
        return '';
      }

      const link = removeTrailingSlashes(`${DOCUMENTATION_BASE_URL}${DOCUMENTATION_LINKS[key]}`);
      const linkMessage = `<a href="${link}" target="_blank"><strong>${link}</strong></a>`;

      return t('pageHeaders.learnMore', { link: linkMessage });
    });

    const message = computed(() => {
      const key = pageName.value;

      const baseMessage = hasMessage.value && key
        ? t(`pageHeaders.${key}.message`)
        : '';

      return learMoreMessage.value ? `${baseMessage}\n${learMoreMessage.value}` : baseMessage;
    });

    if (!messageWasHidden.value) {
      shownMessage.value = true;
    }

    /**
     * Toggles the expanded/collapsed state of the page header help message block (help icon target).
     */
    const toggleMessageVisibility = () => shownMessage.value = !shownMessage.value;

    /**
     * Hides the message panel: if the user has not dismissed this page tour yet, persists completion
     * via `finishTourByName` for `pageName`, then clears the loading state and collapses the message.
     */
    const hideMessage = async () => {
      isHidePending.value = true;

      if (!messageWasHidden.value) {
        finishTourByName(pageName.value);
      }

      isHidePending.value = false;
      shownMessage.value = false;
    };

    return {
      pageName,
      hasMessage,
      shownMessage,
      isHidePending,
      messageWasHidden,
      message,
      toggleMessageVisibility,
      hideMessage,
    };
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
