import { sortBy } from 'lodash';

/**
 * Mixin for sharing logic of sorting groups headers
 */
export default {
  props: {
    roles: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    headers() {
      return [
        { text: '', sortable: false },

        ...this.sortedRoles.map(role => ({ text: role._id, sortable: false })),
      ];
    },

    sortedRoles() {
      return sortBy(this.roles, [({ _id: name }) => name.toLowerCase()]);
    },
  },
};
