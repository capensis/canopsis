// LIBS
import { createNamespacedHelpers } from 'vuex';
// MIXINS
import viewMixin from '@/mixins/view';

const { mapGetters, mapActions } = createNamespacedHelpers('view/widget');

export default {
  mixins: [
    viewMixin,
  ],
  computed: {
    ...mapGetters({
      getWidget: 'getItem',
    }),
  },
  methods: {
    ...mapActions({
      saveWidget: 'saveItem',
    }),
  },
};
