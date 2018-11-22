<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-list
      v-list-tile
        v-list-tile-title {{ $t('parameters.interfaceLanguage') }}
        v-list-tile-content
          v-select(:items="languageOptions", :value="$i18n.locale", @input="setLocale")
</template>

<script>
import i18nMixin from '@/mixins/i18n';

export default {
  mixins: [i18nMixin],
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
    };
  },
};
</script>
