<template lang="pug">
  div
    h1.hello {{ $t('common.hello') }}
    h2 {{ msg }}
    v-container(fluid)
      v-layout(row, wrap)
        v-flex(xs6)
          v-subheader Example
        v-flex(xs6)
          v-select(
            :items="locales"
            label="Select language"
            v-model="currentLocaleIndex"
            single-line
          )
</template>

<script>
export default {
  name: 'HelloWorld',
  props: {
    msg: {
      type: String,
      required: true,
    },
  },
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
  watch: {
    currentLocaleIndex() {
      if (this.currentLocaleIndex) {
        this.$store.dispatch('i18n/setLocale', this.currentLocaleIndex.key);
      }
    },
  },
};
</script>
