<template>
  <div v-if="total" :class="isTop ? 'py-1' : 'pa-6'" class="container text-center">
    <ul v-if="isTop" class="v-pagination">
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
    <c-layout v-else align="center" justify="center">
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
    totalVisible: {
      type: Number,
      required: false,
    },
  },

  emits: ['update:page'],

  computed: {
    paginationTotalVisible() {
      return this.totalVisible ?? PAGINATION_TOTAL_VISIBLE;
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
      return 1 + this.limit * (this.page - 1);
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
  opacity: 0.6;
  pointer-events: none;
}

.container {
  flex-direction: column;
}

.v-pagination {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 0;
  margin: 0;
  list-style: none;
}

.v-pagination li {
  display: flex;
  list-style: none;
}

.v-pagination__navigation {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  color: rgba(var(--v-theme-on-surface), 0.87);
  background: transparent;
  border: 0;
  border-radius: 50%;
}

.v-pagination__navigation--disabled {
  opacity: 0.6;
  pointer-events: none;
}

.pagination-details {
  display: inline-flex;
  align-items: center;
}

span {
  padding-left: 1%;
}
</style>
