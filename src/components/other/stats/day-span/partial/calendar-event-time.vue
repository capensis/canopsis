<script>
import { DsCalendarEventTime } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEventTime,
  computed: {
    hasPopover() {
      return !!this.$scopedSlots.eventPopover && this.calendarEvent.data.meta.hasPopover;
    },
  },
};
</script>
