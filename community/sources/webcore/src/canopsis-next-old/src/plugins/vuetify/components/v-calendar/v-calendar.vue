<script>
import VCalendar from 'vuetify/lib/components/VCalendar/VCalendar';
import { isEventHiddenOn } from 'vuetify/lib/components/VCalendar/util/events';
import { getDayIdentifier } from 'vuetify/lib/components/VCalendar/util/timestamp';

const MINUTES_IN_DAY = 1440;

export default {
  extends: VCalendar,
  methods: {
    genTimedEvent({ event, left, width }, day) {
      if (day.timeDelta(event.end) < 0 || day.timeDelta(event.start) >= 1 || isEventHiddenOn(event, day)) {
        return false;
      }

      const dayIdentifier = getDayIdentifier(day);
      const start = event.startIdentifier >= dayIdentifier;
      const end = event.endIdentifier > dayIdentifier;
      const top = start ? day.timeToY(event.start) : 0;
      const bottom = end ? day.timeToY(MINUTES_IN_DAY) : day.timeToY(event.end);
      /* this.eventHeight was changed on 1, because in week and day mode min height was a eventHeight */
      const height = bottom - top;
      const scope = { eventParsed: event, day, start, end, timed: true };

      return this.genEvent(event, scope, true, {
        staticClass: 'v-event-timed',
        style: {
          top: `${top}px`,
          height: `${height}px`,
          left: `${left}%`,
          width: `${width}%`,
        },
      });
    },
  },
};
</script>

<style scoped>
  label {
    cursor: pointer;
  }
</style>
