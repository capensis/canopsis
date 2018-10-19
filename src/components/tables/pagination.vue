<template lang="pug">
  div.container.text-xs-center(v-if="meta.total")
    ul.v-pagination(v-if="type === 'top'")
      li
        button.v-pagination__navigation(
        :disabled="isPreviousPageDisabled",
        :class="{ 'v-pagination__navigation--disabled': isPreviousPageDisabled }",
        @click="previous"
        )
          v-icon chevron_left
      span {{ currentPage }}
      span /
      span {{ totalPages }}
      li
        button.v-pagination__navigation(
        :disabled="currentPage >= totalPages",
        :class="{ 'v-pagination__navigation--disabled': isNextPageDisabled }",
        @click="next"
        )
          v-icon chevron_right
    div(v-else)
      span {{ $t('common.showing') }} {{ first }} {{ $t('common.to') }}
      |  {{ last }} {{ $t('common.of') }} {{ meta.total }} {{ $t('common.entries') }}
      v-pagination(v-model="currentPage", :total-visible="totalVisible" :length="totalPages")
</template>

<script>
import { PAGINATION_TOTAL_VISIBLE } from '@/config';

/**
 * Pagination component
 *
 * @prop {String} type - 'Top' or 'Bottom, to determine
 * if it's a top pagination (with less infos), or a bottom pagination
 * @prop {Object} meta - Object containing meta information (Ex : total number of items)
 * @prop {Number} query - Object containing widget query information
 *
 * @event query#update
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
    query: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      totalVisible: PAGINATION_TOTAL_VISIBLE,
    };
  },
  computed: {
    currentPage: {
      get() {
        return this.query.page || 1;
      },
      set(page) {
        this.$emit('update:query', { ...this.query, page });
      },
    },

    totalPages() {
      if (this.meta.total) {
        return Math.ceil(this.meta.total / this.query.limit);
      }

      return 0;
    },

    /**
     * Calculate first item nb to display on pagination, in case it's not given by the backend
     */
    first() {
      if (this.meta.first) {
        return this.meta.first;
      }

      return 1 + (this.query.limit * (this.query.page - 1));
    },

    /**
     * Calculate last item nb to display on pagination, in case it's not given by the backend
     */
    last() {
      if (this.meta.last) {
        return this.meta.last;
      }

      const calculatedLast = this.query.page * this.query.limit;

      return calculatedLast > this.meta.total ? this.meta.total : calculatedLast;
    },

    isPreviousPageDisabled() {
      return this.currentPage <= 1;
    },

    isNextPageDisabled() {
      return this.currentPage >= this.totalPages;
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
