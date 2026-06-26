<template>
  <div
    v-if="total"
    :class="{ 'py-1': isTop }"
    class="container pa-6 text-center"
  >
    <ul
      v-if="isTop"
      class="v-pagination"
    >
      <li>
        <button
          :disabled="isPreviousPageDisabled"
          :class="{ 'v-pagination__navigation--disabled': isPreviousPageDisabled }"
          class="v-pagination__navigation"
          data-test="paginationPreviewButton"
          @click="previous"
        >
          <c-icon>chevron_left</c-icon>
        </button>
      </li>
      <div class="pagination-details">
        {{ page }} / {{ totalPages }}
      </div>
      <li>
        <button
          :disabled="isNextPageDisabled"
          :class="{ 'v-pagination__navigation--disabled': isNextPageDisabled }"
          class="v-pagination__navigation"
          data-test="paginationNextButton"
          @click="next"
        >
          <c-icon>chevron_right</c-icon>
        </button>
      </li>
    </ul>
    <c-layout
      v-else
      align="center"
      justify="center"
    >
      <span class="text-medium-emphasis">
        {{ $t('common.paginationItems', { first, last, total }) }}
      </span>
      <v-pagination
        :model-value="page"
        :total-visible="paginationTotalVisible"
        :length="totalPages"
        data-test="vPagination"
        @update:model-value="updatePage"
      />
    </c-layout>
  </div>
</template>

<script>
import { PAGINATION_LIMIT, PAGINATION_TOTAL_VISIBLE } from '@/config';

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

  emits: ['update:page'],

  computed: {
    paginationTotalVisible() {
      return PAGINATION_TOTAL_VISIBLE;
    },

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
      this.$emit('update:page', page);
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
