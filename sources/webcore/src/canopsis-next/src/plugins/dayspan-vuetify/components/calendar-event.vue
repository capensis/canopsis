<script>
import { get } from 'lodash';
import { DsCalendarEvent } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEvent,
  props: {
    enableResize: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasPopover() {
      return !!this.$scopedSlots.eventPopover && get(this.calendarEvent, 'data.meta.hasPopover');
    },
  },
};
</script>
