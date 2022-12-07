export const modalInnerItemsMixin = {
  computed: {
    /**
     * Function returns items by entity type and entity ids
     *
     * @returns {Array}
     */
    items() {
      return this.config.items ?? [];
    },
  },
};
