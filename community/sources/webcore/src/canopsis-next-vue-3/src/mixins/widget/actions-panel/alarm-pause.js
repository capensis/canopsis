import { ALARM_LIST_ACTIONS_TYPES, ALARM_LIST_PAUSE_METHODS } from '@/constants';

import {
  entityPauseFormToRequest,
  entityResumeFormToRequest,
} from '@/helpers/entities/entity/pause';
import { getAlarmListPauseMode } from '@/helpers/entities/widget/forms/alarm';

import { entityPauseMixin } from '@/mixins/entities/entity-pause';

export const widgetActionsPanelAlarmPauseMixin = {
  mixins: [entityPauseMixin],

  computed: {
    isPauseActionMode() {
      return getAlarmListPauseMode(this.widget.parameters) === ALARM_LIST_PAUSE_METHODS.pause;
    },
  },

  methods: {
    getFastPbehaviorAction(pbehaviorsParameters, method, { keepMenuOpen = false } = {}) {
      const action = {
        type: ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd,
        title: this.$t('alarm.actions.titles.fastPbehaviorAdd'),
      };

      if (pbehaviorsParameters.length > 1) {
        action.items = pbehaviorsParameters.map(parameters => ({
          title: parameters.name_prefix,
          ...(keepMenuOpen && { closeOnClick: false }),
          method: () => method(parameters),
        }));
      } else {
        action.method = () => method(pbehaviorsParameters[0]);
      }

      return action;
    },

    getEntityPauseAction(pauseTypes, method) {
      const action = {
        type: ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd,
        icon: 'pause',
        title: this.$t('alarm.actions.titles.fastPbehaviorAdd'),
      };

      if (pauseTypes.length > 1) {
        action.items = pauseTypes.map(pauseType => ({
          title: pauseType.name,
          method: () => method(pauseType),
        }));
      } else {
        action.method = () => method(pauseTypes[0]);
      }

      return action;
    },

    getEntityResumeAction(method) {
      return {
        type: ALARM_LIST_ACTIONS_TYPES.fastPbehaviorRemove,
        icon: 'play_arrow',
        title: this.$t('alarm.actions.titles.fastPbehaviorRemove'),
        method,
      };
    },

    async pauseAlarmEntities(entities, pauseParameters = {}, { bulk = false } = {}) {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd, true);

      try {
        const pause = entityPauseFormToRequest(pauseParameters);
        const data = entities.map(({ _id: id }) => ({ ...pause, _id: id }));

        if (!bulk && data.length === 1) {
          await this.pauseEntity({ data: data[0] });
        } else {
          await this.bulkPauseEntities({ data });
        }

        this.$popups.success({ text: this.$t('entity.popups.pauseOperation') });
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
        this.$popups.error({ text: this.$t('alarm.popups.pauseFailed') });
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd, false);
      }
    },

    async resumeAlarmEntities(entities, { bulk = false } = {}) {
      this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorRemove, true);

      try {
        const resume = entityResumeFormToRequest();
        const data = entities.map(({ _id: id }) => ({ ...resume, _id: id }));

        if (!bulk && data.length === 1) {
          await this.resumeEntity({ data: data[0] });
        } else {
          await this.bulkResumeEntities({ data });
        }

        this.$popups.success({ text: this.$t('entity.popups.resumeOperation') });
        await this.afterSubmit();
      } catch (err) {
        console.error(err);
        this.$popups.error({ text: this.$t('alarm.popups.unpauseFailed') });
      } finally {
        this.setActionPending(ALARM_LIST_ACTIONS_TYPES.fastPbehaviorRemove, false);
      }
    },
  },
};
