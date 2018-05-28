<template lang="pug">
  div.container(v-show="meta.total")
    span {{ $t('common.showing') }} {{ meta.first }} {{ $t('common.to') }}
         |  {{ meta.last }} {{ $t('common.of') }} {{ meta.total }} {{ $t('common.entries') }}
    v-pagination(v-model="currentPage", :length="totalPages")
</template>

<script>
export default {
  name: 'alarm-list-pagination',
  props: {
    meta: {
      type: Object,
      default() {
        return {
          total: 0,
          first: 0,
          last: 0,
        };
      },
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
};
</script>

<style scoped>
  .container{
   flex-direction: column;
  }
  span {
    padding-left: 1%;
  }
</style>
