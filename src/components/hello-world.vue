<template lang="pug">
  div
    v-btn(color="error", @click="addErrorPopup({ text: '<h1>ERROR</h1>', autoClose: false })") Error
    v-btn(color="warning", @click="addWarningPopup({ text: 'Warning' })") Warning
    v-btn(color="info", @click="addInfoPopup({ text: 'Info', autoClose: 10000 })") Info
    v-btn(color="success", @click="addSuccessPopup({ text: 'Success' })") Success
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
import PopupMixin from '@/mixins/popup';

export default {
  name: 'HelloWorld',
  mixins: [PopupMixin],
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
