import { omit } from 'lodash';

export default function (comments) {
  return comments.map(comment => omit(comment, ['key', 'ts']));
}
