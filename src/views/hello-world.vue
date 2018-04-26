<template lang="pug">
  div
    h1.hello {{ $t('common.hello') }}
    v-container(fluid)
      v-layout(row, wrap, v-for="item in items", :key="item._id")
        v-flex(xs6)
          v-subheader {{Object.keys(item)}}
        v-flex(xs6)
          v-subheader {{item.v.display_name}}
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  name: 'HelloWorld',
  mounted() {
    this.$store.dispatch('entities/alarm/fetchList');
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
  computed: {
    ...mapGetters('entities/alarm', [
      'items',
      'meta',
    ]),
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
