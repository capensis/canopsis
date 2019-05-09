import moment from 'moment';
import 'moment-duration-format';

export default function (value) {
  return moment.duration(value, 'seconds').format('HH[h] mm[m] ss[s]', { trim: 'both final' }) || 0;
}
