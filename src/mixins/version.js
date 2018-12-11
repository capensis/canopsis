import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('version');

export default {
  computed: {
    ...mapGetters({
      getVersion: 'version',
    }),
  },
  methods: {
    ...mapActions(['fetchVersion']),
  },
};
