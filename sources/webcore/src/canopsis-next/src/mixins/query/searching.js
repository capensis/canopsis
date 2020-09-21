import { omit } from 'lodash';

export default {
  data() {
    return {
      searchingText: '',
    };
  },
  methods: {
    submitSearch() {
      this.$emit('update:query', {
        ...this.query,

        page: 1,
        search: this.searchingText,
      });
    },

    clearSearch() {
      this.searchingText = '';

      this.$emit('update:query', omit(this.query, ['search']));
    },
  },
};
