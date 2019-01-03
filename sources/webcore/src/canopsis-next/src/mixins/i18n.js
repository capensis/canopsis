import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('i18n');

export default {
  methods: {
    ...mapActions(['setLocale']),
  },
};
