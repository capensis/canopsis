import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('version');

export default {
  computed: {
    ...mapGetters(['version']),
  },
  methods: {
    ...mapActions(['fetchVersion']),
  },
};
