<template lang="pug">
  v-layout(justify-space-between, align-center)
    v-flex.pa-2(v-for="(icon, index) in mainIcons", :key="index")
      v-icon(color="white", small) {{ icon }}
    v-flex.pl-1.white--text.subheading(xs12)
      v-layout(align-center)
        div.mr-1.entity-name(v-resize-text="{ maxFontSize: '16px' }") {{ entityName }}
        v-btn.mx-1.white(
          v-for="icon in extraIcons",
          :key="icon.icon",
          :color="icon.color",
          small,
          dark,
          icon
        )
          v-icon(small) {{ icon.icon }}
        c-no-events-icon(:value="entity.idle_since", color="white", top)
</template>

<script>
import { get } from 'lodash';

import {
  ENTITIES_STATUSES,
  EVENT_ENTITY_STYLE,
  EVENT_ENTITY_TYPES,
  WEATHER_ICONS,
} from '@/constants';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
    entityNameField: {
      type: String,
      default: 'name',
    },
    active: {
      type: Boolean,
      default: false,
    },
    paused: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    entityName() {
      return get({ entity: this.entity }, this.entityNameField, this.entityNameField);
    },

    mainIcons() {
      const mainIcons = [];

      if (!this.paused && !this.active) {
        mainIcons.push(WEATHER_ICONS[this.entity.icon]);
      }

      mainIcons.push(...this.entity.pbehaviors.map(({ type }) => type.icon_name));

      return mainIcons;
    },

    extraIcons() {
      const extraIcons = [];

      if (this.entity.ack) {
        extraIcons.push({
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.fastAck].icon,
          color: 'purple',
        });
      }

      if (this.entity.ticket) {
        extraIcons.push({
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          color: 'blue',
        });
      }

      if (this.entity.status && this.entity.status.val === ENTITIES_STATUSES.cancelled) {
        extraIcons.push({
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          color: 'grey darken-1',
        });
      }

      return extraIcons;
    },
  },
};
</script>

<style lang="scss" scoped>
.entity-name {
  line-height: 1.5em;
  word-break: break-all;
}
</style>
