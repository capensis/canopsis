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
      |  {{ meta.last || lastItem }} {{ $t('common.of') }} {{ meta.total }} {{ $t('common.entries') }}
      v-pagination(v-model="currentPage", :length="totalPages")
</template>

<script>
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
    first: {
      type: Number,
      default: () => 0,
    },
    last: {
      type: Number,
      default: () => 0,
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
    lastItem() {
      if (this.last > this.meta.total) {
        return this.meta.total;
      }
      return this.last;
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
