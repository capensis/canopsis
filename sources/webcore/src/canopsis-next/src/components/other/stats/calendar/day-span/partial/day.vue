<script>
import { DsDay } from 'dayspan-vuetify/src/components';

export default {
  extends: DsDay,
  computed: {
    classesDay() {
      return {
        'ds-past': this.day.currentOffset < 0,
        'ds-today': this.day.currentDay,
        'ds-first-day-day': this.day.dayOfMonth === 1,
        'ds-out-calendar': !this.day.inCalendar,
      };
    },
  },
};
</script>
