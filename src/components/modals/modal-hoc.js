import Vue from 'vue';
import { mapGetters } from 'vuex';

export default WrappedComponent => Vue.extend({
  ...WrappedComponent,

  data() {
    const data = WrappedComponent.data ? WrappedComponent.data() : {};

    return {
      opened: false,
      ...data,
    };
  },

  computed: {
    name() {
      return this.$options.name;
    },

    ...mapGetters('modal', {
      modalComponent: 'component',
      modalConfig: 'config',
    }),

    ...WrappedComponent.computed,
  },

  methods: {
    hideModal() {
      this.opened = false;

      setTimeout(() => {
        if (this.errors) {
          this.errors.clear();
        }

        if (this.modalComponent === this.$options.name) {
          this.$store.dispatch('modal/hide');
        }
      }, 300); // TODO: see it
    },

    ...WrappedComponent.methods,
  },

  watch: {
    modalComponent(value) {
      this.opened = value === this.$options.name;
    },

    ...WrappedComponent.watch,
  },
});
