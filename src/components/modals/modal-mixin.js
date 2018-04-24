import { mapGetters } from 'vuex';

export default {
  data() {
    return {
      opened: false,
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
  },

  watch: {
    modalComponent(value) {
      this.opened = value === this.$options.name;
    },
  },
};
