<script>
import { get } from 'lodash';
import { DsCalendarEventTime } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEventTime,
  computed: {
    hasPopover() {
      return !!this.$scopedSlots.eventPopover && get(this.calendarEvent, 'data.meta.hasPopover');
    },
  },
};
</script>
