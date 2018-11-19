<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-list
      v-list-tile
        v-list-tile-title {{ $t('parameters.interfaceLanguage') }}
        v-list-tile-content
          v-select(:items="languageOptions", v-model="selectedLanguage", @input="updateLocale")
</template>

<script>
export default {
  data() {
    return {
      selectedLanguage: this.$i18n.locale,
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
  methods: {
    updateLocale(event) {
      this.$i18n.locale = event;
    },
  },
};
</script>
