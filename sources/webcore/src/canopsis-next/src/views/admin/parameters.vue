<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-list
      v-list-tile
        v-list-tile-title {{ $t('parameters.interfaceLanguage') }}
        v-list-tile-content
          v-select(:items="languageOptions", :value="$i18n.locale", @input="changeLocale")
      v-list-tile
        v-list-tile-title {{ $t('parameters.groupsNavigationType.title') }}
        v-list-tile-content
          v-select(
          :items="groupsNavigationOptions",
          :value="groupsNavigationType",
          @change="setGroupsNavigationType"
          )
</template>

<script>
import { GROUPS_NAVIGATION_TYPES } from '@/constants';

import appMixin from '@/mixins/app';
import i18nMixin from '@/mixins/i18n';
import authMixin from '@/mixins/auth';
import entitiesUserMixin from '@/mixins/entities/user';

export default {
  mixins: [appMixin, i18nMixin, authMixin, entitiesUserMixin],
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
    async changeLocale(locale) {
      this.setLocale(locale);

      const user = { ...this.currentUser, ui_language: locale };

      await this.createUser({ data: user });
      await this.fetchCurrentUser();
    },
  },
};
</script>
