<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-list
      v-list-tile
        v-list-tile-title {{ $t('parameters.interfaceLanguage') }}
        v-list-tile-content
          v-select(:items="languageOptions", :value="$i18n.locale", @input="setLocale")
      v-list-tile
        v-list-tile-title {{ $t('parameters.groupsNavigationType.title') }}
        v-list-tile-content
          v-select(
          :items="groupsNavigationOptions",
          :value="groupsNavigationType",
          @change="setGroupsNavigationType"
          )
      v-list-tile
        v-list-tile-title {{ $t('parameters.ldapAuthentication.title') }}
        v-list-tile-content
          v-btn(@click="showLDAPConfigModal", dark, color="secondary") {{ $t('parameters.configuration') }}
      v-list-tile
        v-list-tile-title {{ $t('parameters.casAuthentication.title') }}
        v-list-tile-content
          v-btn(@click="showCASConfigModal", dark, color="secondary") {{ $t('parameters.configuration') }}
</template>

<script>
import { GROUPS_NAVIGATION_TYPES, MODALS } from '@/constants';

import appMixin from '@/mixins/app';
import i18nMixin from '@/mixins/i18n';
import modalMixin from '@/mixins/modal';

export default {
  mixins: [appMixin, i18nMixin, modalMixin],
  data() {
    return {
      languageOptions: [
        {
          text: 'Français',
          value: 'fr',
        },
        {
          text: 'English',
          value: 'en',
        },
      ],
      groupsNavigationOptions: [
        {
          text: this.$t('parameters.groupsNavigationType.items.sideBar'),
          value: GROUPS_NAVIGATION_TYPES.sideBar,
        },
        {
          text: this.$t('parameters.groupsNavigationType.items.topBar'),
          value: GROUPS_NAVIGATION_TYPES.topBar,
        },
      ],
    };
  },
  methods: {
    showLDAPConfigModal() {
      this.showModal({
        name: MODALS.ldapConfiguration,
      });
    },
    showCASConfigModal() {
      this.showModal({
        name: MODALS.casConfiguration,
      });
    },
  },
};
</script>
