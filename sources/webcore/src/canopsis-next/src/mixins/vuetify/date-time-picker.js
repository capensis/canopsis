import formBaseMixin from '../form/base';
import { modelPropKeyComputed } from '../form/internal/computed-properties';

export default {
  mixins: [formBaseMixin],
  methods: {
    updateTime(time = '00:00:00') {
      const value = this[this[modelPropKeyComputed]];
      const newValue = new Date(value ? value.getTime() : null);
      const [hours = 0, minutes = 0, seconds = 0] = time.split(':');

      newValue.setHours(parseInt(hours, 10) || 0, parseInt(minutes, 10) || 0, parseInt(seconds, 10) || 0, 0);

      this.updateModel(newValue);
    },

    updateDate(date) {
      const value = this[this[modelPropKeyComputed]];
      const newValue = new Date(value ? value.getTime() : null);
      const [year, month, day] = date.split('-');

      newValue.setFullYear(parseInt(year, 10), parseInt(month, 10) - 1, parseInt(day, 10));

      if (!value) {
        newValue.setHours(0, 0, 0, 0);
      } else if (this.useSeconds) {
        newValue.setMilliseconds(0);
      } else {
        newValue.setSeconds(0, 0);
      }

      this.updateModel(newValue);
    },
  },
};
