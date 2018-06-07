import omit from 'lodash/omit';

export function getQueryAlarm() {
  const query = omit(this.$route.query, ['page']);

  query.skip = ((this.$route.query.page - 1) * this.limit) || 0;

  return query;
}

export function getQueryContext() {
  const query = omit(this.$route.query, ['page']);

  query.limit = this.limit;

  return query;
}
