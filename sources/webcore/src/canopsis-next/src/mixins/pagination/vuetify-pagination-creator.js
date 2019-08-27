import vuetifyPagination from './vuetify-pagination';

export default function (itemsKey) {
  return {
    mixins: [vuetifyPagination],
    watch: {
      [itemsKey](value) {
        this.$set(this.pagination, 'totalItems', value.length);
      },
    },
  };
}
