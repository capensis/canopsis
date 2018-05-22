import omit from 'lodash/omit';

export default function getQuery() {
  const query = omit(this.$route.query, ['page']);

  query.limit = this.limit;
  query.skip = ((this.$route.query.page - 1) * this.limit) || 0;

  return query;
}

