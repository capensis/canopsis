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
    v-layout
      v-card
        v-card-title Login footer
        v-card-text
          text-editor(:value="frontendServiceItem.login_footer")
          v-btn(@click="updateLoginFooter")

</template>

<script>
import { GROUPS_NAVIGATION_TYPES } from '@/constants';

import appMixin from '@/mixins/app';
import i18nMixin from '@/mixins/i18n';

import entitiesFrontendServiceMixin from '@/mixins/entities/frontend-service';

import TextEditor from '@/components/other/text-editor/text-editor.vue';

export default {
  components: { TextEditor },
  mixins: [appMixin, i18nMixin, entitiesFrontendServiceMixin],
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
    updateLoginFooter() {
      this.updateFrontendService({
        data: {
          ...this.frontendServiceItem,

        },
      });
    },
  },
};
</script>
