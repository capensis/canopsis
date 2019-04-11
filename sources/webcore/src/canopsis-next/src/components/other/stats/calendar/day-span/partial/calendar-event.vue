<script>
import { get } from 'lodash';
import { DsCalendarEvent } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEvent,
  computed: {
    hasPopover() {
      return !!this.$scopedSlots.eventPopover && get(this.calendarEvent, 'data.meta.hasPopover');
    },
  },
};
</script>
