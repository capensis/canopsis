<template lang="pug">
    v-toolbar(
      dense
      fixed
      clipped-left
      app
      color='blue darken-4'
    )
      v-toolbar-side-icon(
          class='white--text'
          @click.stop="handleClick"
        )
      v-spacer
      v-select(
        class='languageSelect'
        solo
        :items="locales"
        label="Select language"
        v-model="currentLocaleIndex"
        single-line
        dense
      )
</template>

<script>
export default {
  name: 'TopBar',
  data() {
    const locales = [
      { key: 'fr', text: 'Français' },
      { key: 'en', text: 'English' },
    ];

    return {
      currentLocaleIndex: locales.find(({ key }) => key === this.$i18n.locale),
      locales,
    };
  },
  methods: {
    handleClick() {
      this.$store.dispatch('app/toggleSideBar');
    },
  },
  watch: {
    currentLocaleIndex() {
      if (this.currentLocaleIndex) {
        this.$store.dispatch('i18n/setLocale', this.currentLocaleIndex.key);
      }
    },
  },
};
</script>

<style scoped>
  .languageSelect {
    max-width: 9em;
  }
</style>

