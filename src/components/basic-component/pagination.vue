<template lang="pug">
  div.container(v-if="meta.total")
    span {{ $t('common.showing') }} {{ meta.first }} {{ $t('common.to') }}
      |  {{ meta.last }} {{ $t('common.of') }} {{ meta.total }} {{ $t('common.entries') }}
    ul.pagination(v-if="type === 'top'")
      li
        button.pagination__navigation(:disabled="currentPage <= 1", @click="previous")
          v-icon chevron_left
      li
        button.pagination__navigation(:disabled="currentPage >= totalPages", @click="next")
          v-icon chevron_right
    v-pagination(v-else, v-model="currentPage", :length="totalPages")
</template>

<script>
export default {
  props: {
    type: {
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
