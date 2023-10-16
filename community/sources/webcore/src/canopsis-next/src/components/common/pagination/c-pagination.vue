<template>
  <div
    class="container text-center"
    v-if="total"
    :class="{ 'py-1': isTop }"
  >
    <ul
      class="v-pagination"
      v-if="isTop"
    >
      <li>
        <button
          class="v-pagination__navigation"
          data-test="paginationPreviewButton"
          :disabled="isPreviousPageDisabled"
          :class="{ 'v-pagination__navigation--disabled': isPreviousPageDisabled }"
          @click="previous"
        >
          <v-icon>chevron_left</v-icon>
        </button>
      </li>
      <div class="pagination-details">
        {{ page }} / {{ totalPages }}
      </div>
      <li>
        <button
          class="v-pagination__navigation"
          data-test="paginationNextButton"
          :disabled="isNextPageDisabled"
          :class="{ 'v-pagination__navigation--disabled': isNextPageDisabled }"
          @click="next"
        >
          <v-icon>chevron_right</v-icon>
        </button>
      </li>
    </ul>
    <v-layout
      v-else
      align-center
      justify-center
    >
      <span class="text--secondary">{{ $t('common.paginationItems', { first, last, total }) }}</span>
      <v-pagination
        data-test="vPagination"
        :value="page"
        :total-visible="$config.PAGINATION_TOTAL_VISIBLE"
        :length="totalPages"
        @input="updatePage"
      />
    </v-layout>
  </div>
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';

/**
 * Pagination component
 *
 * @prop {Object} page - Current page
 * @prop {Number} limit - Elements per page
 * @prop {Number} total - Total items
 * @prop {string} type - 'Top' or 'Bottom, to determine
 * if it's a top pagination (with less infos), or a bottom pagination
 *
 * @event input
 */
export default {
  model: {
    prop: 'page',
    event: 'input',
  },
  props: {
    page: {
      type: Number,
      default: 1,
    },
    limit: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
    total: {
      type: Number,
      default: 0,
    },
    type: {
      type: String,
      validator: value => ['top', 'bottom'].indexOf(value) !== -1,
      default: 'bottom',
    },
  },
  computed: {
    isTop() {
      return this.type === 'top';
    },

    totalPages() {
      return Math.ceil(this.total / this.limit) || 0;
    },

    /**
     * Calculate first item nb to display on pagination, in case it's not given by the backend
     */
    first() {
      return 1 + (this.limit * (this.page - 1));
    },

    /**
     * Calculate last item nb to display on pagination, in case it's not given by the backend
     */
    last() {
      const calculatedLast = this.page * this.limit;

      return calculatedLast > this.total ? this.total : calculatedLast;
    },

    isPreviousPageDisabled() {
      return this.page <= 1;
    },

    isNextPageDisabled() {
      return this.page >= this.totalPages;
    },
  },
  methods: {
    previous() {
      this.updatePage(this.page - 1);
    },

    next() {
      this.updatePage(this.page + 1);
    },

    updatePage(page) {
      this.$emit('input', page);
    },
  },
};
</script>

<style scoped>
  .pagination__navigation:disabled {
    opacity: .6;
    pointer-events: none;
  }

  .container {
    flex-direction: column;
  }

  span {
    padding-left: 1%;
  }
</style>
