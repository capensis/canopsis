import uid from '@/helpers/uid';

export default function (pbehavior = {}) {
  const exdate = pbehavior.exdate || [];

  return exdate.map(unix => ({
    value: new Date(unix * 1000),
    key: uid(),
  }));
}
