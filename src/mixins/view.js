import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

export default {
  methods: {
    ...mapActions([
      'loadView',
    ]),
  },
  computed: {
    ...mapGetters([
      'loadedView',
    ]),
  },
};
