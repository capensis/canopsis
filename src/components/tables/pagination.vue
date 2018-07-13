<template lang="pug">
  div.container.text-xs-center(v-if="meta.total")
    ul.pagination(v-if="type === 'top'")
      li
        button.pagination__navigation(:disabled="currentPage <= 1", @click="previous")
          v-icon chevron_left
      span {{ currentPage }}
      span /
      span {{ totalPages }}
      li
        button.pagination__navigation(:disabled="currentPage >= totalPages", @click="next")
          v-icon chevron_right
    div(v-else)
      span {{ $t('common.showing') }} {{ meta.first || first }} {{ $t('common.to') }}
      |  {{ meta.last || last }} {{ $t('common.of') }} {{ meta.total }} {{ $t('common.entries') }}
      v-pagination(v-model="currentPage", :length="totalPages")
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

/**
* Pagination component
*
* @prop {String} [type] - 'Top' or 'Bottom, to determine
* if it's a top pagination (with less infos), or a bottom pagination
* @prop {Object} [meta] - Object containing meta informations (Ex : total number of items)
* @prop {Number} [limit] - Number of items per pages
*/
export default {
  props: {
    type: {
      type: String,
      validator: value => ['top', 'bottom'].indexOf(value) !== -1,
      default: 'bottom',
    },
    meta: {
      type: Object,
      default: () => ({
        total: 0,
        first: 0,
        last: 0,
      }),
    },
    limit: {
      type: Number,
      required: true,
    },
  },
  computed: {
    currentPage: {
      get() {
        return parseInt(this.$route.query.page, 10) || 1;
      },
      set(page) {
        this.$router.push({
          query: {
            ...this.$route.query,
            page,
          },
        });
      },
    },
    totalPages() {
      if (this.meta.total) {
        return Math.ceil(this.meta.total / this.limit);
      }
      return 0;
    },
    /**
     * Calculate first item nb to display on pagination, in case it's not given by the backend
     */
    first() {
      const { page } = this.$route.query;
      if (page === 1 || !this.$route.query.page) {
        return 1;
      }
      if (this.$route.query.limit) {
        return 1 + (this.$route.query.limit * (page - 1));
      }

      return 1 + (PAGINATION_LIMIT * (page - 1));
    },
    /**
     * Calculate last item nb to display on pagination, in case it's not given by the backend
     */
    last() {
      let last;

      if (this.meta.last) {
        return this.meta.last;
      }

      if (this.$route.query.page === 1 || !this.$route.query.page) {
        last = this.$route.query.limit || PAGINATION_LIMIT;
      } else if (this.$route.query.limit) {
        last = this.$route.query.page * this.$route.query.limit;
      } else {
        last = this.$route.query.page * PAGINATION_LIMIT;
      }

      if (last > this.meta.total) {
        return this.meta.total;
      }

      return last;
    },
  },
  methods: {
    previous() {
      this.currentPage = this.currentPage - 1;
    },
    next() {
      this.currentPage = this.currentPage + 1;
    },
  },
};
</script>

<style scoped>
  .pagination__navigation:disabled {
    opacity: .6;
    pointer-events: none;
  }

  .container{
    flex-direction: column;
  }
  span {
    padding-left: 1%;
  }
</style>
