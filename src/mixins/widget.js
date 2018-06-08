import { createNamespacedHelpers } from 'vuex';
import viewMixin from './view';

const { mapGetters, mapActions } = createNamespacedHelpers('view/widget');

export default {
  mixins: [
    viewMixin,
  ],
  computed: {
    ...mapGetters([
      'getWidget',
    ]),
  },
  methods: {
    ...mapActions({
      saveWidget: 'save',
    }),
  },
};
