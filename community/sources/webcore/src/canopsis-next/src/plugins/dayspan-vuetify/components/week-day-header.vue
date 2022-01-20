<script>
import { DsWeekDayHeader } from 'dayspan-vuetify/src/components';

export default {
  extends: DsWeekDayHeader,
  methods: {
    isVisible(calendarEvent) {
      if (!calendarEvent.fullDay) {
        return false;
      }

      return !(
        this.$dayspan.features.hideOnMove
        && this.placeholder
        && this.placeholder.event.id === calendarEvent.event.id
      );
    },
  },
};
</script>
