<script>
import { DsDay } from 'dayspan-vuetify/src/components';

export default {
  extends: DsDay,
  methods: {
    isVisible(calendarEvent) {
      return !(
        this.$dayspan.features.hideOnMove
        && this.placeholder
        && this.placeholder.event.id === calendarEvent.event.id
      );
    },
  },
};
</script>
