<template>
  <v-layout
    justify-space-between
    align-center
  >
    <v-flex @click.stop="">
      <v-simple-checkbox
        :value="selected"
        :disabled="!selectable"
        class="service-entity-header__checkbox ma-0 pa-0"
        dark
        @input="$emit('update:selected', $event)"
      />
    </v-flex>
    <v-flex class="pa-2">
      <v-icon class="color--inherit" small>
        {{ entity.icon }}
      </v-icon>
    </v-flex>
    <v-flex
      class="pl-1 white--text text-subtitle-1"
      xs12
    >
      <v-layout align-center>
        <div class="mr-1 service-entity-header__name">
          {{ entityName }}
        </div>
        <v-btn
          v-for="icon in extraIcons"
          :key="icon.icon"
          :style="{ backgroundColor: icon.color }"
          class="service-entity-header__extra-icon mx-1"
          small
          dark
          icon
        >
          <v-icon small>
            {{ icon.icon }}
          </v-icon>
        </v-btn>
        <c-simple-tooltip
          v-for="link in entityLinks"
          :key="link.rule_id"
          :content="link.label"
          top
        >
          <template #activator="{ on }">
            <v-btn
              :style="{ backgroundColor: linkIconBackgroundColor }"
              :aria-label="link.label"
              class="service-entity-header__extra-icon mx-1"
              small
              dark
              icon
              v-on="on"
              @click.stop="linkClick(link)"
            >
              <v-icon small>
                {{ link.icon_name }}
              </v-icon>
            </v-btn>
          </template>
        </c-simple-tooltip>
        <c-no-events-icon
          :value="entity.idle_since"
          color="white"
          top
        />
        <div @click.stop="">
          <v-alert
            v-if="lastActionUnavailable"
            :value="lastActionUnavailable"
            class="service-entity-header__alert ma-0 px-2 py-1"
            color="black"
            dismissible
            @input="hideAlert"
          >
            {{ $t('serviceWeather.cannotBeApplied') }}
          </v-alert>
        </div>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';
import { get } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import {
  ALARM_STATUSES,
  BUSINESS_USER_PERMISSIONS_ACTIONS_MAP,
  EVENT_ENTITY_TYPES,
  LINK_RULE_ACTIONS,
  WEATHER_ACTIONS_TYPES,
} from '@/constants';

import { writeTextToClipboard } from '@/helpers/clipboard';
import { getEntityEventIcon } from '@/helpers/entities/entity/icons';
import { harmonizeLinks } from '@/helpers/entities/link/list';

import { useCurrentUserPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    lastActionUnavailable: {
      type: Boolean,
      default: false,
    },
    entityNameField: {
      type: String,
      default: 'name',
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const popups = usePopups();
    const { checkAccess } = useCurrentUserPermissions();

    const entityLinksPermission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.weather[
      WEATHER_ACTIONS_TYPES.entityLinks
    ];

    const linkIconBackgroundColor = CSS_COLORS_VARS.secondary;

    const entityName = computed(() => (
      get({ entity: props.entity }, props.entityNameField, props.entityNameField)
    ));

    const extraIcons = computed(() => {
      const icons = [];

      if (props.entity.ack) {
        icons.push({
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.fastAck),
          color: 'purple',
        });
      }

      if (props.entity.ticket) {
        icons.push({
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.assocTicket),
          color: 'blue',
        });
      }

      if (props.entity.status?.val === ALARM_STATUSES.cancelled) {
        icons.push({
          icon: getEntityEventIcon(EVENT_ENTITY_TYPES.delete),
          color: 'grey darken-1',
        });
      }

      if (props.entity.pbh_origin_icon) {
        icons.push({
          icon: props.entity.pbh_origin_icon,
          color: CSS_COLORS_VARS.secondary,
        });
      }

      return icons;
    });

    const entityLinks = computed(() => (
      checkAccess(entityLinksPermission)
        ? harmonizeLinks(props.entity.links)
        : []
    ));

    /**
     * Runs link rule action: clipboard copy with popup feedback, otherwise opens URL in a new tab.
     *
     * @param {Object} link
     */
    const linkClick = async (link) => {
      const action = link.action ?? LINK_RULE_ACTIONS.open;

      if (action === LINK_RULE_ACTIONS.copy) {
        try {
          await writeTextToClipboard(link.url);

          popups.success({ text: t('popups.copySuccess') });
        } catch (err) {
          console.error(err);

          popups.error({ text: t('popups.copyError') });
        }

        return;
      }

      window.open(link.url, '_blank');
    };

    /**
     * Notifies parent to dismiss the last-action-unavailable alert.
     */
    const hideAlert = () => emit('remove:unavailable');

    return {
      entityName,
      extraIcons,
      entityLinks,
      linkIconBackgroundColor,
      linkClick,
      hideAlert,
    };
  },
};
</script>

<style lang="scss">
.service-entity-header__name {
  line-height: 1.5em;
  word-break: break-all;
}
.service-entity-header__alert {
  border: none;
  background-color: rgba(255, 255, 255, 0.2) !important;
  border-radius: 5px;

  & ::v-deep .v-alert__dismissible .v-icon {
    margin-left: 0;
    font-size: 18px;
  }
}

.service-entity-header__extra-icon * {
  color: white !important;
}

.service-entity-header__checkbox i {
  color: inherit !important;
}
</style>
