import uid from '@/helpers/uid';

export default function (pbehavior = {}) {
  const comments = pbehavior.comments || [];

  return comments.map(comment => ({
    ...comment,

    key: uid(),
  }));
}
