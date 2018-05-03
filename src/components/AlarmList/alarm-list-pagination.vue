<template lang="pug">
  v-pagination( @input="changePage"
                :length="totalPages"
                v-model="currentPage" )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters } = createNamespacedHelpers('entities/alarm');

export default {
  name: 'alarm-list-pagination',
  data() {
    return {
      currentPage: 1,
    };
  },
  props: {
    itemsPerPage: {
      type: Number,
      required: true,
    },
  },
  computed: {
    ...mapGetters(['meta']),
    totalPages() {
      if (this.meta.total) {
        return Math.ceil(this.meta.total / this.itemsPerPage);
      }
      return 0;
    },
  },
  methods: {
    changePage() {
      this.$emit('changedPage', {
        params: {
          skip: (this.currentPage - 1) * this.itemsPerPage,
        },
      });
    },
  },
};
</script>

<style scoped>

</style>
