import { isEmpty, pick } from 'lodash';

import { PBEHAVIOR_ORIGINS } from '@/constants';

import { createDowntimePbehavior } from '@/helpers/entities/pbehavior/form';

import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';

export const widgetActionsPanelCommonMixin = {
  mixins: [entitiesPbehaviorMixin],
  methods: {
    showPbehaviorResponsePopups(response, remove = false) {
      if (response?.length) {
        let hasErrors;

        response.forEach(({ error, errors }) => {
          if (error || !isEmpty(errors)) {
            hasErrors = true;
            this.$popups.error({ text: error || Object.values(errors).join('\n') });
          }
        });

        if (!hasErrors) {
          this.$popups.success({ text: this.$t(`modals.createPbehavior.success.${remove ? 'remove' : 'create'}`) });
        }
      }
    },

    async createDowntimePbehavior(entities, payload) {
      const response = await this.createEntityPbehaviors({
        data: entities.map(entity => createDowntimePbehavior({
          entity,
          ...pick(payload, ['comment', 'reason', 'type', 'prefix', 'origin']),
        }), []),
      });

      this.showPbehaviorResponsePopups(response);
    },

    async removeDowntimePbehavior(entities, origin = PBEHAVIOR_ORIGINS.serviceWeather) {
      const response = await this.removeEntityPbehaviors({
        data: entities.map(({ _id: id }) => ({
          origin,

          entity: id,
        })),
      });

      this.showPbehaviorResponsePopups(response, true);
    },
  },
};
